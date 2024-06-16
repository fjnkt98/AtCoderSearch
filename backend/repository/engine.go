package repository

import (
	"context"

	"github.com/goark/errs"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, errs.New("failed to parse database connection config", errs.WithCause(err))
	}
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, errs.New("failed to create connection pool from config", errs.WithCause(err))
	}
	return pool, nil
}
