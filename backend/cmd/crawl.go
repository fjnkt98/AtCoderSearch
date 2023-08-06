package cmd

import (
	"fjnkt98/atcodersearch/problem"
	"fjnkt98/atcodersearch/user"
	"log"

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
			log.Fatalf("failed to save contest information: %s", err.Error())
		}

		difficultyCrawler := problem.NewDifficultyCrawler(db)
		if err := difficultyCrawler.Run(); err != nil {
			log.Fatalf("failed to save difficulty information: %s", err.Error())
		}

		all, err := cmd.Flags().GetBool("all")
		if err != nil {
			log.Fatalf("failed to get flag `all`: %s", err.Error())
		}
		duration, err := cmd.Flags().GetInt("duration")
		if err != nil {
			log.Fatalf("failed to get flag `duration`: %s", err.Error())
		}
		problemCrawler := problem.NewProblemCrawler(db)
		if err := problemCrawler.Run(all, duration); err != nil {
			log.Fatalf("failed to save problem information: %s", err.Error())
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
			log.Fatalf("failed to get value of `duration` flag: %s", err.Error())
		}

		if err := crawler.Run(duration); err != nil {
			log.Fatalf("failed to save contest information: %s", err.Error())
		}
	},
}

func init() {
	crawlProblemCmd.Flags().BoolP("all", "a", false, "When true, crawl all problems")
	crawlCmd.PersistentFlags().Int("duration", 1000, "Duration[ms] in crawling problem")

	crawlCmd.AddCommand(crawlProblemCmd)
	crawlCmd.AddCommand(crawlUserCmd)
	rootCmd.AddCommand(crawlCmd)
}
