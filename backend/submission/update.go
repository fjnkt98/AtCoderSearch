package submission

import (
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

func Update(cfg UpdateConfig, db *sqlx.DB, core solr.SolrCore[any, any]) error {
	generator := NewDocumentGenerator(db, cfg.SaveDir)
	if err := generator.Run(cfg.ChunkSize, cfg.GenerateConcurrent); err != nil {
		return failure.Wrap(err)
	}

	uploader := acs.NewDefaultDocumentUploader(core, cfg.SaveDir)
	if err := uploader.PostDocument(cfg.Optimize, cfg.PostConcurrent); err != nil {
		return failure.Wrap(err)
	}
	return nil
}
