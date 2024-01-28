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

type UserUpdater interface {
	batch.Batch
}

type userUpdater struct {
	crawler   crawl.UserCrawler
	generator generate.UserGenerator
	uploader  upload.DocumentUploader
	repo      repository.UpdateHistoryRepository
	skipFetch bool
}

func NewUserUpdater(
	crawler crawl.UserCrawler,
	generator generate.UserGenerator,
	uploader upload.DocumentUploader,
	repo repository.UpdateHistoryRepository,
	skipFetch bool,
) UserUpdater {
	return &userUpdater{
		crawler:   crawler,
		generator: generator,
		uploader:  uploader,
		repo:      repo,
		skipFetch: skipFetch,
	}
}

func (u *userUpdater) Name() string {
	return "UserUpdater"
}

func (u *userUpdater) Config() any {
	config := map[string]any{
		"crawl":      u.crawler.Config(),
		"generate":   u.generator.Config(),
		"uploader":   u.uploader.Config(),
		"skip_fetch": u.skipFetch,
	}
	return config
}

func (u *userUpdater) Run(ctx context.Context) error {
	config, err := json.Marshal(u.Config())
	if err != nil {
		return errs.New(
			"failed to encode update config",
			errs.WithCause(err),
		)
	}

	history := repository.NewUpdateHistory("user", string(config))
	defer u.repo.Cancel(ctx, &history)

	slog.Info("Start to update user index.")
	if u.skipFetch {
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
