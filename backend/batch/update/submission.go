package update

import (
	"context"
	"encoding/json"
	"fjnkt98/atcodersearch/batch"
	"fjnkt98/atcodersearch/batch/generate"
	"fjnkt98/atcodersearch/batch/upload"
	"fjnkt98/atcodersearch/config"
	"fjnkt98/atcodersearch/repository"
	"log/slog"

	"github.com/goark/errs"
)

type SubmissionUpdater interface {
	batch.Batch
}

type submissionUpdater struct {
	cfg       config.SubmissionConfig
	generator generate.SubmissionGenerator
	uploader  upload.DocumentUploader
	repo      repository.UpdateHistoryRepository
}

func NewSubmissionUpdater(
	cfg config.SubmissionConfig,
	generator generate.SubmissionGenerator,
	uploader upload.DocumentUploader,
	repo repository.UpdateHistoryRepository,
) SubmissionUpdater {
	return &submissionUpdater{
		cfg:       cfg,
		generator: generator,
		uploader:  uploader,
		repo:      repo,
	}
}

func (u *submissionUpdater) Name() string {
	return "submissionUpdater"
}

func (u *submissionUpdater) Run(ctx context.Context) error {
	cfg, err := json.Marshal(u.cfg)
	if err != nil {
		return errs.New(
			"failed to encode update config",
			errs.WithCause(err),
		)
	}

	history := repository.NewUpdateHistory("submission", string(cfg))
	defer u.repo.Cancel(ctx, &history)

	slog.Info("Start to update submission index.")
	if err := u.generator.GenerateSubmission(ctx); err != nil {
		return errs.Wrap(err)
	}

	if err := u.uploader.Upload(ctx); err != nil {
		return errs.Wrap(err)
	}

	if err := u.repo.Finish(ctx, &history); err != nil {
		return errs.Wrap(err)
	}
	slog.Info("Finished updating submission index successfully.")
	return nil
}
