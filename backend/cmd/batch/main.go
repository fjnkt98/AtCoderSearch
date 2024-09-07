package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

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
		&cli.Command{
			Name: "crawl",
			Subcommands: []*cli.Command{
				&cli.Command{
					Name: "problem",
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name:    "all",
							Aliases: []string{"a"},
							Value:   false,
						},
					},
					Action: func(ctx *cli.Context) error {
						panic("not implemented")
					},
				},
				&cli.Command{
					Name: "user",
					Action: func(ctx *cli.Context) error {
						panic("not implemented")
					},
				},
				&cli.Command{
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
						panic("not implemented")
					},
				},
			},
			Flags: []cli.Flag{
				&cli.DurationFlag{
					Name:    "duration",
					Aliases: []string{"d"},
					Value:   1000 * time.Millisecond,
				},
			},
		},
		&cli.Command{
			Name: "update",
			Subcommands: []*cli.Command{
				&cli.Command{
					Name:  "problem",
					Flags: []cli.Flag{},
					Action: func(ctx *cli.Context) error {
						panic("not implemented")
					},
				},
				&cli.Command{
					Name:  "user",
					Flags: []cli.Flag{},
					Action: func(ctx *cli.Context) error {
						panic("not implemented")
					},
				},
				&cli.Command{
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
