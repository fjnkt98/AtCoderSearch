//go:generate mockgen -source=$GOFILE -destination=./mock/mock_$GOFILE -package=$GOPACKAGE

package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"log/slog"

	"github.com/goark/errs"
	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
)

type SubmissionCrawlHistoryRepository interface {
	Save(ctx context.Context, history SubmissionCrawlHistory) error
	GetLatestHistory(ctx context.Context, contestID string) (SubmissionCrawlHistory, error)
}

type SubmissionCrawlHistory struct {
	bun.BaseModel `bun:"table:submission_crawl_history,alias:h"`
	StartedAt     int    `bun:"started_at"`
	ContestID     string `bun:"contest_id,type:text"`
}

func NewSubmissionCrawlHistory(contestID string) SubmissionCrawlHistory {
	return SubmissionCrawlHistory{
		StartedAt: int(time.Now().Unix()),
		ContestID: contestID,
	}
}

type submissionCrawlHistoryRepository struct {
	db *bun.DB
}

func NewSubmissionCrawlHistoryRepository(db *bun.DB) SubmissionCrawlHistoryRepository {
	return &submissionCrawlHistoryRepository{
		db: db,
	}
}

func (r *submissionCrawlHistoryRepository) Save(ctx context.Context, history SubmissionCrawlHistory) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errs.New(
			"failed to start transaction",
			errs.WithCause(err),
		)
	}
	defer tx.Rollback()

	if _, err := tx.NewInsert().Model(&history).Exec(ctx); err != nil {
		return errs.New(
			"failed to execute sql",
			errs.WithCause(err),
		)
	}

	if err := tx.Commit(); err != nil {
		return errs.New(
			"failed to commit transaction",
			errs.WithCause(err),
		)
	}

	return nil
}

func (r *submissionCrawlHistoryRepository) GetLatestHistory(ctx context.Context, contestID string) (SubmissionCrawlHistory, error) {
	var history SubmissionCrawlHistory
	err := r.db.NewSelect().
		Model(&history).
		Where("? = ?", bun.Ident("contest_id"), contestID).
		Order("started_at DESC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Info(fmt.Sprintf("`submission_crawl_history` table is empty in term of contest id `%s`", contestID))
			return SubmissionCrawlHistory{StartedAt: 0, ContestID: contestID}, nil
		} else {
			return SubmissionCrawlHistory{}, errs.New(
				"failed to get latest crawl history",
				errs.WithCause(err),
			)
		}
	}

	return history, nil
}

type UpdateHistoryRepository interface {
	Finish(ctx context.Context, history *UpdateHistory) error
	Cancel(ctx context.Context, history *UpdateHistory) error
	GetLatest(ctx context.Context, domain string) (UpdateHistory, error)
}

type UpdateHistory struct {
	bun.BaseModel `bun:"table:update_history,alias:h"`
	Domain        string    `bun:"domain,type:text"`
	StartedAt     time.Time `bun:"started_at"`
	FinishedAt    time.Time `bun:"finished_at"`
	Status        string    `bun:"status,type:text"`
	Options       string    `bun:"options"`
	wasSaved      bool
}

func NewUpdateHistory(domain string, options string) UpdateHistory {
	return UpdateHistory{
		Domain:    domain,
		StartedAt: time.Now(),
		Options:   options,
		wasSaved:  false,
	}
}

type updateHistoryRepository struct {
	db *bun.DB
}

func (r *updateHistoryRepository) save(ctx context.Context, history *UpdateHistory, status string) error {
	if history.wasSaved {
		return nil
	}

	history.wasSaved = true
	history.FinishedAt = time.Now()
	history.Status = status

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errs.New(
			"failed to start transaction",
			errs.WithCause(err),
		)
	}
	defer tx.Rollback()

	if _, err := tx.NewInsert().Model(history).Exec(ctx); err != nil {
		return errs.New(
			"failed to execute sql",
			errs.WithCause(err),
		)
	}

	if err := tx.Commit(); err != nil {
		return errs.New(
			"failed to commit transaction",
			errs.WithCause(err),
		)
	}

	return nil
}

func (r *updateHistoryRepository) Finish(ctx context.Context, history *UpdateHistory) error {
	return r.save(ctx, history, "finished")
}

func (r *updateHistoryRepository) Cancel(ctx context.Context, history *UpdateHistory) error {
	return r.save(ctx, history, "canceled")
}

func (r *updateHistoryRepository) GetLatest(ctx context.Context, domain string) (UpdateHistory, error) {
	var history UpdateHistory
	err := r.db.NewSelect().
		Model(&history).
		Where("? = ?", bun.Ident("domain"), domain).
		Where("? = ?", bun.Ident("status"), "finished").
		Order("started_at DESC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Info("`update_history` table is empty.")
			return UpdateHistory{Domain: domain, Status: "finished"}, nil
		} else {
			return UpdateHistory{}, errs.New(
				"failed to get latest update history",
				errs.WithCause(err),
			)
		}
	}

	return history, nil
}

func NewUpdateHistoryRepository(db *bun.DB) UpdateHistoryRepository {
	return &updateHistoryRepository{
		db: db,
	}
}
