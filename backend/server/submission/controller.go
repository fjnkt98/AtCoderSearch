package submission

import (
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/server/utility"
	"log/slog"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
)

type SubmissionController interface {
	HandleGET(*gin.Context)
	HandlePOST(*gin.Context)
}

type submissionController struct {
	uc      SubmissionUsecase
	pr      SubmissionPresenter
	decoder *schema.Decoder
}

func NewSubmissionController(
	uc SubmissionUsecase,
	pr SubmissionPresenter,
) SubmissionController {
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	decoder.RegisterConverter([]string{}, func(input string) reflect.Value {
		return reflect.ValueOf(strings.Split(input, ","))
	})

	return &submissionController{
		uc:      uc,
		pr:      pr,
		decoder: decoder,
	}
}

func (c *submissionController) HandleGET(ctx *gin.Context) {
	startTime := time.Now()
	raw := ctx.Request.URL.RawQuery
	query, err := url.ParseQuery(raw)
	if err != nil {
		slog.Error(
			"failed to parse query string",
			slog.Any("error", err),
		)
		ctx.JSON(http.StatusBadRequest, c.pr.Error("failed to parse query string"))
		return
	}

	var params SearchParams
	if err := c.decoder.Decode(&params, query); err != nil {
		slog.Error(
			"failed to decode request parameter",
			slog.Any("error", err),
		)
		ctx.JSON(http.StatusBadRequest, c.pr.Error("failed to decode request parameter"))
		return
	}

	if !params.Validate() {
		ctx.JSON(http.StatusBadRequest, c.pr.Error("validation error"))
		return
	}

	res, err := c.uc.Search(ctx, params)
	if err != nil {
		slog.Error(
			"failed to execute search query",
			slog.Any("error", err),
		)
		ctx.JSON(http.StatusInternalServerError, c.pr.Error("internal server error"))
		return
	}

	t := int(time.Since(startTime).Milliseconds())

	ctx.JSON(http.StatusOK, c.pr.Format(params, res, t))
}

func (c *submissionController) HandlePOST(ctx *gin.Context) {
	startTime := time.Now()

	var params SearchParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		slog.Error(
			"failed to bind request body",
			slog.Any("error", err),
		)
		ctx.JSON(http.StatusBadRequest, c.pr.Error("failed to bind request body"))
		return
	}

	if !params.Validate() {
		ctx.JSON(http.StatusBadRequest, c.pr.Error("validation error"))
		return
	}

	res, err := c.uc.Search(ctx, params)
	if err != nil {
		slog.Error(
			"failed to execute search query",
			slog.Any("error", err),
		)
		ctx.JSON(http.StatusInternalServerError, c.pr.Error("internal server error"))
		return
	}

	t := int(time.Since(startTime).Milliseconds())

	ctx.JSON(http.StatusOK, c.pr.Format(params, res, t))
}

type SearchParams struct {
	utility.SearchParams[FilterParams, FacetParams]
}

var SortValues = mapset.NewSet[string]("execution_time", "-execution_time", "submitted_at", "-submitted_at", "point", "-point", "length", "-length")

func (p SearchParams) Validate() bool {
	if p.Limit != nil && *p.Limit > 1000 {
		return false
	}
	if p.Page > 1000 {
		return false
	}
	if !SortValues.Contains(p.Sort...) {
		return false
	}
	if !p.Filter.Validate() {
		return false
	}
	if !p.Facet.Validate() {
		return false
	}

	return true
}

type FilterParams struct {
	EpochSecond   utility.IntegerRange `json:"epoch_second" schema:"epoch_second"`
	ProblemID     []string             `json:"problem_id" schema:"problem_id"`
	ContestID     []string             `json:"contest_id" schema:"contest_id"`
	Category      []string             `json:"category" schema:"category"`
	UserID        []string             `json:"user_id" schema:"user_id"`
	Language      []string             `json:"language" schema:"language"`
	LanguageGroup []string             `json:"language_group" schema:"language_group"`
	Point         utility.FloatRange   `json:"point" schema:"point"`
	Length        utility.IntegerRange `json:"length" schema:"length"`
	Result        []string             `json:"result" schema:"result"`
	ExecutionTime utility.IntegerRange `json:"execution_time" schema:"execution_time"`
}

func (p FilterParams) Validate() bool {
	return true
}

type FacetParams struct {
	Term          []string                `json:"term" schema:"term" facet:"problem_id:problem_id,user_id:user_id,language:language,language_group:language_group,result:result,contest_id:contest_id"`
	Length        utility.RangeFacetParam `json:"length" schema:"length" facet:"length:length"`
	ExecutionTime utility.RangeFacetParam `json:"execution_time" schema:"execution_time" facet:"execution_time:execution_time"`
}

var TermValues = mapset.NewSet[string]("problem_id", "user_id", "language", "language_group", "result", "contest_id")

func (p FacetParams) Validate() bool {
	return TermValues.Contains(p.Term...)
}

func (p *SearchParams) Query() url.Values {
	return solr.NewLuceneQueryBuilder().
		Facet(p.GetFacet()).
		Fl(strings.Join(utility.FieldList(new(Submission)), ",")).
		Fq(p.GetFilter()).
		Op("AND").
		Q("*:*").
		Rows(p.GetRows()).
		Sort(p.GetSort()).
		Start(p.GetStart()).
		Build()
}
