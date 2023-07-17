package cmd

import (
	"fjnkt98/atcodersearch/atcodersearch/common"
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
			log.Fatal(err.Error())
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
		if err := uploader.PostDocument(optimize, 1); err != nil {
			log.Fatal(err.Error())
		}
	},
}

func init() {
	postCmd.PersistentFlags().BoolP("optimize", "o", false, "When true, send optimize request to Solr")
	postCmd.PersistentFlags().String("save-dir", "", "Directory path at which generated documents will be saved")
	postCmd.AddCommand(postProblemCmd)
	rootCmd.AddCommand(postCmd)
}
