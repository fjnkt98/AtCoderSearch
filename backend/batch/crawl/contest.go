package crawl

import (
	"context"
	"fjnkt98/atcodersearch/pkg/atcoder"
	"fjnkt98/atcodersearch/repository"

	"log/slog"

	"github.com/goark/errs"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ContestCrawler interface {
	CrawlContest(ctx context.Context) error
}

type contestCrawler struct {
	client atcoder.AtCoderProblemsClient
	pool   *pgxpool.Pool
}

func NewContestCrawler(client atcoder.AtCoderProblemsClient, pool *pgxpool.Pool) ContestCrawler {
	return &contestCrawler{
		client: client,
		pool:   pool,
	}
}

func (c *contestCrawler) CrawlContest(ctx context.Context) error {
	slog.Info("Start to fetch contests.")
	contests, err := c.client.FetchContests(ctx)
	if err != nil {
		return errs.Wrap(err)
	}
	slog.Info("Finish fetching contests.")

	slog.Info("Start to save contests.")
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return errs.New("failed to start transaction", errs.WithCause(err))
	}
	q := repository.New(c.pool).WithTx(tx)
	for _, contest := range contests {
		if _, err := q.InsertContest(ctx, repository.InsertContestParams{
			ContestID:        contest.ID,
			StartEpochSecond: contest.StartEpochSecond,
			DurationSecond:   contest.DurationSecond,
			Title:            contest.Title,
			RateChange:       contest.RateChange,
			Category:         contest.Categorize(),
		}); err != nil {
			return errs.Wrap(err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return errs.New("failed to commit transaction", errs.WithCause(err))
	}

	slog.Info("Finish saving contest list.")

	return nil
}
