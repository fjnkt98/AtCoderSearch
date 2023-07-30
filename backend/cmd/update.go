/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fjnkt98/atcodersearch/atcodersearch/common"
	"fjnkt98/atcodersearch/atcodersearch/problem"
	"fjnkt98/atcodersearch/solr"
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
		saveDir, err := GetSaveDir(cmd, "problem")
		if err != nil {
			log.Fatal(err.Error())
		}

		db := GetDB()
		skipFetch, err := cmd.Flags().GetBool("skip-fetch")
		if err != nil {
			log.Fatalf("failed to get flag `--skip-fetch`: %s", err.Error())
		}

		if !skipFetch {
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
		}

		generator := problem.NewProblemDocumentGenerator(db, saveDir)
		chunkSize, err := cmd.Flags().GetInt("chunk-size")
		if err != nil {
			log.Fatalf("failed to get flag `--chunk-size`: %s", err.Error())
		}
		generateConcurrency, err := cmd.Flags().GetInt("generate-concurrent")
		if err != nil {
			log.Fatalf("failed to get flag `--generate-concurrent`: %s", err.Error())
		}

		if err := generator.Run(chunkSize, generateConcurrency); err != nil {
			log.Fatalf("failed to generate document: %s", err.Error())
		}

		postConcurrency, err := cmd.Flags().GetInt("post-concurrent")
		if err != nil {
			log.Fatalf("failed to get flag `--post-concurrent`: %s", err.Error())
		}

		solrURL := os.Getenv("SOLR_HOST")
		if solrURL == "" {
			log.Fatalln("environment variable `SOLR_HOST` must be set.")
		}
		core, err := solr.NewSolrCore[any, any]("problem", solrURL)
		if err != nil {
			log.Fatalf("failed to create `problem` core: %s", err.Error())
		}

		uploader := common.NewDefaultDocumentUploader(core, saveDir)
		optimize, err := cmd.Flags().GetBool("optimize")
		if err != nil {
			log.Fatalf("failed to get value of `optimize` flag: %s", err.Error())
		}
		if err := uploader.PostDocument(optimize, postConcurrency); err != nil {
			log.Fatalf("failed to post documents: %s", err.Error())
		}
	},
}

func init() {
	updateCmd.PersistentFlags().String("save-dir", "", "Directory path at which generated documents will be saved.")
	updateCmd.PersistentFlags().BoolP("skip-fetch", "f", false, "Skip crawling if true.")
	updateCmd.PersistentFlags().BoolP("optimize", "o", false, "Optimize index if true.")
	updateCmd.PersistentFlags().Int("chunk-size", 1000, "Number of documents to write in 1 file.")
	updateCmd.PersistentFlags().Int("generate-concurrent", 6, "Number of concurrent document generation processes")
	updateCmd.PersistentFlags().Int("post-concurrent", 4, "Number of concurrent document upload processes")

	updateProblemCmd.Flags().Int("duration", 1000, "Duration[ms] in crawling problem")
	updateProblemCmd.Flags().BoolP("all", "a", false, "Crawl all problems if true.")

	updateCmd.AddCommand(updateProblemCmd)
	rootCmd.AddCommand(updateCmd)
}
