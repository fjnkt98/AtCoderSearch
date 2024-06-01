package crawl

import (
	"context"
	"fjnkt98/atcodersearch/pkg/atcoder"
	"fjnkt98/atcodersearch/repository"
	"fmt"
	"time"

	"log/slog"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/goark/errs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

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

func (c *SubmissionCrawler) crawlContest(ctx context.Context, contestID string, lastCrawled int64) error {
	submissions := make([]atcoder.Submission, 0)
loop:
	for i := 1; i <= 1_000_000_000; i++ {
		slog.Info("fetch submissions", slog.String("contest id", contestID), slog.Int("page", i))
		subs, err := c.client.FetchSubmissions(ctx, contestID, i)
		if err != nil {
		retryLoop:
			for j := 0; err != nil && j < c.retry; j++ {
				select {
				case <-ctx.Done():
					return nil
				default:
					slog.Error("failed to crawl submission", slog.String("contestID", contestID), slog.String("error", fmt.Sprintf("%+v", err)))
					slog.Info("retry to crawl submission after 1 minutes...")
					time.Sleep(time.Duration(60) * time.Second)
					subs, err = c.client.FetchSubmissions(ctx, contestID, i)
					if err == nil {
						break retryLoop
					}
				}
			}

			if err != nil {
				return errs.New(
					"failed to crawl submissions",
					errs.WithCause(err),
					errs.WithContext("contest id", contestID),
				)
			}
		}

		if len(subs) == 0 {
			slog.Info("There is no more submissions", slog.String("contest id", contestID))
			break loop
		}

		submissions = append(submissions, subs...)

		if subs[0].EpochSecond < lastCrawled {
			slog.Info("Break crawling since all submissions after have been crawled.", slog.String("contest id", contestID), slog.Int("page", i))
			time.Sleep(c.duration)
			break loop
		}

		time.Sleep(c.duration)
	}

	if len(submissions) == 0 {
		slog.Info("There is no submissions to save", slog.String("contest id", contestID))
		return nil
	}

	count, err := repository.BulkUpdate(ctx, c.pool, "submissions", convertSubmissions(dedupSubmissions(submissions)))
	if err != nil {
		return errs.New("failed to bulk update submissions", errs.WithCause(err))
	}
	slog.Info("Save submissions successfully", slog.String("contest id", contestID), slog.Int64("count", count))
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
		return errs.New("failed to fetch contest categories", errs.WithCause(err), errs.WithContext("targets", c.targets))
	}

	for _, id := range ids {
		lastCrawled, err := q.FetchLatestCrawlHistory(ctx, id)
		if err != nil && !errs.Is(err, pgx.ErrNoRows) {
			return errs.New("failed to fetch latest crawl history", errs.WithCause(err), errs.WithContext("contest id", id))
		}

		startedAt := time.Now().Unix()
		slog.Info("Start to crawl", slog.String("contest id", id), slog.Time("last crawled", time.Unix(lastCrawled, 0)))
		if err := c.crawlContest(ctx, id, lastCrawled); err != nil {
			return errs.Wrap(err)
		}

		if _, err := q.CreateCrawlHistory(ctx, repository.CreateCrawlHistoryParams{StartedAt: startedAt, ContestID: id}); err != nil {
			return errs.New("failed to create crawl history", errs.WithCause(err), errs.WithContext("contest id", id))
		}
		time.Sleep(c.duration)
	}
	return nil
}

func dedupSubmissions(submissions []atcoder.Submission) []atcoder.Submission {
	result := make([]atcoder.Submission, 0, len(submissions))
	ids := mapset.NewSet[int64]()
	for _, s := range submissions {
		if ids.Contains(s.ID) {
			continue
		}
		ids.Add(s.ID)
		result = append(result, s)
	}
	return result
}

func convertSubmissions(submissions []atcoder.Submission) []repository.Submission {
	result := make([]repository.Submission, len(submissions))
	for i, s := range submissions {
		result[i] = repository.Submission{
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
		}
	}
	return result
}
