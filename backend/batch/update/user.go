package update

import (
	"context"
	"encoding/json"
	"fjnkt98/atcodersearch/batch"
	"fjnkt98/atcodersearch/batch/crawl"
	"fjnkt98/atcodersearch/batch/generate"
	"fjnkt98/atcodersearch/batch/upload"
	"fjnkt98/atcodersearch/config"
	"fjnkt98/atcodersearch/repository"
	"log/slog"

	"github.com/goark/errs"
)

type UserUpdater interface {
	batch.Batch
}

type userUpdater struct {
	cfg       config.UserConfig
	crawler   crawl.UserCrawler
	generator generate.UserGenerator
	uploader  upload.DocumentUploader
	repo      repository.UpdateHistoryRepository
}

func NewUserUpdater(
	cfg config.UserConfig,
	crawler crawl.UserCrawler,
	generator generate.UserGenerator,
	uploader upload.DocumentUploader,
	repo repository.UpdateHistoryRepository,
) UserUpdater {
	return &userUpdater{
		cfg:       cfg,
		crawler:   crawler,
		generator: generator,
		uploader:  uploader,
		repo:      repo,
	}
}

func (u *userUpdater) Name() string {
	return "UserUpdater"
}

func (u *userUpdater) Run(ctx context.Context) error {
	cfg, err := json.Marshal(u.cfg)
	if err != nil {
		return errs.New(
			"failed to encode update config",
			errs.WithCause(err),
		)
	}

	history := repository.NewUpdateHistory("user", string(cfg))
	defer u.repo.Cancel(ctx, &history)

	slog.Info("Start to update user index.")
	if u.cfg.Update.SkipFetch {
		slog.Info("Skip to crawl.")
	} else {
		if err := u.crawler.CrawlUser(ctx); err != nil {
			return errs.Wrap(err)
		}
	}

	if err := u.generator.GenerateUser(ctx); err != nil {
		return errs.Wrap(err)
	}

	if err := u.uploader.Upload(ctx); err != nil {
		return errs.Wrap(err)
	}

	if err := u.repo.Finish(ctx, &history); err != nil {
		return errs.Wrap(err)
	}
	slog.Info("Finished updating user index successfully.")
	return nil
}
