package update

import (
	"context"
	"encoding/json"
	"fjnkt98/atcodersearch/batch"
	"fjnkt98/atcodersearch/batch/crawl"
	"fjnkt98/atcodersearch/batch/generate"
	"fjnkt98/atcodersearch/batch/repository"
	"fjnkt98/atcodersearch/batch/upload"
	"fjnkt98/atcodersearch/config"
	"log/slog"

	"github.com/goark/errs"
)

type ProblemUpdater interface {
	batch.Batch
}

type problemUpdater struct {
	cfg               config.ProblemConfig
	problemCrawler    crawl.ProblemCrawler
	contestCrawler    crawl.ContestCrawler
	difficultyCrawler crawl.DifficultyCrawler
	generator         generate.ProblemGenerator
	uploader          upload.DocumentUploader
	repo              repository.UpdateHistoryRepository
}

func NewProblemUpdater(
	cfg config.ProblemConfig,
	problemCrawler crawl.ProblemCrawler,
	contestCrawler crawl.ContestCrawler,
	difficultyCrawler crawl.DifficultyCrawler,
	generator generate.ProblemGenerator,
	uploader upload.DocumentUploader,
	repo repository.UpdateHistoryRepository,
) ProblemUpdater {
	return &problemUpdater{
		cfg:               cfg,
		problemCrawler:    problemCrawler,
		contestCrawler:    contestCrawler,
		difficultyCrawler: difficultyCrawler,
		generator:         generator,
		uploader:          uploader,
		repo:              repo,
	}
}

func (u *problemUpdater) Name() string {
	return "ProblemUpdater"
}

func (u *problemUpdater) Run(ctx context.Context) error {
	cfg, err := json.Marshal(u.cfg)
	if err != nil {
		return errs.New(
			"failed to encode update config",
			errs.WithCause(err),
		)
	}

	history := repository.NewUpdateHistory("problem", string(cfg))
	defer u.repo.Cancel(ctx, &history)

	slog.Info("Start to update problem index.")
	if u.cfg.Update.SkipFetch {
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
