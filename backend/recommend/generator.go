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
	ctx  context.Context
	data Data
}

type Data struct {
	ProblemID   string `db:"problem_id"`
	SolvedCount int    `db:"solved_count"`
}

func (r Row) ToDocument() (Document, error) {
	return Document{
		ProblemID:   r.data.ProblemID,
		SolvedCount: r.data.SolvedCount,
	}, nil
}

type Document struct {
	ProblemID   string `json:"problem_id"`
	SolvedCount int    `json:"solved_count"`
}

type RowReader[R acs.ToDocument[D], D any] struct {
	db *sqlx.DB
}

func (r *RowReader[R, D]) ReadRows(ctx context.Context, tx chan<- Row) error {
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
	rows, err := r.db.Queryx(sql)
	if err != nil {
		return failure.Translate(err, acs.DBError, failure.Context{"sql": sql}, failure.Message("failed to read rows"))
	}
	defer rows.Close()
	defer close(tx)

	for rows.Next() {
		select {
		case <-ctx.Done():
			slog.Info("ReadRows canceled.")
			return nil
		default:
			var data Data
			err := rows.StructScan(&data)
			if err != nil {
				return failure.Translate(err, acs.DBError, failure.Message("failed to scan row"))
			}

			row := Row{
				db:   r.db,
				data: data,
				ctx:  ctx,
			}
			tx <- row
		}
	}

	return nil
}

type DocumentGenerator struct {
	saveDir string
	reader  *RowReader[Row, Document]
}

func NewDocumentGenerator(db *sqlx.DB, saveDir string) DocumentGenerator {
	return DocumentGenerator{
		saveDir: saveDir,
		reader:  &RowReader[Row, Document]{db: db},
	}
}

func (g *DocumentGenerator) Clean() error {
	if err := acs.CleanDocument(g.saveDir); err != nil {
		return failure.Translate(err, acs.FileOperationError, failure.Context{"directory": g.saveDir}, failure.Message("failed to delete problem document files in `%s`"))
	}
	return nil
}

func (g *DocumentGenerator) Generate(chunkSize int, concurrent int) error {
	if err := acs.GenerateDocument[Row, Document](g.reader, g.saveDir, chunkSize, concurrent); err != nil {
		return failure.Wrap(err)
	}
	return nil
}

func (g *DocumentGenerator) Run(chunkSize int, concurrent int) error {
	if err := g.Clean(); err != nil {
		return failure.Wrap(err)
	}

	if err := g.Generate(chunkSize, concurrent); err != nil {
		return failure.Wrap(err)
	}
	return nil
}
