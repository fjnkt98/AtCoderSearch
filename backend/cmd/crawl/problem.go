package crawl

import (
	"fjnkt98/atcodersearch/batch/crawl"
	"fjnkt98/atcodersearch/pkg/atcoder"
	"fjnkt98/atcodersearch/repository"
	"time"

	"github.com/goark/errs"
	"github.com/urfave/cli/v2"
)

func newCrawlProblemCmd() *cli.Command {
	return &cli.Command{
		Name: "problem",
		Flags: []cli.Flag{
			&cli.DurationFlag{
				Name:    "duration",
				Aliases: []string{"d"},
				Value:   1000 * time.Millisecond,
			},
			&cli.BoolFlag{
				Name:    "all",
				Aliases: []string{"a"},
				Value:   false,
			},
		},
		Action: func(ctx *cli.Context) error {
			problemsClient := atcoder.NewAtCoderProblemsClient()
			atcoderClient, err := atcoder.NewAtCoderClient()
			if err != nil {
				return errs.Wrap(err)
			}
			pool, err := repository.NewPool(ctx.Context, ctx.String("database-url"))
			if err != nil {
				return errs.Wrap(err)
			}

			{
				crawler := crawl.NewContestCrawler(problemsClient, pool)
				if err := crawler.CrawlContest(ctx.Context); err != nil {
					return errs.Wrap(err)
				}
			}
			{
				crawler := crawl.NewDifficultyCrawler(problemsClient, pool)
				if err := crawler.CrawlDifficulty(ctx.Context); err != nil {
					return errs.Wrap(err)
				}
			}
			{
				crawler := crawl.NewProblemCrawler(
					problemsClient,
					atcoderClient,
					pool,
					ctx.Duration("duration"),
					ctx.Bool("all"),
				)
				if err := crawler.CrawlProblem(ctx.Context); err != nil {
					return errs.Wrap(err)
				}
			}

			return nil
		},
	}
}
