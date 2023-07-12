package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate <domain>",
	Short: "Generate document JSON files",
	Long:  "Generate document JSON files",
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
	generateCmd.Flags().String("save-dir", "", "Directory path at which generated documents will be saved")
	rootCmd.AddCommand(generateCmd)
}
