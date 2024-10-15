package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	HistoryStatusWorking   = "working"
	HistoryStatusCompleted = "completed"
	HistoryStatusAborted   = "aborted"
)

func NewBatchHistory(ctx context.Context, pool *pgxpool.Pool, name string, option []byte) (*BatchHistory, error) {
	q := New(pool)

	h, err := q.CreateBatchHistory(ctx, CreateBatchHistoryParams{Name: name, Options: option})
	if err != nil {
		return nil, fmt.Errorf("create batch history: %w", err)
	}
	return &h, nil
}

func (h *BatchHistory) Complete(ctx context.Context, pool *pgxpool.Pool) error {
	q := New(pool)

	returned, err := q.CompleteBatchHistory(ctx, h.ID)
	if err != nil {
		return fmt.Errorf("complete batch history: %w", err)
	}
	*h = returned

	return nil
}

var ErrHistoryConfirmed = errors.New("the batch history already confirmed")

func (h *BatchHistory) Abort(ctx context.Context, pool *pgxpool.Pool) error {
	if h.Status != HistoryStatusWorking {
		return ErrHistoryConfirmed
	}
	q := New(pool)

	returned, err := q.AbortBatchHistory(ctx, h.ID)
	if err != nil {
		return fmt.Errorf("abort batch history: %w", err)
	}
	*h = returned

	return nil
}

func NewCrawlHistory(ctx context.Context, pool *pgxpool.Pool, contestID string) (*SubmissionCrawlHistory, error) {
	q := New(pool)

	h, err := q.CreateCrawlHistory(ctx, contestID)
	if err != nil {
		return nil, fmt.Errorf("create crawl history: %w", err)
	}
	return &h, nil
}

func (h *SubmissionCrawlHistory) Complete(ctx context.Context, pool *pgxpool.Pool) error {
	q := New(pool)

	returned, err := q.CompleteCrawlHistory(ctx, h.ID)
	if err != nil {
		return fmt.Errorf("complete crawl history: %w", err)
	}
	*h = returned

	return nil
}

func (h *SubmissionCrawlHistory) Abort(ctx context.Context, pool *pgxpool.Pool) error {
	if h.Status != HistoryStatusWorking {
		return ErrHistoryConfirmed
	}
	q := New(pool)

	returned, err := q.AbortCrawlHistory(ctx, h.ID)
	if err != nil {
		return fmt.Errorf("abort crawl history: %w", err)
	}
	*h = returned

	return nil
}

func FetchLatestCrawlHistory(ctx context.Context, pool *pgxpool.Pool, contestID string) (*SubmissionCrawlHistory, error) {
	q := New(pool)

	h, err := q.FetchLatestCrawlHistory(ctx, contestID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			d := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
			h = SubmissionCrawlHistory{
				ID:         -1,
				ContestID:  contestID,
				StartedAt:  d,
				Status:     HistoryStatusCompleted,
				FinishedAt: &d,
			}
		} else {
			return nil, err
		}
	}
	return &h, nil
}
