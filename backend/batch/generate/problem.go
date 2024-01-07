package generate

import (
	"context"
	"fjnkt98/atcodersearch/batch"
	"fjnkt98/atcodersearch/pkg/solr"
	"fmt"
	"io"
	"strings"
	"time"

	"log/slog"

	"github.com/PuerkitoBio/goquery"
	"github.com/goark/errs"
	"github.com/uptrace/bun"
)

type ProblemGenerator interface {
	batch.Batch
	GenerateProblem(ctx context.Context) error
}

type problemGenerator struct {
	defaultGenerator
}

func NewProblemGenerator(reader RowReader, saveDir string, chunkSize int, concurrent int) ProblemGenerator {
	return &problemGenerator{
		defaultGenerator{
			config: defaultGeneratorConfig{
				SaveDir:    saveDir,
				ChunkSize:  chunkSize,
				Concurrent: concurrent,
			},
			reader: reader,
		},
	}
}

func (g *problemGenerator) Name() string {
	return "ProblemGenerator"
}

func (g *problemGenerator) Run(ctx context.Context) error {
	return g.GenerateProblem(ctx)
}

func (g *problemGenerator) GenerateProblem(ctx context.Context) error {
	if err := g.Generate(ctx); err != nil {
		return errs.Wrap(err)
	}
	return nil
}

type ProblemRow struct {
	ProblemID      string `bun:"problem_id"`
	ProblemTitle   string `bun:"problem_title"`
	ProblemURL     string `bun:"problem_url"`
	ContestID      string `bun:"contest_id"`
	ContestTitle   string `bun:"contest_title"`
	StartAt        int64  `bun:"start_at"`
	Duration       int64  `bun:"duration"`
	RateChange     string `bun:"rate_change"`
	Category       string `bun:"category"`
	HTML           string `bun:"html"`
	Difficulty     *int   `bun:"difficulty"`
	IsExperimental bool   `bun:"is_experimental"`
}

type ProblemDocument struct {
	ProblemID      string                `solr:"problem_id"`
	ProblemTitle   string                `solr:"problem_title"`
	ProblemURL     string                `solr:"problem_url"`
	ContestID      string                `solr:"contest_id"`
	ContestTitle   string                `solr:"contest_title"`
	ContestURL     string                `solr:"contest_url"`
	Color          string                `solr:"color"`
	StartAt        solr.IntoSolrDateTime `solr:"start_at"`
	Duration       int64                 `solr:"duration"`
	RateChange     string                `solr:"rate_change"`
	Category       string                `solr:"category"`
	Difficulty     *int                  `solr:"difficulty"`
	IsExperimental bool                  `solr:"is_experimental"`
	StatementJa    []string              `solr:"statement_ja"`
	StatementEn    []string              `solr:"statement_en"`
}

func (r *ProblemRow) Document(ctx context.Context) (map[string]any, error) {
	statementJa, statementEn, err := ExtractStatements(strings.NewReader(r.HTML))
	if err != nil {
		return nil, errs.New(
			"failed to extract statement from problem",
			errs.WithCause(err),
			errs.WithContext("problem id", r.ProblemID),
		)
	}

	contestURL := fmt.Sprintf("https://atcoder.jp/contests/%s", r.ContestID)
	startAt := solr.IntoSolrDateTime(time.Unix(r.StartAt, 0))

	var color string
	if r.Difficulty == nil {
		color = "black"
	} else {
		color = RateToColor(*r.Difficulty)
	}

	return StructToMap(ProblemDocument{
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
	}), nil
}

func ExtractStatements(html io.Reader) ([]string, []string, error) {
	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return nil, nil, errs.New(
			"failed to extract statements",
			errs.WithCause(err),
		)
	}

	textJa := make([]string, 0)
	textEn := make([]string, 0)

	doc.Find("section").Each(func(_ int, section *goquery.Selection) {
		// For modern contest problem format
		if strings.Contains(section.Find("h3").Text(), "問題") {
			textJa = append(textJa, strings.SplitAfter(section.Text(), "。")...)
		}

		// For legacy contest problem format
		if prev := section.Prev(); goquery.NodeName(prev) == "h3" {
			if text := prev.Text(); strings.Contains(text, "問題") {
				textJa = append(textJa, strings.SplitAfter(section.Text(), "。")...)
			}
		}
	})

	doc.Find("span.lang-en").Find("section").Each(func(_ int, section *goquery.Selection) {
		if strings.Contains(section.NextAll().Find("h3").Text(), "Statement") || strings.Contains(section.Find("h3").Text(), "Statement") {
			textEn = append(textEn, strings.SplitAfter(section.Text(), ".")...)
		}
	})
	return textJa, textEn, nil
}

func RateToColor(rate int) string {
	if rate < 0 {
		return "black"
	} else if rate < 400 {
		return "gray"
	} else if rate < 800 {
		return "brown"
	} else if rate < 1200 {
		return "green"
	} else if rate < 1600 {
		return "cyan"
	} else if rate < 2000 {
		return "blue"
	} else if rate < 2400 {
		return "yellow"
	} else if rate < 2800 {
		return "orange"
	} else if rate < 3200 {
		return "red"
	} else if rate < 3600 {
		return "silver"
	} else {
		return "gold"
	}
}

type problemRowReader struct {
	db *bun.DB
}

func NewProblemRowReader(db *bun.DB) RowReader {
	return &problemRowReader{
		db: db,
	}
}

func (r *problemRowReader) ReadRows(ctx context.Context, tx chan<- Documenter) error {
	rows, err := r.db.NewSelect().
		ColumnExpr("?.? AS ?", bun.Ident("p"), bun.Ident("problem_id"), bun.Ident("problem_id")).
		ColumnExpr("?.? AS ?", bun.Ident("p"), bun.Ident("title"), bun.Ident("problem_title")).
		ColumnExpr("?.? AS ?", bun.Ident("p"), bun.Ident("url"), bun.Ident("problem_url")).
		ColumnExpr("?.? AS ?", bun.Ident("p"), bun.Ident("html"), bun.Ident("html")).
		ColumnExpr("?.? AS ?", bun.Ident("c"), bun.Ident("contest_id"), bun.Ident("contest_id")).
		ColumnExpr("?.? AS ?", bun.Ident("c"), bun.Ident("title"), bun.Ident("contest_title")).
		ColumnExpr("?.? AS ?", bun.Ident("c"), bun.Ident("start_epoch_second"), bun.Ident("start_at")).
		ColumnExpr("?.? AS ?", bun.Ident("c"), bun.Ident("duration_second"), bun.Ident("duration")).
		ColumnExpr("?.? AS ?", bun.Ident("c"), bun.Ident("rate_change"), bun.Ident("rate_change")).
		ColumnExpr("?.? AS ?", bun.Ident("c"), bun.Ident("category"), bun.Ident("category")).
		ColumnExpr("?.? AS ?", bun.Ident("d"), bun.Ident("difficulty"), bun.Ident("difficulty")).
		ColumnExpr("COALESCE(?.?, FALSE) AS ?", bun.Ident("d"), bun.Ident("is_experimental"), bun.Ident("is_experimental")).
		TableExpr("? AS ?", bun.Ident("problems"), bun.Ident("p")).
		Join("JOIN ? AS ? ON ?.? = ?.?", bun.Ident("contests"), bun.Ident("c"), bun.Ident("p"), bun.Ident("contest_id"), bun.Ident("c"), bun.Ident("contest_id")).
		Join("LEFT JOIN ? AS ? ON ?.? = ?.?", bun.Ident("difficulties"), bun.Ident("d"), bun.Ident("p"), bun.Ident("problem_id"), bun.Ident("d"), bun.Ident("problem_id")).
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
			var row ProblemRow
			if err := r.db.ScanRow(ctx, rows, &row); err != nil {
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
