package generate

import (
	"fjnkt98/atcodersearch/batch/generate"
	"fjnkt98/atcodersearch/repository"

	"github.com/goark/errs"
	"github.com/urfave/cli/v2"
)

func newGenerateSubmissionCmd() *cli.Command {
	return &cli.Command{
		Name: "submission",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:    "save-dir",
				EnvVars: []string{"SUBMISSION_SAVE_DIR"},
			},
			&cli.IntFlag{
				Name:  "chunk-size",
				Value: 10000,
			},
			&cli.IntFlag{
				Name:  "concurrent",
				Value: 4,
			},
			&cli.IntFlag{
				Name:  "interval",
				Value: 90,
			},
			&cli.BoolFlag{
				Name:  "all",
				Value: false,
			},
		},
		Action: func(ctx *cli.Context) error {
			pool, err := repository.NewPool(ctx.Context, ctx.String("database-url"))
			if err != nil {
				return errs.Wrap(err)
			}
			reader := generate.NewSubmissionRowReader(pool, ctx.Int("interval"), ctx.Bool("all"))

			err = generate.GenerateSubmissionDocument(
				ctx.Context,
				reader,
				ctx.String("save-dir"),
				generate.WithChunkSize(ctx.Int("chunk-size")),
				generate.WithConcurrent(ctx.Int("concurrent")),
			)
			if err != nil {
				return errs.Wrap(err)
			}
			return nil
		},
	}
}
