package update

import (
	"time"

	"github.com/urfave/cli/v2"
)

func newUpdateSubmissionCmd() *cli.Command {
	return &cli.Command{
		Name: "submission",
		Flags: []cli.Flag{
			&cli.DurationFlag{
				Name:     "duration",
				Value:    1000 * time.Millisecond,
				Category: "crawl",
			},
			&cli.IntFlag{
				Name:     "retry",
				Value:    0,
				Category: "crawl",
			},
			&cli.StringFlag{
				Name:     "target",
				Category: "crawl",
			},
			&cli.BoolFlag{
				Name:     "all",
				Value:    false,
				Category: "generate",
			},
			&cli.PathFlag{
				Name:     "save-dir",
				Category: "generate, post",
			},
			&cli.IntFlag{
				Name:     "chunk-size",
				Value:    1000,
				Category: "generate",
			},
			&cli.IntFlag{
				Name:     "generate-concurrent",
				Value:    4,
				Category: "generate",
			},
			&cli.IntFlag{
				Name:     "post-concurrent",
				Value:    4,
				Category: "post",
			},
			&cli.BoolFlag{
				Name:     "optimize",
				Value:    false,
				Category: "post",
			},
			&cli.BoolFlag{
				Name:     "truncate",
				Value:    false,
				Category: "post",
			},
			&cli.BoolFlag{
				Name:     "skip-fetch",
				Value:    false,
				Category: "post",
			},
		},
		Action: func(c *cli.Context) error {
			panic(0)
		},
	}
}
