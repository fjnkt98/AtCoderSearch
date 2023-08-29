package list

import (
	"database/sql"
	"net/http"

	"golang.org/x/exp/slog"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type Searcher struct {
	db *sqlx.DB
}

func NewSearcher(db *sqlx.DB) Searcher {
	searcher := Searcher{
		db: db,
	}
	return searcher
}

func (s *Searcher) HandleCategory(c echo.Context) error {
	rows, err := s.db.Query(`
	SELECT
		DISTINCT "category"
	FROM
		"contests"
	ORDER BY
		"category"
	`)
	if err != nil {
		slog.Error("failed to get a list of contest category", slog.String("error", err.Error()))
		return c.String(http.StatusInternalServerError, "failed to get contest category list")
	}

	categories := make([]string, 0, 12)
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			slog.Error("failed to scan row for fetching a list of category", slog.String("error", err.Error()))
			return c.String(http.StatusInternalServerError, "failed to scan category")
		}
		categories = append(categories, category)
	}

	return c.JSON(http.StatusOK, categories)
}

func (s *Searcher) HandleLanguage(c echo.Context) error {
	rows, err := s.db.Query(`
	SELECT
		"language"
	FROM
		"languages"
	ORDER BY
		"language"
	`)
	if err != nil {
		slog.Error("failed to get a list of language", slog.String("error", err.Error()))
		return c.String(http.StatusInternalServerError, "failed to get language list")
	}

	languages := make([]string, 0, 12)
	for rows.Next() {
		var language string
		if err := rows.Scan(&language); err != nil {
			slog.Error("failed to scan row for fetching a list of language", slog.String("error", err.Error()))
			return c.String(http.StatusInternalServerError, "failed to scan language")
		}
		languages = append(languages, language)
	}

	return c.JSON(http.StatusOK, languages)
}

func (s *Searcher) HandleContest(c echo.Context) error {
	category := c.QueryParam("category")

	var rows *sql.Rows
	var err error
	if category == "" {
		rows, err = s.db.Query(`
		SELECT
			"contest_id"
		FROM
			"contests"
		ORDER BY
			"start_epoch_second" DESC
		`)
	} else {
		rows, err = s.db.Query(`
		SELECT
			"contest_id"
		FROM
			"contests"
		WHERE
			"category" = $1::text
		ORDER BY
			"start_epoch_second" DESC
		`,
			category,
		)
	}

	if err != nil {
		slog.Error("failed to get a list of contests ", slog.String("error", err.Error()))
		return c.String(http.StatusInternalServerError, "failed to get contests")
	}

	contests := make([]string, 0)
	for rows.Next() {
		var contest string
		if err := rows.Scan(&contest); err != nil {
			slog.Error("failed to scan row for fetching a list of contests", slog.String("error", err.Error()))
			return c.String(http.StatusInternalServerError, "failed to scan contest")
		}
		contests = append(contests, contest)
	}

	return c.JSON(http.StatusOK, contests)
}

func (s *Searcher) HandleProblem(c echo.Context) error {
	id := c.QueryParam("contest_id")

	var rows *sql.Rows
	var err error
	if id == "" {
		rows, err = s.db.Query(`
		SELECT
			"problem_id"
		FROM
			"problems"
		ORDER BY
			"problem_id"
		`)
	} else {
		rows, err = s.db.Query(`
		SELECT
			"problem_id"
		FROM
			"problems"
		WHERE
			"contest_id" = $1::text
		ORDER BY
			"problem_id"
		`,
			id,
		)
	}

	if err != nil {
		slog.Error("failed to get a list of problems ", slog.String("error", err.Error()))
		return c.String(http.StatusInternalServerError, "failed to get problems")
	}

	problems := make([]string, 0, 6)
	for rows.Next() {
		var problem string
		if err := rows.Scan(&problem); err != nil {
			slog.Error("failed to scan row for fetching a list of problems", slog.String("error", err.Error()))
			return c.String(http.StatusInternalServerError, "failed to scan problem")
		}
		problems = append(problems, problem)
	}

	return c.JSON(http.StatusOK, problems)
}
