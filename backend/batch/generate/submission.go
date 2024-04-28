package generate

import (
	"context"
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/repository"
	"fmt"
	"time"

	"github.com/goark/errs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
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
	SubmissionID  int                   `json:"submissionId"`
	EpochSecond   int64                 `json:"epochSecond"`
	SubmittedAt   solr.IntoSolrDateTime `json:"submittedAt"`
	SubmissionURL string                `json:"submissionUrl"`
	ProblemID     string                `json:"problemId"`
	ProblemTitle  string                `json:"problemTitle"`
	Color         string                `json:"color"`
	Difficulty    int                   `json:"difficulty"`
	ContestID     string                `json:"contestId"`
	ContestTitle  string                `json:"contestTitle"`
	Category      string                `json:"category"`
	UserID        string                `json:"userId"`
	Language      string                `json:"language"`
	LanguageGroup string                `json:"languageGroup"`
	Point         float64               `json:"point"`
	Length        int                   `json:"length"`
	Result        string                `json:"result"`
	ExecutionTime *int                  `json:"executionTime"`
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

type SubmissionRowReader struct {
	pool     *pgxpool.Pool
	interval int
	all      bool
}

func NewSubmissionRowReader(
	pool *pgxpool.Pool,
	interval int,
	all bool,
) *SubmissionRowReader {
	return &SubmissionRowReader{
		pool:     pool,
		interval: interval,
		all:      all,
	}
}

func (r *SubmissionRowReader) ReadRows(ctx context.Context, tx chan<- *SubmissionRow) error {
	db := bun.NewDB(stdlib.OpenDBFromPool(r.pool), pgdialect.New())
	query := db.NewSelect().
		ColumnExpr("s.id").
		ColumnExpr("s.epoch_second").
		ColumnExpr("s.problem_id").
		ColumnExpr("COALESCE(p.title, '') AS problem_title").
		ColumnExpr("COALESCE(d.difficulty, 0) AS difficulty").
		ColumnExpr("s.contest_id").
		ColumnExpr("c.title AS contest_title").
		ColumnExpr("c.category").
		ColumnExpr("s.user_id").
		ColumnExpr("s.language").
		ColumnExpr("COALESCE(l.?, '') AS language_group", bun.Ident("group")).
		ColumnExpr("s.point").
		ColumnExpr("s.length").
		ColumnExpr("s.result").
		ColumnExpr("s.execution_time").
		TableExpr("submissions AS s").
		Join("LEFT JOIN contests AS c ON s.contest_id = c.contest_id").
		Join("LEFT JOIN problems AS p ON s.problem_id = p.problem_id").
		Join("LEFT JOIN difficulties AS d ON s.problem_id = d.problem_id").
		Join("LEFT JOIN languages AS l ON s.language = l.language").
		Where("s.epoch_second > EXTRACT(EPOCH FROM CURRENT_DATE - CAST(? || ' day' AS INTERVAL))", r.interval)

	if !r.all {
		row, err := repository.New(r.pool).FetchLatestBatchHistory(ctx, "UpdateSubmission")
		if err != nil && !errs.Is(err, pgx.ErrNoRows) {
			return errs.New(
				"failed to get latest update submission history",
				errs.WithCause(err),
			)
		}
		query = query.Where("s.updated_at > ?", row.StartedAt)
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
			err := db.ScanRow(ctx, rows, &row)
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

func GenerateSubmissionDocument(ctx context.Context, reader RowReader[*SubmissionRow], saveDir string, options ...option) error {
	return GenerateDocument(ctx, reader, saveDir, options...)
}
