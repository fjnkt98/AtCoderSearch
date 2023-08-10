/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fjnkt98/atcodersearch/problem"
	"fjnkt98/atcodersearch/solr"
	"fjnkt98/atcodersearch/submission"
	"fjnkt98/atcodersearch/user"
	"log"
	"os"

	"github.com/spf13/cobra"
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
			log.Fatalf("%+v", err)
		}
		if cfg.SkipFetch, err = cmd.Flags().GetBool("skip-fetch"); err != nil {
			log.Fatalf("failed to get flag `--skip-fetch`: %s", err.Error())
		}
		if cfg.Optimize, err = cmd.Flags().GetBool("optimize"); err != nil {
			log.Fatalf("failed to get value of `optimize` flag: %s", err.Error())
		}
		if cfg.ChunkSize, err = cmd.Flags().GetInt("chunk-size"); err != nil {
			log.Fatalf("failed to get flag `--chunk-size`: %s", err.Error())
		}
		if cfg.GenerateConcurrent, err = cmd.Flags().GetInt("generate-concurrent"); err != nil {
			log.Fatalf("failed to get flag `--generate-concurrent`: %s", err.Error())
		}
		if cfg.PostConcurrent, err = cmd.Flags().GetInt("post-concurrent"); err != nil {
			log.Fatalf("failed to get flag `--post-concurrent`: %s", err.Error())
		}
		if cfg.Duration, err = cmd.Flags().GetInt("duration"); err != nil {
			log.Fatalf("failed to get flag `duration`: %s", err.Error())
		}
		if cfg.All, err = cmd.Flags().GetBool("all"); err != nil {
			log.Fatalf("failed to get flag `all`: %s", err.Error())
		}
		solrURL := os.Getenv("SOLR_HOST")
		if solrURL == "" {
			log.Fatalln("environment variable `SOLR_HOST` must be set.")
		}
		core, err := solr.NewSolrCore[any, any]("problem", solrURL)
		if err != nil {
			log.Fatalf("failed to create `problem` core: %+v", err)
		}

		db := GetDB()

		cfgJSON, err := json.Marshal(cfg)
		if err != nil {
			log.Fatalf("failed to marshal update options: %s", err.Error())
		}

		history := NewUpdateHistory(db, "problem", string(cfgJSON))
		if err := problem.Update(cfg, db, core); err == nil {
			if err := history.Finish(); err != nil {
				log.Fatalf("failed to save update history: %+v", err)
			}
			log.Print("problem update successfully finished.")
		} else {
			if err := history.Cancel(); err != nil {
				log.Printf("failed to save update history: %+v", err)
			}
			log.Fatalf("problem update failed: %+v", err)
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
			log.Fatalf("%+v", err)
		}
		if cfg.SkipFetch, err = cmd.Flags().GetBool("skip-fetch"); err != nil {
			log.Fatalf("failed to get flag `--skip-fetch`: %s", err.Error())
		}
		if cfg.Optimize, err = cmd.Flags().GetBool("optimize"); err != nil {
			log.Fatalf("failed to get value of `optimize` flag: %s", err.Error())
		}
		if cfg.ChunkSize, err = cmd.Flags().GetInt("chunk-size"); err != nil {
			log.Fatalf("failed to get flag `--chunk-size`: %s", err.Error())
		}
		if cfg.GenerateConcurrent, err = cmd.Flags().GetInt("generate-concurrent"); err != nil {
			log.Fatalf("failed to get flag `--generate-concurrent`: %s", err.Error())
		}
		if cfg.PostConcurrent, err = cmd.Flags().GetInt("post-concurrent"); err != nil {
			log.Fatalf("failed to get flag `--post-concurrent`: %s", err.Error())
		}
		if cfg.Duration, err = cmd.Flags().GetInt("duration"); err != nil {
			log.Fatalf("failed to get flag `duration`: %s", err.Error())
		}
		solrURL := os.Getenv("SOLR_HOST")
		if solrURL == "" {
			log.Fatalln("environment variable `SOLR_HOST` must be set.")
		}
		core, err := solr.NewSolrCore[any, any]("user", solrURL)
		if err != nil {
			log.Fatalf("failed to create `user` core: %+v", err)
		}

		db := GetDB()

		options, err := json.Marshal(cfg)
		if err != nil {
			log.Fatalf("failed to marshal update options: %s", err.Error())
		}

		history := NewUpdateHistory(db, "user", string(options))
		if err := user.Update(cfg, db, core); err == nil {
			if err := history.Finish(); err != nil {
				log.Fatalf("failed to save update history: %+v", err)
			}
			log.Print("user update successfully finished.")
		} else {
			if err := history.Cancel(); err != nil {
				log.Printf("failed to save update history: %+v", err)
			}
			log.Fatalf("user update failed: %+v", err)
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
			log.Fatalf("%+v", err)
		}
		if cfg.Optimize, err = cmd.Flags().GetBool("optimize"); err != nil {
			log.Fatalf("failed to get value of `optimize` flag: %s", err.Error())
		}
		if cfg.ChunkSize, err = cmd.Flags().GetInt("chunk-size"); err != nil {
			log.Fatalf("failed to get flag `--chunk-size`: %s", err.Error())
		}
		if cfg.GenerateConcurrent, err = cmd.Flags().GetInt("generate-concurrent"); err != nil {
			log.Fatalf("failed to get flag `--generate-concurrent`: %s", err.Error())
		}
		if cfg.PostConcurrent, err = cmd.Flags().GetInt("post-concurrent"); err != nil {
			log.Fatalf("failed to get flag `--post-concurrent`: %s", err.Error())
		}
		solrURL := os.Getenv("SOLR_HOST")
		if solrURL == "" {
			log.Fatalln("environment variable `SOLR_HOST` must be set.")
		}
		core, err := solr.NewSolrCore[any, any]("submission", solrURL)
		if err != nil {
			log.Fatalf("failed to create `submission` core: %+v", err)
		}

		db := GetDB()

		options, err := json.Marshal(cfg)
		if err != nil {
			log.Fatalf("failed to marshal update options: %s", err.Error())
		}

		history := NewUpdateHistory(db, "submission", string(options))
		if err := submission.Update(cfg, db, core); err == nil {
			if err := history.Finish(); err != nil {
				log.Fatalf("failed to save update history: %+v", err)
			}
			log.Print("submission update successfully finished.")
		} else {
			if err := history.Cancel(); err != nil {
				log.Printf("failed to save update history: %+v", err)
			}
			log.Fatalf("submission update failed: %+v", err)
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
	rootCmd.AddCommand(updateCmd)
}
