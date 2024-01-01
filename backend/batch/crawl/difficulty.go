package crawl

import (
	"context"
	"fjnkt98/atcodersearch/batch"
	"fjnkt98/atcodersearch/pkg/atcoder"
	"fjnkt98/atcodersearch/repository"

	"github.com/goark/errs"
	"golang.org/x/exp/slog"
)

type DifficultyCrawler interface {
	batch.Batch
	CrawlDifficulty(ctx context.Context) error
}

type difficultyCrawler struct {
	client atcoder.AtCoderProblemsClient
	repo   repository.DifficultyRepository
}

func NewDifficultyCrawler(client atcoder.AtCoderProblemsClient, repo repository.DifficultyRepository) DifficultyCrawler {
	return &difficultyCrawler{
		client: client,
		repo:   repo,
	}
}

func (c *difficultyCrawler) Name() string {
	return "DifficultyCrawler"
}

func (c *difficultyCrawler) CrawlDifficulty(ctx context.Context) error {
	slog.Info("Start to crawl difficulties.")
	difficulties, err := c.client.FetchDifficulties(ctx)
	if err != nil {
		return errs.Wrap(err)
	}
	slog.Info("Finish crawling difficulties.")

	slog.Info("Start to save difficulties.")
	if err := c.repo.Save(ctx, convertDifficulties(difficulties)); err != nil {
		return errs.Wrap(err)
	}
	slog.Info("Finish saving difficulties.")

	return nil
}

func (c *difficultyCrawler) Run(ctx context.Context) error {
	return c.CrawlDifficulty(ctx)
}

func convertDifficulties(difficulties map[string]atcoder.Difficulty) []repository.Difficulty {
	result := make([]repository.Difficulty, 0, len(difficulties))
	for problemID, difficulty := range difficulties {
		result = append(result, repository.Difficulty{
			ProblemID:        problemID,
			Slope:            difficulty.Slope,
			Intercept:        difficulty.Intercept,
			Variance:         difficulty.Variance,
			Difficulty:       difficulty.Difficulty,
			Discrimination:   difficulty.Discrimination,
			IrtLogLikelihood: difficulty.IrtLogLikelihood,
			IrtUsers:         difficulty.IrtUsers,
			IsExperimental:   difficulty.IsExperimental,
		})
	}
	return result
}
