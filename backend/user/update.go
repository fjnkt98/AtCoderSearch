package user

import (
	"context"
	"encoding/json"
	"fjnkt98/atcodersearch/acs"
	"fjnkt98/atcodersearch/solr"

	"github.com/jmoiron/sqlx"
	"github.com/morikuni/failure"
)

type UpdateConfig struct {
	SaveDir            string `json:"save-dir"`
	SkipFetch          bool   `json:"skip-fetch"`
	Optimize           bool   `json:"optimize"`
	ChunkSize          int    `json:"chunk-size"`
	GenerateConcurrent int    `json:"generate-concurrent"`
	PostConcurrent     int    `json:"post-concurrent"`
	Duration           int    `json:"duration"`
}

func Update(ctx context.Context, cfg UpdateConfig, db *sqlx.DB, core *solr.Core) error {
	options, err := json.Marshal(cfg)
	if err != nil {
		failure.Translate(err, acs.EncodeError, failure.Message("failed to encode update config"))
	}
	history := acs.NewUpdateHistory(db, "user", string(options))
	defer history.Cancel()

	if !cfg.SkipFetch {
		crawler := NewUserCrawler(db)
		if err := crawler.Run(ctx, cfg.Duration); err != nil {
			return failure.Wrap(err)
		}
	}

	if err := Generate(ctx, db, cfg.SaveDir, cfg.ChunkSize, cfg.GenerateConcurrent); err != nil {
		return failure.Wrap(err)
	}

	if err := acs.PostDocument(ctx, core, cfg.SaveDir, cfg.Optimize, true, cfg.PostConcurrent); err != nil {
		return failure.Wrap(err)
	}

	history.Finish()
	return nil
}
