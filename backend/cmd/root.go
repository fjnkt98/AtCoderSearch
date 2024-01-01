package cmd

import (
	"fjnkt98/atcodersearch/config"
	"fmt"
	"os"

	"golang.org/x/exp/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "atcodersearch",
	Short: "root command for AtCoder Search",
	Long:  "root command for AtCoder Search",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var configFile string

func init() {
	cobra.OnInitialize(func() {
		if configFile == "" {
			configFile = "config/config.yaml"
		}

		if err := config.Load(configFile); err == nil {
			slog.Info(fmt.Sprintf("using config file: %s", viper.ConfigFileUsed()))
		} else {
			slog.Error("failed to load config", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}
	})
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "path to the config file")
}
