package crawl

import (
	"fjnkt98/atcodersearch/batch/crawl"
	"fjnkt98/atcodersearch/pkg/atcoder"
	"fjnkt98/atcodersearch/repository"
	"time"

	"github.com/goark/errs"
	"github.com/urfave/cli/v2"
)

func newCrawlUserCmd() *cli.Command {
	return &cli.Command{
		Name: "user",
		Flags: []cli.Flag{
			&cli.DurationFlag{
				Name:    "duration",
				Aliases: []string{"d"},
				Value:   1000 * time.Millisecond,
			},
		},
		Action: func(ctx *cli.Context) error {
			client, err := atcoder.NewAtCoderClient()
			if err != nil {
				return errs.Wrap(err)
			}
			pool, err := repository.NewPool(ctx.Context, ctx.String("database-url"))
			if err != nil {
				return errs.Wrap(err)
			}
			crawler := crawl.NewUserCrawler(client, pool, ctx.Duration("duration"))
			if err := crawler.CrawlUser(ctx.Context); err != nil {
				return errs.Wrap(err)
			}

			return nil
		},
	}
}
