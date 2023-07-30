package cmd

import (
	"fjnkt98/atcodersearch/atcodersearch/problem"
	"log"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate document JSON files",
	Long:  "Generate document JSON files",
}

var generateProblemCmd = &cobra.Command{
	Use:   "problem",
	Short: "Generate problem document JSON files",
	Long:  "Generate problem document JSON files",
	Run: func(cmd *cobra.Command, args []string) {
		saveDir, err := GetSaveDir(cmd, "problem")
		if err != nil {
			log.Fatal(err.Error())
		}

		generator := problem.NewProblemDocumentGenerator(GetDB(), saveDir)
		concurrent, err := cmd.Flags().GetInt("concurrent")
		if err != nil {
			log.Fatalf("failed to get value of `concurrent` flag: %s", err.Error())
		}
		if err := generator.Run(1000, concurrent); err != nil {
			log.Fatal(err.Error())
		}
	},
}

// var generateUserCmd = &cobra.Command{
// 	Use:   "user",
// 	Short: "Generate user document JSON files",
// 	Long:  "Generate user document JSON files",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		generator := problem.NewProblemDocumentGenerator(getDB(), getSaveDir(cmd))
// 		if err := generator.Run(); err != nil {
// 			log.Fatal(err.Error())
// 		}
// 	},
// }

func init() {
	generateCmd.PersistentFlags().String("save-dir", "", "Directory path at which generated documents will be saved")
	generateCmd.PersistentFlags().Int("concurrent", 10, "Concurrent number of document generation processes")
	generateCmd.AddCommand(generateProblemCmd)

	rootCmd.AddCommand(generateCmd)
}
