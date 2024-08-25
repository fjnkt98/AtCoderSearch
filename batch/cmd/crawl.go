package cmd

import (
	"time"

	"github.com/urfave/cli/v2"
)

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
			// TODO
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
			// TODO
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
			return nil
			// TODO
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
