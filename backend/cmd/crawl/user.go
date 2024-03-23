package crawl

import (
	"time"

	"github.com/urfave/cli/v2"
)

func newCrawlUserCmd() *cli.Command {
	return &cli.Command{
		Name: "user",
		Flags: []cli.Flag{
			&cli.DurationFlag{
				Name:    "duration",
				Aliases: []string{"d"},
				Value:   1000 * time.Millisecond,
			},
		},
		Action: func(c *cli.Context) error {
			panic(0)
		},
	}
}
