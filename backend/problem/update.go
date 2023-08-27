package problem

import (
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
	All                bool   `json:"all"`
}

func Update(cfg UpdateConfig, db *sqlx.DB, core *solr.Core) error {
	options, err := json.Marshal(cfg)
	if err != nil {
		failure.Translate(err, EncodeError, failure.Message("failed to encode update config"))
	}
	history := acs.NewUpdateHistory(db, "problem", string(options))
	defer history.Cancel()

	if !cfg.SkipFetch {
		contestCrawler := NewContestCrawler(db)
		if err := contestCrawler.Run(); err != nil {
			return failure.Wrap(err)
		}

		difficultyCrawler := NewDifficultyCrawler(db)
		if err := difficultyCrawler.Run(); err != nil {
			return failure.Wrap(err)
		}

		problemCrawler := NewProblemCrawler(db)
		if err := problemCrawler.Run(cfg.All, cfg.Duration); err != nil {
			return failure.Wrap(err)
		}
	}

	generator := NewDocumentGenerator(db, cfg.SaveDir)
	if err := generator.Run(cfg.ChunkSize, cfg.GenerateConcurrent); err != nil {
		return failure.Wrap(err)
	}

	uploader := acs.NewDefaultDocumentUploader(core, cfg.SaveDir)
	if err := uploader.PostDocument(cfg.Optimize, true, cfg.PostConcurrent); err != nil {
		return failure.Wrap(err)
	}
	history.Finish()
	return nil
}
