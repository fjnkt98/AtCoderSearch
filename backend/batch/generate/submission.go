package generate

import (
	"context"
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/repository"
	"fmt"
	"time"

	"github.com/goark/errs"
	"github.com/uptrace/bun"
)

type SubmissionRow struct {
	ID            int     `bun:"id"`
	EpochSecond   int64   `bun:"epoch_second"`
	ProblemID     string  `bun:"problem_id"`
	ContestID     string  `bun:"contest_id"`
	UserID        string  `bun:"user_id"`
	Language      string  `bun:"language"`
	Point         float64 `bun:"point"`
	Length        int     `bun:"length"`
	Result        string  `bun:"result"`
	ExecutionTime *int    `bun:"execution_time"`
	LanguageGroup string  `bun:"language_group"`
	Category      string  `bun:"category"`
	ProblemTitle  string  `bun:"problem_title"`
	ContestTitle  string  `bun:"contest_title"`
	Difficulty    int     `bun:"difficulty"`
}

type SubmissionDoc struct {
	SubmissionID  int                   `json:"submission_id"`
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
	LanguageGroup string                `json:"language_group"`
	Point         float64               `json:"point"`
	Length        int                   `json:"length"`
	Result        string                `json:"result"`
	ExecutionTime *int                  `json:"execution_time"`
}

func (r *SubmissionRow) Document(ctx context.Context) (*SubmissionDoc, error) {
	submissionURL := fmt.Sprintf("https://atcoder.jp/contests/%s/submissions/%d", r.ContestID, r.ID)
	color := RateToColor(r.Difficulty)

	return &SubmissionDoc{
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
		LanguageGroup: r.LanguageGroup,
		Point:         r.Point,
		Length:        r.Length,
		Result:        r.Result,
		ExecutionTime: r.ExecutionTime,
	}, nil
}

type submissionRowReader struct {
	db       *bun.DB
	repo     repository.UpdateHistoryRepository
	interval int
	all      bool
}

func NewSubmissionRowReader(
	db *bun.DB,
	interval int,
	all bool,
) RowReader[*SubmissionRow] {
	return &submissionRowReader{
		db:       db,
		repo:     repository.NewUpdateHistoryRepository(db),
		interval: interval,
		all:      all,
	}
}

func (r *submissionRowReader) ReadRows(ctx context.Context, tx chan<- *SubmissionRow) error {
	query := r.db.NewSelect().
		ColumnExpr("?.? AS ?", bun.Ident("s"), bun.Ident("id"), bun.Ident("id")).
		ColumnExpr("?.? AS ?", bun.Ident("s"), bun.Ident("epoch_second"), bun.Ident("epoch_second")).
		ColumnExpr("?.? AS ?", bun.Ident("s"), bun.Ident("problem_id"), bun.Ident("problem_id")).
		ColumnExpr("COALESCE(?.?, '') AS ?", bun.Ident("p"), bun.Ident("title"), bun.Ident("problem_title")).
		ColumnExpr("COALESCE(?.?, 0) AS ?", bun.Ident("d"), bun.Ident("difficulty"), bun.Ident("difficulty")).
		ColumnExpr("?.? AS ?", bun.Ident("s"), bun.Ident("contest_id"), bun.Ident("contest_id")).
		ColumnExpr("?.? AS ?", bun.Ident("c"), bun.Ident("title"), bun.Ident("contest_title")).
		ColumnExpr("?.? AS ?", bun.Ident("c"), bun.Ident("category"), bun.Ident("category")).
		ColumnExpr("?.? AS ?", bun.Ident("s"), bun.Ident("user_id"), bun.Ident("user_id")).
		ColumnExpr("?.? AS ?", bun.Ident("s"), bun.Ident("language"), bun.Ident("language")).
		ColumnExpr("COALESCE(?.?, '') AS ?", bun.Ident("l"), bun.Ident("group"), bun.Ident("language_group")).
		ColumnExpr("?.? AS ?", bun.Ident("s"), bun.Ident("point"), bun.Ident("point")).
		ColumnExpr("?.? AS ?", bun.Ident("s"), bun.Ident("length"), bun.Ident("length")).
		ColumnExpr("?.? AS ?", bun.Ident("s"), bun.Ident("result"), bun.Ident("result")).
		ColumnExpr("?.? AS ?", bun.Ident("s"), bun.Ident("execution_time"), bun.Ident("execution_time")).
		TableExpr("? AS ?", bun.Ident("submissions"), bun.Ident("s")).
		Join("LEFT JOIN ? AS ? ON ?.? = ?.?", bun.Ident("contests"), bun.Ident("c"), bun.Ident("s"), bun.Ident("contest_id"), bun.Ident("c"), bun.Ident("contest_id")).
		Join("LEFT JOIN ? AS ? ON ?.? = ?.?", bun.Ident("problems"), bun.Ident("p"), bun.Ident("s"), bun.Ident("problem_id"), bun.Ident("p"), bun.Ident("problem_id")).
		Join("LEFT JOIN ? AS ? ON ?.? = ?.?", bun.Ident("difficulties"), bun.Ident("d"), bun.Ident("s"), bun.Ident("problem_id"), bun.Ident("d"), bun.Ident("problem_id")).
		Join("LEFT JOIN ? AS ? ON ?.? = ?.?", bun.Ident("languages"), bun.Ident("l"), bun.Ident("s"), bun.Ident("language"), bun.Ident("l"), bun.Ident("language")).
		Where("?.? > EXTRACT(EPOCH FROM CURRENT_DATE - CAST(? || ' day' AS INTERVAL))", bun.Ident("s"), bun.Ident("epoch_second"), r.interval)

	if !r.all {
		latest, err := r.repo.GetLatest(ctx, "submission")
		if err != nil {
			return errs.New(
				"failed to get latest update submission history",
				errs.WithCause(err),
			)
		}

		query = query.Where("?.? > ?", bun.Ident("s"), bun.Ident("crawled_at"), latest.StartedAt)
	}

	rows, err := query.Rows(ctx)
	if err != nil {
		return errs.New(
			"failed to read rows",
			errs.WithCause(err),
		)
	}
	defer rows.Close()
	defer close(tx)

	for rows.Next() {
		select {
		case <-ctx.Done():
			return nil
		default:
			var row SubmissionRow
			err := r.db.ScanRow(ctx, rows, &row)
			if err != nil {
				return errs.New(
					"failed to scan row",
					errs.WithCause(err),
				)
			}
			tx <- &row
		}
	}

	return nil
}

func NewSubmissionGenerator(reader RowReader[*SubmissionRow], saveDir string, chunkSize, concurrent int) DocumentGenerator {
	return NewDocumentGenerator(reader, saveDir, chunkSize, concurrent)
}
