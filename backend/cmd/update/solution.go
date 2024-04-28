package update

import (
	"fjnkt98/atcodersearch/batch/update"
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/repository"

	"github.com/goark/errs"
	"github.com/urfave/cli/v2"
)

func newUpdateSolutionCmd() *cli.Command {
	return &cli.Command{
		Name: "solution",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:     "save-dir",
				Category: "generate, post",
				EnvVars:  []string{"SOLUTION_SAVE_DIR"},
			},
			&cli.IntFlag{
				Name:     "chunk-size",
				Value:    1000,
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
			core, err := solr.NewSolrCore(ctx.String("solr-host"), "solution")
			if err != nil {
				return errs.Wrap(err)
			}

			config := update.UpdateSolutionConfig{
				SaveDir:            ctx.String("save-dir"),
				ChunkSize:          ctx.Int("chunk-size"),
				GenerateConcurrent: ctx.Int("generate-concurrent"),
				PostConcurrent:     ctx.Int("post-concurrent"),
				Optimize:           ctx.Bool("optimize"),
			}

			if err := update.UpdateSolution(ctx.Context, pool, core, config); err != nil {
				return errs.Wrap(err)
			}
			return nil
		},
	}
}
