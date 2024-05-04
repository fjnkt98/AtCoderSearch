package crawl

import (
	"github.com/urfave/cli/v2"
)

func NewCrawlCmd() *cli.Command {
	return &cli.Command{
		Name: "crawl",
		Subcommands: []*cli.Command{
			newCrawlProblemCmd(),
			newCrawlUserCmd(),
			newCrawlSubmissionCmd(),
		},
	}
}
