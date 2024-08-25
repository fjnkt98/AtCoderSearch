package crawl

import (
	"context"
	"errors"
	"fmt"
	"time"

	"log/slog"

	"github.com/fjnkt98/atcodersearch-batch/pkg/atcoder"
	"github.com/fjnkt98/atcodersearch-batch/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"
)

func InterruptibleSleep(ctx context.Context, s int) {
	for i := 0; i < s; i++ {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

type SubmissionCrawler struct {
	client   *atcoder.AtCoderClient
	pool     *pgxpool.Pool
	duration time.Duration
	retry    int
	targets  []string
}

func NewSubmissionCrawler(
	client *atcoder.AtCoderClient,
	pool *pgxpool.Pool,
	duration time.Duration,
	retry int,
	targets []string,
) *SubmissionCrawler {
	return &SubmissionCrawler{
		client:   client,
		pool:     pool,
		duration: duration,
		retry:    retry,
		targets:  targets,
	}
}

func (c *SubmissionCrawler) crawlContest(ctx context.Context, contestID string) ([]atcoder.Submission, error) {
	q := repository.New(c.pool)

	lastCrawled, err := q.FetchLatestCrawlHistory(ctx, contestID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("last crawled history of `%s`: %w", contestID, err)
		}
	}
	startedAt := time.Now().Unix()

	slog.LogAttrs(ctx, slog.LevelInfo, "Start to crawl", slog.String("contestID", contestID), slog.Time("lastCrawled", time.Unix(lastCrawled, 0)))

	submissions := make([]atcoder.Submission, 0)
loop:
	for i := 1; i <= 1_000_000_000; i++ {
		slog.LogAttrs(ctx, slog.LevelInfo, "fetch submissions", slog.String("contestID", contestID), slog.Int("page", i))

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		subs, err := c.client.FetchSubmissions(ctx, contestID, i)
		if err != nil {
		retryLoop:
			for j := 0; err != nil && j < c.retry; j++ {
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				default:
				}

				slog.LogAttrs(ctx, slog.LevelError, "failed to crawl submissions. retry to crawl submission after 1 minutes...", slog.String("contestID", contestID), slog.Any("error", err))

				InterruptibleSleep(ctx, 60)
				subs, err = c.client.FetchSubmissions(ctx, contestID, i)
				if err == nil {
					break retryLoop
				}
			}

			if err != nil {
				return nil, fmt.Errorf("crawl submissions of `%s`: %w", contestID, err)
			}
		}

		if len(subs) == 0 {
			slog.LogAttrs(ctx, slog.LevelInfo, "There is no more submissions", slog.String("contestID", contestID))
			break loop
		}

		submissions = append(submissions, subs...)

		if subs[0].EpochSecond < lastCrawled {
			slog.LogAttrs(ctx, slog.LevelInfo, "Break crawling since all submissions after have been crawled.", slog.String("contestID", contestID), slog.Int("page", i))

			time.Sleep(c.duration)
			break loop
		}

		time.Sleep(c.duration)
	}

	if _, err := q.CreateCrawlHistory(ctx, repository.CreateCrawlHistoryParams{StartedAt: startedAt, ContestID: contestID}); err != nil {
		return nil, fmt.Errorf("save crawl history: %w", err)
	}

	return submissions, nil
}

func (c *SubmissionCrawler) saveSubmissions(ctx context.Context, submissions []atcoder.Submission) error {
	if len(submissions) == 0 {
		return nil
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "Start to save submissions", slog.String("contestID", submissions[0].ContestID))

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("save submissions: %w", err)
	}
	defer tx.Rollback(ctx)

	q := repository.New(tx)

	var count int64 = 0
	for _, s := range submissions {
		result, err := q.InsertSubmission(ctx, repository.InsertSubmissionParams{
			ID:            s.ID,
			EpochSecond:   s.EpochSecond,
			ProblemID:     s.ProblemID,
			ContestID:     &s.ContestID,
			UserID:        &s.UserID,
			Language:      &s.Language,
			Point:         &s.Point,
			Length:        &s.Length,
			Result:        &s.Result,
			ExecutionTime: s.ExecutionTime,
		})
		if err != nil {
			return fmt.Errorf("save submissions: %w", err)
		}

		count += result.RowsAffected()
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("save submissions: %w", err)
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "Save submissions successfully", slog.String("contestID", submissions[0].ContestID), slog.Int64("count", count))

	return nil
}

func (c *SubmissionCrawler) Crawl(ctx context.Context) error {
	q := repository.New(c.pool)

	var ids []string
	var err error
	if len(c.targets) == 0 {
		ids, err = q.FetchContestIDs(ctx)
	} else {
		ids, err = q.FetchContestIDsByCategory(ctx, c.targets)
	}
	if err != nil {
		return fmt.Errorf("crawl submissions: %w", err)
	}

	targets := make(chan string, len(ids))
	for _, id := range ids {
		targets <- id
	}
	close(targets)
	submissions := make(chan []atcoder.Submission)

	eg, ctx := errgroup.WithContext(ctx)

	// crawl and send submissions of it
	eg.Go(func() error {
		defer close(submissions)
	loop:
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case contestID, ok := <-targets:
				if !ok {
					break loop
				}

				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
				}

				s, err := c.crawlContest(ctx, contestID)
				if err != nil {
					return err
				}
				submissions <- s
			}
		}
		return nil
	})

	// receive and save the submissions
	eg.Go(func() error {
	loop:
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case s, ok := <-submissions:
				if !ok {
					break loop
				}

				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
				}

				if err := c.saveSubmissions(ctx, s); err != nil {
					return err
				}
			}
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return err
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "Finish crawling successfully")
	return nil
}
