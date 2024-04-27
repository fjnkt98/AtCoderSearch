package generate

import (
	"fjnkt98/atcodersearch/batch/generate"
	"fjnkt98/atcodersearch/repository"

	"github.com/goark/errs"
	"github.com/urfave/cli/v2"
)

func newGenerateProblemCmd() *cli.Command {
	return &cli.Command{
		Name: "problem",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:    "save-dir",
				EnvVars: []string{"PROBLEM_SAVE_DIR"},
			},
			&cli.IntFlag{
				Name:  "chunk-size",
				Value: 1000,
			},
			&cli.IntFlag{
				Name:  "concurrent",
				Value: 4,
			},
		},
		Action: func(ctx *cli.Context) error {
			pool, err := repository.NewPool(ctx.Context, ctx.String("database-url"))
			if err != nil {
				return errs.Wrap(err)
			}
			reader := generate.NewProblemRowReader(pool)
			generator := generate.NewProblemGenerator(
				reader,
				ctx.String("save-dir"),
				ctx.Int("chunk-size"),
				ctx.Int("concurrent"),
			)

			if err := generator.Generate(ctx.Context); err != nil {
				return errs.Wrap(err)
			}
			return nil
		},
	}
}
