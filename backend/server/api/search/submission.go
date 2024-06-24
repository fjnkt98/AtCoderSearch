package search

import (
	"fjnkt98/atcodersearch/server/api"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

type SubmissionParameter struct {
	api.ParameterBase
	Sort              []string `json:"sort" query:"sort"`
	EpochSecondFrom   *int     `json:"epochSecondFrom" query:"epochSecondFrom"`
	EpochSecondTo     *int     `json:"epochSecondTo" query:"epochSecondTo"`
	ProblemID         []string `json:"problemId" query:"problemId"`
	ContestID         []string `json:"contestId" query:"contestId"`
	Category          []string `json:"category" query:"category"`
	UserID            []string `json:"userId" query:"userId"`
	Language          []string `json:"language" query:"language"`
	LanguageGroup     []string `json:"languageGroup" query:"languageGroup"`
	PointFrom         *float64 `json:"pointFrom" query:"pointFrom"`
	PointTo           *float64 `json:"pointTo" query:"pointTo"`
	LengthFrom        *int     `json:"lengthFrom" query:"lengthFrom"`
	LengthTo          *int     `json:"lengthTo" query:"lengthTo"`
	Result            []string `json:"result" query:"result"`
	ExecutionTimeFrom *int     `json:"executionTimeFrom" query:"executionTimeFrom"`
	ExecutionTimeTo   *int     `json:"executionTimeTo" query:"executionTimeTo"`
}

func (p SubmissionParameter) Validate() error {
	return validation.ValidateStruct(
		&p,
		validation.Field(&p.Limit, validation.Min(0), validation.Max(200)),
		validation.Field(&p.Page, validation.Min(0)),
		validation.Field(&p.Sort, validation.Each(validation.In(
			"execution_time",
			"-execution_time",
			"epoch_second",
			"-epoch_second",
			"point",
			"-point",
			"length",
			"-length",
		))),
	)
}

func (p *SubmissionParameter) Query(pool *pgxpool.Pool) *bun.SelectQuery {
	db := bun.NewDB(stdlib.OpenDBFromPool(pool), pgdialect.New())

	q := db.NewSelect().
		ColumnExpr("s.id AS submission_id").
		ColumnExpr("TO_TIMESTAMP(s.epoch_second) AS submitted_at").
		ColumnExpr("FORMAT('https://atcoder.jp/contests/%s/submissions/%s', s.contest_id, s.id) AS submission_url").
		ColumnExpr("s.problem_id").
		ColumnExpr("p.title AS problem_title").
		ColumnExpr("s.contest_id").
		ColumnExpr("c.title AS contest_title").
		ColumnExpr("c.category").
		ColumnExpr("d.difficulty").
		ColumnExpr("CASE WHEN d.difficulty < 0 THEN 'black' WHEN d.difficulty < 400 THEN 'gray' WHEN d.difficulty < 800 THEN 'brown' WHEN d.difficulty < 1200 THEN 'green' WHEN d.difficulty < 1600 THEN 'cyan' WHEN d.difficulty < 2000 THEN 'blue' WHEN d.difficulty < 2400 THEN 'yellow' WHEN d.difficulty < 2800 THEN 'orange' WHEN d.difficulty < 3200 THEN 'red' WHEN d.difficulty < 3600 THEN 'silver' ELSE 'gold' END AS color").
		ColumnExpr("s.user_id").
		ColumnExpr("s.language").
		ColumnExpr("l.group AS language_group").
		ColumnExpr("s.point").
		ColumnExpr("s.length").
		ColumnExpr("s.result").
		ColumnExpr("s.execution_time").
		TableExpr("submissions AS s").
		Join("LEFT JOIN contests AS c ON s.contest_id = c.contest_id").
		Join("LEFT JOIN problems AS p ON s.problem_id = p.problem_id").
		Join("LEFT JOIN difficulties AS d ON s.problem_id = d.problem_id").
		Join("LEFT JOIN languages AS l ON s.language = l.language").
		Order(api.ParseSort(p.Sort)...).
		Limit(p.Rows()).
		Offset(p.Start())

	if p.EpochSecondFrom != nil {
		q = q.Where("s.epoch_second > ?", p.EpochSecondFrom)
	}
	if p.EpochSecondTo != nil {
		q = q.Where("s.epoch_second < ?", p.EpochSecondTo)
	}
	if len(p.ProblemID) > 0 {
		q = q.Where("s.problem_id in (?)", bun.In(p.ProblemID))
	}
	if len(p.ContestID) > 0 {
		q = q.Where("s.contest_id in (?)", bun.In(p.ContestID))
	}
	if len(p.Category) > 0 {
		q = q.Where("c.category in (?)", bun.In(p.Category))
	}
	if len(p.UserID) > 0 {
		q = q.Where("s.user_id in (?)", bun.In(p.UserID))
	}
	if len(p.Language) > 0 {
		q = q.Where("s.language in (?)", bun.In(p.Language))
	}
	if len(p.LanguageGroup) > 0 {
		q = q.Where("l.group in (?)", bun.In(p.LanguageGroup))
	}
	if p.PointFrom != nil {
		q = q.Where("s.point > ?", p.PointFrom)
	}
	if p.PointTo != nil {
		q = q.Where("s.point < ?", p.PointTo)
	}
	if p.LengthFrom != nil {
		q = q.Where("s.length > ?", p.LengthFrom)
	}
	if p.LengthTo != nil {
		q = q.Where("s.length < ?", p.LengthTo)
	}
	if len(p.Result) > 0 {
		q = q.Where("s.result in (?)", bun.In(p.Result))
	}
	if p.ExecutionTimeFrom != nil {
		q = q.Where("s.execution_time > ?", p.ExecutionTimeFrom)
	}
	if p.ExecutionTimeTo != nil {
		q = q.Where("s.execution_time < ?", p.ExecutionTimeTo)
	}

	return q
}

type SubmissionResponse struct {
	SubmissionID  int64     `bun:"submission_id" json:"submissionId"`
	SubmittedAt   time.Time `bun:"submitted_at" json:"submittedAt"`
	SubmissionURL string    `bun:"submission_url" json:"submissionUrl"`
	ProblemID     string    `bun:"problem_id" json:"problemId"`
	ProblemTitle  string    `bun:"problem_title" json:"problemTitle"`
	ContestID     string    `bun:"contest_id" json:"contestId"`
	ContestTitle  string    `bun:"contest_title" json:"contestTitle"`
	Category      string    `bun:"category" json:"category"`
	Difficulty    int       `bun:"difficulty" json:"difficulty"`
	Color         string    `bun:"color" json:"color"`
	UserID        string    `bun:"user_id" json:"userId"`
	Language      string    `bun:"language" json:"language"`
	LanguageGroup string    `bun:"language_group" json:"languageGroup"`
	Point         float64   `bun:"point" json:"point"`
	Length        int64     `bun:"length" json:"length"`
	Result        string    `bun:"result" json:"result"`
	ExecutionTime *int64    `bun:"execution_time" json:"executionTime"`
}

type SearchSubmissionHandler struct {
	pool *pgxpool.Pool
}

func NewSearchSubmissionHandler(pool *pgxpool.Pool) *SearchSubmissionHandler {
	return &SearchSubmissionHandler{
		pool: pool,
	}
}

func (h *SearchSubmissionHandler) SearchSubmission(ctx echo.Context) error {
	startAt := time.Now()

	var p SubmissionParameter
	if err := ctx.Bind(&p); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: api.NewErrorResponse("bad request", nil)}
	}
	if err := ctx.Validate(p); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: api.NewErrorResponse(err.Error(), p), Internal: err}
	}

	q := p.Query(h.pool)
	items := make([]SubmissionResponse, 0, p.Rows())
	err := q.Scan(ctx.Request().Context(), &items)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: api.NewErrorResponse("request failed", p), Internal: err}
	}

	result := api.ResultResponse[SubmissionResponse]{
		Stats: api.ResultStats{
			Time:   int(time.Since(startAt).Milliseconds()),
			Total:  0,
			Index:  p.Page,
			Count:  len(items),
			Pages:  0,
			Params: p,
		},
		Items: items,
	}

	return ctx.JSON(http.StatusOK, result)
}

func (h *SearchSubmissionHandler) Register(e *echo.Echo) {
	e.GET("/api/search/submission", h.SearchSubmission)
	e.POST("/api/search/submission", h.SearchSubmission)
}
