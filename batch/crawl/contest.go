package crawl

import (
	"context"
	"fmt"

	"log/slog"

	"github.com/fjnkt98/atcodersearch-batch/pkg/atcoder"
	"github.com/fjnkt98/atcodersearch-batch/repository"
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
	slog.LogAttrs(ctx, slog.LevelInfo, "Start to fetch contests.")

	contests, err := c.client.FetchContests(ctx)
	if err != nil {
		return err
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "Finish fetching contests. Start to save contests.")

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("save contest: %w", err)
	}
	defer tx.Rollback(ctx)

	q := repository.New(tx)

	var count int64 = 0
	for _, contest := range contests {
		result, err := q.InsertContest(ctx, repository.InsertContestParams{
			ContestID:        contest.ID,
			StartEpochSecond: contest.StartEpochSecond,
			DurationSecond:   contest.DurationSecond,
			Title:            contest.Title,
			RateChange:       contest.RateChange,
			Category:         contest.Categorize(),
		})
		if err != nil {
			return fmt.Errorf("save contest: %w", err)
		}

		count += result.RowsAffected()
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("save contests: %w", err)
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "Finish saving contests.", slog.Int64("count", count))

	return nil
}
