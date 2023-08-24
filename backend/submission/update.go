package submission

import (
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

func Update(cfg UpdateConfig, db *sqlx.DB, core *solr.Core) error {
	options, err := json.Marshal(cfg)
	if err != nil {
		failure.Translate(err, EncodeError, failure.Message("failed to encode update config"))
	}
	history := acs.NewUpdateHistory(db, "submission", string(options))
	defer history.Cancel()

	period, err := history.GetLatest()
	if err != nil {
		return failure.Wrap(err)
	}

	var generator DocumentGenerator
	if cfg.All {
		generator = NewDocumentGenerator(db, cfg.SaveDir, time.Time{})
	} else {
		generator = NewDocumentGenerator(db, cfg.SaveDir, period)
	}
	if err := generator.Run(cfg.ChunkSize, cfg.GenerateConcurrent); err != nil {
		return failure.Wrap(err)
	}

	uploader := acs.NewDefaultDocumentUploader(core, cfg.SaveDir)
	if err := uploader.PostDocument(cfg.Optimize, cfg.PostConcurrent); err != nil {
		return failure.Wrap(err)
	}

	history.Finish()
	return nil
}
