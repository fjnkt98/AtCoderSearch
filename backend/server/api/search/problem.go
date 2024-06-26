package search

import (
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/server/api"
	"fjnkt98/atcodersearch/settings"
	"fmt"
	"net/http"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type ProblemParameter struct {
	api.ParameterBase
	Q                string   `json:"q" query:"q"`
	Sort             []string `json:"sort" query:"sort"`
	Facet            []string `json:"facet" query:"facet"`
	Category         []string `json:"category" query:"category"`
	DifficultyFrom   *int     `json:"difficultyFrom" query:"difficultyFrom"`
	DifficultyTo     *int     `json:"difficultyTo" query:"difficultyTo"`
	Color            []string `json:"color" query:"color"`
	UserID           string   `json:"userId" query:"userId"`
	Difficulty       *int     `json:"difficulty" query:"difficulty"`
	ExcludeSolved    bool     `json:"excludeSolved" query:"excludeSolved"`
	Experimental     *bool    `json:"experimental" query:"experimental"`
	PrioritizeRecent bool     `json:"prioritizeRecent" query:"prioritizeRecent"`
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

func (p *ProblemParameter) Query(core *solr.SolrCore) *solr.SelectQuery {
	q := core.NewSelect().
		Rows(p.Rows()).
		Start(p.Start()).
		Sort(strings.Join(api.ParseSort(p.Sort, "problemId asc"), ",")).
		Fl(strings.Join(solr.FieldList(new(ProblemResponse)), ",")).
		Q(api.ParseQ(p.Q)).
		Op("AND").
		Qf("text_ja text_en text_reading text_unigram^0 text_ws^10").
		Bq(p.Bq()...).
		Sow(true).
		QAlt("*:*").
		Fq(
			api.TermsFilter(p.Category, "category", api.LocalParam("tag", "category")),
			api.TermsFilter(p.Color, "color", api.LocalParam("tag", "color")),
			api.IntegerRangeFilter(p.DifficultyFrom, p.DifficultyTo, "difficulty", api.LocalParam("tag", "difficulty")),
			api.PointerBoolFilter(p.Experimental, "isExperimental"),
		)

	if p.ExcludeSolved && p.UserID != "" {
		q = q.Fq(
			fmt.Sprintf(`-{!join fromIndex=%s from=problemId to=problemId v='userId:"%s"'}`, settings.SOLUTION_CORE_NAME, solr.Sanitize(p.UserID)),
		)
	}

	if p.PrioritizeRecent {
		q = q.Bq(
			fmt.Sprintf("{!boost b=%d}{!func}pow(2,mul(-1,div(ms(NOW,startAt),2592000000)))", 7),
		)
	}
	if p.Difficulty != nil {
		q = q.Bq(
			fmt.Sprintf("{!boost b=%d}{!func}pow(2.71828182846,mul(-1,div(pow(sub(%d,difficulty),2),20000)))", 10, *p.Difficulty),
		)
	}

	jsonFacet := solr.NewJSONFacetQuery()
	for _, f := range p.Facet {
		switch f {
		case "category":
			jsonFacet.Terms(solr.NewTermsFacetQuery(f).Limit(-1).MinCount(0).Sort("index").ExcludeTags(f))
		case "difficulty":
			jsonFacet.Range(solr.NewRangeFacetQuery("difficulty", 0, 4000, 400).Other("all").ExcludeTags("difficulty"))
		}
	}
	q = q.JsonFacet(jsonFacet)

	return q
}

func (p *ProblemParameter) Bq() []string {
	parsed := solr.Parse(p.Q)
	res := make([]string, 0, len(parsed))
	for _, w := range parsed {
		if w.Negative {
			continue
		}

		w.Phrase = true
		res = append(
			res,
			fmt.Sprintf("text_ja:%s", w.String()),
			fmt.Sprintf("text_en:%s", w.String()),
		)
	}
	return res
}

type ProblemResponse struct {
	ProblemID    string                `json:"problemId"`
	ProblemTitle string                `json:"problemTitle"`
	ProblemURL   string                `json:"problemUrl"`
	ContestID    string                `json:"contestId"`
	ContestTitle string                `json:"contestTitle"`
	ContestURL   string                `json:"contestUrl"`
	Difficulty   *int                  `json:"difficulty"`
	Color        *string               `json:"color"`
	StartAt      solr.FromSolrDateTime `json:"startAt"`
	Duration     int                   `json:"duration"`
	RateChange   string                `json:"rateChange"`
	Category     string                `json:"category"`
	Score        float64               `json:"score"`
}

type SearchProblemHandler struct {
	core *solr.SolrCore
	pool *pgxpool.Pool
}

func NewSearchProblemHandler(core *solr.SolrCore, pool *pgxpool.Pool) *SearchProblemHandler {
	return &SearchProblemHandler{
		core: core,
		pool: pool,
	}
}

func (h *SearchProblemHandler) SearchProblem(ctx echo.Context) error {
	startAt := time.Now()

	var p ProblemParameter
	if err := ctx.Bind(&p); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: api.NewErrorResponse("bad request", nil)}
	}
	if err := ctx.Validate(p); err != nil {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: api.NewErrorResponse(err.Error(), p), Internal: err}
	}

	q := p.Query(h.core)
	res, err := q.Exec(ctx.Request().Context())
	if err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: api.NewErrorResponse("request failed", p), Internal: err}
	}

	facet, err := res.Facet()
	if err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: api.NewErrorResponse("request failed", p), Internal: err}
	}

	var items []ProblemResponse
	if err := res.Scan(&items); err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: api.NewErrorResponse("request failed", p), Internal: err}
	}

	result := api.ResultResponse[ProblemResponse]{
		Stats: api.ResultStats{
			Time:   int(time.Since(startAt).Milliseconds()),
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

func (h *SearchProblemHandler) Register(e *echo.Echo) {
	e.GET("/api/search/problem", h.SearchProblem)
	e.POST("/api/search/problem", h.SearchProblem)
}
