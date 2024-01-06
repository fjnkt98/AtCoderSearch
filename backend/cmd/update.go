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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
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

func newUpdateProblemCmd(args []string, runFunc func(cmd *cobra.Command, args []string)) *cobra.Command {
	updateProblemCmd := &cobra.Command{
		Use:   "problem",
		Short: "update problem index",
		Long:  "update problem index",
		PreRun: func(cmd *cobra.Command, args []string) {
			cmd.Flags().IntP("duration", "d", 1000, "Duration[ms] in crawling problem.")
			viper.BindPFlag("crawl.problem.duration", cmd.Flags().Lookup("duration"))

			cmd.Flags().BoolP("all", "a", false, "When true, crawl all problems. Otherwise, crawl the problems which doesn't have been crawled.")
			viper.BindPFlag("crawl.problem.all", cmd.Flags().Lookup("all"))

			cmd.Flags().String("save-dir", "", "Directory path at which generated documents will be saved.")
			viper.BindPFlag("generate.problem.save_dir", cmd.Flags().Lookup("save-dir"))
			viper.BindPFlag("upload.problem.save_dir", cmd.Flags().Lookup("save-dir"))

			cmd.Flags().Int("chunk-size", 1000, "Number of documents to write in 1 file.")
			viper.BindPFlag("generate.problem.chunk_size", cmd.Flags().Lookup("chunk-size"))

			cmd.Flags().Int("generate-concurrent", 10, "Concurrent number of document generation processes.")
			viper.BindPFlag("generate.problem.concurrent", cmd.Flags().Lookup("generate-concurrent"))

			cmd.Flags().BoolP("optimize", "o", false, "When true, send optimize request to Solr")
			viper.BindPFlag("upload.problem.optimize", cmd.Flags().Lookup("optimize"))

			cmd.Flags().BoolP("truncate", "t", false, "When true, truncate index before upload")
			viper.BindPFlag("upload.problem.truncate", cmd.Flags().Lookup("truncate"))

			cmd.Flags().Int("upload-concurrent", 3, "Concurrent number of document upload processes")
			viper.BindPFlag("upload.problem.concurrent", cmd.Flags().Lookup("upload-concurrent"))

			cmd.Flags().Bool("skip-fetch", false, "When true, skip to crawl problems.")
			viper.BindPFlag("update.problem.skip_fetch", cmd.Flags().Lookup("skip-fetch"))
		},
		Run: func(cmd *cobra.Command, args []string) {
			db := repository.MustGetDB(Config.DataBaseURL)
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
				Config.Crawl.Problem.Duration,
				Config.Crawl.Problem.All,
			)
			generator := generate.NewProblemGenerator(
				generate.NewProblemRowReader(db),
				Config.Generate.Problem.SaveDir,
				Config.Generate.Problem.ChunkSize,
				Config.Generate.Problem.Concurrent,
			)

			core, err := solr.NewSolrCore(Config.SolrHost, Config.ProblemCoreName)
			if err != nil {
				slog.Error("failed to create core", slog.Any("error", err))
				panic("failed to create core")
			}
			uploader := upload.NewDocumentUploader(
				core,
				Config.Upload.Problem.SaveDir,
				Config.Upload.Problem.Concurrent,
				Config.Upload.Problem.Optimize,
				Config.Upload.Problem.Truncate,
			)

			updater := update.NewProblemUpdater(
				problemCrawler,
				contestCrawler,
				difficultyCrawler,
				generator,
				uploader,
				repository.NewUpdateHistoryRepository(db),
				Config.Update.Problem.SkipFetch,
			)

			batch.RunBatch(updater)
		},
	}

	updateProblemCmd.SetArgs(args)
	if runFunc != nil {
		updateProblemCmd.Run = runFunc
	}

	return updateProblemCmd
}

func newUpdateUserCmd(args []string, runFunc func(cmd *cobra.Command, args []string)) *cobra.Command {
	updateUserCmd := &cobra.Command{
		Use:   "user",
		Short: "update user index",
		Long:  "update user index",
		PreRun: func(cmd *cobra.Command, args []string) {
			cmd.Flags().IntP("duration", "d", 1000, "Duration[ms] in crawling user.")
			viper.BindPFlag("crawl.user.duration", cmd.Flags().Lookup("duration"))

			cmd.Flags().String("save-dir", "", "Directory path at which generated documents will be saved.")
			viper.BindPFlag("generate.user.save_dir", cmd.Flags().Lookup("save-dir"))
			viper.BindPFlag("upload.user.save_dir", cmd.Flags().Lookup("save-dir"))

			cmd.Flags().Int("chunk-size", 1000, "Number of documents to write in 1 file.")
			viper.BindPFlag("generate.user.chunk_size", cmd.Flags().Lookup("chunk-size"))

			cmd.Flags().Int("generate-concurrent", 10, "Concurrent number of document generation processes.")
			viper.BindPFlag("generate.user.concurrent", cmd.Flags().Lookup("generate-concurrent"))

			cmd.Flags().BoolP("optimize", "o", false, "When true, send optimize request to Solr")
			viper.BindPFlag("upload.user.optimize", cmd.Flags().Lookup("optimize"))

			cmd.Flags().BoolP("truncate", "t", false, "When true, truncate index before upload")
			viper.BindPFlag("upload.user.truncate", cmd.Flags().Lookup("truncate"))

			cmd.Flags().Int("upload-concurrent", 3, "Concurrent number of document upload processes")
			viper.BindPFlag("upload.user.concurrent", cmd.Flags().Lookup("upload-concurrent"))

			cmd.Flags().Bool("skip-fetch", false, "When true, skip to crawl users.")
			viper.BindPFlag("update.user.skip_fetch", cmd.Flags().Lookup("skip-fetch"))
		},
		Run: func(cmd *cobra.Command, args []string) {
			db := repository.MustGetDB(Config.DataBaseURL)

			client, err := atcoder.NewAtCoderClient()
			if err != nil {
				slog.Error("failed to instantiate atcoder client", slog.Any("error", err))
				panic("failed to instantiate atcoder client")
			}
			crawler := crawl.NewUserCrawler(
				client,
				repository.NewUserRepository(db),
				Config.Crawl.User.Duration,
			)

			generator := generate.NewUserGenerator(
				generate.NewUserRowReader(db),
				Config.Generate.User.SaveDir,
				Config.Generate.User.ChunkSize,
				Config.Generate.User.Concurrent,
			)

			core, err := solr.NewSolrCore(Config.SolrHost, Config.UserCoreName)
			if err != nil {
				slog.Error("failed to create core", slog.Any("error", err))
				panic("failed to create core")
			}

			uploader := upload.NewDocumentUploader(
				core,
				Config.Upload.User.SaveDir,
				Config.Upload.User.Concurrent,
				Config.Upload.User.Optimize,
				Config.Upload.User.Truncate,
			)

			updater := update.NewUserUpdater(
				crawler,
				generator,
				uploader,
				repository.NewUpdateHistoryRepository(db),
				Config.Update.User.SkipFetch,
			)

			batch.RunBatch(updater)
		},
	}

	updateUserCmd.SetArgs(args)
	if runFunc != nil {
		updateUserCmd.Run = runFunc
	}

	return updateUserCmd
}

func newUpdateSubmissionCmd(args []string, runFunc func(cmd *cobra.Command, args []string)) *cobra.Command {
	updateSubmissionCmd := &cobra.Command{
		Use:   "submission",
		Short: "update submission index",
		Long:  "update submission index",
		PreRun: func(cmd *cobra.Command, args []string) {
			cmd.Flags().IntP("duration", "d", 1000, "Duration[ms] in crawling user.")
			viper.BindPFlag("crawl.submission.duration", cmd.Flags().Lookup("duration"))

			cmd.Flags().IntP("retry", "r", 0, "Limit of the number of retry when an error occurred in crawling submissions.")
			viper.BindPFlag("crawl.submission.retry", cmd.Flags().Lookup("retry"))

			cmd.Flags().String("target", "", "Target category to crawl. Multiple categories can be specified by separating tem with comma. If not specified, all categories will be crawled.")
			viper.BindPFlag("crawl.submission.target", cmd.Flags().Lookup("target"))

			cmd.Flags().String("save-dir", "", "Directory path at which generated documents will be saved.")
			viper.BindPFlag("generate.submission.save_dir", cmd.Flags().Lookup("save-dir"))
			viper.BindPFlag("upload.submission.save_dir", cmd.Flags().Lookup("save-dir"))

			cmd.Flags().Int("chunk-size", 1000, "Number of documents to write in 1 file.")
			viper.BindPFlag("generate.submission.chunk_size", cmd.Flags().Lookup("chunk-size"))

			cmd.Flags().Int("generate-concurrent", 10, "Concurrent number of document generation processes.")
			viper.BindPFlag("generate.submission.concurrent", cmd.Flags().Lookup("generate-concurrent"))

			cmd.Flags().Int("interval", 10, "The latest N days' submissions shall be considered.")
			viper.BindPFlag("generate.submission.interval", cmd.Flags().Lookup("interval"))

			cmd.Flags().BoolP("all", "a", false, "When false, crawl only the submissions which doesn't have been crawled (in the interval).")
			viper.BindPFlag("generate.submission.all", cmd.Flags().Lookup("all"))

			cmd.Flags().BoolP("optimize", "o", false, "When true, send optimize request to Solr")
			viper.BindPFlag("upload.submission.optimize", cmd.Flags().Lookup("optimize"))

			cmd.Flags().BoolP("truncate", "t", false, "When true, truncate index before upload")
			viper.BindPFlag("upload.submission.truncate", cmd.Flags().Lookup("truncate"))

			cmd.Flags().Int("upload-concurrent", 3, "Concurrent number of document upload processes")
			viper.BindPFlag("upload.submission.concurrent", cmd.Flags().Lookup("upload-concurrent"))
		},
		Run: func(cmd *cobra.Command, args []string) {
			db := repository.MustGetDB(Config.DataBaseURL)

			generator := generate.NewSubmissionGenerator(
				generate.NewSubmissionRowReader(
					db,
					Config.Generate.Submission.Interval,
					Config.Generate.Submission.All,
				),
				Config.Generate.Submission.SaveDir,
				Config.Generate.Submission.ChunkSize,
				Config.Generate.Submission.Concurrent,
			)

			core, err := solr.NewSolrCore(Config.SolrHost, Config.SubmissionCoreName)
			if err != nil {
				slog.Error("failed to create core", slog.Any("error", err))
				panic("failed to create core")
			}
			uploader := upload.NewDocumentUploader(
				core,
				Config.Upload.Submission.SaveDir,
				Config.Upload.Submission.Concurrent,
				Config.Upload.Submission.Optimize,
				Config.Upload.Submission.Truncate,
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

	return updateSubmissionCmd
}

func newUpdateLanguageCmd(args []string, runFunc func(cmd *cobra.Command, args []string)) *cobra.Command {
	updateLanguageCmd := &cobra.Command{
		Use:   "language",
		Short: "update language index",
		Long:  "update language index",
		PreRun: func(cmd *cobra.Command, args []string) {
		},
		Run: func(cmd *cobra.Command, args []string) {
			db := repository.MustGetDB(Config.DataBaseURL)

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
