package repository

import (
	"context"

	"github.com/goark/errs"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewBatchHistory(ctx context.Context, pool *pgxpool.Pool, name string, option []byte) (*BatchHistory, error) {
	q := New(pool)
	h, err := q.CreateBatchHistory(ctx, CreateBatchHistoryParams{Name: name, Options: option})
	if err != nil {
		return nil, errs.Wrap(err, errs.WithContext("name", name), errs.WithContext("option", option))
	}
	return &h, nil
}

func overwrite(l *BatchHistory, r BatchHistory) {
	l.ID = r.ID
	l.Name = r.Name
	l.StartedAt = r.StartedAt
	l.FinishedAt = r.FinishedAt
	l.Status = r.Status
	l.Options = r.Options
}

func (h *BatchHistory) Finish(ctx context.Context, pool *pgxpool.Pool) error {
	q := New(pool)

	updated, err := q.UpdateBatchHistory(ctx, UpdateBatchHistoryParams{Status: "finished", ID: h.ID})
	if err != nil {
		return errs.Wrap(err, errs.WithContext("name", h.Name))
	}
	overwrite(h, updated)

	return nil
}

func (h *BatchHistory) Fail(ctx context.Context, pool *pgxpool.Pool) error {
	if h.Status != "working" {
		return nil
	}
	q := New(pool)

	updated, err := q.UpdateBatchHistory(ctx, UpdateBatchHistoryParams{Status: "failed", ID: h.ID})
	if err != nil {
		return errs.Wrap(err, errs.WithContext("name", h.Name))
	}
	overwrite(h, updated)

	return nil
}
