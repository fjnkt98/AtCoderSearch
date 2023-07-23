package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

func GetSaveDir(cmd *cobra.Command, domain string) (string, error) {
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

func GetDB() *sqlx.DB {
	dsn := os.Getenv("DATABASE_URL")
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to connect database: %w", err)
	}

	return db
}
