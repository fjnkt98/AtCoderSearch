package submission

import (
	"context"
	"encoding/json"
	"fjnkt98/atcodersearch/acs"
	"fjnkt98/atcodersearch/solr"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/morikuni/failure"
)

type UpdateConfig struct {
	SaveDir            string `json:"save-dir"`
	Optimize           bool   `json:"optimize"`
	ChunkSize          int    `json:"chunk-size"`
	GenerateConcurrent int    `json:"generate-concurrent"`
	PostConcurrent     int    `json:"post-concurrent"`
	All                bool   `json:"all"`
}

func Update(ctx context.Context, cfg UpdateConfig, db *sqlx.DB, core *solr.Core) error {
	options, err := json.Marshal(cfg)
	if err != nil {
		failure.Translate(err, acs.EncodeError, failure.Message("failed to encode update config"))
	}
	history := acs.NewUpdateHistory(db, "submission", string(options))
	defer history.Cancel()

	period, err := history.GetLatest()
	if err != nil {
		return failure.Wrap(err)
	}

	if cfg.All {
		period = time.Time{}
	}
	if err := Generate(ctx, db, cfg.SaveDir, cfg.ChunkSize, cfg.GenerateConcurrent, period); err != nil {
		return failure.Wrap(err)
	}

	if err := acs.PostDocument(ctx, core, cfg.SaveDir, cfg.Optimize, cfg.All, cfg.PostConcurrent); err != nil {
		return failure.Wrap(err)
	}

	history.Finish()
	return nil
}
