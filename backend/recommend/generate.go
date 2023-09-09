package recommend

import (
	"context"
	"fjnkt98/atcodersearch/acs"

	"github.com/jmoiron/sqlx"
	"github.com/morikuni/failure"
	"golang.org/x/exp/slog"
)

type Row struct {
	db   *sqlx.DB
	data Data
}

type Data struct {
	ProblemID   string `db:"problem_id"`
	SolvedCount int    `db:"solved_count"`
}

func ToDocument(ctx context.Context, r Row) (Document, error) {
	return Document{
		ProblemID:   r.data.ProblemID,
		SolvedCount: r.data.SolvedCount,
	}, nil
}

type Document struct {
	ProblemID   string `json:"problem_id"`
	SolvedCount int    `json:"solved_count"`
}

type RowReader struct {
	db *sqlx.DB
}

func (r *RowReader) ReadRows(ctx context.Context, tx chan<- Row) error {
	sql := `
	WITH "solved_counts" AS (
		SELECT
			"problem_id",
			COUNT(1) AS "solved_count"
		FROM
			"submissions"
		WHERE
			"result" = 'AC'
		GROUP BY
			"problem_id"
	)
	SELECT
		"solved_counts"."problem_id",
		"solved_counts"."solved_count"
	FROM
		"solved_counts"
		LEFT JOIN "difficulties" USING("problem_id")
	WHERE
		"difficulty" IS NOT NULL
	`
	rows, err := r.db.QueryxContext(ctx, sql)
	if err != nil {
		return failure.Translate(err, acs.DBError, failure.Context{"sql": sql}, failure.Message("failed to read rows"))
	}
	defer rows.Close()
	defer close(tx)

	for rows.Next() {
		select {
		case <-ctx.Done():
			slog.Info("ReadRows canceled.")
			return failure.New(acs.Interrupt, failure.Message("ReadRows canceled"))
		default:
			var data Data
			err := rows.StructScan(&data)
			if err != nil {
				return failure.Translate(err, acs.DBError, failure.Message("failed to scan row"))
			}

			row := Row{
				db:   r.db,
				data: data,
			}
			tx <- row
		}
	}

	return nil
}

func Generate(ctx context.Context, db *sqlx.DB, saveDir string, chunkSize int, concurrent int) error {
	if err := acs.CleanDocument(saveDir); err != nil {
		return failure.Translate(err, acs.GenerateError, failure.Message("failed to clean recommend save directory"))
	}

	reader := RowReader{db: db}

	if err := acs.GenerateDocument[Row, Document](ctx, saveDir, chunkSize, concurrent, reader.ReadRows, ToDocument); err != nil {
		return failure.Translate(err, acs.GenerateError, failure.Message("failed to generate recommend document"))
	}
	return nil
}
