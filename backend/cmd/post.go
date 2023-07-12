package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var postCmd = &cobra.Command{
	Use:   "post <domain>",
	Short: "Post document JSON files into Solr core",
	Long:  "Post document JSON files into Solr core",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("arg `domain` is mandatory")
		}

		domain := args[0]
		switch domain {
		case "problems":
			log.Println("generate problems")
		case "users":
			log.Println("generate users")
		case "recommends":
			log.Println("generate recommends")
		default:
			log.Fatalf(
				"%s is invalid.\nvariety of domain is below:\n- problems\n- users\n- recommends",
				domain,
			)
		}
	},
}

func init() {
	postCmd.Flags().BoolP("optimize", "o", false, "When true, send optimize request to Solr")
	postCmd.Flags().String("save-dir", "", "Directory path at which generated documents will be saved")
	rootCmd.AddCommand(postCmd)
}
