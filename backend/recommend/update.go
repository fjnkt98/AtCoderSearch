package recommend

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
	Optimize           bool   `json:"optimize"`
	ChunkSize          int    `json:"chunk-size"`
	GenerateConcurrent int    `json:"generate-concurrent"`
	PostConcurrent     int    `json:"post-concurrent"`
}

func Update(ctx context.Context, cfg UpdateConfig, db *sqlx.DB, core *solr.Core) error {
	options, err := json.Marshal(cfg)
	if err != nil {
		return failure.Translate(err, acs.EncodeError, failure.Message("failed to encode update config"))
	}
	history := acs.NewUpdateHistory(db, "recommend", string(options))
	defer history.Cancel()

	if err := Generate(ctx, db, cfg.SaveDir, cfg.ChunkSize, cfg.GenerateConcurrent); err != nil {
		return failure.Translate(err, acs.UpdateIndexError, failure.Message("failed to update recommend index"))
	}

	if err := acs.PostDocument(ctx, core, cfg.SaveDir, cfg.Optimize, true, cfg.PostConcurrent); err != nil {
		return failure.Translate(err, acs.UpdateIndexError, failure.Message("failed to update recommend index"))
	}

	history.Finish()
	return nil
}
