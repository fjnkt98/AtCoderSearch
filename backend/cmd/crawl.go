package cmd

import (
	"fjnkt98/atcodersearch/batch"
	"fjnkt98/atcodersearch/batch/crawl"
	"fjnkt98/atcodersearch/pkg/atcoder"
	"fjnkt98/atcodersearch/repository"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

func newCrawlCmd(args []string, sub ...*cobra.Command) *cobra.Command {
	crawlCmd := &cobra.Command{
		Use:   "crawl",
		Short: "Crawl and save",
		Long:  "Crawl and save",
	}

	crawlCmd.SetArgs(args)
	crawlCmd.AddCommand(sub...)

	return crawlCmd
}

func newCrawlProblemCmd(args []string, runFunc func(cmd *cobra.Command, args []string)) *cobra.Command {
	crawlProblemCmd := &cobra.Command{
		Use:   "problem",
		Short: "Crawl and save problem information",
		Long:  "Crawl and save problem information",
		PreRun: func(cmd *cobra.Command, args []string) {
			cmd.Flags().IntP("duration", "d", 1000, "Duration[ms] in crawling problem.")
			viper.BindPFlag("crawl.problem.duration", cmd.Flags().Lookup("duration"))
			cmd.Flags().BoolP("all", "a", false, "When true, crawl all problems. Otherwise, crawl the problems which doesn't have been crawled.")
			viper.BindPFlag("crawl.problem.all", cmd.Flags().Lookup("all"))
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
			batch.RunBatch(contestCrawler)

			difficultyCrawler := crawl.NewDifficultyCrawler(
				atcoderProblemsClient,
				repository.NewDifficultyRepository(db),
			)
			batch.RunBatch(difficultyCrawler)

			problemCrawler := crawl.NewProblemCrawler(
				atcoder.NewAtCoderProblemsClient(),
				atcoderClient,
				repository.NewProblemRepository(db),
				Config.Crawl.Problem.Duration,
				Config.Crawl.Problem.All,
			)
			batch.RunBatch(problemCrawler)
		},
	}

	crawlProblemCmd.SetArgs(args)
	if runFunc != nil {
		crawlProblemCmd.Run = runFunc
	}

	return crawlProblemCmd
}

func newCrawlUserCmd(args []string, runFunc func(cmd *cobra.Command, args []string)) *cobra.Command {
	crawlUserCmd := &cobra.Command{
		Use:   "user",
		Short: "Crawl and save user information",
		Long:  "Crawl and save user information",
		PreRun: func(cmd *cobra.Command, args []string) {
			cmd.Flags().IntP("duration", "d", 1000, "Duration[ms] in crawling user.")
			viper.BindPFlag("crawl.user.duration", cmd.Flags().Lookup("duration"))
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

			batch.RunBatch(crawler)
		},
	}

	crawlUserCmd.SetArgs(args)
	if runFunc != nil {
		crawlUserCmd.Run = runFunc
	}

	return crawlUserCmd
}

func newCrawlSubmissionCmd(args []string, runFunc func(cmd *cobra.Command, args []string)) *cobra.Command {
	crawlSubmissionCmd := &cobra.Command{
		Use:   "submission",
		Short: "Crawl and save submissions",
		Long:  "Crawl and save submissions",
		PreRun: func(cmd *cobra.Command, args []string) {
			cmd.Flags().IntP("duration", "d", 1000, "Duration[ms] in crawling user.")
			viper.BindPFlag("crawl.submission.duration", cmd.Flags().Lookup("duration"))
			cmd.Flags().IntP("retry", "r", 0, "Limit of the number of retry when an error occurred in crawling submissions.")
			viper.BindPFlag("crawl.submission.retry", cmd.Flags().Lookup("retry"))
			cmd.Flags().String("target", "", "Target category to crawl. Multiple categories can be specified by separating tem with comma. If not specified, all categories will be crawled.")
			viper.BindPFlag("crawl.submission.target", cmd.Flags().Lookup("target"))
		},
		Run: func(cmd *cobra.Command, args []string) {
			db := repository.MustGetDB(Config.DataBaseURL)

			client, err := atcoder.NewAtCoderClient()
			if err != nil {
				slog.Error("failed to instantiate atcoder client", slog.Any("error", err))
				os.Exit(1)
			}

			targets := strings.Split(Config.Crawl.Submission.Targets, ",")

			crawler := crawl.NewSubmissionCrawler(
				client,
				repository.NewSubmissionRepository(db),
				repository.NewContestRepository(db),
				repository.NewSubmissionCrawlHistoryRepository(db),
				Config.Crawl.Submission.Duration,
				Config.Crawl.Submission.Retry,
				targets,
			)

			batch.RunBatch(crawler)
		},
	}
	crawlSubmissionCmd.SetArgs(args)
	if runFunc != nil {
		crawlSubmissionCmd.Run = runFunc
	}
	return crawlSubmissionCmd
}
