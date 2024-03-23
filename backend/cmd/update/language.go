package update

import "github.com/urfave/cli/v2"

func newUpdateLanguageCmd() *cli.Command {
	return &cli.Command{
		Name:  "language",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) error {
			panic(0)
		},
	}
}
