package crawl

import (
	"context"
	"fmt"
	"time"

	"log/slog"

	"github.com/fjnkt98/atcodersearch-batch/pkg/atcoder"
	"github.com/fjnkt98/atcodersearch-batch/repository"
	"github.com/jackc/pgx/v5/pgxpool"
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
	slog.LogAttrs(ctx, slog.LevelInfo, "Start to fetch users.")

	users := make([]atcoder.User, 0)
loop:
	for i := 1; ; i++ {
		slog.LogAttrs(ctx, slog.LevelInfo, "Crawl users", slog.Int("page", i))

		us, err := c.client.FetchUsers(ctx, i)
		if err != nil {
			return fmt.Errorf("fetch users: %w", err)
		}

		if len(us) == 0 {
			slog.LogAttrs(ctx, slog.LevelInfo, "There is no more crawl target.")
			break loop
		}

		users = append(users, us...)

		time.Sleep(c.duration)
	}

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("save users: %w", err)
	}
	defer tx.Rollback(ctx)

	q := repository.New(tx)

	var count int64 = 0
	for _, user := range users {
		result, err := q.InsertUser(ctx, repository.InsertUserParams{
			UserID:        user.UserID,
			Rating:        user.Rating,
			HighestRating: user.HighestRating,
			Affiliation:   user.Affiliation,
			BirthYear:     user.BirthYear,
			Country:       user.Country,
			Crown:         user.Crown,
			JoinCount:     user.JoinCount,
			Rank:          user.Rank,
			ActiveRank:    user.ActiveRank,
			Wins:          user.Wins,
		})
		if err != nil {
			return fmt.Errorf("save users: %w", err)
		}

		count += result.RowsAffected()
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("save users: %w", err)
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "Finish saving users.", slog.Int64("count", count))

	return nil
}
