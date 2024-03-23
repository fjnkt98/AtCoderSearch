package generate

import "github.com/urfave/cli/v2"

func newGenerateSubmissionCmd() *cli.Command {
	return &cli.Command{
		Name: "submission",
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
			&cli.IntFlag{
				Name:  "interval",
				Value: 90,
			},
			&cli.BoolFlag{
				Name:  "all",
				Value: false,
			},
		},
		Action: func(c *cli.Context) error {
			panic(0)
		},
	}
}
