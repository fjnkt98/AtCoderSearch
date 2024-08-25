package cmd

import (
	"time"

	"github.com/urfave/cli/v2"
)

func NewUpdateCmd() *cli.Command {
	updateLanguageCmd := &cli.Command{
		Name: "language",
		Action: func(ctx *cli.Context) error {
			return nil
			// TODO
		},
	}
	updateProblemCmd := &cli.Command{
		Name: "problem",
		Flags: []cli.Flag{
			&cli.DurationFlag{
				Name:  "duration",
				Value: 1000 * time.Millisecond,
			},
			&cli.BoolFlag{
				Name:  "all",
				Value: false,
			},
			&cli.BoolFlag{
				Name:  "skip-fetch",
				Value: false,
			},
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
			return nil
			// TODO
		},
	}
	updateSolutionCmd := &cli.Command{
		Name: "solution",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "all",
				Value: false,
			},
			&cli.IntFlag{
				Name:  "interval",
				Value: 180,
			},
			&cli.IntFlag{
				Name:  "chunk-size",
				Value: 50000,
			},
			&cli.IntFlag{
				Name:  "concurrent",
				Value: 4,
			},
		},
		Action: func(ctx *cli.Context) error {
			// TODO
			return nil
		},
	}
	updateUserCmd := &cli.Command{
		Name: "user",
		Flags: []cli.Flag{
			&cli.DurationFlag{
				Name:     "duration",
				Value:    1000 * time.Millisecond,
				Category: "crawl",
			},
			&cli.BoolFlag{
				Name:     "skip-fetch",
				Value:    false,
				Category: "crawl",
			},
			&cli.IntFlag{
				Name:     "chunk-size",
				Value:    10000,
				Category: "generate",
			},
			&cli.IntFlag{
				Name:  "concurrent",
				Value: 4,
			},
		},
		Action: func(ctx *cli.Context) error {
			// TODO
			return nil
		},
	}

	return &cli.Command{
		Name: "crawl",
		Subcommands: []*cli.Command{
			updateLanguageCmd,
			updateProblemCmd,
			updateSolutionCmd,
			updateUserCmd,
		},
	}
}
