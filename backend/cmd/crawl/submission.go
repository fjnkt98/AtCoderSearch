package crawl

import (
	"fjnkt98/atcodersearch/batch/crawl"
	"fjnkt98/atcodersearch/pkg/atcoder"
	"fjnkt98/atcodersearch/repository"
	"strings"
	"time"

	"github.com/goark/errs"
	"github.com/urfave/cli/v2"
)

func newCrawlSubmissionCmd() *cli.Command {
	return &cli.Command{
		Name: "submission",
		Flags: []cli.Flag{
			&cli.DurationFlag{
				Name:    "duration",
				Aliases: []string{"d"},
				Value:   1000 * time.Millisecond,
			},
			&cli.IntFlag{
				Name:    "retry",
				Aliases: []string{"r"},
				Value:   0,
			},
			&cli.StringFlag{
				Name: "target",
			},
			&cli.StringFlag{
				Name:    "atcoder-username",
				EnvVars: []string{"ATCODER_USERNAME"},
			},
			&cli.StringFlag{
				Name:    "atcoder-password",
				EnvVars: []string{"ATCODER_PASSWORD"},
			},
		},
		Action: func(ctx *cli.Context) error {
			client, err := atcoder.NewAtCoderClient()
			if err != nil {
				return errs.Wrap(err)
			}
			username := ctx.String("atcoder-username")
			password := ctx.String("atcoder-password")
			if username == "" || password == "" {
				return errs.New("atcoder-username or atcoder-password is empty")
			}
			if err := client.Login(ctx.Context, username, password); err != nil {
				return errs.Wrap(err)
			}

			pool, err := repository.NewPool(ctx.Context, ctx.String("database-url"))
			if err != nil {
				return errs.Wrap(err)
			}
			var targets []string
			if target := ctx.String("target"); target != "" {
				targets = strings.Split(target, ",")
			}
			crawler := crawl.NewSubmissionCrawler(
				client,
				pool,
				ctx.Duration("duration"),
				ctx.Int("retry"),
				targets,
			)

			if err := crawler.Crawl(ctx.Context); err != nil {
				return errs.Wrap(err)
			}
			return nil
		},
	}
}
