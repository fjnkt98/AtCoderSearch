/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fjnkt98/atcodersearch/batch"
	"fjnkt98/atcodersearch/batch/crawl"
	"fjnkt98/atcodersearch/batch/generate"
	"fjnkt98/atcodersearch/batch/update"
	"fjnkt98/atcodersearch/batch/upload"
	"fjnkt98/atcodersearch/pkg/atcoder"
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/repository"

	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newUpdateCmd(args []string, sub ...*cobra.Command) *cobra.Command {
	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "update index",
		Long:  "update index",
	}

	updateCmd.SetArgs(args)
	updateCmd.AddCommand(sub...)

	return updateCmd
}

func newUpdateProblemCmd(args []string, config *RootConfig, runFunc func(cmd *cobra.Command, args []string)) *cobra.Command {
	updateProblemCmd := &cobra.Command{
		Use:   "problem",
		Short: "update problem index",
		Long:  "update problem index",
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlag("crawl.problem.duration", cmd.Flags().Lookup("duration"))
			viper.BindPFlag("crawl.problem.all", cmd.Flags().Lookup("all"))
			viper.BindPFlag("generate.problem.save_dir", cmd.Flags().Lookup("save-dir"))
			viper.BindPFlag("upload.problem.save_dir", cmd.Flags().Lookup("save-dir"))
			viper.BindPFlag("generate.problem.chunk_size", cmd.Flags().Lookup("chunk-size"))
			viper.BindPFlag("generate.problem.concurrent", cmd.Flags().Lookup("generate-concurrent"))
			viper.BindPFlag("upload.problem.optimize", cmd.Flags().Lookup("optimize"))
			viper.BindPFlag("upload.problem.truncate", cmd.Flags().Lookup("truncate"))
			viper.BindPFlag("upload.problem.concurrent", cmd.Flags().Lookup("upload-concurrent"))
			viper.BindPFlag("update.problem.skip_fetch", cmd.Flags().Lookup("skip-fetch"))

			MustLoadConfigFromFlags(cmd.Flags(), config)
		},
		Run: func(cmd *cobra.Command, args []string) {
			db := repository.MustGetDB(config.DataBaseURL)
			atcoderClient, err := atcoder.NewAtCoderClient()
			if err != nil {
				slog.Error("failed to instantiate atcoder client", slog.Any("error", err))
				panic("failed to instantiate atcoder client")
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
				atcoder.NewAtCoderProblemsClient(),
				atcoderClient,
				repository.NewProblemRepository(db),
				config.Crawl.Problem.Duration,
				config.Crawl.Problem.All,
			)
			generator := generate.NewProblemGenerator(
				generate.NewProblemRowReader(db),
				config.Generate.Problem.SaveDir,
				config.Generate.Problem.ChunkSize,
				config.Generate.Problem.Concurrent,
			)

			core, err := solr.NewSolrCore(config.SolrHost, config.ProblemCoreName)
			if err != nil {
				slog.Error("failed to create core", slog.Any("error", err))
				panic("failed to create core")
			}
			uploader := upload.NewDocumentUploader(
				core,
				config.Upload.Problem.SaveDir,
				config.Upload.Problem.Concurrent,
				config.Upload.Problem.Optimize,
				config.Upload.Problem.Truncate,
			)

			updater := update.NewProblemUpdater(
				problemCrawler,
				contestCrawler,
				difficultyCrawler,
				generator,
				uploader,
				repository.NewUpdateHistoryRepository(db),
				config.Update.Problem.SkipFetch,
			)

			batch.RunBatch(updater)
		},
	}

	updateProblemCmd.SetArgs(args)
	if runFunc != nil {
		updateProblemCmd.Run = runFunc
	}
	updateProblemCmd.Flags().IntP("duration", "d", 1000, "Duration[ms] in crawling problem.")
	updateProblemCmd.Flags().BoolP("all", "a", false, "When true, crawl all problems. Otherwise, crawl the problems which doesn't have been crawled.")
	updateProblemCmd.Flags().String("save-dir", "", "Directory path at which generated documents will be saved.")
	updateProblemCmd.Flags().Int("chunk-size", 1000, "Number of documents to write in 1 file.")
	updateProblemCmd.Flags().Int("generate-concurrent", 10, "Concurrent number of document generation processes.")
	updateProblemCmd.Flags().BoolP("optimize", "o", false, "When true, send optimize request to Solr")
	updateProblemCmd.Flags().BoolP("truncate", "t", false, "When true, truncate index before upload")
	updateProblemCmd.Flags().Int("upload-concurrent", 3, "Concurrent number of document upload processes")
	updateProblemCmd.Flags().Bool("skip-fetch", false, "When true, skip to crawl problems.")

	return updateProblemCmd
}

func newUpdateUserCmd(args []string, config *RootConfig, runFunc func(cmd *cobra.Command, args []string)) *cobra.Command {
	updateUserCmd := &cobra.Command{
		Use:   "user",
		Short: "update user index",
		Long:  "update user index",
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlag("crawl.user.duration", cmd.Flags().Lookup("duration"))
			viper.BindPFlag("generate.user.save_dir", cmd.Flags().Lookup("save-dir"))
			viper.BindPFlag("upload.user.save_dir", cmd.Flags().Lookup("save-dir"))
			viper.BindPFlag("generate.user.chunk_size", cmd.Flags().Lookup("chunk-size"))
			viper.BindPFlag("generate.user.concurrent", cmd.Flags().Lookup("generate-concurrent"))
			viper.BindPFlag("upload.user.optimize", cmd.Flags().Lookup("optimize"))
			viper.BindPFlag("upload.user.truncate", cmd.Flags().Lookup("truncate"))
			viper.BindPFlag("upload.user.concurrent", cmd.Flags().Lookup("upload-concurrent"))
			viper.BindPFlag("update.user.skip_fetch", cmd.Flags().Lookup("skip-fetch"))

			MustLoadConfigFromFlags(cmd.Flags(), config)
		},
		Run: func(cmd *cobra.Command, args []string) {
			db := repository.MustGetDB(config.DataBaseURL)

			client, err := atcoder.NewAtCoderClient()
			if err != nil {
				slog.Error("failed to instantiate atcoder client", slog.Any("error", err))
				panic("failed to instantiate atcoder client")
			}
			crawler := crawl.NewUserCrawler(
				client,
				repository.NewUserRepository(db),
				config.Crawl.User.Duration,
			)

			generator := generate.NewUserGenerator(
				generate.NewUserRowReader(db),
				config.Generate.User.SaveDir,
				config.Generate.User.ChunkSize,
				config.Generate.User.Concurrent,
			)

			core, err := solr.NewSolrCore(config.SolrHost, config.UserCoreName)
			if err != nil {
				slog.Error("failed to create core", slog.Any("error", err))
				panic("failed to create core")
			}

			uploader := upload.NewDocumentUploader(
				core,
				config.Upload.User.SaveDir,
				config.Upload.User.Concurrent,
				config.Upload.User.Optimize,
				config.Upload.User.Truncate,
			)

			updater := update.NewUserUpdater(
				crawler,
				generator,
				uploader,
				repository.NewUpdateHistoryRepository(db),
				config.Update.User.SkipFetch,
			)

			batch.RunBatch(updater)
		},
	}

	updateUserCmd.SetArgs(args)
	if runFunc != nil {
		updateUserCmd.Run = runFunc
	}
	updateUserCmd.Flags().IntP("duration", "d", 1000, "Duration[ms] in crawling user.")
	updateUserCmd.Flags().String("save-dir", "", "Directory path at which generated documents will be saved.")
	updateUserCmd.Flags().Int("chunk-size", 1000, "Number of documents to write in 1 file.")
	updateUserCmd.Flags().Int("generate-concurrent", 10, "Concurrent number of document generation processes.")
	updateUserCmd.Flags().BoolP("optimize", "o", false, "When true, send optimize request to Solr")
	updateUserCmd.Flags().BoolP("truncate", "t", false, "When true, truncate index before upload")
	updateUserCmd.Flags().Int("upload-concurrent", 3, "Concurrent number of document upload processes")
	updateUserCmd.Flags().Bool("skip-fetch", false, "When true, skip to crawl users.")

	return updateUserCmd
}

func newUpdateSubmissionCmd(args []string, config *RootConfig, runFunc func(cmd *cobra.Command, args []string)) *cobra.Command {
	updateSubmissionCmd := &cobra.Command{
		Use:   "submission",
		Short: "update submission index",
		Long:  "update submission index",
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlag("crawl.submission.duration", cmd.Flags().Lookup("duration"))
			viper.BindPFlag("crawl.submission.retry", cmd.Flags().Lookup("retry"))
			viper.BindPFlag("crawl.submission.targets", cmd.Flags().Lookup("target"))
			viper.BindPFlag("generate.submission.save_dir", cmd.Flags().Lookup("save-dir"))
			viper.BindPFlag("upload.submission.save_dir", cmd.Flags().Lookup("save-dir"))
			viper.BindPFlag("generate.submission.chunk_size", cmd.Flags().Lookup("chunk-size"))
			viper.BindPFlag("generate.submission.concurrent", cmd.Flags().Lookup("generate-concurrent"))
			viper.BindPFlag("generate.submission.interval", cmd.Flags().Lookup("interval"))
			viper.BindPFlag("generate.submission.all", cmd.Flags().Lookup("all"))
			viper.BindPFlag("upload.submission.optimize", cmd.Flags().Lookup("optimize"))
			viper.BindPFlag("upload.submission.truncate", cmd.Flags().Lookup("truncate"))
			viper.BindPFlag("upload.submission.concurrent", cmd.Flags().Lookup("upload-concurrent"))

			MustLoadConfigFromFlags(cmd.Flags(), config)
		},
		Run: func(cmd *cobra.Command, args []string) {
			db := repository.MustGetDB(config.DataBaseURL)

			generator := generate.NewSubmissionGenerator(
				generate.NewSubmissionRowReader(
					db,
					config.Generate.Submission.Interval,
					config.Generate.Submission.All,
				),
				config.Generate.Submission.SaveDir,
				config.Generate.Submission.ChunkSize,
				config.Generate.Submission.Concurrent,
			)

			core, err := solr.NewSolrCore(config.SolrHost, config.SubmissionCoreName)
			if err != nil {
				slog.Error("failed to create core", slog.Any("error", err))
				panic("failed to create core")
			}
			uploader := upload.NewDocumentUploader(
				core,
				config.Upload.Submission.SaveDir,
				config.Upload.Submission.Concurrent,
				config.Upload.Submission.Optimize,
				config.Upload.Submission.Truncate,
			)

			updater := update.NewSubmissionUpdater(
				generator,
				uploader,
				repository.NewUpdateHistoryRepository(db),
			)

			batch.RunBatch(updater)
		},
	}

	updateSubmissionCmd.SetArgs(args)
	if runFunc != nil {
		updateSubmissionCmd.Run = runFunc
	}
	updateSubmissionCmd.Flags().IntP("duration", "d", 1000, "Duration[ms] in crawling user.")
	updateSubmissionCmd.Flags().IntP("retry", "r", 0, "Limit of the number of retry when an error occurred in crawling submissions.")
	updateSubmissionCmd.Flags().StringSlice("target", nil, "Target category to crawl. Multiple categories can be specified by separating tem with comma. If not specified, all categories will be crawled.")
	updateSubmissionCmd.Flags().String("save-dir", "", "Directory path at which generated documents will be saved.")
	updateSubmissionCmd.Flags().Int("chunk-size", 1000, "Number of documents to write in 1 file.")
	updateSubmissionCmd.Flags().Int("generate-concurrent", 10, "Concurrent number of document generation processes.")
	updateSubmissionCmd.Flags().Int("interval", 10, "The latest N days' submissions shall be considered.")
	updateSubmissionCmd.Flags().BoolP("all", "a", false, "When false, crawl only the submissions which doesn't have been crawled (in the interval).")
	updateSubmissionCmd.Flags().BoolP("optimize", "o", false, "When true, send optimize request to Solr")
	updateSubmissionCmd.Flags().BoolP("truncate", "t", false, "When true, truncate index before upload")
	updateSubmissionCmd.Flags().Int("upload-concurrent", 3, "Concurrent number of document upload processes")

	return updateSubmissionCmd
}

func newUpdateLanguageCmd(args []string, config *RootConfig, runFunc func(cmd *cobra.Command, args []string)) *cobra.Command {
	updateLanguageCmd := &cobra.Command{
		Use:   "language",
		Short: "update language index",
		Long:  "update language index",
		PreRun: func(cmd *cobra.Command, args []string) {
			MustLoadConfigFromFlags(cmd.Flags(), config)
		},
		Run: func(cmd *cobra.Command, args []string) {
			db := repository.MustGetDB(config.DataBaseURL)

			updater := update.NewLanguageUpdater(
				repository.NewSubmissionRepository(db),
				repository.NewLanguageRepository(db),
			)

			batch.RunBatch(updater)
		},
	}
	updateLanguageCmd.SetArgs(args)
	if runFunc != nil {
		updateLanguageCmd.Run = runFunc
	}

	return updateLanguageCmd
}
