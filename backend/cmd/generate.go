package cmd

import (
	"fjnkt98/atcodersearch/atcodersearch/problem"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

func getSaveDir(cmd *cobra.Command) string {
	var saveDir string
	if s, err := cmd.Flags().GetString("save-dir"); err != nil {
		saveDir = s
	} else if dir := os.Getenv("DOCUMENT_SAVE_DIRECTORY"); dir != "" {
		saveDir = dir
	} else {
		log.Fatal("couldn't determine document save directory")
	}

	return saveDir
}

func getDB() *sqlx.DB {
	dsn := os.Getenv("DATABASE_URL")
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open database: %s", err.Error())
	}
	return db
}

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
		generator := problem.NewProblemDocumentGenerator(getDB(), getSaveDir(cmd))
		if err := generator.Run(1000, 10); err != nil {
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
	generateCmd.AddCommand(generateProblemCmd)
	rootCmd.AddCommand(generateCmd)
}
