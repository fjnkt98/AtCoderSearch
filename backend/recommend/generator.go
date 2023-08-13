package recommend

import (
	"context"
	"fjnkt98/atcodersearch/acs"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/morikuni/failure"
)

type Row struct {
	db   *sqlx.DB
	ctx  context.Context
	data Data
}

type Data struct {
	ProblemID      string `db:"problem_id"`
	Category       string `db:"category"`
	Difficulty     *int   `db:"difficulty"`
	IsExperimental bool   `db:"is_experimental"`
	SolvedCount    int    `db:"solved_count"`
}

type Correlation struct {
	ProblemID   string  `db:"problem_id"`
	Correlation float64 `db:"correlation"`
	Weight      float64 `db:"weight"`
}

func (r *Row) correlations() (string, string, error) {
	sql := `
	WITH "difficulty_correlations" AS (
		SELECT
			"problem_id",
			"contest_id",
			CAST (
				ROUND(
					EXP(
						- POW(($1::integer - "difficulty"), 2.0) / 57707.8
					),
					6
				) AS DOUBLE PRECISION
			) AS "correlation"
		FROM
			"problems"
			LEFT JOIN "difficulties" USING("problem_id")
		WHERE
			"problems"."problem_id" <> $2::text
			AND "difficulty" IS NOT NULL
		ORDER BY
			"correlation" DESC
		LIMIT
			100
	)
	SELECT
		"problem_id",
		"correlation",
		COALESCE("weight", 0.0) AS "weight"
	FROM
		"difficulty_correlations"
	LEFT JOIN "contests" USING("contest_id")
	LEFT JOIN (SELECT "to", "weight" FROM "category_relationships" WHERE "from" = $3::text) AS "relations" ON "contests"."category" = "relations"."to"
	`
	rows, err := r.db.Queryx(
		sql,
		r.data.Difficulty,
		r.data.ProblemID,
		r.data.Category,
	)

	if err != nil {
		return "", "", failure.Translate(err, DBError, failure.Context{"sql": sql}, failure.Message("failed to read rows"))
	}
	defer rows.Close()

	diff := make([]string, 0)
	cate := make([]string, 0)
	for rows.Next() {
		select {
		case <-r.ctx.Done():
			log.Println("ToDocument canceled.")
			return "", "", nil
		default:
			var c Correlation
			if err := rows.StructScan(&c); err != nil {
				return "", "", failure.Translate(err, DBError, failure.Message("failed to scan row"))
			}
			if c.Correlation != 0.0 {
				diff = append(diff, fmt.Sprintf("%s|%.6f", c.ProblemID, c.Correlation))
				cate = append(cate, fmt.Sprintf("%s|%.6f", c.ProblemID, c.Weight))
			}
		}
	}

	return strings.Join(diff, " "), strings.Join(cate, " "), nil
}

func (r Row) ToDocument() (Document, error) {
	d, c, err := r.correlations()
	if err != nil {
		return Document{}, failure.Wrap(err)
	}

	return Document{
		ProblemID:             r.data.ProblemID,
		DifficultyCorrelation: d,
		CategoryCorrelation:   c,
		Difficulty:            r.data.Difficulty,
		IsExperimental:        r.data.IsExperimental,
		SolvedCount:           r.data.SolvedCount,
	}, nil
}

type Document struct {
	ProblemID             string `json:"problem_id"`
	DifficultyCorrelation string `json:"difficulty_correlation"`
	CategoryCorrelation   string `json:"category_correlation"`
	Difficulty            *int   `json:"difficulty"`
	IsExperimental        bool   `json:"is_experimental"`
	SolvedCount           int    `json:"solved_count"`
}

type RowReader[R acs.ToDocument[D], D any] struct {
	db *sqlx.DB
}

func (r *RowReader[R, D]) ReadRows(ctx context.Context, tx chan<- Row) error {
	sql := `
	WITH "solved_counts" AS (
		SELECT
			"problem_id",
			COUNT(1) AS "solved_count"
		FROM
			"submissions"
		WHERE
			"result" = 'AC'
		GROUP BY
			"problem_id"
	)
	SELECT
		"problems"."problem_id" AS "problem_id",
		"contests"."category" AS "category",
		"difficulties"."difficulty" AS "difficulty",
		"difficulties"."is_experimental" AS "is_experimental",
		"solved_count"
	FROM
		"problems"
		LEFT JOIN "difficulties" ON "problems"."problem_id" = "difficulties"."problem_id"
		LEFT JOIN "contests" ON "problems"."contest_id" = "contests"."contest_id"
		LEFT JOIN "solved_counts" ON "problems"."problem_id" = "solved_counts"."problem_id"
	WHERE
		"difficulty" IS NOT NULL
	`
	rows, err := r.db.Queryx(sql)
	if err != nil {
		return failure.Translate(err, DBError, failure.Context{"sql": sql}, failure.Message("failed to read rows"))
	}
	defer rows.Close()
	defer close(tx)

	for rows.Next() {
		select {
		case <-ctx.Done():
			log.Println("ReadRows canceled.")
			return nil
		default:
			var data Data
			err := rows.StructScan(&data)
			if err != nil {
				return failure.Translate(err, DBError, failure.Message("failed to scan row"))
			}

			row := Row{
				db:   r.db,
				data: data,
				ctx:  ctx,
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

func NewDocumentGenerator(db *sqlx.DB, saveDir string) DocumentGenerator {
	return DocumentGenerator{
		saveDir: saveDir,
		reader:  &RowReader[Row, Document]{db: db},
	}
}

func (g *DocumentGenerator) Clean() error {
	if err := acs.CleanDocument(g.saveDir); err != nil {
		return failure.Translate(err, FileOperationError, failure.Context{"directory": g.saveDir}, failure.Message("failed to delete problem document files in `%s`"))
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
