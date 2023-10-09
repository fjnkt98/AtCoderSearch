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

func ToDocument(ctx context.Context, r Row) (Document, error) {
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

type RowReader struct {
	db       *sqlx.DB
	period   time.Time
	interval int
}

func (r *RowReader) ReadRows(ctx context.Context, tx chan<- Row) error {
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
		"submissions"."epoch_second" > EXTRACT(EPOCH FROM CURRENT_DATE - CAST($1::integer || ' day' AS INTERVAL))
		AND "submissions"."crawled_at" > $2::timestamp with time zone
	`
	rows, err := r.db.QueryxContext(ctx, sql, r.interval, r.period)
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
			var row Row
			err := rows.StructScan(&row)
			if err != nil {
				return failure.Translate(err, acs.DBError, failure.Message("failed to scan row"))
			}
			tx <- row
		}
	}

	return nil
}

func Generate(ctx context.Context, db *sqlx.DB, saveDir string, chunkSize int, concurrent int, period time.Time, interval int) error {
	if err := acs.CleanDocument(saveDir); err != nil {
		return failure.Translate(err, acs.GenerateError, failure.Message("failed to clean submission save directory"))
	}

	reader := RowReader{db: db, period: period, interval: interval}

	if err := acs.GenerateDocument[Row, Document](ctx, saveDir, chunkSize, concurrent, reader.ReadRows, ToDocument); err != nil {
		return failure.Translate(err, acs.GenerateError, failure.Message("failed to generate submission document"))
	}
	return nil
}
