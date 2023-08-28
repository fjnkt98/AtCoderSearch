package list

import (
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
		slog.Error("failed to get a list of contest language", slog.String("error", err.Error()))
		return c.String(http.StatusInternalServerError, "failed to get contest language list")
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
