/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fjnkt98/atcodersearch/batch"
	"fjnkt98/atcodersearch/batch/crawl"
	"fjnkt98/atcodersearch/batch/generate"
	"fjnkt98/atcodersearch/batch/repository"
	"fjnkt98/atcodersearch/batch/update"
	"fjnkt98/atcodersearch/batch/upload"
	"fjnkt98/atcodersearch/config"
	"fjnkt98/atcodersearch/pkg/atcoder"
	"fjnkt98/atcodersearch/pkg/solr"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update index",
	Long:  "update index",
}

var updateProblemCmd = &cobra.Command{
	Use:   "problem",
	Short: "update problem index",
	Long:  "update problem index",
	Run: func(cmd *cobra.Command, args []string) {
		db := GetDB(GetEngine())
		atcoderClient, err := atcoder.NewAtCoderClient()
		if err != nil {
			slog.Error("failed to instantiate atcoder client", slog.Any("error", err))
			os.Exit(1)
		}
		atcoderProblemsClient := atcoder.NewAtCoderProblemsClient()

		contestCrawler := crawl.NewContestCrawler(
			atcoderProblemsClient,
			repository.NewContestRepository(db),
		)
		difficultyCrawler := crawl.NewDifficultyCrawler(
			atcoderProblemsClient,
			repository.NewDifficultyRepository(db),
		)
		problemCrawler := crawl.NewProblemCrawler(
			atcoderProblemsClient,
			atcoderClient,
			repository.NewProblemRepository(db),
			config.Config.Problem.Crawl,
		)
		generator := generate.NewProblemGenerator(
			config.Config.Problem.Generate,
			generate.NewProblemRowReader(db),
		)
		core, err := solr.NewSolrCore(config.Config.SolrHost, config.Config.Problem.CoreName)
		if err != nil {
			slog.Error("failed to create core", slog.Any("error", err))
			os.Exit(1)
		}

		uploader := upload.NewProblemUploader(
			config.Config.Problem.Upload,
			core,
		)

		updater := update.NewProblemUpdater(
			config.Config.Problem,
			problemCrawler,
			contestCrawler,
			difficultyCrawler,
			generator,
			uploader,
			repository.NewUpdateHistoryRepository(db),
		)

		batch.RunBatch(updater)
	},
}

var updateUserCmd = &cobra.Command{
	Use:   "user",
	Short: "update user index",
	Long:  "update user index",
	Run: func(cmd *cobra.Command, args []string) {
		db := GetDB(GetEngine())
		atcoderClient, err := atcoder.NewAtCoderClient()
		if err != nil {
			slog.Error("failed to instantiate atcoder client", slog.Any("error", err))
			os.Exit(1)
		}

		crawler := crawl.NewUserCrawler(
			atcoderClient,
			repository.NewUserRepository(db),
			config.Config.User.Crawl,
		)
		generator := generate.NewUserGenerator(
			config.Config.User.Generate,
			generate.NewUserRowReader(db),
		)
		core, err := solr.NewSolrCore(config.Config.SolrHost, config.Config.User.CoreName)
		if err != nil {
			slog.Error("failed to create core", slog.Any("error", err))
			os.Exit(1)
		}
		uploader := upload.NewUserUploader(
			config.Config.User.Upload,
			core,
		)

		updater := update.NewUserUpdater(
			config.Config.User,
			crawler,
			generator,
			uploader,
			repository.NewUpdateHistoryRepository(db),
		)

		batch.RunBatch(updater)
	},
}

var updateSubmissionCmd = &cobra.Command{
	Use:   "submission",
	Short: "update submission index",
	Long:  "update submission index",
	Run: func(cmd *cobra.Command, args []string) {
		db := GetDB(GetEngine())

		historyRepo := repository.NewUpdateHistoryRepository(db)
		generator := generate.NewSubmissionGenerator(
			config.Config.Submission.Generate,
			generate.NewSubmissionRowReader(
				db,
				historyRepo,
				config.Config.Submission.Read,
			),
		)
		core, err := solr.NewSolrCore(config.Config.SolrHost, config.Config.Submission.CoreName)
		if err != nil {
			slog.Error("failed to create core", slog.Any("error", err))
			os.Exit(1)
		}
		uploader := upload.NewSubmissionUploader(
			config.Config.Submission.Upload,
			core,
		)

		updater := update.NewSubmissionUpdater(
			config.Config.Submission,
			generator,
			uploader,
			historyRepo,
		)

		batch.RunBatch(updater)
	},
}

var updateRecommendCmd = &cobra.Command{
	Use:   "recommend",
	Short: "update recommend index",
	Long:  "update recommend index",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var updateLanguageCmd = &cobra.Command{
	Use:   "language",
	Short: "update language index",
	Long:  "update language index",
	Run: func(cmd *cobra.Command, args []string) {
		db := GetDB(GetEngine())
		updater := update.NewLanguageUpdater(
			repository.NewSubmissionRepository(db),
			repository.NewLanguageRepository(db),
		)
		batch.RunBatch(updater)
	},
}

func init() {
	updateCmd.PersistentFlags().Bool("migrate", false, "Execute database migration before update index.")
	viper.BindPFlag("do_migrate", updateCmd.PersistentFlags().Lookup("migrate"))

	updateCmd.PersistentFlags().String("save-dir", "", "Directory path at which generated documents will be saved.")
	viper.BindPFlag("problem.generate.save_dir", updateCmd.PersistentFlags().Lookup("save-dir"))
	viper.BindPFlag("problem.upload.save_dir", updateCmd.PersistentFlags().Lookup("save-dir"))
	viper.BindPFlag("user.generate.save_dir", updateCmd.PersistentFlags().Lookup("save-dir"))
	viper.BindPFlag("user.upload.save_dir", updateCmd.PersistentFlags().Lookup("save-dir"))
	viper.BindPFlag("submission.generate.save_dir", updateCmd.PersistentFlags().Lookup("save-dir"))
	viper.BindPFlag("submission.upload.save_dir", updateCmd.PersistentFlags().Lookup("save-dir"))

	updateCmd.PersistentFlags().BoolP("optimize", "o", false, "Optimize index if true.")
	viper.BindPFlag("problem.upload.optimize", updateCmd.PersistentFlags().Lookup("optimize"))
	viper.BindPFlag("user.upload.optimize", updateCmd.PersistentFlags().Lookup("optimize"))
	viper.BindPFlag("submission.upload.optimize", updateCmd.PersistentFlags().Lookup("optimize"))

	updateCmd.PersistentFlags().Int("chunk-size", 1000, "Number of documents to write in 1 file.")
	viper.BindPFlag("problem.upload.chunk_size", updateCmd.PersistentFlags().Lookup("chunk-size"))
	viper.BindPFlag("user.upload.chunk_size", updateCmd.PersistentFlags().Lookup("chunk-size"))
	viper.BindPFlag("submission.upload.chunk_size", updateCmd.PersistentFlags().Lookup("chunk-size"))

	updateCmd.PersistentFlags().Int("generate-concurrent", 6, "Number of concurrent document generation processes")
	viper.BindPFlag("problem.generate.concurrent", updateCmd.PersistentFlags().Lookup("generate-concurrent"))
	viper.BindPFlag("user.generate.concurrent", updateCmd.PersistentFlags().Lookup("generate-concurrent"))
	viper.BindPFlag("submission.generate.concurrent", updateCmd.PersistentFlags().Lookup("generate-concurrent"))

	updateCmd.PersistentFlags().Int("upload-concurrent", 4, "Number of concurrent document upload processes")
	viper.BindPFlag("problem.upload.concurrent", updateCmd.PersistentFlags().Lookup("upload-concurrent"))
	viper.BindPFlag("user.upload.concurrent", updateCmd.PersistentFlags().Lookup("upload-concurrent"))
	viper.BindPFlag("submission.upload.concurrent", updateCmd.PersistentFlags().Lookup("upload-concurrent"))

	updateProblemCmd.Flags().BoolP("all", "a", false, "Crawl all problems if true.")
	viper.BindPFlag("problem.crawl.all", updateProblemCmd.Flags().Lookup("all"))
	updateProblemCmd.Flags().BoolP("skip-fetch", "f", false, "Skip crawling if true.")
	viper.BindPFlag("problem.update.skip_fetch", updateProblemCmd.Flags().Lookup("skip-fetch"))
	updateProblemCmd.Flags().Int("duration", 1000, "Interval time[ms] for crawling.")
	viper.BindPFlag("problem.crawl.duration", updateProblemCmd.Flags().Lookup("duration"))

	updateUserCmd.Flags().BoolP("skip-fetch", "f", false, "Skip crawling if true.")
	viper.BindPFlag("user.update.skip_fetch", updateUserCmd.Flags().Lookup("skip-fetch"))
	updateUserCmd.Flags().Int("duration", 1000, "Interval time[ms] for crawling.")
	viper.BindPFlag("user.crawl.duration", updateUserCmd.Flags().Lookup("duration"))

	updateSubmissionCmd.Flags().BoolP("all", "a", false, "Update all submissions.")
	viper.BindPFlag("submission.read.all", updateSubmissionCmd.Flags().Lookup("all"))
	updateSubmissionCmd.Flags().Int("interval", 90, "Indexing submissions for the past in N days.")
	viper.BindPFlag("submission.read.interval", updateSubmissionCmd.Flags().Lookup("interval"))

	updateCmd.AddCommand(updateProblemCmd)
	updateCmd.AddCommand(updateUserCmd)
	updateCmd.AddCommand(updateSubmissionCmd)
	updateCmd.AddCommand(updateRecommendCmd)
	updateCmd.AddCommand(updateLanguageCmd)

	rootCmd.AddCommand(updateCmd)
}
