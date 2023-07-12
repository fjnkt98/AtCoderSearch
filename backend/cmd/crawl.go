package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var crawlCmd = &cobra.Command{
	Use:   "crawl",
	Short: "Crawl and save web resource",
	Long:  "Crawl and save web resource",
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
	crawlCmd.Flags().BoolP("all", "a", false, "When true, crawl all problems")
	rootCmd.AddCommand(crawlCmd)
}
