/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fjnkt98/atcodersearch/problem"
	"fjnkt98/atcodersearch/recommend"
	"fjnkt98/atcodersearch/solr"
	"fjnkt98/atcodersearch/submission"
	"fjnkt98/atcodersearch/user"
	"fmt"
	"os"

	"github.com/spf13/cobra"
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
			slog.Error("environment variable `SOLR_HOST` must be set.", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}
		core, err := solr.NewSolrCore[any, any]("problem", solrURL)
		if err != nil {
			slog.Error("failed to create `problem` core", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}

		db := GetDB()

		cfgJSON, err := json.Marshal(cfg)
		if err != nil {
			slog.Error("failed to marshal update options", slog.String("error", err.Error()))
			os.Exit(1)
		}

		history := NewUpdateHistory(db, "problem", string(cfgJSON))
		if err := problem.Update(cfg, db, core); err == nil {
			if err := history.Finish(); err != nil {
				slog.Error("failed to save update history", slog.String("error", fmt.Sprintf("%+v", err)))
				os.Exit(1)
			}
			slog.Info("problem update successfully finished.")
		} else {
			if err := history.Cancel(); err != nil {
				slog.Error("failed to save update history", slog.String("error", fmt.Sprintf("%+v", err)))
				os.Exit(1)
			}
			slog.Error("problem update failed", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}
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
			slog.Error("environment variable `SOLR_HOST` must be set.", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}
		core, err := solr.NewSolrCore[any, any]("user", solrURL)
		if err != nil {
			slog.Error("failed to create `user` core", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}

		db := GetDB()

		options, err := json.Marshal(cfg)
		if err != nil {
			slog.Error("failed to marshal update options", slog.String("error", err.Error()))
			os.Exit(1)
		}

		history := NewUpdateHistory(db, "user", string(options))
		if err := user.Update(cfg, db, core); err == nil {
			if err := history.Finish(); err != nil {
				slog.Error("failed to save update history", slog.String("error", fmt.Sprintf("%+v", err)))
				os.Exit(1)
			}
			slog.Info("user update successfully finished.")
		} else {
			if err := history.Cancel(); err != nil {
				slog.Error("failed to save update history", slog.String("error", fmt.Sprintf("%+v", err)))
				os.Exit(1)
			}
			slog.Error("user update failed", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}
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

		solrURL := os.Getenv("SOLR_HOST")
		if solrURL == "" {
			slog.Error("environment variable `SOLR_HOST` must be set.", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}
		core, err := solr.NewSolrCore[any, any]("submission", solrURL)
		if err != nil {
			slog.Error("failed to create `submission` core", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}

		db := GetDB()

		options, err := json.Marshal(cfg)
		if err != nil {
			slog.Error("failed to marshal update options", slog.String("error", err.Error()))
			os.Exit(1)
		}

		history := NewUpdateHistory(db, "submission", string(options))
		if err := submission.Update(cfg, db, core); err == nil {
			if err := history.Finish(); err != nil {
				slog.Error("failed to save update history: %+v", slog.String("error", fmt.Sprintf("%+v", err)))
				os.Exit(1)
			}
			slog.Info("submission update successfully finished.")
		} else {
			if err := history.Cancel(); err != nil {
				slog.Error("failed to save update history", slog.String("error", fmt.Sprintf("%+v", err)))
				os.Exit(1)
			}
			slog.Error("submission update failed", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}
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
			slog.Error("environment variable `SOLR_HOST` must be set.", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}
		core, err := solr.NewSolrCore[any, any]("recommend", solrURL)
		if err != nil {
			slog.Error("failed to create `recommend` core", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}

		db := GetDB()

		options, err := json.Marshal(cfg)
		if err != nil {
			slog.Error("failed to marshal update options: %s", slog.String("error", err.Error()))
			os.Exit(1)
		}

		history := NewUpdateHistory(db, "recommend", string(options))
		if err := recommend.Update(cfg, db, core); err == nil {
			if err := history.Finish(); err != nil {
				slog.Error("failed to save update history", slog.String("error", fmt.Sprintf("%+v", err)))
				os.Exit(1)
			}
			slog.Info("recommend update successfully finished.")
		} else {
			if err := history.Cancel(); err != nil {
				slog.Error("failed to save update history", slog.String("error", fmt.Sprintf("%+v", err)))
				os.Exit(1)
			}
			slog.Error("recommend update failed", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}
	},
}

func init() {
	updateCmd.PersistentFlags().String("save-dir", "", "Directory path at which generated documents will be saved.")
	updateCmd.PersistentFlags().BoolP("skip-fetch", "f", false, "Skip crawling if true.")
	updateCmd.PersistentFlags().BoolP("optimize", "o", true, "Optimize index if true.")
	updateCmd.PersistentFlags().Int("chunk-size", 1000, "Number of documents to write in 1 file.")
	updateCmd.PersistentFlags().Int("generate-concurrent", 6, "Number of concurrent document generation processes")
	updateCmd.PersistentFlags().Int("post-concurrent", 4, "Number of concurrent document upload processes")
	updateCmd.PersistentFlags().Int("duration", 1000, "Interval time[ms] for crawling.")

	updateProblemCmd.Flags().BoolP("all", "a", false, "Crawl all problems if true.")

	updateCmd.AddCommand(updateProblemCmd)
	updateCmd.AddCommand(updateUserCmd)
	updateCmd.AddCommand(updateSubmissionCmd)
	updateCmd.AddCommand(updateRecommendCmd)

	rootCmd.AddCommand(updateCmd)
}
