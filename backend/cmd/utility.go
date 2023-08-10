package cmd

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/morikuni/failure"
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
				return "", failure.Translate(err, FileOperationError, failure.Context{"directory": saveDir}, failure.Message("failed to create directory"))
			} else {
				log.Printf("The directory `%s` was successfully created.", saveDir)
				return saveDir, nil
			}
		} else {
			return "", failure.Translate(err, FileOperationError, failure.Context{"directory": saveDir}, failure.Message("failed to get stat directory"))
		}
	} else {
		return saveDir, nil
	}
}

func GetDB() *sqlx.DB {
	dsn := os.Getenv("DATABASE_URL")
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open database: %s", err.Error())
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to connect database: %s", err.Error())
	}

	return db
}

type UpdateHistory struct {
	db        *sqlx.DB
	StartedAt time.Time
	Domain    string
	Options   string
}

func NewUpdateHistory(db *sqlx.DB, domain string, options string) UpdateHistory {
	return UpdateHistory{
		db:        db,
		StartedAt: time.Now(),
		Domain:    domain,
		Options:   options,
	}
}

func (h *UpdateHistory) save(status string) error {
	tx, err := h.db.Beginx()
	if err != nil {
		return failure.Translate(err, DBError, failure.Message("failed to start transaction to save update history"))
	}
	defer tx.Rollback()

	if _, err := tx.Exec(
		`INSERT INTO "update_history" ("domain", "started_at", "finished_at", "status", "options") VALUES ($1::text, $2::timestamp, $3::timestamp, $4::text, $5::json);`,
		h.Domain,
		h.StartedAt,
		time.Now(),
		status,
		h.Options,
	); err != nil {
		return failure.Translate(err, DBError, failure.Message("failed to exec sql to save update history"))
	}

	if err := tx.Commit(); err != nil {
		return failure.Translate(err, DBError, failure.Message("failed to commit transaction to save update history"))
	}

	return nil
}

func (h *UpdateHistory) Finish() error {
	return h.save("finished")
}

func (h *UpdateHistory) Cancel() error {
	return h.save("canceled")
}
