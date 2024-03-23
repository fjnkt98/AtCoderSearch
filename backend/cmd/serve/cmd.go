package serve

import (
	"github.com/urfave/cli/v2"
)

func NewServeCmd() *cli.Command {
	return &cli.Command{
		Name: "serve",
		Action: func(c *cli.Context) error {
			panic(0)
		},
	}
}
