package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestNewPool(t *testing.T) {
	ctx := context.Background()

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

	t.Run("success", func(t *testing.T) {
		dsn := fmt.Sprintf(
			"postgres://atcodersearch:atcodersearch@%s:%d/atcodersearch?sslmode=disable",
			host,
			port.Int(),
		)
		pool, err := NewPool(ctx, dsn)
		if err != nil {
			t.Errorf("couldn't create connection pool: %s", err)
		}
		if err := pool.Ping(ctx); err != nil {
			t.Errorf("ping failed: %s", err)
		}
	})
}
