package cmd

import (
	"context"
	"fjnkt98/atcodersearch/acs"
	"fjnkt98/atcodersearch/atcoder"
	"fjnkt98/atcodersearch/problem"
	"fjnkt98/atcodersearch/submission"
	"fjnkt98/atcodersearch/user"
	"fmt"
	"os"
	"os/signal"

	"github.com/morikuni/failure"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
	"golang.org/x/sync/errgroup"
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

		ctx, cancel := context.WithCancel(context.Background())
		eg, ctx := errgroup.WithContext(ctx)

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)

		done := make(chan Msg, 1)

		eg.Go(func() error {
			contestCrawler := problem.NewContestCrawler(db)
			if err := contestCrawler.Run(ctx); err != nil {
				return failure.Wrap(err)
			}

			difficultyCrawler := problem.NewDifficultyCrawler(db)
			if err := difficultyCrawler.Run(ctx); err != nil {
				return failure.Wrap(err)
			}

			all := GetBool(cmd, "all")
			duration := GetInt(cmd, "duration")

			problemCrawler := problem.NewProblemCrawler(db)
			if err := problemCrawler.Run(ctx, all, duration); err != nil {
				return failure.Wrap(err)
			}

			done <- Msg{}

			return nil
		})

		eg.Go(func() error {
			select {
			case <-quit:
				defer cancel()
				return failure.New(acs.Interrupt, failure.Message("problem crawling has been interrupted"))
			case <-ctx.Done():
				return nil
			case <-done:
				return nil
			}
		})

		if err := eg.Wait(); err != nil {
			if failure.Is(err, acs.Interrupt) {
				slog.Error("problem crawling has been interrupted", slog.String("error", fmt.Sprintf("%+v", err)))
				return
			} else {
				slog.Error("failed to crawl problems", slog.String("error", fmt.Sprintf("%+v", err)))
				os.Exit(1)
			}
		}
		slog.Info("finished crawl problems successfully.")
	},
}

var crawlUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Crawl and save user information",
	Long:  "Crawl and save user information",
	Run: func(cmd *cobra.Command, args []string) {
		db := GetDB()

		ctx, cancel := context.WithCancel(context.Background())
		eg, ctx := errgroup.WithContext(ctx)

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)

		done := make(chan Msg, 1)

		eg.Go(func() error {
			crawler := user.NewUserCrawler(db)
			duration := GetInt(cmd, "duration")

			if err := crawler.Run(ctx, duration); err != nil {
				return failure.Wrap(err)
			}

			return nil
		})

		eg.Go(func() error {
			select {
			case <-quit:
				defer cancel()
				return failure.New(acs.Interrupt, failure.Message("problem crawling has been interrupted"))
			case <-ctx.Done():
				return nil
			case <-done:
				return nil
			}
		})

		if err := eg.Wait(); err != nil {
			if failure.Is(err, acs.Interrupt) {
				slog.Error("user crawling has been interrupted", slog.String("error", fmt.Sprintf("%+v", err)))
				return
			} else {
				slog.Error("failed to crawl users", slog.String("error", fmt.Sprintf("%+v", err)))
				os.Exit(1)
			}
		}
		slog.Info("finished crawl problems successfully.")
	},
}

var crawlSubmissionCmd = &cobra.Command{
	Use:   "submission",
	Short: "Crawl and save submissions",
	Long:  "Crawl and save submissions",
	Run: func(cmd *cobra.Command, args []string) {
		db := GetDB()

		endless := GetBool(cmd, "endless")

		ctx, cancel := context.WithCancel(context.Background())
		eg, ctx := errgroup.WithContext(ctx)

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)

		done := make(chan Msg, 1)

		username := os.Getenv("ATCODER_USER_NAME")
		password := os.Getenv("ATCODER_USER_PASSWORD")

		eg.Go(func() error {
			slog.Info("Login to AtCoder...")
			client, err := atcoder.NewAtCoderClient(ctx, username, password)
			if err != nil {
				return failure.Wrap(err)
			}
			slog.Info("Successfully logged in to AtCoder.")

			crawler := submission.NewCrawler(client, db)
			duration := GetInt(cmd, "duration")

			slog.Info("Start to crawl submissions")
			if endless {
				for {
					select {
					case <-quit:
						return failure.New(acs.Interrupt, failure.Message("crawling problems has been interrupted."))
					default:
						if err := crawler.Run(ctx, duration); err != nil {
							slog.Error("failed to crawl submissions", slog.String("error", fmt.Sprintf("%+v", err)))
						}
					}
				}
			} else {
				if err := crawler.Run(ctx, duration); err != nil {
					return failure.Wrap(err)
				}
			}

			done <- Msg{}

			return nil
		})

		eg.Go(func() error {
			select {
			case <-quit:
				defer cancel()
				return failure.New(acs.Interrupt, failure.Message("problem crawling has been interrupted"))
			case <-ctx.Done():
				return nil
			case <-done:
				return nil
			}
		})

		if err := eg.Wait(); err != nil {
			if failure.Is(err, acs.Interrupt) {
				slog.Error("submissions crawling has been interrupted", slog.String("error", fmt.Sprintf("%+v", err)))
				return
			} else {
				slog.Error("failed to crawl submissions", slog.String("error", fmt.Sprintf("%+v", err)))
				os.Exit(1)
			}
		}
		slog.Info("finished crawl submissions successfully.")
	},
}

func init() {
	crawlProblemCmd.Flags().BoolP("all", "a", false, "When true, crawl all problems")
	crawlSubmissionCmd.Flags().BoolP("endless", "e", false, "When true, crawler will continue crawling even if an error occurred.")
	crawlCmd.PersistentFlags().Int("duration", 1000, "Duration[ms] in crawling problem")

	crawlCmd.AddCommand(crawlProblemCmd)
	crawlCmd.AddCommand(crawlUserCmd)
	crawlCmd.AddCommand(crawlSubmissionCmd)
	rootCmd.AddCommand(crawlCmd)
}
