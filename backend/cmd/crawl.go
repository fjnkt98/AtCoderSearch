package cmd

import (
	"fjnkt98/atcodersearch/atcoder"
	"fjnkt98/atcodersearch/problem"
	"fjnkt98/atcodersearch/submission"
	"fjnkt98/atcodersearch/user"
	"fmt"
	"os"

	"github.com/spf13/cobra"
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
		db := GetDB()
		contestCrawler := problem.NewContestCrawler(db)
		if err := contestCrawler.Run(); err != nil {
			slog.Error("failed to save contest information", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}

		difficultyCrawler := problem.NewDifficultyCrawler(db)
		if err := difficultyCrawler.Run(); err != nil {
			slog.Error("failed to save difficulty information", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}

		all := GetBool(cmd, "all")
		duration := GetInt(cmd, "duration")

		problemCrawler := problem.NewProblemCrawler(db)
		if err := problemCrawler.Run(all, duration); err != nil {
			slog.Error("failed to save problem information", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}
	},
}

var crawlUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Crawl and save user information",
	Long:  "Crawl and save user information",
	Run: func(cmd *cobra.Command, args []string) {
		db := GetDB()
		crawler := user.NewUserCrawler(db)
		duration := GetInt(cmd, "duration")

		if err := crawler.Run(duration); err != nil {
			slog.Error("failed to save user information", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}
	},
}

var crawlSubmissionCmd = &cobra.Command{
	Use:   "submission",
	Short: "Crawl and save submissions",
	Long:  "Crawl and save submissions",
	Run: func(cmd *cobra.Command, args []string) {
		db := GetDB()

		username := os.Getenv("ATCODER_USER_NAME")
		password := os.Getenv("ATCODER_USER_PASSWORD")

		slog.Info("Login to AtCoder...")
		client, err := atcoder.NewAtCoderClient(username, password)
		if err != nil {
			slog.Error("failed to login atcoder", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}
		slog.Info("Successfully logged in to AtCoder.")

		crawler := submission.NewCrawler(client, db)
		duration := GetInt(cmd, "duration")

		slog.Info("Start to crawl submissions")
		if err := crawler.Run(duration); err != nil {
			slog.Error("failed to save submissions", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}
	},
}

func init() {
	crawlProblemCmd.Flags().BoolP("all", "a", false, "When true, crawl all problems")
	crawlCmd.PersistentFlags().Int("duration", 1000, "Duration[ms] in crawling problem")

	crawlCmd.AddCommand(crawlProblemCmd)
	crawlCmd.AddCommand(crawlUserCmd)
	crawlCmd.AddCommand(crawlSubmissionCmd)
	rootCmd.AddCommand(crawlCmd)
}
