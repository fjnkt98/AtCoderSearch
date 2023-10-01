package list

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
	"golang.org/x/exp/slog"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type Searcher struct {
	db *sqlx.DB
	c  *cache.Cache
}

func NewSearcher(db *sqlx.DB) Searcher {
	searcher := Searcher{
		db: db,
		c:  cache.New(1*time.Minute, 5*time.Minute),
	}
	return searcher
}

func (s *Searcher) HandleCategory(c echo.Context) error {
	if categories, ok := s.c.Get("categories"); ok {
		return c.JSON(http.StatusOK, categories)
	} else {
		rows, err := s.db.Query(
			`
			SELECT
				DISTINCT "category"
			FROM
				"contests"
			ORDER BY
				"category"
			`,
		)
		if err != nil {
			slog.Error("failed to get a list of contest category", slog.String("error", err.Error()))
			return c.String(http.StatusInternalServerError, "failed to get contest category list")
		}
		defer rows.Close()

		categories := make([]string, 0, 12)
		for rows.Next() {
			var category string
			if err := rows.Scan(&category); err != nil {
				slog.Error("failed to scan row for fetching a list of category", slog.String("error", err.Error()))
				return c.String(http.StatusInternalServerError, "failed to scan category")
			}
			categories = append(categories, category)
		}

		s.c.Set("categories", categories, cache.DefaultExpiration)

		return c.JSON(http.StatusOK, categories)
	}
}

func (s *Searcher) HandleLanguage(c echo.Context) error {
	if languages, ok := s.c.Get("languages"); ok {
		return c.JSON(http.StatusOK, languages)
	} else {
		rows, err := s.db.Query(
			`
			SELECT
				"language"
			FROM
				"languages"
			ORDER BY
				"language"
			`,
		)
		if err != nil {
			slog.Error("failed to get a list of language", slog.String("error", err.Error()))
			return c.String(http.StatusInternalServerError, "failed to get language list")
		}
		defer rows.Close()

		languages := make([]string, 0, 12)
		for rows.Next() {
			var language string
			if err := rows.Scan(&language); err != nil {
				slog.Error("failed to scan row for fetching a list of language", slog.String("error", err.Error()))
				return c.String(http.StatusInternalServerError, "failed to scan language")
			}
			languages = append(languages, language)
		}
		s.c.Set("languages", languages, cache.DefaultExpiration)

		return c.JSON(http.StatusOK, languages)
	}
}

func (s *Searcher) HandleContest(c echo.Context) error {
	category := c.QueryParam("category")

	if contests, ok := s.c.Get(fmt.Sprintf("contests_of_%s", category)); ok {
		return c.JSON(http.StatusOK, contests)
	} else {
		var rows *sql.Rows
		var err error
		if category == "" {
			rows, err = s.db.Query(
				`
				SELECT
					"contest_id"
				FROM
					"contests"
				ORDER BY
					"start_epoch_second" DESC
				`,
			)
		} else {
			rows, err = s.db.Query(
				`
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
		defer rows.Close()

		contests := make([]string, 0)
		for rows.Next() {
			var contest string
			if err := rows.Scan(&contest); err != nil {
				slog.Error("failed to scan row for fetching a list of contests", slog.String("error", err.Error()))
				return c.String(http.StatusInternalServerError, "failed to scan contest")
			}
			contests = append(contests, contest)
		}
		s.c.Set(fmt.Sprintf("contests_of_%s", category), contests, cache.DefaultExpiration)

		return c.JSON(http.StatusOK, contests)
	}

}

func (s *Searcher) HandleProblem(c echo.Context) error {
	id := c.QueryParam("contest_id")

	if problems, ok := s.c.Get(fmt.Sprintf("problems_in_%s", id)); ok {
		return c.JSON(http.StatusOK, problems)
	} else {
		var rows *sql.Rows
		var err error
		if id == "" {
			rows, err = s.db.Query(
				`
				SELECT
					"problem_id"
				FROM
					"problems"
				ORDER BY
					"problem_id"
				`,
			)
		} else {
			rows, err = s.db.Query(
				`
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
		defer rows.Close()

		problems := make([]string, 0, 6)
		for rows.Next() {
			var problem string
			if err := rows.Scan(&problem); err != nil {
				slog.Error("failed to scan row for fetching a list of problems", slog.String("error", err.Error()))
				return c.String(http.StatusInternalServerError, "failed to scan problem")
			}
			problems = append(problems, problem)
		}
		s.c.Set(fmt.Sprintf("problems_in_%s", id), problems, cache.DefaultExpiration)

		return c.JSON(http.StatusOK, problems)
	}
}
