package crawl

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/fjnkt98/atcodersearch-batch/pkg/atcoder"
	"github.com/fjnkt98/atcodersearch-batch/repository"
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
	slog.LogAttrs(ctx, slog.LevelInfo, "Start to fetch difficulties.")

	difficulties, err := c.client.FetchDifficulties(ctx)
	if err != nil {
		return err
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "Finish fetching difficulties. Start to save difficulties.")

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("save difficulty: %w", err)
	}
	defer tx.Rollback(ctx)

	q := repository.New(tx)

	var count int64 = 0
	for problemID, d := range difficulties {
		result, err := q.InsertDifficulty(ctx, repository.InsertDifficultyParams{
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
		if err != nil {
			return fmt.Errorf("save difficulty: %w", err)
		}
		count += result.RowsAffected()
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("save difficulty: %w", err)
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "Finish saving difficulties.", slog.Int64("count", count))

	return nil
}
