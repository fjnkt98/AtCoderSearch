package cmd

import (
	"fjnkt98/atcodersearch/batch"
	"fjnkt98/atcodersearch/batch/crawl"
	"fjnkt98/atcodersearch/batch/repository"
	"fjnkt98/atcodersearch/config"
	"fjnkt98/atcodersearch/pkg/atcoder"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

var crawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "Crawl and save",
	Long:  "Crawl and save",
}

var crawlProblemCmd = &cobra.Command{
	Use:   "problem",
	Short: "Crawl and save problem information",
	Long:  "Crawl and save problem information",
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
			config.Config.Problem.Crawl,
		)
		batch.RunBatch(problemCrawler)
	},
}

var crawlUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Crawl and save user information",
	Long:  "Crawl and save user information",
	Run: func(cmd *cobra.Command, args []string) {
		db := GetDB(GetEngine())

		client, err := atcoder.NewAtCoderClient()
		if err != nil {
			slog.Error("failed to instantiate atcoder client", slog.Any("error", err))
			os.Exit(1)
		}

		crawler := crawl.NewUserCrawler(
			client,
			repository.NewUserRepository(db),
			config.Config.User.Crawl,
		)

		batch.RunBatch(crawler)
	},
}

var crawlSubmissionCmd = &cobra.Command{
	Use:   "submission",
	Short: "Crawl and save submissions",
	Long:  "Crawl and save submissions",
	Run: func(cmd *cobra.Command, args []string) {
		db := GetDB(GetEngine())

		client, err := atcoder.NewAtCoderClient()
		if err != nil {
			slog.Error("failed to instantiate atcoder client", slog.Any("error", err))
			os.Exit(1)
		}

		crawler := crawl.NewSubmissionCrawler(
			client,
			repository.NewSubmissionRepository(db),
			repository.NewContestRepository(db),
			repository.NewSubmissionCrawlHistoryRepository(db),
			config.Config.Submission.Crawl,
		)

		batch.RunBatch(crawler)
	},
}

func init() {
	crawlCmd.PersistentFlags().IntP("duration", "d", 1000, "Duration[ms] in crawling problem")
	viper.BindPFlag("problem.crawl.duration", crawlCmd.PersistentFlags().Lookup("duration"))
	viper.BindPFlag("user.crawl.duration", crawlCmd.PersistentFlags().Lookup("duration"))
	viper.BindPFlag("submission.crawl.duration", crawlCmd.PersistentFlags().Lookup("duration"))

	crawlProblemCmd.Flags().BoolP("all", "a", false, "When true, crawl all problems")
	viper.BindPFlag("problem.crawl.all", crawlProblemCmd.Flags().Lookup("all"))

	crawlSubmissionCmd.Flags().IntP("retry", "r", 0, "Limit of the number of retry when an error occurred in crawling submissions.")
	viper.BindPFlag("submission.crawl.retry", crawlSubmissionCmd.Flags().Lookup("retry"))
	crawlSubmissionCmd.Flags().String("target", "", "Target category to crawl. Multiple categories can be specified by separating tem with comma. If not specified, all categories will be crawled.")
	viper.BindPFlag("submission.crawl.targets", crawlSubmissionCmd.Flags().Lookup("target"))

	crawlCmd.AddCommand(crawlProblemCmd)
	crawlCmd.AddCommand(crawlUserCmd)
	crawlCmd.AddCommand(crawlSubmissionCmd)
	rootCmd.AddCommand(crawlCmd)
}
