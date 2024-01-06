package cmd

import (
	"fmt"
	"os"

	"golang.org/x/exp/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newRootCmd(args []string, sub ...*cobra.Command) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "atcodersearch",
		Short: "root command for AtCoder Search",
		Long:  "root command for AtCoder Search",
	}

	rootCmd.SetArgs(args)

	var configFile string
	cobra.OnInitialize(func() {
		if configFile == "" {
			configFile = "config/config.yaml"
		}

		if err := LoadConfig(configFile, &Config); err == nil {
			slog.Error("failed to load config", slog.String("error", fmt.Sprintf("%+v", err)))
			panic("failed to load config.")
		} else {
			slog.Info(fmt.Sprintf("using config file: %s", viper.ConfigFileUsed()))
		}
	})
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "path to the config file")
	rootCmd.AddCommand(sub...)

	return rootCmd
}

func Execute() error {
	args := os.Args[1:]
	rootCmd := newRootCmd(args)
	rootCmd.AddCommand(
		newCrawlCmd(
			args,
			newCrawlProblemCmd(args, nil),
			newCrawlUserCmd(args, nil),
			newCrawlSubmissionCmd(args, nil),
		),
		newGenerateCmd(
			args,
			newGenerateProblemCmd(args, nil),
			newGenerateUserCmd(args, nil),
			newGenerateSubmissionCmd(args, nil),
		),
		newUploadCmd(
			args,
			newUploadProblemCmd(args, nil),
			newUploadUserCmd(args, nil),
			newUploadSubmissionCmd(args, nil),
		),
		newUpdateCmd(
			args,
			newUpdateProblemCmd(args, nil),
			newUpdateUserCmd(args, nil),
			newUpdateSubmissionCmd(args, nil),
			newUpdateLanguageCmd(args, nil),
		),
		newServerCmd(args, nil),
		newMigrateCmd(args, nil),
	)

	return rootCmd.Execute()
}
