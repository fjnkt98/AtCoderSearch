package cmd

import (
	"fjnkt98/atcodersearch/cmd/crawl"
	"fjnkt98/atcodersearch/cmd/generate"
	"fjnkt98/atcodersearch/cmd/post"
	"fjnkt98/atcodersearch/cmd/serve"
	"fjnkt98/atcodersearch/cmd/update"

	"github.com/urfave/cli/v2"
)

func NewApp() *cli.App {
	cli.NewApp()
	return &cli.App{
		Name: "atcodersearch",
		Commands: []*cli.Command{
			crawl.NewCrawlCmd(),
			generate.NewGenerateCmd(),
			post.NewPostCmd(),
			update.NewUpdateCmd(),
			serve.NewServeCmd(),
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "database-url",
				EnvVars: []string{"DATABASE_URL"},
			},
			&cli.StringFlag{
				Name:    "solr-host",
				EnvVars: []string{"SOLR_HOST"},
			},
		},
	}
}
