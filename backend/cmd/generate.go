package cmd

import (
	"errors"
	"fjnkt98/atcodersearch/atcodersearch/problem"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

func getSaveDir(cmd *cobra.Command, domain string) (string, error) {
	var saveDir string
	if s, err := cmd.Flags().GetString("save-dir"); err == nil && s != "" {
		return s, nil
	} else if dir := os.Getenv("DOCUMENT_SAVE_DIRECTORY"); dir != "" {
		saveDir = dir
	} else {
		return "", errors.New("couldn't determine document save directory")
	}

	saveDir = filepath.Join(saveDir, domain)
	if _, err := os.Stat(saveDir); err != nil {
		if os.IsNotExist(err) {
			log.Printf("The directory `%s` doesn't exists, so attempt to create it.", saveDir)
			if err := os.Mkdir(saveDir, os.ModePerm); err != nil {
				return "", fmt.Errorf("failed to create directory `%s`: %s", saveDir, err.Error())
			} else {
				log.Printf("The directory `%s` was successfully created.", saveDir)
				return saveDir, nil
			}
		} else {
			return "", fmt.Errorf("failed to get stat directory `%s`: %s", saveDir, err.Error())
		}
	} else {
		return saveDir, nil
	}
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
		saveDir, err := getSaveDir(cmd, "problem")
		if err != nil {
			log.Fatal(err.Error())
		}

		generator := problem.NewProblemDocumentGenerator(getDB(), saveDir)
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
