package problem

import (
	"context"
	"fjnkt98/atcodersearch/atcodersearch/common"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

var extractor = NewFullTextExtractor()

type Row struct {
	ProblemID      string
	ProblemTitle   string
	ProblemURL     string
	ContestID      string
	ContestTitle   string
	StartAt        int64
	Duration       int64
	RateChange     string
	Category       string
	HTML           string
	Difficulty     *int
	IsExperimental bool
}

func (r Row) ToDocument() (common.Document, error) {
	statementJa, statementEn, err := extractor.Extract(strings.NewReader(r.HTML))
	if err != nil {
		log.Printf("failed to extract statement at problem `%s`: %s", r.ProblemID, err.Error())
		return ProblemDocument{}, err
	}

	contestURL := fmt.Sprintf("https://atcoder.jp/contests/%s", r.ContestID)
	startAt := time.Unix(r.StartAt, 0)

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
	ProblemID      string
	ProblemTitle   string
	ProblemURL     string
	ContestID      string
	ContestURL     string
	ContestTitle   string
	Color          string
	StartAt        time.Time
	Duration       int64
	RateChange     string
	Category       string
	Difficulty     *int
	IsExperimental bool
	StatementJa    []string
	StatementEn    []string
}

type ProblemRowReader struct {
	db *sqlx.DB
}

func (r *ProblemRowReader) ReadRows(ctx context.Context, tx chan<- common.ToDocument) error {
	rows, err := r.db.Query(`
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
		"difficulties"."is_experimental" AS "is_experimental"
	FROM
		"problems"
		JOIN "contests" ON "problems"."contest_id" = "contests"."contest_id"
		LEFT JOIN "difficulties" ON "problems"."problem_id" = "difficulties"."problem_id"`,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

loop:
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if !rows.Next() {
				break loop
			}
			var row Row
			if err := rows.Scan(&row); err != nil {
				tx <- row
			}
		}
	}

	return nil
}

type ProblemDocumentGenerator struct {
	common.DefaultDocumentGenerator
}

func NewProblemDocumentGenerator(db *sqlx.DB, saveDir string) ProblemDocumentGenerator {
	return ProblemDocumentGenerator{
		common.NewDefaultDocumentGenerator(saveDir, &ProblemRowReader{db: db}),
	}
}

func (g *ProblemDocumentGenerator) Run() error {
	if err := g.Clean(); err != nil {
		log.Printf("failed to clean existing document files: %s", err.Error())
		return err
	}

	if err := g.Generate(1000); err != nil {
		log.Printf("failed to generate documents: %s", err.Error())
		return err
	}
	return nil
}
