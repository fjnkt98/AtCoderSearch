package problem

import (
	"context"
	"fjnkt98/atcodersearch/atcodersearch/common"
	"fjnkt98/atcodersearch/solr"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
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

func (r *Row) ToDocument() (any, error) {
	statementJa, statementEn, err := extractor.Extract(strings.NewReader(r.HTML))
	if err != nil {
		log.Printf("failed to extract statement at problem `%s`: %s", r.ProblemID, err.Error())
		return ProblemDocument{}, err
	}

	contestURL := fmt.Sprintf("https://atcoder.jp/contests/%s", r.ContestID)
	startAt := solr.IntoSolrDateTime(time.Unix(r.StartAt, 0))

	var color string
	if r.Difficulty == nil {
		color = "black"
	} else {
		color = common.RateToColor(*r.Difficulty)
	}

	return ProblemDocument{
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

type ProblemDocument struct {
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

type ProblemRowReader struct {
	db *sqlx.DB
}

func (r *ProblemRowReader) ReadRows(ctx context.Context, tx chan<- common.ToDocument) error {
	rows, err := r.db.Queryx(`
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
		LEFT JOIN "difficulties" ON "problems"."problem_id" = "difficulties"."problem_id"`,
	)
	if err != nil {
		return err
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
				log.Printf("failed to scan row: %s", err.Error())
				return err
			}
			tx <- &row
		}
	}

	return nil
}

type ProblemDocumentGenerator struct {
	saveDir string
	reader  *ProblemRowReader
}

func NewProblemDocumentGenerator(db *sqlx.DB, saveDir string) ProblemDocumentGenerator {
	return ProblemDocumentGenerator{
		saveDir: saveDir,
		reader:  &ProblemRowReader{db: db},
	}
}

func (g *ProblemDocumentGenerator) Clean() error {
	if err := common.CleanDocument(g.saveDir); err != nil {
		return fmt.Errorf("failed to delete problem document files in `%s`: %s", g.saveDir, err.Error())
	}
	return nil
}

func (g *ProblemDocumentGenerator) Generate(chunkSize int, concurrent int) error {
	if err := common.GenerateDocument(g.reader, g.saveDir, chunkSize, concurrent); err != nil {
		return fmt.Errorf("failed to generate problem document files: %s", err.Error())
	}
	return nil
}

func (g *ProblemDocumentGenerator) Run(chunkSize int, concurrent int) error {
	if err := g.Clean(); err != nil {
		return err
	}

	if err := g.Generate(chunkSize, concurrent); err != nil {
		return err
	}
	return nil
}
