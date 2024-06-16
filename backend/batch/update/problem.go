package update

import (
	"context"
	"encoding/json"
	"fjnkt98/atcodersearch/batch/crawl"
	"fjnkt98/atcodersearch/batch/generate"
	"fjnkt98/atcodersearch/batch/post"
	"fjnkt98/atcodersearch/pkg/atcoder"
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/repository"
	"fjnkt98/atcodersearch/settings"
	"log/slog"
	"time"

	"github.com/goark/errs"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UpdateProblemConfig struct {
	Duration           time.Duration `json:"duration"`
	All                bool          `json:"all"`
	SkipFetch          bool          `json:"skip-fetch"`
	SaveDir            string        `json:"save-dir"`
	ChunkSize          int           `json:"chunk-size"`
	GenerateConcurrent int           `json:"generate-concurrent"`
	PostConcurrent     int           `json:"post-concurrent"`
	Optimize           bool          `json:"optimize"`
}

func UpdateProblem(ctx context.Context, pool *pgxpool.Pool, core *solr.SolrCore, config UpdateProblemConfig) error {
	slog.Info("Start Batch", slog.String("name", settings.UPDATE_PROBLEM_BATCH_NAME), slog.Any("config", config))
	options, err := json.Marshal(config)
	if err != nil {
		return errs.New("failed to marshal update problem config", errs.WithCause(err))
	}

	h, err := repository.NewBatchHistory(ctx, pool, settings.UPDATE_PROBLEM_BATCH_NAME, options)
	if err != nil {
		return errs.Wrap(err, errs.WithCause(err), errs.WithContext("name", settings.UPDATE_PROBLEM_BATCH_NAME), errs.WithContext("config", config))
	}
	defer h.Fail(ctx, pool)

	if !config.SkipFetch {
		problemsClient := atcoder.NewAtCoderProblemsClient()
		atcoderClient, err := atcoder.NewAtCoderClient()
		if err != nil {
			return errs.Wrap(err, errs.WithContext("name", settings.UPDATE_PROBLEM_BATCH_NAME), errs.WithContext("config", config))
		}
		if err := crawl.NewContestCrawler(problemsClient, pool).Crawl(ctx); err != nil {
			return errs.Wrap(err, errs.WithContext("name", settings.UPDATE_PROBLEM_BATCH_NAME), errs.WithContext("config", config))
		}
		if err := crawl.NewDifficultyCrawler(problemsClient, pool).Crawl(ctx); err != nil {
			return errs.Wrap(err, errs.WithContext("name", settings.UPDATE_PROBLEM_BATCH_NAME), errs.WithContext("config", config))
		}
		if err := crawl.NewProblemCrawler(problemsClient, atcoderClient, pool, config.Duration, config.All).Crawl(ctx); err != nil {
			return errs.Wrap(err, errs.WithContext("name", settings.UPDATE_PROBLEM_BATCH_NAME), errs.WithContext("config", config))
		}
	}

	if err := generate.GenerateProblemDocument(
		ctx,
		generate.NewProblemRowReader(pool),
		config.SaveDir,
		generate.WithChunkSize(config.ChunkSize),
		generate.WithConcurrent(config.GenerateConcurrent),
	); err != nil {
		return errs.Wrap(err, errs.WithContext("name", settings.UPDATE_PROBLEM_BATCH_NAME), errs.WithContext("config", config))
	}

	if err := post.PostDocument(
		ctx,
		core,
		config.SaveDir,
		post.WithConcurrent(config.PostConcurrent),
		post.WithOptimize(config.Optimize),
		post.WithTruncate(true),
	); err != nil {
		if errs.Is(err, post.ErrNoFiles) {
			slog.Info("there is no files to post", slog.Any("detail", err))
		} else {
			return errs.Wrap(err, errs.WithContext("name", settings.UPDATE_PROBLEM_BATCH_NAME), errs.WithContext("config", config))
		}
	}

	if err := h.Finish(ctx, pool); err != nil {
		return errs.Wrap(err, errs.WithCause(err), errs.WithContext("name", settings.UPDATE_PROBLEM_BATCH_NAME), errs.WithContext("config", config))
	}

	slog.Info("Finish Batch", slog.String("name", settings.UPDATE_PROBLEM_BATCH_NAME), slog.Any("config", config))
	return nil
}
