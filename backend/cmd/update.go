/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fjnkt98/atcodersearch/acs"
	"fjnkt98/atcodersearch/problem"
	"fjnkt98/atcodersearch/solr"
	"fjnkt98/atcodersearch/user"
	"fmt"
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
		type Options struct {
			SaveDir            string `json:"save-dir"`
			SkipFetch          bool   `json:"skip-fetch"`
			Optimize           bool   `json:"optimize"`
			ChunkSize          int    `json:"chunk-size"`
			GenerateConcurrent int    `json:"generate-concurrent"`
			PostConcurrent     int    `json:"post-concurrent"`
			Duration           int    `json:"duration"`
			All                bool   `json:"all"`
		}

		opt := Options{}
		var err error

		if opt.SaveDir, err = GetSaveDir(cmd, "problem"); err != nil {
			log.Fatal(err.Error())
		}
		if opt.SkipFetch, err = cmd.Flags().GetBool("skip-fetch"); err != nil {
			log.Fatalf("failed to get flag `--skip-fetch`: %s", err.Error())
		}
		if opt.Optimize, err = cmd.Flags().GetBool("optimize"); err != nil {
			log.Fatalf("failed to get value of `optimize` flag: %s", err.Error())
		}
		if opt.ChunkSize, err = cmd.Flags().GetInt("chunk-size"); err != nil {
			log.Fatalf("failed to get flag `--chunk-size`: %s", err.Error())
		}
		if opt.GenerateConcurrent, err = cmd.Flags().GetInt("generate-concurrent"); err != nil {
			log.Fatalf("failed to get flag `--generate-concurrent`: %s", err.Error())
		}
		if opt.PostConcurrent, err = cmd.Flags().GetInt("post-concurrent"); err != nil {
			log.Fatalf("failed to get flag `--post-concurrent`: %s", err.Error())
		}
		if opt.Duration, err = cmd.Flags().GetInt("duration"); err != nil {
			log.Fatalf("failed to get flag `duration`: %s", err.Error())
		}
		if opt.All, err = cmd.Flags().GetBool("all"); err != nil {
			log.Fatalf("failed to get flag `all`: %s", err.Error())
		}
		solrURL := os.Getenv("SOLR_HOST")
		if solrURL == "" {
			log.Fatalln("environment variable `SOLR_HOST` must be set.")
		}
		core, err := solr.NewSolrCore[any, any]("problem", solrURL)
		if err != nil {
			log.Fatalf("failed to create `problem` core: %s", err.Error())
		}

		db := GetDB()

		f := func() error {
			if !opt.SkipFetch {
				contestCrawler := problem.NewContestCrawler(db)
				if err := contestCrawler.Run(); err != nil {
					return fmt.Errorf("failed to save contest information: %w", err)
				}

				difficultyCrawler := problem.NewDifficultyCrawler(db)
				if err := difficultyCrawler.Run(); err != nil {
					return fmt.Errorf("failed to save difficulty information: %w", err)
				}

				problemCrawler := problem.NewProblemCrawler(db)
				if err := problemCrawler.Run(opt.All, opt.Duration); err != nil {
					return fmt.Errorf("failed to save problem information: %w", err)
				}
			}

			generator := problem.NewDocumentGenerator(db, opt.SaveDir)
			if err := generator.Run(opt.ChunkSize, opt.GenerateConcurrent); err != nil {
				return fmt.Errorf("failed to generate document: %w", err)
			}

			uploader := acs.NewDefaultDocumentUploader(core, opt.SaveDir)
			if err := uploader.PostDocument(opt.Optimize, opt.PostConcurrent); err != nil {
				return fmt.Errorf("failed to post documents: %w", err)
			}
			return nil
		}

		options, err := json.Marshal(opt)
		if err != nil {
			log.Fatalf("failed to marshal update options: %s", err.Error())
		}

		history := NewUpdateHistory(db, "problem", string(options))
		if err := f(); err == nil {
			if err := history.Finish(); err != nil {
				log.Fatalf("failed to save update history: %s", err.Error())
			}
			log.Print("problem update successfully finished.")
		} else {
			if err := history.Cancel(); err != nil {
				log.Printf("failed to save update history: %s", err.Error())
			}
			log.Fatalf("problem update failed: %s", err.Error())
		}
	},
}

var updateUserCmd = &cobra.Command{
	Use:   "user",
	Short: "update user index",
	Long:  "update user index",
	Run: func(cmd *cobra.Command, args []string) {
		type Options struct {
			SaveDir            string `json:"save-dir"`
			SkipFetch          bool   `json:"skip-fetch"`
			Optimize           bool   `json:"optimize"`
			ChunkSize          int    `json:"chunk-size"`
			GenerateConcurrent int    `json:"generate-concurrent"`
			PostConcurrent     int    `json:"post-concurrent"`
			Duration           int    `json:"duration"`
		}

		opt := Options{}
		var err error

		if opt.SaveDir, err = GetSaveDir(cmd, "user"); err != nil {
			log.Fatal(err.Error())
		}
		if opt.SkipFetch, err = cmd.Flags().GetBool("skip-fetch"); err != nil {
			log.Fatalf("failed to get flag `--skip-fetch`: %s", err.Error())
		}
		if opt.Optimize, err = cmd.Flags().GetBool("optimize"); err != nil {
			log.Fatalf("failed to get value of `optimize` flag: %s", err.Error())
		}
		if opt.ChunkSize, err = cmd.Flags().GetInt("chunk-size"); err != nil {
			log.Fatalf("failed to get flag `--chunk-size`: %s", err.Error())
		}
		if opt.GenerateConcurrent, err = cmd.Flags().GetInt("generate-concurrent"); err != nil {
			log.Fatalf("failed to get flag `--generate-concurrent`: %s", err.Error())
		}
		if opt.PostConcurrent, err = cmd.Flags().GetInt("post-concurrent"); err != nil {
			log.Fatalf("failed to get flag `--post-concurrent`: %s", err.Error())
		}
		if opt.Duration, err = cmd.Flags().GetInt("duration"); err != nil {
			log.Fatalf("failed to get flag `duration`: %s", err.Error())
		}
		solrURL := os.Getenv("SOLR_HOST")
		if solrURL == "" {
			log.Fatalln("environment variable `SOLR_HOST` must be set.")
		}
		core, err := solr.NewSolrCore[any, any]("user", solrURL)
		if err != nil {
			log.Fatalf("failed to create `user` core: %s", err.Error())
		}

		db := GetDB()

		f := func() error {
			if !opt.SkipFetch {
				crawler := user.NewUserCrawler(db)
				if err := crawler.Run(opt.Duration); err != nil {
					return fmt.Errorf("failed to save user information: %w", err)
				}
			}

			generator := user.NewDocumentGenerator(db, opt.SaveDir)
			if err := generator.Run(opt.ChunkSize, opt.GenerateConcurrent); err != nil {
				return fmt.Errorf("failed to generate document: %w", err)
			}

			uploader := acs.NewDefaultDocumentUploader(core, opt.SaveDir)
			if err := uploader.PostDocument(opt.Optimize, opt.PostConcurrent); err != nil {
				return fmt.Errorf("failed to post documents: %w", err)
			}
			return nil
		}

		options, err := json.Marshal(opt)
		if err != nil {
			log.Fatalf("failed to marshal update options: %s", err.Error())
		}

		history := NewUpdateHistory(db, "user", string(options))
		if err := f(); err == nil {
			if err := history.Finish(); err != nil {
				log.Fatalf("failed to save update history: %s", err.Error())
			}
			log.Print("user update successfully finished.")
		} else {
			if err := history.Cancel(); err != nil {
				log.Printf("failed to save update history: %s", err.Error())
			}
			log.Fatalf("user update failed: %s", err.Error())
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
	rootCmd.AddCommand(updateCmd)
}
