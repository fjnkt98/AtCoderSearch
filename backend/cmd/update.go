/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fjnkt98/atcodersearch/acs"
	"fjnkt98/atcodersearch/problem"
	"fjnkt98/atcodersearch/recommend"
	"fjnkt98/atcodersearch/solr"
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
		cfg := problem.UpdateConfig{}
		var err error

		if cfg.SaveDir, err = GetSaveDir(cmd, "problem"); err != nil {
			slog.Error("failed to get save dir", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}
		cfg.SkipFetch = GetBool(cmd, "skip-fetch")
		cfg.Optimize = GetBool(cmd, "optimize")
		cfg.ChunkSize = GetInt(cmd, "chunk-size")
		cfg.GenerateConcurrent = GetInt(cmd, "generate-concurrent")
		cfg.PostConcurrent = GetInt(cmd, "post-concurrent")
		cfg.Duration = GetInt(cmd, "duration")
		cfg.All = GetBool(cmd, "all")

		solrURL := os.Getenv("SOLR_HOST")
		if solrURL == "" {
			slog.Error("environment variable `SOLR_HOST` must be set.")
			os.Exit(1)
		}
		core, err := solr.NewSolrCore("problem", solrURL)
		if err != nil {
			slog.Error("failed to create `problem` core", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}

		if doMigrate := GetBool(cmd, "migrate"); doMigrate {
			DoMigrate()
		}

		db := GetDB()
		ctx, cancel := context.WithCancel(context.Background())
		eg, ctx := errgroup.WithContext(ctx)

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)

		done := make(chan Msg, 1)

		eg.Go(func() error {
			if err := problem.Update(ctx, cfg, db, core); err != nil {
				failure.Wrap(err)
			}
			done <- Msg{}
			return nil
		})

		eg.Go(func() error {
			select {
			case <-quit:
				defer cancel()
				return failure.New(acs.Interrupt, failure.Message("updating problem index has been canceled"))
			case <-ctx.Done():
				return nil
			case <-done:
				return nil
			}
		})

		if err := eg.Wait(); err != nil {
			if failure.Is(err, acs.Interrupt) {
				slog.Error("updating problem index has been canceled", slog.String("error", fmt.Sprintf("%+v", err)))
				return
			} else {
				slog.Error("failed to update problem index", slog.String("error", fmt.Sprintf("%+v", err)))
				os.Exit(1)
			}
		}

		slog.Info("updating problem index successfully finished.")
	},
}

var updateUserCmd = &cobra.Command{
	Use:   "user",
	Short: "update user index",
	Long:  "update user index",
	Run: func(cmd *cobra.Command, args []string) {

		cfg := user.UpdateConfig{}
		var err error

		if cfg.SaveDir, err = GetSaveDir(cmd, "user"); err != nil {
			slog.Error("failed to get save dir", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}

		cfg.SkipFetch = GetBool(cmd, "skip-fetch")
		cfg.Optimize = GetBool(cmd, "optimize")
		cfg.ChunkSize = GetInt(cmd, "chunk-size")
		cfg.GenerateConcurrent = GetInt(cmd, "generate-concurrent")
		cfg.PostConcurrent = GetInt(cmd, "post-concurrent")
		cfg.Duration = GetInt(cmd, "duration")

		solrURL := os.Getenv("SOLR_HOST")
		if solrURL == "" {
			slog.Error("environment variable `SOLR_HOST` must be set.")
			os.Exit(1)
		}
		core, err := solr.NewSolrCore("user", solrURL)
		if err != nil {
			slog.Error("failed to create `user` core", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}

		if doMigrate := GetBool(cmd, "migrate"); doMigrate {
			DoMigrate()
		}

		db := GetDB()
		ctx, cancel := context.WithCancel(context.Background())
		eg, ctx := errgroup.WithContext(ctx)

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)

		done := make(chan Msg, 1)

		eg.Go(func() error {
			if err := user.Update(ctx, cfg, db, core); err != nil {
				failure.Wrap(err)
			}
			done <- Msg{}
			return nil
		})

		eg.Go(func() error {
			select {
			case <-quit:
				defer cancel()
				return failure.New(acs.Interrupt, failure.Message("updating user index has been canceled"))
			case <-ctx.Done():
				return nil
			case <-done:
				return nil
			}
		})

		if err := eg.Wait(); err != nil {
			if failure.Is(err, acs.Interrupt) {
				slog.Error("updating user index has been canceled", slog.String("error", fmt.Sprintf("%+v", err)))
				return
			} else {
				slog.Error("failed to update user index", slog.String("error", fmt.Sprintf("%+v", err)))
				os.Exit(1)
			}
		}

		slog.Info("updating user index successfully finished.")
	},
}

var updateSubmissionCmd = &cobra.Command{
	Use:   "submission",
	Short: "update submission index",
	Long:  "update submission index",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := submission.UpdateConfig{}
		var err error

		if cfg.SaveDir, err = GetSaveDir(cmd, "submission"); err != nil {
			slog.Error("failed to get save dir", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}

		cfg.Optimize = GetBool(cmd, "optimize")
		cfg.ChunkSize = GetInt(cmd, "chunk-size")
		cfg.GenerateConcurrent = GetInt(cmd, "generate-concurrent")
		cfg.PostConcurrent = GetInt(cmd, "post-concurrent")
		cfg.All = GetBool(cmd, "all")

		solrURL := os.Getenv("SOLR_HOST")
		if solrURL == "" {
			slog.Error("environment variable `SOLR_HOST` must be set.")
			os.Exit(1)
		}
		core, err := solr.NewSolrCore("submission", solrURL)
		if err != nil {
			slog.Error("failed to create `submission` core", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}

		if doMigrate := GetBool(cmd, "migrate"); doMigrate {
			DoMigrate()
		}

		db := GetDB()
		ctx, cancel := context.WithCancel(context.Background())
		eg, ctx := errgroup.WithContext(ctx)

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)

		done := make(chan Msg, 1)

		eg.Go(func() error {
			if err := submission.Update(ctx, cfg, db, core); err != nil {
				return failure.Wrap(err)
			}
			done <- Msg{}
			return nil
		})

		eg.Go(func() error {
			select {
			case <-quit:
				defer cancel()
				return failure.New(acs.Interrupt, failure.Message("updating submission index has been canceled"))
			case <-ctx.Done():
				return nil
			case <-done:
				return nil
			}
		})

		if err := eg.Wait(); err != nil {
			if failure.Is(err, acs.Interrupt) {
				slog.Error("updating submission index has been canceled", slog.String("error", fmt.Sprintf("%+v", err)))
				return
			} else {
				slog.Error("failed to update submission index", slog.String("error", fmt.Sprintf("%+v", err)))
				os.Exit(1)
			}
		}

		slog.Info("updating submission index successfully finished.")
	},
}

var updateRecommendCmd = &cobra.Command{
	Use:   "recommend",
	Short: "update recommend index",
	Long:  "update recommend index",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := recommend.UpdateConfig{}
		var err error

		if cfg.SaveDir, err = GetSaveDir(cmd, "recommend"); err != nil {
			slog.Error("failed to get save dir", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}

		cfg.Optimize = GetBool(cmd, "optimize")
		cfg.ChunkSize = GetInt(cmd, "chunk-size")
		cfg.GenerateConcurrent = GetInt(cmd, "generate-concurrent")
		cfg.PostConcurrent = GetInt(cmd, "post-concurrent")

		solrURL := os.Getenv("SOLR_HOST")
		if solrURL == "" {
			slog.Error("environment variable `SOLR_HOST` must be set.")
			os.Exit(1)
		}
		core, err := solr.NewSolrCore("recommend", solrURL)
		if err != nil {
			slog.Error("failed to create `recommend` core", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}

		if doMigrate := GetBool(cmd, "migrate"); doMigrate {
			DoMigrate()
		}

		db := GetDB()
		ctx, cancel := context.WithCancel(context.Background())
		eg, ctx := errgroup.WithContext(ctx)

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)

		done := make(chan Msg, 1)

		eg.Go(func() error {
			if err := recommend.Update(ctx, cfg, db, core); err != nil {
				return failure.Wrap(err)
			}
			done <- Msg{}
			return nil
		})

		eg.Go(func() error {
			select {
			case <-quit:
				defer cancel()
				return failure.New(acs.Interrupt, failure.Message("updating recommend index has been canceled"))
			case <-ctx.Done():
				return nil
			case <-done:
				return nil
			}
		})

		if err := eg.Wait(); err != nil {
			if failure.Is(err, acs.Interrupt) {
				slog.Error("updating recommend index has been canceled", slog.String("error", fmt.Sprintf("%+v", err)))
				return
			} else {
				slog.Error("failed to update recommend index", slog.String("error", fmt.Sprintf("%+v", err)))
				os.Exit(1)
			}
		}

		slog.Info("updating recommend index successfully finished.")
	},
}

func init() {
	updateCmd.PersistentFlags().Bool("migrate", false, "Execute database migration before update index.")
	updateCmd.PersistentFlags().String("save-dir", "", "Directory path at which generated documents will be saved.")
	updateCmd.PersistentFlags().BoolP("optimize", "o", true, "Optimize index if true.")
	updateCmd.PersistentFlags().Int("chunk-size", 1000, "Number of documents to write in 1 file.")
	updateCmd.PersistentFlags().Int("generate-concurrent", 6, "Number of concurrent document generation processes")
	updateCmd.PersistentFlags().Int("post-concurrent", 4, "Number of concurrent document upload processes")

	updateProblemCmd.Flags().BoolP("all", "a", false, "Crawl all problems if true.")
	updateProblemCmd.Flags().BoolP("skip-fetch", "f", false, "Skip crawling if true.")
	updateProblemCmd.Flags().Int("duration", 1000, "Interval time[ms] for crawling.")

	updateUserCmd.Flags().BoolP("skip-fetch", "f", false, "Skip crawling if true.")
	updateUserCmd.Flags().Int("duration", 1000, "Interval time[ms] for crawling.")

	updateSubmissionCmd.Flags().BoolP("all", "a", false, "Update all submissions.")

	updateCmd.AddCommand(updateProblemCmd)
	updateCmd.AddCommand(updateUserCmd)
	updateCmd.AddCommand(updateSubmissionCmd)
	updateCmd.AddCommand(updateRecommendCmd)

	rootCmd.AddCommand(updateCmd)
}
