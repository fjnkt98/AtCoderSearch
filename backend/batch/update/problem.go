package update

import (
	"context"
	"encoding/json"
	"fjnkt98/atcodersearch/batch"
	"fjnkt98/atcodersearch/batch/crawl"
	"fjnkt98/atcodersearch/batch/generate"
	"fjnkt98/atcodersearch/batch/upload"
	"fjnkt98/atcodersearch/repository"
	"log/slog"

	"github.com/goark/errs"
)

type ProblemUpdater interface {
	batch.Batch
}

type problemUpdater struct {
	problemCrawler    crawl.ProblemCrawler
	contestCrawler    crawl.ContestCrawler
	difficultyCrawler crawl.DifficultyCrawler
	generator         generate.ProblemGenerator
	uploader          upload.DocumentUploader
	repo              repository.UpdateHistoryRepository
	skipFetch         bool
}

func NewProblemUpdater(
	problemCrawler crawl.ProblemCrawler,
	contestCrawler crawl.ContestCrawler,
	difficultyCrawler crawl.DifficultyCrawler,
	generator generate.ProblemGenerator,
	uploader upload.DocumentUploader,
	repo repository.UpdateHistoryRepository,
	skipFetch bool,
) ProblemUpdater {
	return &problemUpdater{
		problemCrawler:    problemCrawler,
		contestCrawler:    contestCrawler,
		difficultyCrawler: difficultyCrawler,
		generator:         generator,
		uploader:          uploader,
		repo:              repo,
		skipFetch:         skipFetch,
	}
}

func (u *problemUpdater) Name() string {
	return "ProblemUpdater"
}

func (u *problemUpdater) Config() any {
	config := map[string]any{
		"crawl": map[string]any{
			"problem":    u.problemCrawler.Config(),
			"contest":    u.contestCrawler.Config(),
			"difficulty": u.difficultyCrawler.Config(),
		},
		"generate":   u.generator.Config(),
		"upload":     u.uploader.Config(),
		"skip_fetch": u.skipFetch,
	}

	return config
}

func (u *problemUpdater) Run(ctx context.Context) error {
	config, err := json.Marshal(u.Config())
	if err != nil {
		return errs.New(
			"failed to encode update config",
			errs.WithCause(err),
		)
	}

	history := repository.NewUpdateHistory("problem", string(config))
	defer u.repo.Cancel(ctx, &history)

	slog.Info("Start to update problem index.")
	if u.skipFetch {
		slog.Info("Skip to crawl.")
	} else {
		if err := u.problemCrawler.CrawlProblem(ctx); err != nil {
			return errs.Wrap(err)
		}

		if err := u.contestCrawler.CrawlContest(ctx); err != nil {
			return errs.Wrap(err)
		}

		if err := u.difficultyCrawler.CrawlDifficulty(ctx); err != nil {
			return errs.Wrap(err)
		}
	}

	if err := u.generator.GenerateProblem(ctx); err != nil {
		return errs.Wrap(err)
	}

	if err := u.uploader.Upload(ctx); err != nil {
		return errs.Wrap(err)
	}

	if err := u.repo.Finish(ctx, &history); err != nil {
		return errs.Wrap(err)
	}
	slog.Info("Finished updating problem index successfully.")
	return nil
}
