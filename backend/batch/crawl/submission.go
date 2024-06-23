package crawl

import (
	"context"
	"fjnkt98/atcodersearch/pkg/atcoder"
	"fjnkt98/atcodersearch/repository"
	"time"

	"log/slog"

	"github.com/goark/errs"
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

func (c *SubmissionCrawler) crawl(ctx context.Context, contestID string) ([]atcoder.Submission, error) {
	q := repository.New(c.pool)
	lastCrawled, err := q.FetchLatestCrawlHistory(ctx, contestID)
	if err != nil && !errs.Is(err, pgx.ErrNoRows) {
		return nil, errs.New("failed to fetch latest crawl history", errs.WithCause(err), errs.WithContext("contestID", contestID))
	}
	startedAt := time.Now().Unix()

	slog.Info("Start to crawl", slog.String("contestID", contestID), slog.Time("lastCrawled", time.Unix(lastCrawled, 0)))

	submissions := make([]atcoder.Submission, 0)
loop:
	for i := 1; i <= 1_000_000_000; i++ {
		slog.Info("fetch submissions", slog.String("contestID", contestID), slog.Int("page", i))

		subs, err := c.client.FetchSubmissions(ctx, contestID, i)
		if err != nil {
		retryLoop:
			for j := 0; err != nil && j < c.retry; j++ {
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				default:
					slog.Error("failed to crawl submissions. retry to crawl submission after 1 minutes...", slog.String("contestID", contestID), slog.Any("error", err))

					InterruptibleSleep(ctx, 60)
					subs, err = c.client.FetchSubmissions(ctx, contestID, i)
					if err == nil {
						break retryLoop
					}
				}
			}

			if err != nil {
				return nil, errs.New(
					"failed to crawl submissions",
					errs.WithCause(err),
					errs.WithContext("contestID", contestID),
				)
			}
		}

		if len(subs) == 0 {
			slog.Info("There is no more submissions", slog.String("contestID", contestID))
			break loop
		}

		submissions = append(submissions, subs...)

		if subs[0].EpochSecond < lastCrawled {
			slog.Info("Break crawling since all submissions after have been crawled.", slog.String("contestID", contestID), slog.Int("page", i))

			time.Sleep(c.duration)
			break loop
		}

		time.Sleep(c.duration)
	}

	if _, err := q.CreateCrawlHistory(ctx, repository.CreateCrawlHistoryParams{StartedAt: startedAt, ContestID: contestID}); err != nil {
		return nil, errs.New("failed to create crawl history", errs.WithCause(err), errs.WithContext("contestID", contestID))
	}

	return submissions, nil
}

func (c *SubmissionCrawler) save(ctx context.Context, submissions []atcoder.Submission) error {
	if len(submissions) == 0 {
		return nil
	}

	slog.Info("Start to save submissions", slog.String("contestID", submissions[0].ContestID))
	count, err := repository.BulkUpdate(ctx, c.pool, "submissions", repository.NewSubmissions(submissions))
	if err != nil {
		return errs.New("failed to bulk update submissions", errs.WithCause(err))
	}
	slog.Info("Save submissions successfully", slog.String("contestID", submissions[0].ContestID), slog.Int64("count", count))
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

				s, err := c.crawl(ctx, contestID)
				if err != nil {
					return errs.Wrap(err)
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

				if err := c.save(ctx, s); err != nil {
					return errs.Wrap(err)
				}
			}
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return errs.Wrap(err)
	}

	slog.Info("Finish crawling successfully")
	return nil
}
