//go:build test_generate

// docker run --rm -d -p 5432:5432 --name postgres -e POSTGRES_DB=test_atcodersearch -e POSTGRES_USER=test_atcodersearch -e POSTGRES_PASSWORD=test_atcodersearch --mount type=bind,src=./schema.sql,dst=/docker-entrypoint-initdb.d/schema.sql postgres:15

package generate

import (
	"context"
	"database/sql"
	"errors"
	"fjnkt98/atcodersearch/batch"
	"os"
	"testing"
	"time"

	"github.com/goark/errs"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

func getTestDB() (*bun.DB, error) {
	os.Setenv("PGSSLMODE", "disable")
	engine, err := sql.Open("postgres", "postgres://test_atcodersearch:test_atcodersearch@localhost/test_atcodersearch")
	if err != nil {
		return nil, errs.New("failed to open database", errs.WithCause(err))
	}

	if err := engine.Ping(); err != nil {
		return nil, errs.New("failed to connect database", errs.WithCause(err))
	}

	db := bun.NewDB(engine, pgdialect.New())
	db.AddQueryHook(
		bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
			bundebug.FromEnv("BUNDEBUG"),
		),
	)

	return db, nil
}

func TestReadProblems(t *testing.T) {
	db, err := getTestDB()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	reader := NewProblemRowReader(db)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	ch := make(chan Documenter, 1)

	err = reader.ReadRows(ctx, ch)
	if err != nil && !errors.Is(err, batch.ErrInterrupt) {
		t.Errorf("an error occurred in read problem rows: %s", err.Error())
	}
}
