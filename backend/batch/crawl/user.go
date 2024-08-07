package crawl

import (
	"context"
	"fjnkt98/atcodersearch/pkg/atcoder"
	"fjnkt98/atcodersearch/repository"
	"time"

	"log/slog"

	"github.com/goark/errs"
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
	users := make([]atcoder.User, 0)
loop:
	for i := 1; ; i++ {
		slog.Info("Crawl users", slog.Int("page", i))
		us, err := c.client.FetchUsers(ctx, i)
		if err != nil {
			return errs.Wrap(err)
		}

		if len(us) == 0 {
			slog.Info("There is no more crawl target.")
			break loop
		}

		users = append(users, us...)

		time.Sleep(c.duration)
	}
	count, err := repository.BulkUpdate(ctx, c.pool, "users", repository.NewUsers(users))
	if err != nil {
		return errs.New("failed to bulk update users", errs.WithCause(err))
	}
	slog.Info("Finish crawling users successfully.", slog.Int64("count", count))
	return nil
}
