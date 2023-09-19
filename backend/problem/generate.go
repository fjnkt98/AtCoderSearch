package problem

import (
	"context"
	"fjnkt98/atcodersearch/acs"
	"fjnkt98/atcodersearch/solr"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/morikuni/failure"
	"golang.org/x/exp/slog"
)

var extractor = NewFullTextExtractor()

type Row struct {
	ProblemID      string `db:"problem_id"`
	ProblemTitle   string `db:"problem_title"`
	ProblemURL     string `db:"problem_url"`
	ContestID      string `db:"contest_id"`
	ContestTitle   string `db:"contest_title"`
	StartAt        int64  `db:"start_at"`
	Duration       int64  `db:"duration"`
	RateChange     string `db:"rate_change"`
	Category       string `db:"category"`
	HTML           string `db:"html"`
	Difficulty     *int   `db:"difficulty"`
	IsExperimental bool   `db:"is_experimental"`
}

func ToDocument(ctx context.Context, r Row) (Document, error) {
	statementJa, statementEn, err := extractor.Extract(strings.NewReader(r.HTML))
	if err != nil {
		return Document{}, failure.Translate(err, acs.ExtractError, failure.Context{"problemID": r.ProblemID}, failure.Message("failed to extract statement at problem `%s`"))
	}

	contestURL := fmt.Sprintf("https://atcoder.jp/contests/%s", r.ContestID)
	startAt := solr.IntoSolrDateTime(time.Unix(r.StartAt, 0))

	var color string
	if r.Difficulty == nil {
		color = "black"
	} else {
		color = acs.RateToColor(*r.Difficulty)
	}

	return Document{
		ProblemID:      r.ProblemID,
		ProblemURL:     r.ProblemURL,
		ProblemTitle:   r.ProblemTitle,
		ContestID:      r.ContestID,
		ContestURL:     contestURL,
		ContestTitle:   r.ContestTitle,
		Color:          color,
		StartAt:        startAt,
		Duration:       r.Duration,
		RateChange:     r.RateChange,
		Category:       r.Category,
		Difficulty:     r.Difficulty,
		IsExperimental: r.IsExperimental,
		StatementJa:    statementJa,
		StatementEn:    statementEn,
	}, nil
}

type Document struct {
	ProblemID      string                `json:"problem_id"`
	ProblemTitle   string                `json:"problem_title"`
	ProblemURL     string                `json:"problem_url"`
	ContestID      string                `json:"contest_id"`
	ContestTitle   string                `json:"contest_title"`
	ContestURL     string                `json:"contest_url"`
	Color          string                `json:"color"`
	StartAt        solr.IntoSolrDateTime `json:"start_at"`
	Duration       int64                 `json:"duration"`
	RateChange     string                `json:"rate_change"`
	Category       string                `json:"category"`
	Difficulty     *int                  `json:"difficulty"`
	IsExperimental bool                  `json:"is_experimental"`
	StatementJa    []string              `json:"statement_ja"`
	StatementEn    []string              `json:"statement_en"`
}

type RowReader struct {
	db *sqlx.DB
}

func (r *RowReader) ReadRows(ctx context.Context, tx chan<- Row) error {
	sql := `
	SELECT
		"problems"."problem_id" AS "problem_id",
		"problems"."title" AS "problem_title",
		"problems"."url" AS "problem_url",
		"contests"."contest_id" AS "contest_id",
		"contests"."title" AS "contest_title",
		"contests"."start_epoch_second" AS "start_at",
		"contests"."duration_second" AS "duration",
		"contests"."rate_change" AS "rate_change",
		"contests"."category" AS "category",
		"problems"."html" AS "html",
		"difficulties"."difficulty" AS "difficulty",
		COALESCE("difficulties"."is_experimental", FALSE) AS "is_experimental"
	FROM
		"problems"
		JOIN "contests" ON "problems"."contest_id" = "contests"."contest_id"
		LEFT JOIN "difficulties" ON "problems"."problem_id" = "difficulties"."problem_id"
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

func Generate(ctx context.Context, db *sqlx.DB, saveDir string, chunkSize int, concurrent int) error {
	if err := acs.CleanDocument(saveDir); err != nil {
		return failure.Translate(err, acs.GenerateError, failure.Message("failed to clean problem save directory"))
	}

	reader := RowReader{db: db}

	if err := acs.GenerateDocument[Row, Document](ctx, saveDir, chunkSize, concurrent, reader.ReadRows, ToDocument); err != nil {
		return failure.Translate(err, acs.GenerateError, failure.Message("failed to generate problem document"))
	}
	return nil
}