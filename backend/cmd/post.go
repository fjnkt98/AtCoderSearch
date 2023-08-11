package cmd

import (
	"fjnkt98/atcodersearch/acs"
	"fjnkt98/atcodersearch/solr"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var postCmd = &cobra.Command{
	Use:   "post",
	Short: "Post document JSON files into Solr core",
	Long:  "Post document JSON files into Solr core",
}

var postProblemCmd = &cobra.Command{
	Use:   "problem",
	Short: "Post document JSON files into problem core",
	Long:  "Post document JSON files into problem core",
	Run: func(cmd *cobra.Command, args []string) {
		saveDir, err := GetSaveDir(cmd, "problem")
		if err != nil {
			log.Fatalf("%+v", err)
		}

		solrURL := os.Getenv("SOLR_HOST")
		if solrURL == "" {
			log.Fatalln("environment variable `SOLR_HOST` must be set.")
		}
		core, err := solr.NewSolrCore[any, any]("problem", solrURL)
		if err != nil {
			log.Fatalf("failed to create `problem` core: %+v", err)
		}

		uploader := acs.NewDefaultDocumentUploader(core, saveDir)
		optimize, err := cmd.Flags().GetBool("optimize")
		if err != nil {
			log.Fatalf("failed to get value of `optimize` flag: %+v", err)
		}
		concurrent, err := cmd.Flags().GetInt("concurrent")
		if err != nil {
			log.Fatalf("failed to get value of `concurrent` flag: %+v", err)
		}
		if err := uploader.PostDocument(optimize, concurrent); err != nil {
			log.Fatalf("%+v", err)
		}
	},
}

var postUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Post document JSON files into user core",
	Long:  "Post document JSON files into user core",
	Run: func(cmd *cobra.Command, args []string) {
		saveDir, err := GetSaveDir(cmd, "user")
		if err != nil {
			log.Fatalf("%+v", err)
		}

		solrURL := os.Getenv("SOLR_HOST")
		if solrURL == "" {
			log.Fatalln("environment variable `SOLR_HOST` must be set.")
		}
		core, err := solr.NewSolrCore[any, any]("user", solrURL)
		if err != nil {
			log.Fatalf("failed to create `user` core: %+v", err)
		}

		uploader := acs.NewDefaultDocumentUploader(core, saveDir)
		optimize, err := cmd.Flags().GetBool("optimize")
		if err != nil {
			log.Fatalf("failed to get value of `optimize` flag: %+v", err)
		}
		concurrent, err := cmd.Flags().GetInt("concurrent")
		if err != nil {
			log.Fatalf("failed to get value of `concurrent` flag: %+v", err)
		}
		if err := uploader.PostDocument(optimize, concurrent); err != nil {
			log.Fatalf("%+v", err)
		}
	},
}

var postSubmissionCmd = &cobra.Command{
	Use:   "submission",
	Short: "Post document JSON files into submission core",
	Long:  "Post document JSON files into submission core",
	Run: func(cmd *cobra.Command, args []string) {
		saveDir, err := GetSaveDir(cmd, "submission")
		if err != nil {
			log.Fatalf("%+v", err)
		}

		solrURL := os.Getenv("SOLR_HOST")
		if solrURL == "" {
			log.Fatalln("environment variable `SOLR_HOST` must be set.")
		}
		core, err := solr.NewSolrCore[any, any]("submission", solrURL)
		if err != nil {
			log.Fatalf("failed to create `user` core: %+v", err)
		}

		uploader := acs.NewDefaultDocumentUploader(core, saveDir)
		optimize, err := cmd.Flags().GetBool("optimize")
		if err != nil {
			log.Fatalf("failed to get value of `optimize` flag: %+v", err)
		}
		concurrent, err := cmd.Flags().GetInt("concurrent")
		if err != nil {
			log.Fatalf("failed to get value of `concurrent` flag: %+v", err)
		}
		if err := uploader.PostDocument(optimize, concurrent); err != nil {
			log.Fatalf("%+v", err)
		}
	},
}

func init() {
	postCmd.PersistentFlags().BoolP("optimize", "o", false, "When true, send optimize request to Solr")
	postCmd.PersistentFlags().String("save-dir", "", "Directory path at which generated documents will be saved")
	postCmd.PersistentFlags().Int("concurrent", 3, "Concurrent number of document upload processes")
	postCmd.AddCommand(postProblemCmd)
	postCmd.AddCommand(postUserCmd)
	postCmd.AddCommand(postSubmissionCmd)
	rootCmd.AddCommand(postCmd)
}
