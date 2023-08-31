package solved

import (
	"context"
	"fjnkt98/atcodersearch/acs"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/morikuni/failure"
	"golang.org/x/exp/slog"
)

type Row struct {
	ProblemID string `db:"problem_id"`
	UserID    string `db:"user_id"`
}

func (r Row) ToDocument() (Document, error) {
	return Document{
		UniqueKey: fmt.Sprintf("%s-%s", r.ProblemID, r.UserID),
		ProblemID: r.ProblemID,
		UserID:    r.UserID,
	}, nil
}

type Document struct {
	UniqueKey string `json:"unique_key"`
	ProblemID string `json:"problem_id"`
	UserID    string `json:"user_id"`
}

type RowReader[R acs.ToDocument[D], D any] struct {
	db *sqlx.DB
}

func (r *RowReader[R, D]) ReadRows(ctx context.Context, tx chan<- Row) error {
	sql := `
	SELECT
		DISTINCT
		"submissions"."problem_id",
		"submissions"."user_id"
	FROM
		"submissions"
	WHERE
		"epoch_second" > EXTRACT(EPOCH FROM CURRENT_DATE - INTERVAL '30 day')
		AND "result" = 'AC'
	`
	rows, err := r.db.Queryx(sql)
	if err != nil {
		return failure.Translate(err, DBError, failure.Context{"sql": sql}, failure.Message("failed to read rows"))
	}
	defer rows.Close()
	defer close(tx)

	for rows.Next() {
		select {
		case <-ctx.Done():
			slog.Info("ReadRows canceled.")
			return nil
		default:
			var row Row
			err := rows.StructScan(&row)
			if err != nil {
				return failure.Translate(err, DBError, failure.Message("failed to scan row"))
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
		return failure.Translate(err, FileOperationError, failure.Context{"directory": g.saveDir}, failure.Message("failed to delete submission document files"))
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
