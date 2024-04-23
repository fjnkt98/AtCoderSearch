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

type UserCrawler interface {
	CrawlUser(ctx context.Context) error
}

type userCrawler struct {
	client   atcoder.AtCoderClient
	pool     *pgxpool.Pool
	duration time.Duration
}

func NewUserCrawler(
	client atcoder.AtCoderClient,
	pool *pgxpool.Pool,
	duration time.Duration,
) *userCrawler {
	return &userCrawler{
		client:   client,
		pool:     pool,
		duration: duration,
	}
}

func (c *userCrawler) CrawlUser(ctx context.Context) error {
	slog.Info("Start to crawl users.")

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

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return errs.New("failed to start transaction", errs.WithCause(err))
	}
	q := repository.New(tx)
	for _, u := range users {
		if _, err := q.InsertUser(ctx, repository.InsertUserParams{
			UserName:      u.UserName,
			Rating:        u.Rating,
			HighestRating: u.HighestRating,
			Affiliation:   u.Affiliation,
			BirthYear:     u.BirthYear,
			Country:       u.Country,
			Crown:         u.Crown,
			JoinCount:     u.JoinCount,
			Rank:          u.Rank,
			ActiveRank:    u.ActiveRank,
			Wins:          u.Wins,
		}); err != nil {
			return errs.New("failed to insert the user", errs.WithCause(err), errs.WithContext("user", u))
		}
	}

	slog.Info("Finish crawling users successfully.")
	return nil
}
