package repository

import (
	"database/sql"
	"log/slog"

	"github.com/goark/errs"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func GetEngine(dsn string) (*sql.DB, error) {
	engine, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, errs.New(
			"failed to open database",
			errs.WithCause(err),
		)
	}

	if err := engine.Ping(); err != nil {
		return nil, errs.New(
			"failed to connect database",
			errs.WithCause(err),
		)
	}

	return engine, nil
}

func MustGetEngine(dsn string) *sql.DB {
	engine, err := GetEngine(dsn)
	if err != nil {
		slog.Error("failed to get engine instance", slog.Any("error", err))
		panic("failed to get engine instance")
	}
	return engine
}

func GetDB(dsn string) (*bun.DB, error) {
	engine, err := GetEngine(dsn)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	db := bun.NewDB(engine, pgdialect.New())
	return db, nil
}

func MustGetDB(dsn string) *bun.DB {
	db, err := GetDB(dsn)
	if err != nil {
		slog.Error("failed to get db instance", slog.Any("error", err))
		panic("failed to get db instance")
	}

	return db
}
