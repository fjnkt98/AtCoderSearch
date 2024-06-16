package crawl

import (
	"context"
	"fjnkt98/atcodersearch/pkg/atcoder"
	"fjnkt98/atcodersearch/repository"

	"log/slog"

	"github.com/goark/errs"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DifficultyCrawler struct {
	client *atcoder.AtCoderProblemsClient
	pool   *pgxpool.Pool
}

func NewDifficultyCrawler(client *atcoder.AtCoderProblemsClient, pool *pgxpool.Pool) *DifficultyCrawler {
	return &DifficultyCrawler{
		client: client,
		pool:   pool,
	}
}

func (c *DifficultyCrawler) Crawl(ctx context.Context) error {
	slog.Info("Start to crawl difficulties.")
	difficulties, err := c.client.FetchDifficulties(ctx)
	if err != nil {
		return errs.Wrap(err)
	}
	slog.Info("Finish crawling difficulties.")

	slog.Info("Start to save difficulties.")
	count, err := repository.BulkUpdate(ctx, c.pool, "difficulties", repository.NewDifficulties(difficulties))
	if err != nil {
		return errs.New("failed to bulk update difficulties", errs.WithCause(err))
	}
	slog.Info("Finish saving difficulties.", slog.Int64("count", count))

	return nil
}
