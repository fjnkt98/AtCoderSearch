package search

import (
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/repository"
	"fjnkt98/atcodersearch/server/api"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/goark/errs"
	"github.com/jackc/pgx/v5"
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
	Rating           *int     `json:"rating" query:"rating"`
	ExcludeSolved    bool     `json:"excludeSolved" query:"excludeSolved"`
	Experimental     *bool    `json:"experimental" query:"experimental"`
	PrioritizeRecent bool     `json:"prioritizeRecent" query:"prioritizeRecent"`
	UseUserRating    bool     `json:"useUserRating" query:"useUserRating"`
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
		Sort(api.ParseSort(p.Sort)).
		Fl(strings.Join(solr.FieldList(new(ProblemResponse)), ",")).
		Q(api.ParseQ(p.Q)).
		Op("AND").
		Qf("text_ja text_en text_reading").
		Sow(true).
		QAlt("*:*").
		Fq(
			api.TermsFilter(p.Category, "category", api.LocalParam("tag", "category")),
			api.TermsFilter(p.Color, "color", api.LocalParam("tag", "color")),
			api.IntegerRangeFilter(p.DifficultyFrom, p.DifficultyTo, "difficulty", api.LocalParam("tag", "difficulty")),
			api.PointerBoolFilter(p.Experimental, "isExperimental"),
		)

	if p.ExcludeSolved {
		q = q.Fq(
			fmt.Sprintf(`-{!join fromIndex=solution from=problemId to=problemId v='userId:"%s"'}`, solr.Sanitize(p.UserID)),
		)
	}

	if p.PrioritizeRecent {
		q = q.Bq(
			fmt.Sprintf("{!boost b=%d}{!func}pow(2,mul(-1,div(ms(NOW,startAt),2592000000)))", 7),
		)
	}
	if p.Rating != nil {
		q = q.Bq(
			fmt.Sprintf("{!boost b=%d}{!func}pow(2.71828182846,mul(-1,div(pow(sub(%d,difficulty),2),20000)))", 10, *p.Rating),
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
	var p ProblemParameter
	if err := ctx.Bind(&p); err != nil {
		return ctx.JSON(http.StatusBadRequest, api.NewErrorResponse("bad request", nil))
	}
	if err := ctx.Validate(p); err != nil {
		return ctx.JSON(http.StatusBadRequest, api.NewErrorResponse(err.Error(), p))
	}

	var message string
	if p.UseUserRating && p.UserID != "" {
		rating, err := repository.New(h.pool).FetchRatingByUserID(ctx.Request().Context(), p.UserID)
		if err != nil {
			if errs.Is(err, pgx.ErrNoRows) {
				message = "specified user not found"
			} else {
				slog.Error("request failed", slog.Any("error", err))
				return ctx.JSON(http.StatusInternalServerError, api.NewErrorResponse("request failed", p))
			}
		} else {
			rating := int(rating)
			p.Rating = &rating
		}
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

	var items []ProblemResponse
	if err := res.Scan(&items); err != nil {
		slog.Error("request failed", slog.Any("error", err))
		return ctx.JSON(http.StatusInternalServerError, api.NewErrorResponse("request failed", p))
	}

	result := api.ResultResponse[ProblemResponse]{
		Stats: api.ResultStats{
			Total:  res.Raw.Response.NumFound,
			Index:  (res.Raw.Response.Start / p.Rows()) + 1,
			Count:  len(items),
			Pages:  (res.Raw.Response.NumFound + p.Rows() - 1) / p.Rows(),
			Params: p,
			Facet:  api.NewFacetCount(facet),
		},
		Items:   items,
		Message: message,
	}

	return ctx.JSON(http.StatusOK, result)
}

func (h *SearchProblemHandler) Register(e *echo.Echo) {
	e.GET("/api/search/problem", h.SearchProblem)
	e.POST("/api/search/problem", h.SearchProblem)
}
