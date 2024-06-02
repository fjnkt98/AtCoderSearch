package search

import (
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/server/api"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type UserParameter struct {
	api.ParameterBase
	Q             string   `json:"q" query:"q"`
	Sort          []string `json:"sort" query:"sort"`
	Facet         []string `json:"facet" query:"facet"`
	UserID        []string `json:"userId" query:"userId"`
	RatingFrom    *int     `json:"ratingFrom" query:"ratingFrom"`
	RatingTo      *int     `json:"ratingTo" query:"ratingTo"`
	BirthYearFrom *int     `json:"birthYearFrom" query:"birthYearFrom"`
	BirthYearTo   *int     `json:"birthYearTo" query:"birthYearTo"`
	JoinCountFrom *int     `json:"joinCountFrom" query:"joinCountFrom"`
	JoinCountTo   *int     `json:"joinCountTo" query:"joinCountTo"`
	Country       []string `json:"country" query:"country"`
	Color         []string `json:"color" query:"color"`
}

func (p UserParameter) Validate() error {
	return validation.ValidateStruct(
		&p,
		validation.Field(&p.Q, validation.RuneLength(0, 200)),
		validation.Field(&p.Limit, validation.Min(0), validation.Max(200)),
		validation.Field(&p.Page, validation.Min(0)),
		validation.Field(&p.Sort, validation.Each(validation.In("-score", "rating", "-rating", "birthYear", "-birthYear"))),
		validation.Field(&p.Facet, validation.Each(validation.In("country", "rating", "birthYear", "joinCount"))),
	)
}

func (p *UserParameter) Query(core *solr.SolrCore) *solr.SelectQuery {
	q := core.NewSelect().
		Rows(p.Rows()).
		Start(p.Start()).
		Sort(strings.Join(api.ParseSort(p.Sort), ",")).
		Fl(strings.Join(solr.FieldList(new(UserResponse)), ",")).
		Q(api.ParseQ(p.Q)).
		Op("AND").
		Qf("text_unigram").
		QAlt("*:*").
		Fq(
			api.TermsFilter(p.UserID, "userId"),
			api.IntegerRangeFilter(p.RatingFrom, p.RatingTo, "rating", api.LocalParam("tag", "rating")),
			api.IntegerRangeFilter(p.BirthYearFrom, p.BirthYearTo, "birthYear", api.LocalParam("tag", "birthYear")),
			api.IntegerRangeFilter(p.JoinCountFrom, p.JoinCountTo, "joinCount", api.LocalParam("tag", "joinCount")),
			api.TermsFilter(p.Country, "country", api.LocalParam("tag", "country")),
			api.TermsFilter(p.Color, "color", api.LocalParam("tag", "color")),
		)

	jsonFacet := solr.NewJSONFacetQuery()
	for _, f := range p.Facet {
		switch f {
		case "country":
			jsonFacet.Terms(solr.NewTermsFacetQuery(f).Limit(-1).MinCount(0).Sort("index").ExcludeTags(f))
		case "rating":
			jsonFacet.Range(solr.NewRangeFacetQuery("rating", 0, 4000, 400).Other("all").ExcludeTags("rating"))
		case "birthYear":
			jsonFacet.Range(solr.NewRangeFacetQuery("birthYear", 1970, 2020, 10).Other("all").ExcludeTags("birthYear"))
		case "joinCount":
			jsonFacet.Range(solr.NewRangeFacetQuery("joinCount", 0, 100, 20).Other("all").ExcludeTags("joinCount"))
		}
	}
	q = q.JsonFacet(jsonFacet)

	return q
}

type UserResponse struct {
	UserID        string  `json:"userId"`
	Rating        int     `json:"rating"`
	HighestRating int     `json:"highestRating"`
	Affiliation   *string `json:"affiliation"`
	BirthYear     *int    `json:"birthYear"`
	Country       *string `json:"country"`
	Crown         *string `json:"crown"`
	JoinCount     int     `json:"joinCount"`
	Rank          int     `json:"rank"`
	ActiveRank    *int    `json:"activeRank"`
	Wins          int     `json:"wins" `
	Color         string  `json:"color"`
	UserURL       string  `json:"userUrl"`
	Score         float64 `json:"score"`
}

type SearchUserHandler struct {
	core *solr.SolrCore
	pool *pgxpool.Pool
}

func NewSearchUserHandler(core *solr.SolrCore, pool *pgxpool.Pool) *SearchUserHandler {
	return &SearchUserHandler{
		core: core,
		pool: pool,
	}
}

func (h *SearchUserHandler) SearchUser(ctx echo.Context) error {
	var p UserParameter
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

	var items []UserResponse
	if err := res.Scan(&items); err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: api.NewErrorResponse("request failed", p), Internal: err}
	}

	result := api.ResultResponse[UserResponse]{
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

func (h *SearchUserHandler) Register(e *echo.Echo) {
	e.GET("/api/search/user", h.SearchUser)
	e.POST("/api/search/user", h.SearchUser)
}
