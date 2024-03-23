package crawl

import (
	"time"

	"github.com/urfave/cli/v2"
)

func newCrawlSubmissionCmd() *cli.Command {
	return &cli.Command{
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
		},
		Action: func(c *cli.Context) error {
			panic(0)
		},
	}
}
