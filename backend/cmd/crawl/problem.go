package crawl

import (
	"time"

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
		Action: func(c *cli.Context) error {
			panic(0)
		},
	}
}
