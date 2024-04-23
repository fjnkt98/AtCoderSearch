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
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return errs.New("failed to start transaction", errs.WithCause(err))
	}

	q := repository.New(c.pool).WithTx(tx)
	for problemID, d := range difficulties {
		if _, err := q.InsertDifficulty(ctx, repository.InsertDifficultyParams{
			ProblemID:        problemID,
			Slope:            d.Slope,
			Intercept:        d.Intercept,
			Variance:         d.Variance,
			Difficulty:       d.Difficulty,
			Discrimination:   d.Discrimination,
			IrtLoglikelihood: d.IrtLogLikelihood,
			IrtUsers:         d.IrtUsers,
			IsExperimental:   d.IsExperimental,
		}); err != nil {
			return errs.Wrap(err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return errs.New("failed to commit transaction", errs.WithCause(err))
	}
	slog.Info("Finish saving difficulties.")

	return nil
}
