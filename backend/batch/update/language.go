package update

import (
	"context"
	"fjnkt98/atcodersearch/repository"
	"log/slog"

	"github.com/goark/errs"
	"github.com/jackc/pgx/v5/pgxpool"
)

func UpdateLanguage(ctx context.Context, pool *pgxpool.Pool) error {
	slog.Info("Start UpdateLanguage")

	q := repository.New(pool)

	id, err := q.CreateBatchHistory(ctx, repository.CreateBatchHistoryParams{Name: "UpdateLanguage", Options: []byte("{}")})
	if err != nil {
		return errs.New("failed to create batch history", errs.WithCause(err))
	}

	if res, err := q.UpdateLanguages(ctx); err != nil {
		return errs.New("failed to update languages", errs.WithCause(err))
	} else {
		slog.Info("Languages updated successfully.", slog.Int64("count", res.RowsAffected()))
	}

	if err := q.UpdateBatchHistory(ctx, repository.UpdateBatchHistoryParams{ID: id, Status: "finished"}); err != nil {
		return errs.New("failed to update batch history", errs.WithCause(err))
	}

	slog.Info("Finish UpdateProblem")
	return nil
}
