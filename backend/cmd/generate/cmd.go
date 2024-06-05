package generate

import (
	"github.com/urfave/cli/v2"
)

func NewGenerateCmd() *cli.Command {
	return &cli.Command{
		Name: "generate",
		Subcommands: []*cli.Command{
			newGenerateProblemCmd(),
			newGenerateUserCmd(),
			// newGenerateSolutionCmd(),
		},
	}
}
