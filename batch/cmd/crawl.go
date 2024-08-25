package cmd

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/fjnkt98/atcodersearch-batch/crawl"
	"github.com/fjnkt98/atcodersearch-batch/pkg/atcoder"
	"github.com/fjnkt98/atcodersearch-batch/repository"
	"github.com/urfave/cli/v2"
)

var ErrEmptyAuth = errors.New("atcoder-username or atcoder-password is empty")

func NewCrawlCmd() *cli.Command {
	crawlProblemCmd := &cli.Command{
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
			problemsClient := atcoder.NewAtCoderProblemsClientWithHTTPClient(
				&http.Client{
					Timeout: 5 * time.Second,
				},
			)
			atcoderClient, err := atcoder.NewAtCoderClientWithHTTPClient(
				&http.Client{
					Timeout: 5 * time.Second,
				},
			)
			if err != nil {
				return err
			}

			pool, err := repository.NewPool(ctx.Context, ctx.String("database-url"))
			if err != nil {
				return err
			}

			if err := crawl.NewContestCrawler(problemsClient, pool).Crawl(ctx.Context); err != nil {
				return err
			}
			if err := crawl.NewDifficultyCrawler(problemsClient, pool).Crawl(ctx.Context); err != nil {
				return err
			}
			if err := crawl.NewProblemCrawler(
				problemsClient,
				atcoderClient,
				pool,
				ctx.Duration("duration"),
				ctx.Bool("all"),
			).Crawl(ctx.Context); err != nil {
				return err
			}

			return nil
		},
	}
	crawlUserCmd := &cli.Command{
		Name: "user",
		Flags: []cli.Flag{
			&cli.DurationFlag{
				Name:    "duration",
				Aliases: []string{"d"},
				Value:   1000 * time.Millisecond,
			},
		},
		Action: func(ctx *cli.Context) error {
			client, err := atcoder.NewAtCoderClientWithHTTPClient(
				&http.Client{
					Timeout: 5 * time.Second,
				},
			)
			if err != nil {
				return err
			}

			pool, err := repository.NewPool(ctx.Context, ctx.String("database-url"))
			if err != nil {
				return err
			}

			if err := crawl.NewUserCrawler(client, pool, ctx.Duration("duration")).Crawl(ctx.Context); err != nil {
				return err
			}

			return nil
		},
	}
	crawlSubmissionCmd := &cli.Command{
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
			client, err := atcoder.NewAtCoderClientWithHTTPClient(
				&http.Client{
					Timeout: 5 * time.Second,
				},
			)
			if err != nil {
				return err
			}

			username := ctx.String("atcoder-username")
			password := ctx.String("atcoder-password")
			if username == "" || password == "" {
				return ErrEmptyAuth
			}

			if err := client.Login(ctx.Context, username, password); err != nil {
				return err
			}

			pool, err := repository.NewPool(ctx.Context, ctx.String("database-url"))
			if err != nil {
				return err
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
				return err
			}

			return nil
		},
	}

	return &cli.Command{
		Name: "crawl",
		Subcommands: []*cli.Command{
			crawlProblemCmd,
			crawlUserCmd,
			crawlSubmissionCmd,
		},
	}
}
