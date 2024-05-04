package search

import (
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/server/api"
	"log/slog"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
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
			"executionTime",
			"-executionTime",
			"submittedAt",
			"-submittedAt",
			"point",
			"-point",
			"length",
			"-length",
		))),
	)
}

func (p *SubmissionParameter) Query(core *solr.SolrCore) *solr.SelectQuery {
	q := core.NewSelect().
		Rows(p.Rows()).
		Start(p.Start()).
		Sort(api.ParseSort(p.Sort)).
		Fl(strings.Join(solr.FieldList(new(SubmissionResponse)), ",")).
		Q("*:*").
		Df("null").
		Fq(
			api.IntegerRangeFilter(p.EpochSecondFrom, p.EpochSecondTo, "epochSecond", api.LocalParam("tag", "epochSecond")),
			api.TermsFilter(p.ProblemID, "problemId", api.LocalParam("tag", "problemId")),
			api.TermsFilter(p.ContestID, "contestId", api.LocalParam("tag", "contestId")),
			api.TermsFilter(p.Category, "category", api.LocalParam("tag", "category")),
			api.TermsFilter(p.UserID, "userId", api.LocalParam("tag", "userId")),
			api.TermsFilter(p.Language, "language", api.LocalParam("tag", "language")),
			api.TermsFilter(p.LanguageGroup, "languageGroup", api.LocalParam("tag", "languageGroup")),
			api.FloatRangeFilter(p.PointFrom, p.PointTo, "point", api.LocalParam("tag", "point")),
			api.IntegerRangeFilter(p.LengthFrom, p.LengthTo, "length", api.LocalParam("tag", "length")),
			api.TermsFilter(p.Result, "result", api.LocalParam("tag", "result")),
			api.IntegerRangeFilter(p.ExecutionTimeFrom, p.ExecutionTimeTo, "executionTime", api.LocalParam("tag", "executionTime")),
		)

	return q
}

type SubmissionResponse struct {
	SubmissionID  int64                 `json:"submissionId"`
	SubmittedAt   solr.FromSolrDateTime `json:"submittedAt"`
	SubmissionURL string                `json:"submissionUrl"`
	ProblemID     string                `json:"problemId"`
	ProblemTitle  string                `json:"problemTitle"`
	ContestID     string                `json:"contestId"`
	ContestTitle  string                `json:"contestTitle"`
	Category      string                `json:"category"`
	Difficulty    int                   `json:"difficulty"`
	Color         string                `json:"color"`
	UserID        string                `json:"userId"`
	Language      string                `json:"language"`
	Point         float64               `json:"point"`
	Length        int64                 `json:"length"`
	Result        string                `json:"result"`
	ExecutionTime *int64                `json:"executionTime"`
}

type SearchSubmissionHandler struct {
	core *solr.SolrCore
}

func NewSearchSubmissionHandler(core *solr.SolrCore) *SearchSubmissionHandler {
	return &SearchSubmissionHandler{
		core: core,
	}
}

func (h *SearchSubmissionHandler) SearchSubmission(ctx echo.Context) error {
	var p SubmissionParameter
	if err := ctx.Bind(&p); err != nil {
		return ctx.JSON(http.StatusBadRequest, api.NewErrorResponse("bad request", nil))
	}
	if err := ctx.Validate(p); err != nil {
		return ctx.JSON(http.StatusBadRequest, api.NewErrorResponse(err.Error(), p))
	}

	q := p.Query(h.core)
	res, err := q.Exec(ctx.Request().Context())
	if err != nil {
		slog.Error("request failed", slog.Any("error", err))
		return ctx.JSON(http.StatusInternalServerError, api.NewErrorResponse("request failed", p))
	}

	facet, err := res.Facet()
	if err != nil {
		slog.Error("request failed", slog.Any("error", err))
		return ctx.JSON(http.StatusInternalServerError, api.NewErrorResponse("request failed", p))
	}

	var items []SubmissionResponse
	if err := res.Scan(&items); err != nil {
		slog.Error("request failed", slog.Any("error", err))
		return ctx.JSON(http.StatusInternalServerError, api.NewErrorResponse("request failed", p))
	}

	result := api.ResultResponse[SubmissionResponse]{
		Stats: api.ResultStats{
			Total:  res.Raw.Response.NumFound,
			Index:  (res.Raw.Response.Start / p.Rows()) + 1,
			Count:  len(items),
			Pages:  (res.Raw.Response.NumFound + p.Rows() - 1) / p.Rows(),
			Params: p,
			Facet:  api.NewFacetCount(facet),
		},
		Items: items,
	}

	return ctx.JSON(http.StatusOK, result)
}

func (h *SearchSubmissionHandler) Register(e *echo.Echo) {
	e.GET("/api/search/submission", h.SearchSubmission)
	e.POST("/api/search/submission", h.SearchSubmission)
}
