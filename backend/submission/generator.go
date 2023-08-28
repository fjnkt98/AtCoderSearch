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
	Category     string `db:"category"`
	ProblemTitle string `db:"problem_title"`
	ContestTitle string `db:"contest_title"`
	Difficulty   int    `db:"difficulty"`
}

func (r Row) ToDocument() (Document, error) {
	submissionURL := fmt.Sprintf("https://atcoder.jp/contests/%s/submissions/%d", r.ContestID, r.ID)
	color := acs.RateToColor(r.Difficulty)

	return Document{
		SubmissionID:  r.ID,
		EpochSecond:   r.EpochSecond,
		SubmittedAt:   solr.IntoSolrDateTime(time.Unix(r.EpochSecond, 0)),
		SubmissionURL: submissionURL,
		ProblemID:     r.ProblemID,
		ProblemTitle:  r.ProblemTitle,
		Color:         color,
		Difficulty:    r.Difficulty,
		ContestID:     r.ContestID,
		ContestTitle:  r.ContestTitle,
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
	SubmissionID  int64                 `json:"submission_id"`
	EpochSecond   int64                 `json:"epoch_second"`
	SubmittedAt   solr.IntoSolrDateTime `json:"submitted_at"`
	SubmissionURL string                `json:"submission_url"`
	ProblemID     string                `json:"problem_id"`
	ProblemTitle  string                `json:"problem_title"`
	Color         string                `json:"color"`
	Difficulty    int                   `json:"difficulty"`
	ContestID     string                `json:"contest_id"`
	ContestTitle  string                `json:"contest_title"`
	Category      string                `json:"category"`
	UserID        string                `json:"user_id"`
	Language      string                `json:"language"`
	Point         float64               `json:"point"`
	Length        int64                 `json:"length"`
	Result        string                `json:"result"`
	ExecutionTime *int64                `json:"execution_time"`
}

type RowReader[R acs.ToDocument[D], D any] struct {
	db     *sqlx.DB
	period time.Time
}

func (r *RowReader[R, D]) ReadRows(ctx context.Context, tx chan<- Row) error {
	sql := `
	SELECT
		"submissions"."id",
		"submissions"."epoch_second",
		"submissions"."problem_id",
		COALESCE("problems"."title", '') AS "problem_title",
		COALESCE("difficulties"."difficulty", 0) AS "difficulty",
		"submissions"."contest_id",
		"contests"."title" AS "contest_title",
		"contests"."category",
		"submissions"."user_id",
		"submissions"."language",
		"submissions"."point",
		"submissions"."length",
		"submissions"."result",
		"submissions"."execution_time"
	FROM
		"submissions"
		LEFT JOIN "contests" USING("contest_id")
		LEFT JOIN "problems" ON "submissions"."problem_id" = "problems"."problem_id"
		LEFT JOIN "difficulties" ON "submissions"."problem_id" = "difficulties"."problem_id"
	WHERE
		"submissions"."crawled_at" > $1::timestamp with time zone
	LIMIT
		100000
	`
	rows, err := r.db.Queryx(sql, r.period)
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

func NewDocumentGenerator(db *sqlx.DB, saveDir string, period time.Time) DocumentGenerator {
	return DocumentGenerator{
		saveDir: saveDir,
		reader:  &RowReader[Row, Document]{db: db, period: period},
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
