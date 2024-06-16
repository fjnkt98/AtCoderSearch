package crawl

import (
	"context"
	"fjnkt98/atcodersearch/pkg/atcoder"
	"fjnkt98/atcodersearch/repository"

	"log/slog"

	"github.com/goark/errs"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ContestCrawler struct {
	client *atcoder.AtCoderProblemsClient
	pool   *pgxpool.Pool
}

func NewContestCrawler(client *atcoder.AtCoderProblemsClient, pool *pgxpool.Pool) *ContestCrawler {
	return &ContestCrawler{
		client: client,
		pool:   pool,
	}
}

func (c *ContestCrawler) Crawl(ctx context.Context) error {
	slog.Info("Start to fetch contests.")
	contests, err := c.client.FetchContests(ctx)
	if err != nil {
		return errs.Wrap(err)
	}
	slog.Info("Finish fetching contests.")

	slog.Info("Start to save contests.")
	count, err := repository.BulkUpdate(ctx, c.pool, "contests", repository.NewContests(contests))
	if err != nil {
		return errs.New("failed to bulk update contests", errs.WithCause(err))
	}
	slog.Info("Finish saving contest list.", slog.Int64("count", count))

	return nil
}
