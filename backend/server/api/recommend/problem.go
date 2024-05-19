package recommend

import (
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/server/api"
	"log/slog"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
)

type ProblemParameter struct {
	api.ParameterBase
	ProblemID string `json:"problemId" query:"problemId"`
}

func (r *ProblemParameter) Rows() int {
	if r.Limit == nil {
		return 5
	}
	return *r.Limit
}
func (p ProblemParameter) Validate() error {
	return validation.ValidateStruct(
		&p,
		validation.Field(&p.ProblemID, validation.RuneLength(0, 200)),
		validation.Field(&p.Limit, validation.Min(0), validation.Max(200)),
		validation.Field(&p.Page, validation.Min(0)),
	)
}

func (p *ProblemParameter) Query(core *solr.SolrCore) *solr.MoreLikeThisQuery {
	q := core.NewMoreLikeThis().
		Rows(p.Rows()).
		Start(p.Start()).
		Fl(strings.Join(solr.FieldList(new(ProblemResponse)), ",")).
		Q(p.ProblemID).
		Qf("text_ja").
		MinTF(2).
		MinDF(5)

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

type RecommendProblemHandler struct {
	core *solr.SolrCore
}

func NewRecommendProblemHandler(core *solr.SolrCore) *RecommendProblemHandler {
	return &RecommendProblemHandler{
		core: core,
	}
}

func (h *RecommendProblemHandler) RecommendProblem(ctx echo.Context) error {
	var p ProblemParameter
	if err := ctx.Bind(&p); err != nil {
		return ctx.JSON(http.StatusBadRequest, api.NewErrorResponse("bad request", nil))
	}
	if err := ctx.Validate(p); err != nil {
		return ctx.JSON(http.StatusBadRequest, api.NewErrorResponse(err.Error(), p))
	}

	q := p.Query(h.core)
	res, err := q.Exec(ctx.Request().Context())
	if err != nil {
		if strings.HasPrefix(res.Raw.Error.Msg, "Error completing MLT request.") {
			return ctx.JSON(http.StatusOK, api.NewEmptyResponse())
		}
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
		},
		Items: items,
	}

	return ctx.JSON(http.StatusOK, result)
}

func (h *RecommendProblemHandler) Register(e *echo.Echo) {
	e.GET("/api/recommend/problem", h.RecommendProblem)
	e.POST("/api/recommend/problem", h.RecommendProblem)
}
