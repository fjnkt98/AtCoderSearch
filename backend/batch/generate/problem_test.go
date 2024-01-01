//go:build test_generate

package generate

import (
	"context"
	"database/sql"
	"errors"
	"fjnkt98/atcodersearch/batch"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

func getTestDB() *bun.DB {
	os.Setenv("PGSSLMODE", "disable")
	engine, err := sql.Open("postgres", "postgres://atcodersearch:atcodersearch@localhost/atcodersearch")
	if err != nil {
		slog.Error("failed to open database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if err := engine.Ping(); err != nil {
		slog.Error("failed to connect database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	db := bun.NewDB(engine, pgdialect.New())
	db.AddQueryHook(
		bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
			bundebug.FromEnv("BUNDEBUG"),
		),
	)

	return db
}

func TestReadProblems(t *testing.T) {
	db := getTestDB()
	reader := NewProblemRowReader(db)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	ch := make(chan Documenter, 1)

	err := reader.ReadRows(ctx, ch)
	if err != nil && !errors.Is(err, batch.ErrInterrupt) {
		t.Errorf("an error occurred in read problem rows: %s", err.Error())
	}
}
