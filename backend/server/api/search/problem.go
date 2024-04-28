package search

import (
	"context"
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

var ErrUserNotFound = errs.New("user not found")

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
	ExcludeSolved    bool     `json:"excludeSolved" query:"excludeSolved"`
	Experimental     *bool    `json:"experimental" query:"experimental"`
	PreferDifficulty string   `json:"preferDifficulty" query:"preferDifficulty"`
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
		validation.Field(&p.PreferDifficulty, validation.In("easy", "normal", "hard")),
	)
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
		return ctx.String(http.StatusBadRequest, "bad request")
	}
	if err := ctx.Validate(p); err != nil {
		return ctx.JSON(http.StatusBadRequest, api.NewErrorResponse(err.Error(), p))
	}

	q, err := h.Query(ctx.Request().Context(), p)
	var message string
	if err != nil {
		if errs.Is(err, ErrUserNotFound) {
			message = "specified user not found"
		} else {
			slog.Error("request failed", slog.Any("error", err))
			return ctx.JSON(http.StatusInternalServerError, api.NewErrorResponse("request failed", p))
		}
	}

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

func (h *SearchProblemHandler) Query(ctx context.Context, p ProblemParameter) (*solr.SelectQuery, error) {
	q := h.core.NewSelect().
		Rows(p.Rows()).
		Start(p.Start()).
		Sort(api.ParseSort(p.Sort)).
		Fl(strings.Join(solr.FieldList(new(ProblemResponse)), ",")).
		Q(api.ParseQ(p.Q)).
		Op("AND").
		Qf("text_ja text_en text_reading").
		QAlt("*:*").
		Fq(
			api.TermsFilter(p.Category, "category", api.LocalParam("tag", "category")),
			api.TermsFilter(p.Color, "color", api.LocalParam("tag", "color")),
			api.RangeFilter(p.DifficultyFrom, p.DifficultyTo, "difficulty", api.LocalParam("tag", "color")),
			api.PointerBoolFilter(p.Experimental, "isExperimental"),
		)

	var err error
	if p.UserID != "" {
		rating, fetchErr := repository.New(h.pool).FetchRatingByUserName(ctx, p.UserID)
		if fetchErr != nil {
			if errs.Is(fetchErr, pgx.ErrNoRows) {
				err = ErrUserNotFound
			} else {
				err = fetchErr
			}
		}

		if p.ExcludeSolved {
			q = q.Fq(
				fmt.Sprintf(`-{!join fromIndex=solution from=problemId to=problemId v="+userId:"%s" +result:AC"}`, solr.Sanitize(p.UserID)),
			)
		}

		weight := 3
		if p.PrioritizeRecent {
			weight = 7
		}
		switch p.PreferDifficulty {
		case "easy":
			rating -= 200
		case "hard":
			rating += 200
		}

		q = q.Bq(
			fmt.Sprintf("{!boost b=%d}{!func}pow(2,mul(-1,div(ms(NOW,startAt),2592000000)))", weight),
			fmt.Sprintf("{!boost b=%d}{!func}pow(2.71828182846,mul(-1,div(pow(sub(%d,difficulty),2),20000)))", 10, rating),
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

	return q, err
}
