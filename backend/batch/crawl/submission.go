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
)

type SubmissionCrawler interface {
	CrawlSubmission(ctx context.Context) error
}

type submissionCrawler struct {
	client         atcoder.AtCoderClient
	submissionRepo repository.SubmissionRepository
	contestRepo    repository.ContestRepository
	historyRepo    repository.SubmissionCrawlHistoryRepository
	duration       time.Duration
	retry          int
	targets        []string
	username       string
	password       string
}

func NewSubmissionCrawler(
	client atcoder.AtCoderClient,
	submissionRepo repository.SubmissionRepository,
	contestRepo repository.ContestRepository,
	historyRepo repository.SubmissionCrawlHistoryRepository,
	duration time.Duration,
	retry int,
	targets []string,
	username string,
	password string,
) SubmissionCrawler {
	return &submissionCrawler{
		client:         client,
		submissionRepo: submissionRepo,
		contestRepo:    contestRepo,
		historyRepo:    historyRepo,
		duration:       duration,
		retry:          retry,
		targets:        targets,
		username:       username,
		password:       password,
	}
}

func (c *submissionCrawler) crawl(ctx context.Context, contestID string, latest repository.SubmissionCrawlHistory) error {
	allSubmissions := make([]atcoder.Submission, 0)
loop:
	for i := 1; i <= 1_000_000_000; i++ {
		slog.Info(fmt.Sprintf("fetch submissions at page %d of the contest `%s`", i, contestID))
		submissions, err := c.client.FetchSubmissions(ctx, contestID, i)
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
					submissions, err = c.client.FetchSubmissions(ctx, contestID, i)
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

		if len(submissions) == 0 {
			slog.Info(fmt.Sprintf("There is no submissions in contest `%s`.", contestID))
			break loop
		}

		allSubmissions = append(allSubmissions, submissions...)

		if submissions[0].EpochSecond < int64(latest.StartedAt) {
			slog.Info(fmt.Sprintf("All submissions after page `%d` have been crawled. Break crawling the contest `%s`", i, contestID))
			time.Sleep(c.duration)
			break loop
		}

		time.Sleep(c.duration)
	}

	if len(allSubmissions) == 0 {
		slog.Info(fmt.Sprintf("No submissions to save for contest `%s`.", contestID))
		return nil
	}

	noDupSubmissions := make([]atcoder.Submission, 0, len(allSubmissions))
	ids := mapset.NewSet[int]()
	for _, s := range allSubmissions {
		if ids.Contains(s.ID) {
			continue
		}
		ids.Add(s.ID)
		noDupSubmissions = append(noDupSubmissions, s)
	}

	if err := c.submissionRepo.Save(ctx, convertSubmissions(noDupSubmissions)); err != nil {
		return errs.New(
			"failed to save submissions",
			errs.WithCause(err),
			errs.WithContext("contest id", contestID),
		)
	}
	slog.Info("Save submissions successfully", slog.String("contest id", contestID))
	return nil
}

func (c *submissionCrawler) CrawlSubmission(ctx context.Context) error {
	if err := c.client.Login(ctx, c.username, c.password); err != nil {
		return errs.Wrap(err)
	}

	ids, err := c.contestRepo.FetchContestIDs(ctx, c.targets)
	if err != nil {
		return errs.Wrap(err)
	}

	for _, id := range ids {
		history := repository.NewSubmissionCrawlHistory(id)
		latest, err := c.historyRepo.GetLatestHistory(ctx, id)
		if err != nil {
			return errs.Wrap(err)
		}

		slog.Info(fmt.Sprintf("Start to crawl contest `%s` since period `%s`", id, time.Unix(int64(latest.StartedAt), 0)))
		if err := c.crawl(ctx, id, latest); err != nil {
			return errs.Wrap(err)
		}
		if err := c.historyRepo.Save(ctx, history); err != nil {
			return errs.Wrap(err)
		}
		time.Sleep(c.duration)
	}
	return nil
}

func (c *submissionCrawler) Run(ctx context.Context) error {
	return c.CrawlSubmission(ctx)

}

func convertSubmission(submission atcoder.Submission) repository.Submission {
	return repository.Submission{
		ID:            submission.ID,
		EpochSecond:   submission.EpochSecond,
		ProblemID:     submission.ProblemID,
		ContestID:     submission.ContestID,
		UserID:        submission.UserID,
		Language:      submission.Language,
		Point:         submission.Point,
		Length:        submission.Length,
		Result:        submission.Result,
		ExecutionTime: submission.ExecutionTime,
	}
}

func convertSubmissions(submissions []atcoder.Submission) []repository.Submission {
	result := make([]repository.Submission, len(submissions))
	for i, submission := range submissions {
		result[i] = convertSubmission(submission)
	}

	return result
}
