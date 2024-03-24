package search

import (
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/server/api"
	"log/slog"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/goark/errs"
	"github.com/labstack/echo/v4"
)

type ProblemParameter struct {
	api.ParameterBase
	Q              string   `json:"q" query:"q"`
	Sort           []string `json:"sort" query:"sort"`
	Facet          []string `json:"facet" query:"facet"`
	Category       []string `json:"category" query:"category"`
	DifficultyFrom *int     `json:"difficultyFrom" query:"difficultyFrom"`
	DifficultyTo   *int     `json:"difficultyTo" query:"difficultyTo"`
	Color          []string `json:"color" query:"color"`
}

func (p ProblemParameter) Validate() error {
	return validation.ValidateStruct(
		&p,
		validation.Field(&p.Q, validation.RuneLength(0, 200)),
		validation.Field(&p.Limit, validation.Min(0), validation.Max(200)),
		validation.Field(&p.Page, validation.Min(0)),
		validation.Field(&p.Sort, validation.Each(validation.In("-score", "startAt", "-startAt", "difficulty", "-difficulty"))),
		validation.Field(&p.Facet, validation.Each(validation.In("category", "difficulty"))),
	)
}

type ProblemResponse struct {
	ProblemID    string                `json:"problemId" solr:"problem_id"`
	ProblemTitle string                `json:"problemTitle" solr:"problem_title"`
	ProblemURL   string                `json:"problemUrl" solr:"problem_url"`
	ContestID    string                `json:"contestId" solr:"contest_id"`
	ContestTitle string                `json:"contestTitle" solr:"contest_title"`
	ContestURL   string                `json:"contestUrl" solr:"contest_url"`
	Difficulty   *int                  `json:"difficulty" solr:"difficulty"`
	Color        *string               `json:"color" solr:"color"`
	StartAt      solr.FromSolrDateTime `json:"startAt" solr:"start_at"`
	Duration     int                   `json:"duration" solr:"duration"`
	RateChange   string                `json:"rateChange" solr:"rate_change"`
	Category     string                `json:"category" solr:"category"`
}

type SearchProblemHandler struct {
	core *solr.SolrCore
}

func NewSearchProblemHandler(core *solr.SolrCore) *SearchProblemHandler {
	return &SearchProblemHandler{
		core: core,
	}
}

func (h *SearchProblemHandler) SearchProblem(c echo.Context) error {
	var p ProblemParameter
	if err := c.Bind(&p); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	if err := c.Validate(p); err != nil {
		return c.JSON(http.StatusBadRequest, api.NewErrorResponse(err.Error(), p))
	}

	q := h.core.NewSelect().
		Rows(p.Rows()).
		Start(p.Start()).
		Sort(api.ParseSort(p.Sort)).
		Fl(strings.Join(solr.FieldList(new(ProblemResponse)), ",")).
		Op("AND").
		Qf("text_ja text_en text_reading")

	jsonFacet := solr.NewJSONFacetQuery()
	for _, f := range p.Facet {
		if f == "difficulty" {
			jsonFacet.Range(solr.NewRangeFacetQuery("difficulty", 0, 4000, 400).Other("all").ExcludeTags("difficulty"))
		}
		jsonFacet.Terms(solr.NewTermsFacetQuery(f).Limit(-1).MinCount(0).Sort("index").ExcludeTags(f))
	}
	q = q.JsonFacet(jsonFacet)

	q = q.Fq(
		api.TermsFilter(p.Category, "category", api.LocalParam("tag", "category")),
		api.TermsFilter(p.Color, "color", api.LocalParam("tag", "color")),
		api.RangeFilter(p.DifficultyFrom, p.DifficultyTo, "difficulty", api.LocalParam("tag", "color")),
	)

	res, err := q.Exec(c.Request().Context())
	if err != nil {
		slog.Error("request failed", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, api.NewErrorResponse("request failed", p))
	}

	facet, err := res.Facet()
	if err != nil && !errs.Is(err, solr.ErrNoFacetCounts) {
		slog.Error("request failed", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, api.NewErrorResponse("request failed", p))
	}

	var items []ProblemResponse
	if err := res.Scan(&items); err != nil {
		slog.Error("request failed", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, api.NewErrorResponse("request failed", p))
	}

	result := api.ResultResponse[ProblemResponse]{
		Stats: api.ResultStats{
			Total:  res.Raw.Response.NumFound,
			Index:  (res.Raw.Response.Start / p.Rows()) + 1,
			Count:  int(len(items)),
			Pages:  (res.Raw.Response.NumFound + p.Rows() - 1) / p.Rows(),
			Params: p,
			Facet:  api.NewFacetCount(facet),
		},
		Items: items,
	}

	return c.JSON(http.StatusOK, result)
}
