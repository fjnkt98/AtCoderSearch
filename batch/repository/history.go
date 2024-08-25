package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewBatchHistory(ctx context.Context, pool *pgxpool.Pool, name string, option []byte) (*BatchHistory, error) {
	q := New(pool)
	h, err := q.CreateBatchHistory(ctx, CreateBatchHistoryParams{Name: name, Options: option})
	if err != nil {
		return nil, fmt.Errorf("batch history: %w", err)
	}
	return &h, nil
}

func (h *BatchHistory) overwrite(r BatchHistory) {
	h.ID = r.ID
	h.Name = r.Name
	h.StartedAt = r.StartedAt
	h.FinishedAt = r.FinishedAt
	h.Status = r.Status
	h.Options = r.Options
}

func (h *BatchHistory) Finish(ctx context.Context, pool *pgxpool.Pool) error {
	q := New(pool)

	updated, err := q.UpdateBatchHistory(ctx, UpdateBatchHistoryParams{Status: "finished", ID: h.ID})
	if err != nil {
		return fmt.Errorf("batch history: %w", err)
	}
	h.overwrite(updated)

	return nil
}

var ErrHistoryConfirmed = errors.New("the batch history already confirmed")

func (h *BatchHistory) Fail(ctx context.Context, pool *pgxpool.Pool) error {
	if h.Status != "working" {
		return ErrHistoryConfirmed
	}
	q := New(pool)

	updated, err := q.UpdateBatchHistory(ctx, UpdateBatchHistoryParams{Status: "failed", ID: h.ID})
	if err != nil {
		return fmt.Errorf("batch history: %w", err)
	}
	h.overwrite(updated)

	return nil
}
