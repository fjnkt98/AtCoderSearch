package testutil

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func GetProjectRoot() string {
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

func CreateDBContainer(files ...string) (container testcontainers.Container, dsn string, stop func() error, err error) {
	ctx := context.Background()

	stop = func() error {
		return nil
	}

	root := GetProjectRoot()
	schema, err := filepath.Abs(filepath.Join(root, "..", "db", "schema.sql"))
	if err != nil {
		return
	}
	r, err := os.Open(schema)
	if err != nil {
		return
	}
	defer r.Close()

	mounts := []testcontainers.ContainerFile{
		{
			Reader:            r,
			ContainerFilePath: "/docker-entrypoint-initdb.d/001_schema.sql",
			FileMode:          0o644,
		},
	}

	for i, file := range files {
		r, err = os.Open(file)
		if err != nil {
			return
		}
		defer r.Close()

		mounts = append(mounts, testcontainers.ContainerFile{
			Reader:            r,
			ContainerFilePath: fmt.Sprintf("/docker-entrypoint-initdb.d/%3d_%s.sql", i+2, filepath.Base(file)),
			FileMode:          0o644,
		})
	}

	container, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "postgres:16-bullseye",
			Env: map[string]string{
				"POSTGRES_PASSWORD":         "atcodersearch",
				"POSTGRES_USER":             "atcodersearch",
				"POSTGRES_DB":               "atcodersearch",
				"POSTGRES_HOST_AUTH_METHOD": "password",
				"TZ":                        "Asia/Tokyo",
			},
			Files:        mounts,
			ExposedPorts: []string{"5432/tcp"},
			WaitingFor:   wait.ForListeningPort("5432/tcp"),
		},
		Started: true,
	})
	if err != nil {
		return
	}
	stop = func() error {
		return container.Terminate(ctx)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return
	}
	port, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return
	}

	dsn = fmt.Sprintf(
		"postgres://atcodersearch:atcodersearch@%s:%d/atcodersearch?sslmode=disable",
		host,
		port.Int(),
	)
	return
}

func CreateEngineContainer() (container testcontainers.Container, url string, key string, stop func() error, err error) {
	ctx := context.Background()

	stop = func() error {
		return nil
	}

	key = "meili-master-key"

	container, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "getmeili/meilisearch:prototype-japanese-184",
			Env: map[string]string{
				"MEILI_MASTER_KEY": key,
			},
			ExposedPorts: []string{"7700/tcp"},
			WaitingFor:   wait.ForListeningPort("7700/tcp"),
		},
		Started: true,
	})
	if err != nil {
		return
	}
	stop = func() error {
		return container.Terminate(ctx)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return
	}
	port, err := container.MappedPort(ctx, "7700/tcp")
	if err != nil {
		return
	}

	url = fmt.Sprintf("http://%s:%d", host, port.Int())

	return
}
