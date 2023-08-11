package cmd

import (
	"fjnkt98/atcodersearch/atcoder"
	"fjnkt98/atcodersearch/problem"
	"fjnkt98/atcodersearch/submission"
	"fjnkt98/atcodersearch/user"
	"log"
	"os"

	"github.com/spf13/cobra"
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
			log.Fatalf("failed to save contest information: %+v", err)
		}

		difficultyCrawler := problem.NewDifficultyCrawler(db)
		if err := difficultyCrawler.Run(); err != nil {
			log.Fatalf("failed to save difficulty information: %+v", err)
		}

		all, err := cmd.Flags().GetBool("all")
		if err != nil {
			log.Fatalf("failed to get flag `all`: %+v", err)
		}
		duration, err := cmd.Flags().GetInt("duration")
		if err != nil {
			log.Fatalf("failed to get flag `duration`: %+v", err)
		}
		problemCrawler := problem.NewProblemCrawler(db)
		if err := problemCrawler.Run(all, duration); err != nil {
			log.Fatalf("failed to save problem information: %+v", err)
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
		duration, err := cmd.Flags().GetInt("duration")
		if err != nil {
			log.Fatalf("failed to get value of `duration` flag: %+v", err)
		}

		if err := crawler.Run(duration); err != nil {
			log.Fatalf("failed to save user information: %+v", err)
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

		log.Println("Login to AtCoder...")
		client, err := atcoder.NewAtCoderClient(username, password)
		if err != nil {
			log.Fatalf("failed to login atcoder: %+v", err)
		}
		log.Println("Successfully logged in to AtCoder.")

		crawler := submission.NewCrawler(client, db)
		duration, err := cmd.Flags().GetInt("duration")
		if err != nil {
			log.Fatalf("failed to get value of `duration` flag: %+v", err)
		}

		log.Println("Start to crawl submissions")
		if err := crawler.Run(duration); err != nil {
			log.Fatalf("failed to save submissions: %+v", err)
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
