package update

import (
	"context"
	"fjnkt98/atcodersearch/repository"
	"fjnkt98/atcodersearch/settings"
	"log/slog"

	"github.com/goark/errs"
	"github.com/jackc/pgx/v5/pgxpool"
)

func UpdateLanguage(ctx context.Context, pool *pgxpool.Pool) error {
	slog.Info("Start Batch", slog.String("name", settings.UPDATE_LANGUAGE_BATCH_NAME))

	h, err := repository.NewBatchHistory(ctx, pool, settings.UPDATE_LANGUAGE_BATCH_NAME, nil)
	if err != nil {
		return errs.Wrap(err, errs.WithCause(err), errs.WithContext("name", settings.UPDATE_LANGUAGE_BATCH_NAME))
	}
	defer h.Fail(ctx, pool)

	q := repository.New(pool)
	if res, err := q.UpdateLanguages(ctx); err != nil {
		return errs.New("failed to update languages", errs.WithCause(err), errs.WithContext("name", settings.UPDATE_LANGUAGE_BATCH_NAME))
	} else {
		slog.Info("Languages updated successfully.", slog.Int64("count", res.RowsAffected()))
	}

	if err := h.Finish(ctx, pool); err != nil {
		return errs.Wrap(err, errs.WithCause(err), errs.WithContext("name", settings.UPDATE_LANGUAGE_BATCH_NAME))
	}

	slog.Info("Finish Batch", slog.String("name", settings.UPDATE_LANGUAGE_BATCH_NAME))
	return nil
}
