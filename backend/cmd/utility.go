package cmd

import (
	"database/sql"
	"fjnkt98/atcodersearch/config"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"golang.org/x/exp/slog"
)

type Msg struct{}

// func GetSaveDir(cmd *cobra.Command, domain string) (string, error) {
// 	var saveDir string
// 	if s, err := cmd.Flags().GetString("save-dir"); err == nil && s != "" {
// 		return s, nil
// 	} else if dir := os.Getenv("DOCUMENT_SAVE_DIRECTORY"); dir != "" {
// 		saveDir = dir
// 	} else {
// 		return "", errors.New("couldn't determine document save directory")
// 	}

// 	saveDir = filepath.Join(saveDir, domain)
// 	if _, err := os.Stat(saveDir); err != nil {
// 		if os.IsNotExist(err) {
// 			slog.Info(fmt.Sprintf("The directory `%s` doesn't exists, so attempt to create it.", saveDir))
// 			if err := os.Mkdir(saveDir, os.ModePerm); err != nil {
// 				return "", failure.Translate(err, acs.FileOperationError, failure.Context{"directory": saveDir}, failure.Message("failed to create directory"))
// 			} else {
// 				slog.Info(fmt.Sprintf("The directory `%s` was successfully created.", saveDir))
// 				return saveDir, nil
// 			}
// 		} else {
// 			return "", failure.Translate(err, acs.FileOperationError, failure.Context{"directory": saveDir}, failure.Message("failed to get stat directory"))
// 		}
// 	} else {
// 		return saveDir, nil
// 	}
// }

func GetEngine() *sql.DB {
	dsn := config.Config.DataBaseURL

	engine, err := sql.Open("postgres", dsn)
	if err != nil {
		slog.Error("failed to open database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if err := engine.Ping(); err != nil {
		slog.Error("failed to connect database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	return engine
}

func GetDB(engine *sql.DB) *bun.DB {
	db := bun.NewDB(engine, pgdialect.New())

	return db
}
