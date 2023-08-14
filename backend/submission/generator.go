package submission

import (
	"context"
	"fjnkt98/atcodersearch/acs"
	"fjnkt98/atcodersearch/atcoder"
	"fjnkt98/atcodersearch/solr"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/morikuni/failure"
	"golang.org/x/exp/slog"
)

type Row struct {
	atcoder.Submission
	Category string `db:"category"`
}

func (r Row) ToDocument() (Document, error) {
	dt := time.Unix(r.EpochSecond, 0)
	id := fmt.Sprintf("%s-%d", r.ContestID, r.ID)

	return Document{
		ID:            id,
		SubmissionID:  r.ID,
		EpochSecond:   r.EpochSecond,
		SubmittedAt:   solr.IntoSolrDateTime(dt),
		ProblemID:     r.ProblemID,
		ContestID:     r.ContestID,
		Category:      r.Category,
		UserID:        r.UserID,
		Language:      r.Language,
		Point:         r.Point,
		Length:        r.Length,
		Result:        r.Result,
		ExecutionTime: r.ExecutionTime,
	}, nil
}

type Document struct {
	ID            string                `json:"id"`
	SubmissionID  int64                 `json:"submission_id"`
	EpochSecond   int64                 `json:"epoch_second"`
	SubmittedAt   solr.IntoSolrDateTime `json:"submitted_at"`
	ProblemID     string                `json:"problem_id"`
	ContestID     string                `json:"contest_id"`
	Category      string                `json:"category"`
	UserID        string                `json:"user_id"`
	Language      string                `json:"language"`
	Point         float64               `json:"point"`
	Length        uint64                `json:"length"`
	Result        string                `json:"result"`
	ExecutionTime *uint64               `json:"execution_time"`
}

type RowReader[R acs.ToDocument[D], D any] struct {
	db *sqlx.DB
}

func (r *RowReader[R, D]) ReadRows(ctx context.Context, tx chan<- Row) error {
	sql := `
	SELECT
		"id",
		"epoch_second",
		"problem_id",
		"contest_id",
		"category",
		"user_id",
		"language",
		"point",
		"length",
		"result",
		"execution_time"
	FROM
		"submissions"
		LEFT JOIN "contests" USING("contest_id")
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
