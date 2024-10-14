package update

import (
	"context"
	"fjnkt98/atcodersearch/repository"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

func UpdateLanguage(ctx context.Context, pool *pgxpool.Pool) error {
	slog.LogAttrs(ctx, slog.LevelInfo, "start update language")

	q := repository.New(pool)

	res, err := q.UpdateLanguages(ctx)
	if err != nil {
		return fmt.Errorf("update language: %w", err)
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "finish update language", slog.Int64("rowsAffected", res.RowsAffected()))
	return nil
}
