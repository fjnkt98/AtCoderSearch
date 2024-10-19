package update

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/meilisearch/meilisearch-go"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

type ProblemIndexer struct {
	client meilisearch.ServiceManager
}

func NewProblemIndexer(client meilisearch.ServiceManager) *ProblemIndexer {
	return &ProblemIndexer{
		client: client,
	}
}

func (ix *ProblemIndexer) Manager() meilisearch.ServiceManager {
	return ix.client
}

func (ix *ProblemIndexer) Settings() *meilisearch.Settings {
	return &meilisearch.Settings{
		Dictionary: []string{
			"ABC",
			"AGC",
			"AHC",
			"ARC",
		},
		DisplayedAttributes: []string{
			"*",
		},
		FilterableAttributes: []string{
			"problemId",
			"contestId",
			"color",
			"difficulty",
			"difficultyFacet",
			"isExperimental",
			"category",
		},
		SearchableAttributes: []string{
			"problemTitle",
			"contestTitle",
			"color",
			"category",
			"statementJa",
			"statementEn",
		},
		SortableAttributes: []string{
			"startAt",
			"difficulty",
		},
		Synonyms: map[string][]string{},
	}
}

func (ix *ProblemIndexer) PrimaryKey() string {
	return "problemId"
}

func (ix *ProblemIndexer) IndexName() string {
	return "problems"
}

type ProblemRowReader struct {
	pool *pgxpool.Pool
}

func NewProblemRowReader(pool *pgxpool.Pool) *ProblemRowReader {
	return &ProblemRowReader{
		pool: pool,
	}
}

func (r *ProblemRowReader) ReadRows(ctx context.Context, tx chan<- ProblemRow) error {
	db := bun.NewDB(stdlib.OpenDBFromPool(r.pool), pgdialect.New())
	rows, err := db.NewSelect().
		ColumnExpr("p.problem_id AS problem_id").
		ColumnExpr("p.title AS problem_title").
		ColumnExpr("p.url AS problem_url").
		ColumnExpr("p.html AS html").
		ColumnExpr("p.contest_id AS contest_id").
		ColumnExpr("c.title AS contest_title").
		ColumnExpr("c.start_epoch_second AS start_at").
		ColumnExpr("c.duration_second AS duration").
		ColumnExpr("c.rate_change AS rate_change").
		ColumnExpr("c.category AS category").
		ColumnExpr("d.difficulty AS difficulty").
		ColumnExpr("COALESCE(d.is_experimental, FALSE) AS is_experimental").
		TableExpr("problems AS p").
		Join("JOIN contests AS c ON p.contest_id = c.contest_id").
		Join("LEFT JOIN difficulties AS d ON p.problem_id = d.problem_id").
		Rows(ctx)

	if err != nil {
		return fmt.Errorf("execute sql: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			var row ProblemRow
			if err := db.ScanRow(ctx, rows, &row); err != nil {
				return fmt.Errorf("scan row: %w", err)
			}
			tx <- row
		}
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

func (r ProblemRow) Document(ctx context.Context) (ProblemDocument, error) {
	statementJa, statementEn, err := ExtractStatements(strings.NewReader(r.HTML))
	if err != nil {
		return ProblemDocument{}, fmt.Errorf("extract statement of the problem `%s`: %w", r.ProblemID, err)
	}

	return ProblemDocument{
		ProblemID:       r.ProblemID,
		ProblemURL:      r.ProblemURL,
		ProblemTitle:    r.ProblemTitle,
		ContestID:       r.ContestID,
		ContestURL:      fmt.Sprintf("https://atcoder.jp/contests/%s", r.ContestID),
		ContestTitle:    r.ContestTitle,
		StartAt:         r.StartAt,
		Duration:        r.Duration,
		RateChange:      r.RateChange,
		Category:        r.Category,
		Difficulty:      r.Difficulty,
		DifficultyFacet: RateToRangeLabel(r.Difficulty),
		IsExperimental:  r.IsExperimental,
		StatementJa:     statementJa,
		StatementEn:     statementEn,
	}, nil
}

type ProblemDocument struct {
	ProblemID       string   `json:"problemId"`
	ProblemTitle    string   `json:"problemTitle"`
	ProblemURL      string   `json:"problemUrl"`
	ContestID       string   `json:"contestId"`
	ContestTitle    string   `json:"contestTitle"`
	ContestURL      string   `json:"contestUrl"`
	StartAt         int64    `json:"startAt"`
	Duration        int64    `json:"duration"`
	RateChange      string   `json:"rateChange"`
	Category        string   `json:"category"`
	Difficulty      *int     `json:"difficulty"`
	DifficultyFacet string   `json:"difficultyFacet,omitempty"`
	IsExperimental  bool     `json:"isExperimental"`
	StatementJa     []string `json:"statementJa"`
	StatementEn     []string `json:"statementEn"`
}

func ExtractStatements(html io.Reader) ([]string, []string, error) {
	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return nil, nil, fmt.Errorf("new document from reader: %w", err)
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

func RateToRangeLabel(rate *int) string {
	if rate == nil || *rate < 0 {
		return "     ~    0"
	} else if *rate < 400 {
		return "   0 ~  400"
	} else if *rate < 800 {
		return " 400 ~  800"
	} else if *rate < 1200 {
		return " 800 ~ 1200"
	} else if *rate < 1600 {
		return "1200 ~ 1600"
	} else if *rate < 2000 {
		return "1600 ~ 2000"
	} else if *rate < 2400 {
		return "2000 ~ 2400"
	} else if *rate < 2800 {
		return "2400 ~ 2800"
	} else if *rate < 3200 {
		return "2800 ~ 3200"
	} else if *rate < 3600 {
		return "3200 ~ 3600"
	} else {
		return "3600 ~     "
	}
}
