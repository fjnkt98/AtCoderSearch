package crawl

import (
	"context"
	"fjnkt98/atcodersearch/batch"
	"fjnkt98/atcodersearch/pkg/atcoder"
	"fjnkt98/atcodersearch/repository"

	"log/slog"

	"github.com/goark/errs"
)

type ContestCrawler interface {
	batch.Batch
	CrawlContest(ctx context.Context) error
}

type contestCrawler struct {
	client atcoder.AtCoderProblemsClient
	repo   repository.ContestRepository
}

func NewContestCrawler(client atcoder.AtCoderProblemsClient, repo repository.ContestRepository) ContestCrawler {
	return &contestCrawler{
		client: client,
		repo:   repo,
	}
}

func (c *contestCrawler) Name() string {
	return "ContestCrawler"
}

func (c *contestCrawler) Config() any {
	return nil
}

func (c *contestCrawler) CrawlContest(ctx context.Context) error {
	slog.Info("Start to fetch contests.")
	contests, err := c.client.FetchContests(ctx)
	if err != nil {
		return errs.Wrap(err)
	}
	slog.Info("Finish fetching contests.")

	slog.Info("Start to save contests.")
	if err := c.repo.Save(ctx, convertContests(contests)); err != nil {
		return errs.Wrap(err)
	}
	slog.Info("Finish saving contest list.")

	return nil
}

func (c *contestCrawler) Run(ctx context.Context) error {
	return c.CrawlContest(ctx)
}

func convertContest(contest atcoder.Contest) repository.Contest {
	return repository.Contest{
		ContestID:        contest.ID,
		StartEpochSecond: contest.StartEpochSecond,
		DurationSecond:   contest.DurationSecond,
		Title:            contest.Title,
		RateChange:       contest.RateChange,
		Category:         contest.Categorize(),
	}
}

func convertContests(contests []atcoder.Contest) []repository.Contest {
	result := make([]repository.Contest, len(contests))
	for i, c := range contests {
		result[i] = convertContest(c)
	}

	return result
}
