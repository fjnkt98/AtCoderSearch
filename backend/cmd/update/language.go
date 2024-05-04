package update

import (
	"fjnkt98/atcodersearch/batch/update"
	"fjnkt98/atcodersearch/repository"

	"github.com/goark/errs"
	"github.com/urfave/cli/v2"
)

func newUpdateLanguageCmd() *cli.Command {
	return &cli.Command{
		Name:  "language",
		Flags: []cli.Flag{},
		Action: func(ctx *cli.Context) error {
			pool, err := repository.NewPool(ctx.Context, ctx.String("database-url"))
			if err != nil {
				return errs.Wrap(err)
			}

			if err := update.UpdateLanguage(ctx.Context, pool); err != nil {
				return errs.Wrap(err)
			}
			return nil
		},
	}
}
