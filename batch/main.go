package main

import (
	"errors"
	"fjnkt98/atcodersearch/crawl"
	"fjnkt98/atcodersearch/pkg/atcoder"
	"fjnkt98/atcodersearch/repository"
	"fjnkt98/atcodersearch/update"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

var ErrNoAuthInfo = errors.New("no auth info")

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	app := cli.NewApp()
	app.Name = "atcodersearch-batch"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "database-url",
			Hidden:  true,
			EnvVars: []string{"DATABASE_URL"},
		},
		&cli.StringFlag{
			Name:    "engine-url",
			Hidden:  true,
			EnvVars: []string{"ENGINE_URL"},
		},
	}
	app.Commands = []*cli.Command{
		{
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
							return err
						}
						problemsClient := atcoder.NewAtCoderProblemsClient()

						pool, err := repository.NewPool(ctx.Context, ctx.String("database-url"))
						if err != nil {
							return err
						}

						crawler := crawl.NewProblemCrawler(atcoderClient, problemsClient, pool, ctx.Duration("duration"), ctx.Bool("all"))
						if err := crawler.CrawlContests(ctx.Context); err != nil {
							return err
						}
						if err := crawler.CrawlDifficulties(ctx.Context); err != nil {
							return err
						}
						if err := crawler.CrawlProblems(ctx.Context); err != nil {
							return err
						}

						return nil
					},
				},
				{
					Name: "user",
					Action: func(ctx *cli.Context) error {
						atcoderClient, err := atcoder.NewAtCoderClient()
						if err != nil {
							return err
						}

						pool, err := repository.NewPool(ctx.Context, ctx.String("database-url"))
						if err != nil {
							return err
						}

						crawler := crawl.NewUserCrawler(atcoderClient, pool, ctx.Duration("duration"))
						if err := crawler.Crawl(ctx.Context); err != nil {
							return err
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
							return err
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
							return err
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
						)

						if err := crawler.Crawl(ctx.Context); err != nil {
							return err
						}
						return nil
					},
				},
			},
		},
		{
			Name: "update",
			Subcommands: []*cli.Command{
				{
					Name: "problem",
					Flags: []cli.Flag{
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
							return fmt.Errorf("new pool: %w", err)
						}

						reader := update.NewProblemRowReader(pool)

						if err := update.UpdateIndex(
							ctx.Context,
							reader,
							nil,
							ctx.Int("chunk-size"),
							ctx.Int("concurrent"),
						); err != nil {
							return fmt.Errorf("update index: %w", err)
						}

						return nil
					},
				},
				{
					Name:  "user",
					Flags: []cli.Flag{},
					Action: func(ctx *cli.Context) error {
						panic("not implemented")
					},
				},
				{
					Name:  "language",
					Flags: []cli.Flag{},
					Action: func(ctx *cli.Context) error {
						panic("not implemented")
					},
				},
			},
			Flags: []cli.Flag{},
		},
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error("command failed", slog.Any("error", err))
		os.Exit(1)
	}
}
