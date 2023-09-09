package cmd

import (
	"errors"
	"fjnkt98/atcodersearch/acs"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	"github.com/morikuni/failure"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

type Msg struct{}

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
			slog.Info(fmt.Sprintf("The directory `%s` doesn't exists, so attempt to create it.", saveDir))
			if err := os.Mkdir(saveDir, os.ModePerm); err != nil {
				return "", failure.Translate(err, acs.FileOperationError, failure.Context{"directory": saveDir}, failure.Message("failed to create directory"))
			} else {
				slog.Info(fmt.Sprintf("The directory `%s` was successfully created.", saveDir))
				return saveDir, nil
			}
		} else {
			return "", failure.Translate(err, acs.FileOperationError, failure.Context{"directory": saveDir}, failure.Message("failed to get stat directory"))
		}
	} else {
		return saveDir, nil
	}
}

func GetDB() *sqlx.DB {
	dsn := os.Getenv("DATABASE_URL")
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		slog.Error("failed to open database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if err := db.Ping(); err != nil {
		slog.Error("failed to connect database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	return db
}

func GetBool(cmd *cobra.Command, flag string) bool {
	v, err := cmd.Flags().GetBool(flag)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to get value of `%s` flag", flag))
		os.Exit(1)
	}
	return v
}

func GetInt(cmd *cobra.Command, flag string) int {
	v, err := cmd.Flags().GetInt(flag)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to get value of `%s` flag", flag))
		os.Exit(1)
	}
	return v
}

func GetString(cmd *cobra.Command, flag string) string {
	v, err := cmd.Flags().GetString(flag)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to get value of `%s` flag", flag))
		os.Exit(1)
	}
	return v
}
