package post

import (
	"fjnkt98/atcodersearch/batch/post"
	"fjnkt98/atcodersearch/pkg/solr"

	"github.com/goark/errs"
	"github.com/urfave/cli/v2"
)

func newPostUserCmd() *cli.Command {
	return &cli.Command{
		Name: "user",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:    "save-dir",
				EnvVars: []string{"USER_SAVE_DIR"},
			},
			&cli.BoolFlag{
				Name:    "optimize",
				Aliases: []string{"o"},
				Value:   false,
			},
			&cli.BoolFlag{
				Name:    "truncate",
				Aliases: []string{"t"},
				Value:   false,
			},
			&cli.IntFlag{
				Name:  "concurrent",
				Value: 2,
			},
		},
		Action: func(ctx *cli.Context) error {
			core, err := solr.NewSolrCore(ctx.String("solr-host"), "user")
			if err != nil {
				return errs.Wrap(err)
			}

			err = post.PostDocument(
				ctx.Context,
				core,
				ctx.String("save-dir"),
				post.WithConcurrent(ctx.Int("concurrent")),
				post.WithOptimize(ctx.Bool("optimize")),
				post.WithTruncate(ctx.Bool("truncate")),
			)
			if err != nil {
				return errs.Wrap(err)
			}
			return nil
		},
	}
}
