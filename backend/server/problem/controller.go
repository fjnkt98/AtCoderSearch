package problem

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
	"golang.org/x/text/unicode/norm"
)

type ProblemController interface {
	HandleGET(*gin.Context)
	HandlePOST(*gin.Context)
}

type problemController struct {
	uc      ProblemUsecase
	pr      ProblemPresenter
	decoder *schema.Decoder
}

func NewProblemController(
	uc ProblemUsecase,
	pr ProblemPresenter,
) ProblemController {
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	decoder.RegisterConverter([]string{}, func(input string) reflect.Value {
		return reflect.ValueOf(strings.Split(input, ","))
	})

	return &problemController{
		uc:      uc,
		pr:      pr,
		decoder: decoder,
	}
}

func (c *problemController) HandleGET(ctx *gin.Context) {
	startTime := time.Now()
	raw := ctx.Request.URL.RawQuery
	query, err := url.ParseQuery(raw)
	if err != nil {
		slog.Error(
			"failed to parse query string",
			slog.Any("error", err),
		)
		ctx.JSON(http.StatusBadRequest, c.pr.Error("failed to parse query string"))
		ctx.Abort()
		return
	}

	var params SearchParams
	if err := c.decoder.Decode(&params, query); err != nil {
		slog.Error(
			"failed to decode request parameter",
			slog.Any("error", err),
		)
		ctx.JSON(http.StatusBadRequest, c.pr.Error("failed to decode request parameter"))
		ctx.Abort()
		return
	}

	if !params.Validate() {
		ctx.JSON(http.StatusBadRequest, c.pr.Error("validation error"))
		ctx.Abort()
		return
	}

	res, err := c.uc.Search(ctx, params)
	if err != nil {
		slog.Error(
			"failed to execute search query",
			slog.Any("error", err),
		)
		ctx.JSON(http.StatusInternalServerError, c.pr.Error("internal server error"))
		ctx.Abort()
		return
	}

	t := int(time.Since(startTime).Milliseconds())

	ctx.JSON(http.StatusOK, c.pr.Format(params, res, t))
}

func (ctr *problemController) HandlePOST(c *gin.Context) {
	panic(0)
}

type SearchParams struct {
	Keyword string `json:"keyword" schema:"keyword"`
	utility.SearchParams[FilterParams, FacetParams]
}

var SortValues = mapset.NewSet[string]("-score", "start_at", "-start_at", "difficulty", "-difficulty", "problem_id", "-problem_id")

func (p SearchParams) Validate() bool {
	if len([]rune(p.Keyword)) > 200 {
		return false
	}
	if p.Limit != nil && *p.Limit > 1000 {
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
	Category   []string             `json:"category" schema:"category" filter:"category,quote"`
	Difficulty utility.IntegerRange `json:"difficulty" schema:"difficulty" filter:"difficulty"`
	Color      []string             `json:"color" schema:"color" filter:"color"`
}

func (p FilterParams) Validate() bool {
	return true
}

type FacetParams struct {
	Term       utility.TermFacetParam  `json:"term" schema:"term" facet:"category:category,color:color"`
	Difficulty utility.RangeFacetParam `json:"difficulty" schema:"difficulty" facet:"difficulty:difficulty"`
}

var TermValues = mapset.NewSet[string]("category", "color")

func (p FacetParams) Validate() bool {
	return TermValues.Contains(p.Term...)
}

func (p *SearchParams) Query() url.Values {
	return solr.NewEDisMaxQueryBuilder().
		Facet(p.GetFacet()).
		Fl(strings.Join(utility.FieldList(new(Problem)), ",")).
		Fq(p.GetFilter()).
		Op("AND").
		Q(solr.Sanitize(norm.NFKC.String(p.Keyword))).
		QAlt("*:*").
		Rows(p.GetRows()).
		Sort(p.GetSort()).
		Sow(true).
		Start(p.GetStart()).
		Build()
}
