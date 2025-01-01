package cmd

import (
	"errors"
	"fjnkt98/atcodersearch/atcoder"
	"fjnkt98/atcodersearch/crawl"
	"fjnkt98/atcodersearch/repository"
	"fmt"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

var ErrNoAuthInfo = errors.New("no auth info")

func NewCrawlCmd() *cli.Command {
	return &cli.Command{
		Name: "crawl",
		Flags: []cli.Flag{
			&cli.DurationFlag{
				Name:    "duration",
				Aliases: []string{"d"},
				Value:   1000 * time.Millisecond,
			},
		},
		Subcommands: []*cli.Command{
			{
				Name: "problem",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "all",
						Aliases: []string{"a"},
						Value:   false,
					},
				},
				Action: func(ctx *cli.Context) error {
					atcoderClient, err := atcoder.NewAtCoderClient()
					if err != nil {
						return fmt.Errorf("crawl problem: %w", err)
					}
					problemsClient := atcoder.NewAtCoderProblemsClient()

					pool, err := repository.NewPool(ctx.Context, ctx.String("database-url"))
					if err != nil {
						return fmt.Errorf("crawl problem: %w", err)
					}

					crawler := crawl.NewProblemCrawler(atcoderClient, problemsClient, pool, ctx.Duration("duration"), ctx.Bool("all"))
					if err := crawler.CrawlContests(ctx.Context); err != nil {
						return fmt.Errorf("crawl problem: %w", err)
					}
					if err := crawler.CrawlDifficulties(ctx.Context); err != nil {
						return fmt.Errorf("crawl problem: %w", err)
					}
					if err := crawler.CrawlProblems(ctx.Context); err != nil {
						return fmt.Errorf("crawl problem: %w", err)
					}

					return nil
				},
			},
			{
				Name: "user",
				Action: func(ctx *cli.Context) error {
					atcoderClient, err := atcoder.NewAtCoderClient()
					if err != nil {
						return fmt.Errorf("crawl user: %w", err)
					}

					pool, err := repository.NewPool(ctx.Context, ctx.String("database-url"))
					if err != nil {
						return fmt.Errorf("crawl user: %w", err)
					}

					crawler := crawl.NewUserCrawler(atcoderClient, pool, ctx.Duration("duration"))
					if err := crawler.Crawl(ctx.Context); err != nil {
						return fmt.Errorf("crawl user: %w", err)
					}

					return nil
				},
			},
			{
				Name: "submission",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "retry",
						Aliases: []string{"r"},
						Value:   0,
					},
					&cli.StringFlag{
						Name: "target",
					},
					&cli.BoolFlag{
						Name: "endless",
					},
					&cli.StringFlag{
						Name:    "atcoder-username",
						Hidden:  true,
						EnvVars: []string{"ATCODER_USERNAME"},
					},
					&cli.StringFlag{
						Name:    "atcoder-password",
						Hidden:  true,
						EnvVars: []string{"ATCODER_PASSWORD"},
					},
				},
				Action: func(ctx *cli.Context) error {
					atcoderClient, err := atcoder.NewAtCoderClient()
					if err != nil {
						return fmt.Errorf("crawl submission: %w", err)
					}
					username := ctx.String("atcoder-username")
					password := ctx.String("atcoder-password")
					if username == "" || password == "" {
						return ErrNoAuthInfo
					}
					if err := atcoderClient.Login(ctx.Context, username, password); err != nil {
						return fmt.Errorf("login: %w", err)
					}
					pool, err := repository.NewPool(ctx.Context, ctx.String("database-url"))
					if err != nil {
						return fmt.Errorf("crawl submission: %w", err)
					}

					var targets []string
					if target := ctx.String("target"); target != "" {
						targets = strings.Split(target, ",")
					}

					crawler := crawl.NewSubmissionCrawler(
						atcoderClient,
						pool,
						ctx.Duration("duration"),
						ctx.Int("retry"),
						1*time.Minute,
						targets,
						ctx.Bool("endless"),
					)

					if err := crawler.Crawl(ctx.Context); err != nil {
						return fmt.Errorf("crawl submission: %w", err)
					}
					return nil
				},
			},
		},
	}
}
