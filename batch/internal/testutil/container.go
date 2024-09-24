package testutil

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func ProjectRoot() string {
	currentDir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		_, err := os.ReadFile(filepath.Join(currentDir, "go.mod"))
		if os.IsNotExist(err) {
			if currentDir == filepath.Dir(currentDir) {
				return ""
			}
			currentDir = filepath.Dir(currentDir)
			continue
		} else if err != nil {
			return ""
		}
		break
	}
	return currentDir
}

func CreateDBContainer() (testcontainers.Container, string, func() error, error) {
	ctx := context.Background()

	stop := func() error {
		return nil
	}

	schema, err := filepath.Abs(filepath.Join(ProjectRoot(), "..", "db", "schema.sql"))
	if err != nil {
		return nil, "", stop, err
	}
	r, err := os.Open(schema)
	if err != nil {
		return nil, "", stop, err
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
		return nil, "", stop, err
	}
	stop = func() error {
		return container.Terminate(ctx)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, "", stop, err
	}
	port, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return nil, "", stop, err
	}

	dsn := fmt.Sprintf(
		"postgres://atcodersearch:atcodersearch@%s:%d/atcodersearch?sslmode=disable",
		host,
		port.Int(),
	)
	return container, dsn, stop, nil
}
