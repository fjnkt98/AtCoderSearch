package post

import (
	"github.com/urfave/cli/v2"
)

func NewPostCmd() *cli.Command {
	return &cli.Command{
		Name: "post",
		Subcommands: []*cli.Command{
			newPostProblemCmd(),
			newPostUserCmd(),
			newPostSubmissionCmd(),
		},
	}
}
