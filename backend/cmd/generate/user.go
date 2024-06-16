package generate

import (
	"fjnkt98/atcodersearch/batch/generate"
	"fjnkt98/atcodersearch/repository"

	"github.com/goark/errs"
	"github.com/urfave/cli/v2"
)

func newGenerateUserCmd() *cli.Command {
	return &cli.Command{
		Name: "user",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:    "save-dir",
				EnvVars: []string{"USER_SAVE_DIR"},
			},
			&cli.IntFlag{
				Name:  "chunk-size",
				Value: 10000,
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
			reader := generate.NewUserRowReader(pool)
			err = generate.GenerateUserDocument(
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
