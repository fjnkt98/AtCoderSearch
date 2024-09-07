package repository

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setup() (testcontainers.Container, *pgxpool.Pool, func() error, error) {
	ctx := context.Background()

	stop := func() error {
		return nil
	}

	schema, err := filepath.Abs(filepath.Join("..", "..", "db", "schema.sql"))
	if err != nil {
		return nil, nil, stop, err
	}
	r, err := os.Open(schema)
	if err != nil {
		return nil, nil, stop, err
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "postgres:16-bullseye",
			Env: map[string]string{
				"POSTGRES_PASSWORD":         "atcodersearch",
				"POSTGRES_USER":             "atcodersearch",
				"POSTGRES_DB":               "atcodersearch",
				"POSTGRES_HOST_AUTH_METHOD": "password",
				"TZ":                        "Asia/Tokyo",
			},
			Files: []testcontainers.ContainerFile{
				{
					Reader:            r,
					ContainerFilePath: "/docker-entrypoint-initdb.d/schema.sql",
					FileMode:          0o666,
				},
			},
			ExposedPorts: []string{"5432/tcp"},
			WaitingFor:   wait.ForListeningPort("5432/tcp"),
		},
		Started: true,
	})
	if err != nil {
		return nil, nil, stop, err
	}
	stop = func() error {
		return container.Terminate(ctx)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, nil, stop, err
	}
	port, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return nil, nil, stop, err
	}

	dsn := fmt.Sprintf(
		"postgres://atcodersearch:atcodersearch@%s:%d/atcodersearch?sslmode=disable",
		host,
		port.Int(),
	)
	pool, err := NewPool(ctx, dsn)
	if err != nil {
		return nil, nil, stop, err
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, nil, stop, err
	}

	return container, pool, stop, nil
}

func TestCreateAndUpdateBatchHistory(t *testing.T) {
	_, pool, stop, err := setup()
	defer stop()

	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	t.Run("create batch history", func(t *testing.T) {
		history, err := NewBatchHistory(ctx, pool, "test", nil)
		if err != nil {
			t.Fatal(err)
		}

		if history.Name != "test" {
			t.Errorf("history.Name = %s, want test", history.Name)
		}
		if history.Status != "working" {
			t.Errorf("history.Status = %s, want working", history.Status)
		}
		if history.FinishedAt != nil {
			t.Errorf("history.FinishedAt = %v, want nil", history.FinishedAt)
		}
	})

	t.Run("complete batch history", func(t *testing.T) {
		history, err := NewBatchHistory(ctx, pool, "test", nil)
		if err != nil {
			t.Fatal(err)
		}

		if err := history.Complete(ctx, pool); err != nil {
			t.Error(err)
		}
		if history.Status != "completed" {
			t.Errorf("history.Status = %s, want completed", history.Status)
		}
		if history.FinishedAt == nil {
			t.Errorf("history.FinishedAt must be registered")
		}
	})

	t.Run("abort batch history", func(t *testing.T) {
		history, err := NewBatchHistory(ctx, pool, "test", nil)
		if err != nil {
			t.Fatal(err)
		}

		if err := history.Abort(ctx, pool); err != nil {
			t.Error(err)
		}
		if history.Status != "aborted" {
			t.Errorf("history.Status = %s, want aborted", history.Status)
		}
		if history.FinishedAt == nil {
			t.Errorf("history.FinishedAt must be registered")
		}
	})

	t.Run("abort completed batch history", func(t *testing.T) {
		history, err := NewBatchHistory(ctx, pool, "test", nil)
		if err != nil {
			t.Fatal(err)
		}

		if err := history.Complete(ctx, pool); err != nil {
			t.Error(err)
		}

		if err := history.Abort(ctx, pool); err != ErrHistoryConfirmed {
			t.Errorf("err = %v, want ErrHistoryConfirmed", err)
		}
	})
}
