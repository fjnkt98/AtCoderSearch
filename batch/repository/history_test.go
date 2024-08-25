package repository

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestBatchHistory(t *testing.T) {
	ctx := context.Background()

	schemaFile, err := filepath.Abs(filepath.Join("..", "..", "db", "schema.sql"))
	if err != nil {
		t.Fatalf("failed to specify schema.sql file: %s", err)
	}
	r, err := os.Open(schemaFile)
	if err != nil {
		t.Fatalf("failed to open schema.sql file: %s", err)
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
		t.Fatalf("couldn't create container: %s", err)
	}

	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("couldn't terminate container: %s", err)
		}
	}()

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("couldn't get container host: %s", err)
	}
	port, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		t.Fatalf("couldn't get container port: %s", err)
	}

	dsn := fmt.Sprintf(
		"postgres://atcodersearch:atcodersearch@%s:%d/atcodersearch?sslmode=disable",
		host,
		port.Int(),
	)
	pool, err := NewPool(ctx, dsn)
	if err != nil {
		t.Fatalf("failed to create connection pool: %s", err)
	}
	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("failed to verify a connection: %s", err)
	}

	t.Run("create and finish batch history", func(t *testing.T) {
		h, err := NewBatchHistory(ctx, pool, "TestBatch", nil)
		if err != nil {
			t.Fatalf("failed to create batch history: %s", err)
		}

		if h.ID != 1 {
			t.Errorf("expected id 1, but got %d", h.ID)
		}
		if h.Name != "TestBatch" {
			t.Errorf("expected id `TestBatch`, but got %s", h.Name)
		}
		if h.Status != "working" {
			t.Errorf("expected id `working`, but got %s", h.Name)
		}

		if err := h.Finish(ctx, pool); err != nil {
			t.Errorf("%s", err)
		}

		if h.Status != "finished" {
			t.Errorf("history status not changed")
		}
	})

	t.Run("create and fail batch history", func(t *testing.T) {
		h, err := NewBatchHistory(ctx, pool, "TestBatch", nil)
		if err != nil {
			t.Fatalf("failed to create batch history: %s", err)
		}

		if h.ID != 2 {
			t.Errorf("expected id 2, but got %d", h.ID)
		}
		if h.Name != "TestBatch" {
			t.Errorf("expected id `TestBatch`, but got %s", h.Name)
		}
		if h.Status != "working" {
			t.Errorf("expected id `working`, but got %s", h.Name)
		}

		if err := h.Fail(ctx, pool); err != nil {
			t.Errorf("%s", err)
		}

		if h.Status != "failed" {
			t.Errorf("history status not changed")
		}
	})
}

func TestOverwrite(t *testing.T) {
	startedAt := time.Date(2024, 6, 1, 0, 0, 0, 0, time.Local)
	finishedAt := time.Date(2024, 6, 1, 3, 0, 0, 0, time.Local)

	h := &BatchHistory{
		ID:         1,
		Name:       "Test",
		StartedAt:  startedAt,
		FinishedAt: nil,
		Status:     "working",
		Options:    nil,
	}

	r := BatchHistory{
		ID:         1,
		Name:       "Test",
		StartedAt:  startedAt,
		FinishedAt: &finishedAt,
		Status:     "finish",
		Options:    nil,
	}

	expected := &BatchHistory{
		ID:         1,
		Name:       "Test",
		StartedAt:  startedAt,
		FinishedAt: &finishedAt,
		Status:     "finish",
		Options:    nil,
	}

	h.overwrite(r)
	if !reflect.DeepEqual(h, expected) {
		t.Errorf("expected \n%+v\n, but got \n%+v\n", expected, h)
	}
}

func TestFailShouldBeSkip(t *testing.T) {
	startedAt := time.Date(2024, 6, 1, 0, 0, 0, 0, time.Local)
	finishedAt := time.Date(2024, 6, 1, 3, 0, 0, 0, time.Local)

	h := &BatchHistory{
		ID:         1,
		Name:       "Test",
		StartedAt:  startedAt,
		FinishedAt: &finishedAt,
		Status:     "finish",
		Options:    nil,
	}
	if err := h.Fail(context.Background(), nil); !errors.Is(err, ErrHistoryConfirmed) {
		t.Fatalf("test failed: %s", err.Error())
	}

	h = &BatchHistory{
		ID:         1,
		Name:       "Test",
		StartedAt:  startedAt,
		FinishedAt: &finishedAt,
		Status:     "canceled",
		Options:    nil,
	}
	if err := h.Fail(context.Background(), nil); !errors.Is(err, ErrHistoryConfirmed) {
		t.Fatalf("test failed: %s", err.Error())
	}
}
