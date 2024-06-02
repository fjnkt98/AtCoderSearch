package update

import (
	"fjnkt98/atcodersearch/batch/update"
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/repository"
	"fjnkt98/atcodersearch/settings"
	"time"

	"github.com/goark/errs"
	"github.com/urfave/cli/v2"
)

func newUpdateSubmissionCmd() *cli.Command {
	return &cli.Command{
		Name: "submission",
		Flags: []cli.Flag{
			&cli.DurationFlag{
				Name:     "duration",
				Value:    1000 * time.Millisecond,
				Category: "crawl",
			},
			&cli.IntFlag{
				Name:     "retry",
				Value:    0,
				Category: "crawl",
			},
			&cli.StringFlag{
				Name:     "target",
				Category: "crawl",
			},
			&cli.BoolFlag{
				Name:     "all",
				Value:    false,
				Category: "generate",
			},
			&cli.IntFlag{
				Name:     "interval",
				Value:    90,
				Category: "generate",
			},
			&cli.PathFlag{
				Name:     "save-dir",
				Category: "generate, post",
				EnvVars:  []string{"SUBMISSION_SAVE_DIR"},
			},
			&cli.IntFlag{
				Name:     "chunk-size",
				Value:    50000,
				Category: "generate",
			},
			&cli.IntFlag{
				Name:     "generate-concurrent",
				Value:    4,
				Category: "generate",
			},
			&cli.IntFlag{
				Name:     "post-concurrent",
				Value:    2,
				Category: "post",
			},
			&cli.BoolFlag{
				Name:     "optimize",
				Value:    false,
				Category: "post",
			},
		},
		Action: func(ctx *cli.Context) error {
			pool, err := repository.NewPool(ctx.Context, ctx.String("database-url"))
			if err != nil {
				return errs.Wrap(err)
			}
			core, err := solr.NewSolrCore(ctx.String("solr-host"), settings.SUBMISSION_CORE_NAME)
			if err != nil {
				return errs.Wrap(err)
			}

			config := update.UpdateSubmissionConfig{
				Duration:           ctx.Duration("duration"),
				Retry:              ctx.Int("retry"),
				All:                ctx.Bool("all"),
				Interval:           ctx.Int("interval"),
				SaveDir:            ctx.String("save-dir"),
				ChunkSize:          ctx.Int("chunk-size"),
				GenerateConcurrent: ctx.Int("generate-concurrent"),
				PostConcurrent:     ctx.Int("post-concurrent"),
				Optimize:           ctx.Bool("optimize"),
			}

			if err := update.UpdateSubmission(ctx.Context, pool, core, config); err != nil {
				return errs.Wrap(err)
			}
			return nil
		},
	}
}
