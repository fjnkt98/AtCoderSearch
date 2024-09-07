package crawl

import (
	"context"
	"database/sql"
	"fjnkt98/atcodersearch/pkg/atcoder"
	"fjnkt98/atcodersearch/repository"
	"fmt"
	"slices"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func SaveContests(ctx context.Context, pool *pgxpool.Pool, contests []atcoder.Contest, timestamp time.Time) (int64, error) {
	if len(contests) == 0 {
		return 0, nil
	}

	tx, err := bun.NewDB(stdlib.OpenDBFromPool(pool), pgdialect.New()).BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return 0, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	var count int64 = 0

	for chunk := range slices.Chunk(slices.Collect(repository.Map(repository.NewContest, slices.Values(contests), timestamp)), 1000) {
		res, err := tx.NewInsert().
			Model(&chunk).
			On("CONFLICT (contest_id) DO UPDATE").
			Set("contest_id = EXCLUDED.contest_id").
			Set("start_epoch_second = EXCLUDED.start_epoch_second").
			Set("duration_second = EXCLUDED.duration_second").
			Set("title = EXCLUDED.title").
			Set("rate_change = EXCLUDED.rate_change").
			Set("category = EXCLUDED.category").
			Set("updated_at = EXCLUDED.updated_at").
			Exec(ctx)
		if err != nil {
			return 0, fmt.Errorf("insert: %w", err)
		}

		if c, err := res.RowsAffected(); err != nil {
			return 0, fmt.Errorf("rows affected: %w", err)
		} else {
			count += c
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit transaction: %w", err)
	}

	return count, nil
}
