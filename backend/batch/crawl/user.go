package crawl

import (
	"context"
	"database/sql"
	"fjnkt98/atcodersearch/pkg/atcoder"
	"fjnkt98/atcodersearch/repository"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

type UserCrawler struct {
	client   *atcoder.AtCoderClient
	pool     *pgxpool.Pool
	duration time.Duration
}

func NewUserCrawler(
	client *atcoder.AtCoderClient,
	pool *pgxpool.Pool,
	duration time.Duration,
) *UserCrawler {
	return &UserCrawler{
		client:   client,
		pool:     pool,
		duration: duration,
	}
}

func (c *UserCrawler) Crawl(ctx context.Context) error {
	slog.LogAttrs(ctx, slog.LevelInfo, "start to crawl users.")

	users := make([]atcoder.User, 0)
loop:
	for i := 1; ; i++ {
		slog.LogAttrs(ctx, slog.LevelInfo, "fetch users", slog.Int("page", i))

		u, err := c.client.FetchUsers(ctx, i)
		if err != nil {
			return fmt.Errorf("fetch users: %w", err)
		}

		if len(u) == 0 {
			slog.LogAttrs(ctx, slog.LevelInfo, "there is no more crawl target.")
			break loop
		}

		users = append(users, u...)

		time.Sleep(c.duration)
	}

	count, err := SaveUsers(ctx, c.pool, users, time.Now())
	if err != nil {
		return fmt.Errorf("save users: %w", err)
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "finish to crawl contests.", slog.Int64("count", count))
	return nil
}

func SaveUsers(ctx context.Context, pool *pgxpool.Pool, users []atcoder.User, timestamp time.Time) (int64, error) {
	if len(users) == 0 {
		return 0, nil
	}

	tx, err := bun.NewDB(stdlib.OpenDBFromPool(pool), pgdialect.New()).BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return 0, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	var count int64 = 0

	for chunk := range slices.Chunk(slices.Collect(repository.Map(repository.NewUser, slices.Values(users), timestamp)), 1000) {
		res, err := tx.NewInsert().
			Model(&chunk).
			On("CONFLICT (user_id) DO UPDATE").
			Set("user_id = EXCLUDED.user_id").
			Set("rating = EXCLUDED.rating").
			Set("highest_rating = EXCLUDED.highest_rating").
			Set("affiliation = EXCLUDED.affiliation").
			Set("birth_year = EXCLUDED.birth_year").
			Set("country = EXCLUDED.country").
			Set("crown = EXCLUDED.crown").
			Set("join_count = EXCLUDED.join_count").
			Set("rank = EXCLUDED.rank").
			Set("active_rank = EXCLUDED.active_rank").
			Set("wins = EXCLUDED.wins").
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
