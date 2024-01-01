package cmd

import (
	"fjnkt98/atcodersearch/batch"
	"fjnkt98/atcodersearch/batch/upload"
	"fjnkt98/atcodersearch/config"
	"fjnkt98/atcodersearch/pkg/solr"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		core, err := solr.NewSolrCore(config.Config.SolrHost, config.Config.Problem.CoreName)
		if err != nil {
			slog.Error("failed to create core", slog.Any("error", err))
			os.Exit(1)
		}

		uploader := upload.NewProblemUploader(
			config.Config.Problem.Upload,
			core,
		)

		batch.RunBatch(uploader)
	},
}

var postUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Post document JSON files into user core",
	Long:  "Post document JSON files into user core",
	Run: func(cmd *cobra.Command, args []string) {
		core, err := solr.NewSolrCore(config.Config.SolrHost, config.Config.User.CoreName)
		if err != nil {
			slog.Error("failed to create core", slog.Any("error", err))
			os.Exit(1)
		}

		uploader := upload.NewUserUploader(
			config.Config.User.Upload,
			core,
		)

		batch.RunBatch(uploader)
	},
}

var postSubmissionCmd = &cobra.Command{
	Use:   "submission",
	Short: "Post document JSON files into submission core",
	Long:  "Post document JSON files into submission core",
	Run: func(cmd *cobra.Command, args []string) {
		core, err := solr.NewSolrCore(config.Config.SolrHost, config.Config.Submission.CoreName)
		if err != nil {
			slog.Error("failed to create core", slog.Any("error", err))
			os.Exit(1)
		}

		uploader := upload.NewSubmissionUploader(
			config.Config.Submission.Upload,
			core,
		)

		batch.RunBatch(uploader)
	},
}

var postRecommendCmd = &cobra.Command{
	Use:   "recommend",
	Short: "Post document JSON files into recommend core",
	Long:  "Post document JSON files into recommend core",
	Run: func(cmd *cobra.Command, args []string) {
		// core, err := solr.NewSolrCore(config.Config.SolrHost, "recommend")
		// if err != nil {
		// 	slog.Error("failed to create core", slog.Any("error", err))
		// 	os.Exit(1)
		// }

		// uploader := upload.NewRecommendUploader(
		// 	config.Config.Recommend.Upload,
		// 	core,
		// )

		// if err := batch.RunBatch(uploader); err != nil {
		// 	slog.Error("batch failed", slog.Any("error", err))
		// 	os.Exit(1)
		// }
	},
}

func init() {
	postCmd.PersistentFlags().BoolP("optimize", "o", false, "When true, send optimize request to Solr")
	viper.BindPFlag("problem.upload.optimize", crawlCmd.PersistentFlags().Lookup("optimize"))
	viper.BindPFlag("user.upload.optimize", crawlCmd.PersistentFlags().Lookup("optimize"))
	viper.BindPFlag("submission.upload.optimize", crawlCmd.PersistentFlags().Lookup("optimize"))

	postCmd.PersistentFlags().BoolP("truncate", "t", false, "When true, truncate index before post")
	viper.BindPFlag("problem.upload.truncate", crawlCmd.PersistentFlags().Lookup("truncate"))
	viper.BindPFlag("user.upload.truncate", crawlCmd.PersistentFlags().Lookup("truncate"))
	viper.BindPFlag("submission.upload.truncate", crawlCmd.PersistentFlags().Lookup("truncate"))

	postCmd.PersistentFlags().String("save-dir", "", "Directory path at which generated documents will be saved")
	viper.BindPFlag("problem.generate.save_dir", crawlCmd.PersistentFlags().Lookup("save-dir"))
	viper.BindPFlag("user.generate.save_dir", crawlCmd.PersistentFlags().Lookup("save-dir"))
	viper.BindPFlag("submission.generate.save_dir", crawlCmd.PersistentFlags().Lookup("save-dir"))

	postCmd.PersistentFlags().Int("concurrent", 3, "Concurrent number of document upload processes")
	viper.BindPFlag("problem.upload.concurrent", crawlCmd.PersistentFlags().Lookup("concurrent"))
	viper.BindPFlag("user.upload.concurrent", crawlCmd.PersistentFlags().Lookup("concurrent"))
	viper.BindPFlag("submission.upload.concurrent", crawlCmd.PersistentFlags().Lookup("concurrent"))

	postCmd.AddCommand(postProblemCmd)
	postCmd.AddCommand(postUserCmd)
	postCmd.AddCommand(postSubmissionCmd)
	postCmd.AddCommand(postRecommendCmd)

	rootCmd.AddCommand(postCmd)
}
