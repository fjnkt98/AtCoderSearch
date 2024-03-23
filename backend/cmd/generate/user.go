package generate

import "github.com/urfave/cli/v2"

func newGenerateUserCmd() *cli.Command {
	return &cli.Command{
		Name: "user",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name: "save-dir",
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
		Action: func(c *cli.Context) error {
			panic(0)
		},
	}
}
