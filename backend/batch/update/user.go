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
	"log/slog"
	"time"

	"github.com/goark/errs"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UpdateUserConfig struct {
	Duration           time.Duration `json:"duration"`
	SkipFetch          bool          `json:"skip-fetch"`
	SaveDir            string        `json:"save-dir"`
	ChunkSize          int           `json:"chunk-size"`
	GenerateConcurrent int           `json:"generate-concurrent"`
	PostConcurrent     int           `json:"post-concurrent"`
	Optimize           bool          `json:"optimize"`
}

func UpdateUser(ctx context.Context, pool *pgxpool.Pool, core *solr.SolrCore, config UpdateUserConfig) error {
	slog.Info("Start UpdateUser")
	options, err := json.Marshal(config)
	if err != nil {
		return errs.New("failed to marshal update problem config", errs.WithCause(err))
	}

	q := repository.New(pool)

	id, err := q.CreateBatchHistory(ctx, repository.CreateBatchHistoryParams{Name: "UpdateUser", Options: options})
	if err != nil {
		return errs.New("failed to create batch history", errs.WithCause(err))
	}

	if !config.SkipFetch {
		client, err := atcoder.NewAtCoderClient()
		if err != nil {
			return errs.Wrap(err)
		}
		if err := crawl.NewUserCrawler(client, pool, config.Duration).Crawl(ctx); err != nil {
			return errs.Wrap(err)
		}
	}

	if err := generate.GenerateUserDocument(
		ctx,
		generate.NewUserRowReader(pool),
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

	slog.Info("Finish UpdateUser")
	return nil
}
