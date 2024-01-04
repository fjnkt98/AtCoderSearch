//go:build test_repository

// docker run --rm -d -p 5432:5432 --name postgres -e POSTGRES_DB=test_atcodersearch -e POSTGRES_USER=test_atcodersearch -e POSTGRES_PASSWORD=test_atcodersearch --mount type=bind,src=./schema.sql,dst=/docker-entrypoint-initdb.d/schema.sql postgres:15

package repository

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

func ptr[T any](v T) *T {
	return &v
}

func getTestDB() *bun.DB {
	os.Setenv("PGSSLMODE", "disable")
	engine, err := sql.Open("postgres", "postgres://test_atcodersearch:test_atcodersearch@localhost/test_atcodersearch")
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
			bundebug.WithVerbose(false),
			bundebug.FromEnv("BUNDEBUG"),
		),
	)

	return db
}

func TestSaveContest(t *testing.T) {
	db := getTestDB()
	repository := NewContestRepository(db)

	contests := []Contest{
		{
			ContestID:        "abc300",
			StartEpochSecond: 1682769600,
			DurationSecond:   6000,
			Title:            "ユニークビジョンプログラミングコンテスト2023 春 (AtCoder Beginner Contest 300)",
			RateChange:       "~1999",
			Category:         "ABC",
		},
	}

	ctx := context.Background()
	if err := repository.Save(ctx, contests); err != nil {
		t.Fatalf("failed to save contests: %s", err.Error())
	}
}

func TestFetchALlContestIDs(t *testing.T) {
	db := getTestDB()
	repository := NewContestRepository(db)

	ctx := context.Background()
	_, err := repository.FetchContestIDs(ctx, nil)
	if err != nil {
		t.Fatalf("failed to fetch contest ids: %s", err.Error())
	}
}

func TestFetchSpecifiedContestIDs(t *testing.T) {
	db := getTestDB()
	repository := NewContestRepository(db)

	ctx := context.Background()
	_, err := repository.FetchContestIDs(ctx, []string{"ABC", "ARC"})
	if err != nil {
		t.Fatalf("failed to fetch contest ids: %s", err.Error())
	}
}

func TestFetchCategories(t *testing.T) {
	db := getTestDB()
	repository := NewContestRepository(db)

	ctx := context.Background()
	_, err := repository.FetchCategories(ctx)
	if err != nil {
		t.Fatalf("failed to fetch categories: %s", err.Error())
	}
}
