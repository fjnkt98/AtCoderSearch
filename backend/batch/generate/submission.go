package generate

import (
	"context"
	"fjnkt98/atcodersearch/batch"
	"fjnkt98/atcodersearch/batch/repository"
	"fjnkt98/atcodersearch/config"
	"fjnkt98/atcodersearch/pkg/solr"
	"fmt"
	"time"

	"github.com/goark/errs"
	"github.com/uptrace/bun"
	"golang.org/x/exp/slog"
)

type SubmissionGenerator interface {
	batch.Batch
	GenerateSubmission(ctx context.Context) error
}

type submissionGenerator struct {
	defaultGenerator
}

func NewSubmissionGenerator(cfg config.GenerateSubmissionConfig, reader RowReader) SubmissionGenerator {
	return &submissionGenerator{
		defaultGenerator{
			cfg: config.GenerateCommonConfig{
				SaveDir:    cfg.SaveDir,
				ChunkSize:  cfg.ChunkSize,
				Concurrent: cfg.Concurrent,
			},
			reader: reader,
		},
	}
}

func (g *submissionGenerator) Name() string {
	return "SubmissionGenerator"
}

func (g *submissionGenerator) Run(ctx context.Context) error {
	return g.GenerateSubmission(ctx)
}

func (g *submissionGenerator) GenerateSubmission(ctx context.Context) error {
	if err := g.Generate(ctx); err != nil {
		return errs.Wrap(err)
	}
	return nil
}

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

type SubmissionDocument struct {
	SubmissionID  int                   `solr:"submission_id"`
	EpochSecond   int64                 `solr:"epoch_second"`
	SubmittedAt   solr.IntoSolrDateTime `solr:"submitted_at"`
	SubmissionURL string                `solr:"submission_url"`
	ProblemID     string                `solr:"problem_id"`
	ProblemTitle  string                `solr:"problem_title"`
	Color         string                `solr:"color"`
	Difficulty    int                   `solr:"difficulty"`
	ContestID     string                `solr:"contest_id"`
	ContestTitle  string                `solr:"contest_title"`
	Category      string                `solr:"category"`
	UserID        string                `solr:"user_id"`
	Language      string                `solr:"language"`
	LanguageGroup string                `solr:"language_group"`
	Point         float64               `solr:"point"`
	Length        int                   `solr:"length"`
	Result        string                `solr:"result"`
	ExecutionTime *int                  `solr:"execution_time"`
}

func (r *SubmissionRow) Document(ctx context.Context) (map[string]any, error) {
	submissionURL := fmt.Sprintf("https://atcoder.jp/contests/%s/submissions/%d", r.ContestID, r.ID)
	color := RateToColor(r.Difficulty)

	return StructToMap(SubmissionDocument{
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
	}), nil
}

type submissionRowReader struct {
	db   *bun.DB
	repo repository.UpdateHistoryRepository
	cfg  config.ReadSubmissionConfig
}

func NewSubmissionRowReader(
	db *bun.DB,
	repo repository.UpdateHistoryRepository,
	cfg config.ReadSubmissionConfig,
) RowReader {
	return &submissionRowReader{
		db:   db,
		repo: repo,
		cfg:  cfg,
	}
}

func (r *submissionRowReader) getPeriod(ctx context.Context) (time.Time, error) {
	if r.cfg.All {
		return time.Time{}, nil
	} else {
		latest, err := r.repo.GetLatest(ctx, "submission")
		if err != nil {
			return time.Time{}, errs.New(
				"failed to get latest update submission history",
				errs.WithCause(err),
			)
		}
		return latest.StartedAt, nil
	}
}

func (r *submissionRowReader) ReadRows(ctx context.Context, tx chan<- Documenter) error {
	period, err := r.getPeriod(ctx)
	if err != nil {
		return errs.Wrap(err)
	}

	rows, err := r.db.NewSelect().
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
		Where("?.? > EXTRACT(EPOCH FROM CURRENT_DATE - CAST(? || ' day' AS INTERVAL))", bun.Ident("s"), bun.Ident("epoch_second"), r.cfg.Interval).
		Where("?.? > ?", bun.Ident("s"), bun.Ident("crawled_at"), period).
		Rows(ctx)

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
			slog.Info("read rows canceled.")
			return batch.ErrInterrupt
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
