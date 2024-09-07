package crawl

import (
	"context"
	"fjnkt98/atcodersearch/pkg/atcoder"
	"fmt"
	"slices"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SaveContests(ctx context.Context, pool *pgxpool.Pool, contests []atcoder.Contest) (int64, error) {
	if len(contests) == 0 {
		return 0, nil
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	dialect := goqu.Dialect("postgres")

	var count int64 = 0
	for chunk := range slices.Chunk(contests, 1000) {
		q := dialect.Insert("contests").Prepared(true).Cols(
			"contest_id",
			"start_epoch_second",
			"duration_second",
			"title",
			"rate_change",
			"category",
			"updated_at",
		)
		for _, c := range chunk {
			q = q.Vals(goqu.Vals{
				c.ID,
				c.StartEpochSecond,
				c.DurationSecond,
				c.Title,
				c.RateChange,
				c.Categorize(),
				goqu.L("NOW()"),
			})
		}

		q = q.OnConflict(
			goqu.DoUpdate(
				"contest_id",
				goqu.Record{
					"contest_id":         goqu.L("EXCLUDED.contest_id"),
					"start_epoch_second": goqu.L("EXCLUDED.start_epoch_second"),
					"duration_second":    goqu.L("EXCLUDED.duration_second"),
					"title":              goqu.L("EXCLUDED.title"),
					"rate_change":        goqu.L("EXCLUDED.rate_change"),
					"category":           goqu.L("EXCLUDED.category"),
					"updated_at":         goqu.L("NOW()"),
				},
			),
		)

		sql, args, err := q.ToSQL()
		if err != nil {
			return 0, fmt.Errorf("to sql: %w", err)
		}

		result, err := tx.Exec(ctx, sql, args...)
		if err != nil {
			return 0, fmt.Errorf("exec: %w", err)
		}

		count += result.RowsAffected()
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit transaction: %w", err)
	}

	return count, nil
}
