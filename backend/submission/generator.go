package submission

import (
	"context"
	"fjnkt98/atcodersearch/acs"
	"fjnkt98/atcodersearch/atcoder"
	"fjnkt98/atcodersearch/solr"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/morikuni/failure"
)

type Row struct {
	atcoder.Submission
}

func (r *Row) ToDocument() (Document, error) {
	dt := time.Unix(r.EpochSecond, 0)

	return Document{
		SubmissionID:  r.ID,
		SubmittedAt:   solr.IntoSolrDateTime(dt),
		ProblemID:     r.ProblemID,
		ContestID:     r.ContestID,
		UserID:        r.UserID,
		Language:      r.Language,
		Point:         r.Point,
		Length:        r.Length,
		Result:        r.Result,
		ExecutionTime: r.ExecutionTime,
	}, nil
}

type Document struct {
	SubmissionID  int64                 `json:"submission_id"`
	SubmittedAt   solr.IntoSolrDateTime `json:"submitted_at"`
	ProblemID     string                `json:"problem_id"`
	ContestID     string                `json:"contest_id"`
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

func (r *RowReader[R, D]) ReadRows(ctx context.Context, tx chan<- *Row) error {
	sql := `
	SELECT
		"id",
		"epoch_second",
		"problem_id",
		"contest_id",
		"user_id",
		"language",
		"point",
		"length",
		"result",
		"execution_time"
	FROM
		"submissions"
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
			log.Println("ReadRows canceled.")
			return nil
		default:
			var row Row
			err := rows.StructScan(&row)
			if err != nil {
				return failure.Translate(err, DBError, failure.Message("failed to scan row"))
			}
			tx <- &row
		}
	}

	return nil
}

type DocumentGenerator struct {
	saveDir string
	reader  *RowReader[*Row, Document]
}

func NewDocumentGenerator(db *sqlx.DB, saveDir string) DocumentGenerator {
	return DocumentGenerator{
		saveDir: saveDir,
		reader:  &RowReader[*Row, Document]{db: db},
	}
}

func (g *DocumentGenerator) Clean() error {
	if err := acs.CleanDocument(g.saveDir); err != nil {
		return failure.Translate(err, FileOperationError, failure.Context{"directory": g.saveDir}, failure.Message("failed to delete submission document files"))
	}
	return nil
}

func (g *DocumentGenerator) Generate(chunkSize int, concurrent int) error {
	if err := acs.GenerateDocument[*Row, Document](g.reader, g.saveDir, chunkSize, concurrent); err != nil {
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
