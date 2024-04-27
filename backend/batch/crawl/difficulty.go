package crawl

import (
	"context"
	"fjnkt98/atcodersearch/pkg/atcoder"
	"fjnkt98/atcodersearch/repository"

	"log/slog"

	"github.com/goark/errs"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DifficultyCrawler interface {
	CrawlDifficulty(ctx context.Context) error
}

type difficultyCrawler struct {
	client atcoder.AtCoderProblemsClient
	pool   *pgxpool.Pool
}

func NewDifficultyCrawler(client atcoder.AtCoderProblemsClient, pool *pgxpool.Pool) DifficultyCrawler {
	return &difficultyCrawler{
		client: client,
		pool:   pool,
	}
}

func (c *difficultyCrawler) CrawlDifficulty(ctx context.Context) error {
	slog.Info("Start to crawl difficulties.")
	difficulties, err := c.client.FetchDifficulties(ctx)
	if err != nil {
		return errs.Wrap(err)
	}
	slog.Info("Finish crawling difficulties.")

	slog.Info("Start to save difficulties.")
	count, err := repository.BulkUpdate(ctx, c.pool, "difficulties", convertDifficulties(difficulties))
	if err != nil {
		return errs.New("failed to bulk update difficulties", errs.WithCause(err))
	}
	slog.Info("Finish saving difficulties.", slog.Int64("count", count))

	return nil
}

func convertDifficulties(difficulties map[string]atcoder.Difficulty) []repository.Difficulty {
	result := make([]repository.Difficulty, 0, len(difficulties))
	for problemID, d := range difficulties {
		result = append(result, repository.Difficulty{
			ProblemID:        problemID,
			Slope:            d.Slope,
			Intercept:        d.Intercept,
			Variance:         d.Variance,
			Difficulty:       d.Difficulty,
			Discrimination:   d.Discrimination,
			IrtLoglikelihood: d.IrtLogLikelihood,
			IrtUsers:         d.IrtUsers,
			IsExperimental:   d.IsExperimental,
		})
	}
	return result
}
