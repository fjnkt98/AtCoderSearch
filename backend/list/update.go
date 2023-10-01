package list

import (
	"context"
	"fjnkt98/atcodersearch/acs"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/morikuni/failure"
	"golang.org/x/exp/slog"
)

func Update(ctx context.Context, db *sqlx.DB) error {
	tx, err := db.Beginx()
	if err != nil {
		return failure.Translate(err, acs.DBError, failure.Message("failed to start transaction"))
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(
		ctx,
		`
		DELETE FROM "languages"
		`,
	); err != nil {
		return failure.Translate(err, acs.DBError, failure.Message("failed to delete languages"))
	}

	affected := 0
	if result, err := tx.ExecContext(
		ctx,
		`
		INSERT INTO "languages" ("language")
		SELECT
			DISTINCT "language"
		FROM
			"submissions"	
		`,
	); err != nil {
		return failure.Translate(err, acs.DBError, failure.Message("failed to insert languages"))
	} else {
		a, _ := result.RowsAffected()
		affected = int(a)
	}

	if err := tx.Commit(); err != nil {
		return failure.Translate(err, acs.DBError, failure.Message("failed to commit transaction to save languages"))
	} else {
		slog.Info(fmt.Sprintf("commit transaction. save %d rows.", affected))
	}

	return nil
}
