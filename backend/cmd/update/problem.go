package update

import (
	"time"

	"github.com/urfave/cli/v2"
)

func newUpdateProblemCmd() *cli.Command {
	return &cli.Command{
		Name: "problem",
		Flags: []cli.Flag{
			&cli.DurationFlag{
				Name:     "duration",
				Value:    1000 * time.Millisecond,
				Category: "crawl",
			},
			&cli.BoolFlag{
				Name:     "all",
				Value:    false,
				Category: "crawl",
			},
			&cli.BoolFlag{
				Name:     "skip-fetch",
				Value:    false,
				Category: "crawl",
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
		},
		Action: func(c *cli.Context) error {
			panic(0)
		},
	}
}
