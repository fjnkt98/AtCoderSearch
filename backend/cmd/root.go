package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func newRootCmd(args []string, sub ...*cobra.Command) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "atcodersearch",
		Short: "root command for AtCoder Search",
		Long:  "root command for AtCoder Search",
	}

	rootCmd.SetArgs(args)

	rootCmd.PersistentFlags().String("config", "cmd/config.yaml", "path to the config file")
	rootCmd.AddCommand(sub...)

	return rootCmd
}

func Execute() error {
	args := os.Args[1:]
	var config RootConfig

	rootCmd := newRootCmd(
		args,
		newCrawlCmd(
			args,
			newCrawlProblemCmd(args, &config, nil),
			newCrawlUserCmd(args, &config, nil),
			newCrawlSubmissionCmd(args, &config, nil),
		),
		newGenerateCmd(
			args,
			newGenerateProblemCmd(args, &config, nil),
			newGenerateUserCmd(args, &config, nil),
			newGenerateSubmissionCmd(args, &config, nil),
		),
		newUploadCmd(
			args,
			newUploadProblemCmd(args, &config, nil),
			newUploadUserCmd(args, &config, nil),
			newUploadSubmissionCmd(args, &config, nil),
		),
		newUpdateCmd(
			args,
			newUpdateProblemCmd(args, &config, nil),
			newUpdateUserCmd(args, &config, nil),
			newUpdateSubmissionCmd(args, &config, nil),
			newUpdateLanguageCmd(args, &config, nil),
		),
		newServerCmd(args, &config, nil),
		newMigrateCmd(args, &config, nil),
	)

	return rootCmd.Execute()
}
