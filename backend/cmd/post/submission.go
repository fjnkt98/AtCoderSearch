package post

import "github.com/urfave/cli/v2"

func newPostSubmissionCmd() *cli.Command {
	return &cli.Command{
		Name: "submission",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name: "save-dir",
			},
			&cli.BoolFlag{
				Name:    "optimize",
				Aliases: []string{"o"},
				Value:   false,
			},
			&cli.BoolFlag{
				Name:    "truncate",
				Aliases: []string{"t"},
				Value:   false,
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
