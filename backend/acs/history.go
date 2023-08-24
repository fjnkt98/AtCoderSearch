package acs

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/morikuni/failure"
	"golang.org/x/exp/slog"
)

type UpdateHistory struct {
	db        *sqlx.DB
	StartedAt time.Time
	Domain    string
	Options   string
	WasSaved  bool
}

func NewUpdateHistory(db *sqlx.DB, domain string, options string) UpdateHistory {
	return UpdateHistory{
		db:        db,
		StartedAt: time.Now(),
		Domain:    domain,
		Options:   options,
		WasSaved:  false,
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
	h.WasSaved = true
	return h.save("finished")
}

func (h *UpdateHistory) Cancel() error {
	if h.WasSaved {
		return nil
	}
	return h.save("canceled")
}

func (h *UpdateHistory) GetLatest() (time.Time, error) {
	rows, err := h.db.Query(
		`SELECT
			"started_at"
		FROM
			"update_history"
		WHERE
			"domain" = $1::text
		ORDER
			BY "started_at" DESC
		LIMIT
			1;`,
		h.Domain,
	)
	if err != nil {
		return time.Time{}, failure.Translate(err, DBError, failure.Message("failed to get latest update history"))
	}

	defer rows.Close()
	var startedAt time.Time
	for rows.Next() {
		if err := rows.Scan(&startedAt); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				slog.Info(fmt.Sprintf("There is no history for domain `%s`", h.Domain))
				return time.Time{}, nil
			} else {
				return time.Time{}, failure.Translate(err, DBError, failure.Message("failed to get latest crawl history"))
			}
		}
	}

	return startedAt, nil
}
