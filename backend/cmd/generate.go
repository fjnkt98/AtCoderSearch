package cmd

import (
	"fjnkt98/atcodersearch/problem"
	"fjnkt98/atcodersearch/recommend"
	"fjnkt98/atcodersearch/submission"
	"fjnkt98/atcodersearch/user"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate document JSON files",
	Long:  "Generate document JSON files",
}

var generateProblemCmd = &cobra.Command{
	Use:   "problem",
	Short: "Generate problem document JSON files",
	Long:  "Generate problem document JSON files",
	Run: func(cmd *cobra.Command, args []string) {
		saveDir, err := GetSaveDir(cmd, "problem")
		if err != nil {
			slog.Error("failed to get save dir", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}

		generator := problem.NewDocumentGenerator(GetDB(), saveDir)
		concurrent := GetInt(cmd, "concurrent")
		if err := generator.Run(1000, concurrent); err != nil {
			slog.Error("generation failed", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}
	},
}

var generateUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Generate user document JSON files",
	Long:  "Generate user document JSON files",
	Run: func(cmd *cobra.Command, args []string) {
		saveDir, err := GetSaveDir(cmd, "user")
		if err != nil {
			slog.Error("failed to get save dir", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}
		generator := user.NewDocumentGenerator(GetDB(), saveDir)
		concurrent := GetInt(cmd, "concurrent")
		if err := generator.Run(1000, concurrent); err != nil {
			slog.Error("generation failed", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}
	},
}

var generateSubmissionCmd = &cobra.Command{
	Use:   "submission",
	Short: "Generate submission document JSON files",
	Long:  "Generate submission document JSON files",
	Run: func(cmd *cobra.Command, args []string) {
		saveDir, err := GetSaveDir(cmd, "submission")
		if err != nil {
			slog.Error("failed to get save dir", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}
		generator := submission.NewDocumentGenerator(GetDB(), saveDir, time.Time{})
		concurrent := GetInt(cmd, "concurrent")
		if err := generator.Run(100000, concurrent); err != nil {
			slog.Error("generation failed", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}
	},
}

var generateRecommendCmd = &cobra.Command{
	Use:   "recommend",
	Short: "Generate recommend document JSON files",
	Long:  "Generate recommend document JSON files",
	Run: func(cmd *cobra.Command, args []string) {
		saveDir, err := GetSaveDir(cmd, "recommend")
		if err != nil {
			slog.Error("failed to get save dir", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}
		generator := recommend.NewDocumentGenerator(GetDB(), saveDir)

		concurrent := GetInt(cmd, "concurrent")
		chunkSize := GetInt(cmd, "chunk-size")
		if err := generator.Run(chunkSize, concurrent); err != nil {
			slog.Error("generation failed", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}
	},
}

func init() {
	generateCmd.PersistentFlags().String("save-dir", "", "Directory path at which generated documents will be saved")
	generateCmd.PersistentFlags().Int("concurrent", 10, "Concurrent number of document generation processes")
	generateCmd.PersistentFlags().Int("chunk-size", 1000, "Number of documents to write in 1 file.")
	generateCmd.AddCommand(generateProblemCmd)
	generateCmd.AddCommand(generateUserCmd)
	generateCmd.AddCommand(generateSubmissionCmd)
	generateCmd.AddCommand(generateRecommendCmd)

	rootCmd.AddCommand(generateCmd)
}
