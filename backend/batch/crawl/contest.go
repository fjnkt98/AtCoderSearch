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
	count, err := repository.BulkUpdate(ctx, c.pool, "contests", convertContests(contests))
	if err != nil {
		return errs.New("failed to bulk update contests", errs.WithCause(err))
	}
	slog.Info("Finish saving contest list.", slog.Int64("count", count))

	return nil
}

func convertContests(contests []atcoder.Contest) []repository.Contest {
	result := make([]repository.Contest, len(contests))
	for i, c := range contests {
		result[i] = repository.Contest{
			ContestID:        c.ID,
			StartEpochSecond: c.StartEpochSecond,
			DurationSecond:   c.DurationSecond,
			Title:            c.Title,
			RateChange:       c.RateChange,
			Category:         c.Categorize(),
		}
	}
	return result
}
