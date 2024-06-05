package update

import (
	"context"
	"encoding/json"
	"fjnkt98/atcodersearch/batch/generate"
	"fjnkt98/atcodersearch/batch/post"
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/repository"
	"fjnkt98/atcodersearch/settings"
	"log/slog"
	"time"

	"github.com/goark/errs"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UpdateSolutionConfig struct {
	All                bool   `json:"all"`
	Interval           int    `json:"interval"`
	SaveDir            string `json:"save-dir"`
	ChunkSize          int    `json:"chunk-size"`
	GenerateConcurrent int    `json:"generate-concurrent"`
	PostConcurrent     int    `json:"post-concurrent"`
	Optimize           bool   `json:"optimize"`
}

func UpdateSolution(ctx context.Context, pool *pgxpool.Pool, core *solr.SolrCore, config UpdateSolutionConfig) error {
	slog.Info("Start Batch", slog.String("name", settings.UPDATE_SOLUTION_BATCH_NAME), slog.Any("config", config))
	options, err := json.Marshal(config)
	if err != nil {
		return errs.New("failed to marshal update solution config", errs.WithCause(err), errs.WithContext("config", config))
	}

	h, err := repository.NewBatchHistory(ctx, pool, settings.UPDATE_SOLUTION_BATCH_NAME, options)
	if err != nil {
		return errs.Wrap(err, errs.WithCause(err), errs.WithContext("name", settings.UPDATE_SOLUTION_BATCH_NAME), errs.WithContext("config", config))
	}
	defer h.Fail(ctx, pool)

	var lastUpdated *time.Time = nil
	if !config.All {
		q := repository.New(pool)
		h, err := q.FetchLatestBatchHistory(ctx, settings.UPDATE_SOLUTION_BATCH_NAME)
		if err != nil {
			return errs.Wrap(err, errs.WithContext("name", settings.UPDATE_SOLUTION_BATCH_NAME), errs.WithContext("config", config))
		}
		lastUpdated = &h.StartedAt
	}

	if err := generate.GenerateSolutionDocument(
		ctx,
		generate.NewSolutionRowReader(pool, config.Interval, lastUpdated),
		config.SaveDir,
		generate.WithChunkSize(config.ChunkSize),
		generate.WithConcurrent(config.GenerateConcurrent),
	); err != nil {
		return errs.Wrap(err, errs.WithContext("name", settings.UPDATE_SOLUTION_BATCH_NAME), errs.WithContext("config", config))
	}

	if err := post.PostDocument(
		ctx,
		core,
		config.SaveDir,
		post.WithConcurrent(config.PostConcurrent),
		post.WithOptimize(config.Optimize),
		post.WithTruncate(config.All),
	); err != nil {
		return errs.Wrap(err, errs.WithContext("name", settings.UPDATE_SOLUTION_BATCH_NAME), errs.WithContext("config", config))
	}

	if err := h.Finish(ctx, pool); err != nil {
		return errs.Wrap(err, errs.WithCause(err), errs.WithContext("name", settings.UPDATE_SOLUTION_BATCH_NAME), errs.WithContext("config", config))
	}

	slog.Info("Finish Batch", slog.String("name", settings.UPDATE_SOLUTION_BATCH_NAME), slog.Any("config", config))
	return nil
}
