package update

import (
	"context"
	"encoding/json"
	"fjnkt98/atcodersearch/batch/generate"
	"fjnkt98/atcodersearch/batch/post"
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/repository"
	"log/slog"
	"time"

	"github.com/goark/errs"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UpdateSubmissionConfig struct {
	Duration           time.Duration `json:"duration"`
	Retry              int           `json:"retry"`
	All                bool          `json:"all"`
	Interval           int           `json:"interval"`
	SaveDir            string        `json:"save-dir"`
	ChunkSize          int           `json:"chunk-size"`
	GenerateConcurrent int           `json:"generate-concurrent"`
	PostConcurrent     int           `json:"post-concurrent"`
	Optimize           bool          `json:"optimize"`
}

func UpdateSubmission(ctx context.Context, pool *pgxpool.Pool, core *solr.SolrCore, config UpdateSubmissionConfig) error {
	slog.Info("Start UpdateSubmission")
	options, err := json.Marshal(config)
	if err != nil {
		return errs.New("failed to marshal update submission config", errs.WithCause(err))
	}

	q := repository.New(pool)

	id, err := q.CreateBatchHistory(ctx, repository.CreateBatchHistoryParams{Name: "UpdateSubmission", Options: options})
	if err != nil {
		return errs.New("failed to create batch history", errs.WithCause(err))
	}

	if err := generate.GenerateSubmissionDocument(
		ctx,
		generate.NewSubmissionRowReader(pool, config.Interval, config.All),
		config.SaveDir,
		generate.WithChunkSize(config.ChunkSize),
		generate.WithConcurrent(config.GenerateConcurrent),
	); err != nil {
		return errs.Wrap(err)
	}

	if err := post.PostDocument(
		ctx,
		core,
		config.SaveDir,
		post.WithConcurrent(config.PostConcurrent),
		post.WithOptimize(config.Optimize),
		post.WithTruncate(true),
	); err != nil {
		return errs.Wrap(err)
	}

	if err := q.UpdateBatchHistory(ctx, repository.UpdateBatchHistoryParams{ID: id, Status: "finished"}); err != nil {
		return errs.New("failed to update batch history", errs.WithCause(err))
	}

	slog.Info("Finish UpdateSubmission")
	return nil
}
