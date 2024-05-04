package update

import (
	"github.com/urfave/cli/v2"
)

func NewUpdateCmd() *cli.Command {
	return &cli.Command{
		Name: "update",
		Subcommands: []*cli.Command{
			newUpdateProblemCmd(),
			newUpdateUserCmd(),
			newUpdateSubmissionCmd(),
			newUpdateSolutionCmd(),
			newUpdateLanguageCmd(),
		},
	}
}
