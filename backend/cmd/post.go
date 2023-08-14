package cmd

import (
	"fjnkt98/atcodersearch/acs"
	"fjnkt98/atcodersearch/solr"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

var postCmd = &cobra.Command{
	Use:   "post",
	Short: "Post document JSON files into Solr core",
	Long:  "Post document JSON files into Solr core",
}

var postProblemCmd = &cobra.Command{
	Use:   "problem",
	Short: "Post document JSON files into problem core",
	Long:  "Post document JSON files into problem core",
	Run: func(cmd *cobra.Command, args []string) {
		saveDir, err := GetSaveDir(cmd, "problem")
		if err != nil {
			slog.Error("failed to get save dir", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}

		solrURL := os.Getenv("SOLR_HOST")
		if solrURL == "" {
			slog.Error("environment variable `SOLR_HOST` must be set.")
			os.Exit(1)
		}
		core, err := solr.NewSolrCore[any, any]("problem", solrURL)
		if err != nil {
			slog.Error("failed to create `problem` core", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}

		uploader := acs.NewDefaultDocumentUploader(core, saveDir)
		optimize := GetBool(cmd, "optimize")
		concurrent := GetInt(cmd, "concurrent")
		if err := uploader.PostDocument(optimize, concurrent); err != nil {
			slog.Error("post failed", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}
	},
}

var postUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Post document JSON files into user core",
	Long:  "Post document JSON files into user core",
	Run: func(cmd *cobra.Command, args []string) {
		saveDir, err := GetSaveDir(cmd, "user")
		if err != nil {
			slog.Error("failed to get save dir", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}

		solrURL := os.Getenv("SOLR_HOST")
		if solrURL == "" {
			slog.Error("environment variable `SOLR_HOST` must be set.")
			os.Exit(1)
		}
		core, err := solr.NewSolrCore[any, any]("user", solrURL)
		if err != nil {
			slog.Error("failed to create `user` core", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}

		uploader := acs.NewDefaultDocumentUploader(core, saveDir)
		optimize := GetBool(cmd, "optimize")
		concurrent := GetInt(cmd, "concurrent")
		if err := uploader.PostDocument(optimize, concurrent); err != nil {
			slog.Error("post failed", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}
	},
}

var postSubmissionCmd = &cobra.Command{
	Use:   "submission",
	Short: "Post document JSON files into submission core",
	Long:  "Post document JSON files into submission core",
	Run: func(cmd *cobra.Command, args []string) {
		saveDir, err := GetSaveDir(cmd, "submission")
		if err != nil {
			slog.Error("failed to get save dir", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}

		solrURL := os.Getenv("SOLR_HOST")
		if solrURL == "" {
			slog.Error("environment variable `SOLR_HOST` must be set.")
			os.Exit(1)
		}
		core, err := solr.NewSolrCore[any, any]("submission", solrURL)
		if err != nil {
			slog.Error("failed to create `user` core", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}

		uploader := acs.NewDefaultDocumentUploader(core, saveDir)
		optimize := GetBool(cmd, "optimize")
		concurrent := GetInt(cmd, "concurrent")
		if err := uploader.PostDocument(optimize, concurrent); err != nil {
			slog.Error("post failed", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}
	},
}

var postRecommendCmd = &cobra.Command{
	Use:   "recommend",
	Short: "Post document JSON files into recommend core",
	Long:  "Post document JSON files into recommend core",
	Run: func(cmd *cobra.Command, args []string) {
		saveDir, err := GetSaveDir(cmd, "recommend")
		if err != nil {
			slog.Error("failed to get save dir", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}

		solrURL := os.Getenv("SOLR_HOST")
		if solrURL == "" {
			slog.Error("environment variable `SOLR_HOST` must be set.")
			os.Exit(1)
		}
		core, err := solr.NewSolrCore[any, any]("recommend", solrURL)
		if err != nil {
			slog.Error("failed to create `user` core", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}

		uploader := acs.NewDefaultDocumentUploader(core, saveDir)
		optimize := GetBool(cmd, "optimize")
		concurrent := GetInt(cmd, "concurrent")
		if err := uploader.PostDocument(optimize, concurrent); err != nil {
			slog.Error("post failed", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}
	},
}

func init() {
	postCmd.PersistentFlags().BoolP("optimize", "o", false, "When true, send optimize request to Solr")
	postCmd.PersistentFlags().String("save-dir", "", "Directory path at which generated documents will be saved")
	postCmd.PersistentFlags().Int("concurrent", 3, "Concurrent number of document upload processes")

	postCmd.AddCommand(postProblemCmd)
	postCmd.AddCommand(postUserCmd)
	postCmd.AddCommand(postSubmissionCmd)
	postCmd.AddCommand(postRecommendCmd)

	rootCmd.AddCommand(postCmd)
}
