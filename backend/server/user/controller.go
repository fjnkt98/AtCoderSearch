package user

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

type UserController interface {
	HandleGET(*gin.Context)
	HandlePOST(*gin.Context)
}

type userController struct {
	uc      UserUsecase
	pr      UserPresenter
	decoder *schema.Decoder
}

func NewUserController(
	uc UserUsecase,
	pr UserPresenter,
) UserController {
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	decoder.RegisterConverter([]string{}, func(input string) reflect.Value {
		return reflect.ValueOf(strings.Split(input, ","))
	})

	return &userController{
		uc:      uc,
		pr:      pr,
		decoder: decoder,
	}
}

func (c *userController) HandleGET(ctx *gin.Context) {
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

func (c *userController) HandlePOST(ctx *gin.Context) {
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
	Keyword string `json:"keyword" schema:"keyword"`
	utility.SearchParams[FilterParams, FacetParams]
}

var SortValues = mapset.NewSet[string]("-score", "rating", "-rating", "birth_year", "-birth_year")

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
	UserID    []string             `json:"user_id" schema:"user_id"`
	Rating    utility.IntegerRange `json:"rating" schema:"rating"`
	BirthYear utility.IntegerRange `json:"birth_year" schema:"birth_year"`
	JoinCount utility.IntegerRange `json:"join_count" schema:"join_count"`
	Country   []string             `json:"country" schema:"country"`
	Color     []string             `json:"color" schema:"color"`
}

func (p FilterParams) Validate() bool {
	return true
}

type FacetParams struct {
	Term      []string                `json:"term" schema:"term" validate:"dive,oneof=country" facet:"country:country"`
	Rating    utility.RangeFacetParam `json:"rating" schema:"rating" facet:"rating:rating"`
	BirthYear utility.RangeFacetParam `json:"birth_year" schema:"birth_year" facet:"birth_year:birth_year"`
	JoinCount utility.RangeFacetParam `json:"join_count" schema:"join_count" facet:"join_count:join_count"`
}

var TermValues = mapset.NewSet[string]("country")

func (p FacetParams) Validate() bool {
	return TermValues.Contains(p.Term...)
}

func (p *SearchParams) Query() url.Values {
	return solr.NewEDisMaxQueryBuilder().
		Facet(p.GetFacet()).
		Fl(strings.Join(utility.FieldList(new(User)), ",")).
		Fq(p.GetFilter()).
		Op("AND").
		Q(solr.Sanitize(norm.NFKC.String(p.Keyword))).
		QAlt("*:*").
		Qf("text_unigram").
		Rows(p.GetRows()).
		Sort(p.GetSort()).
		Sow(true).
		Start(p.GetStart()).
		Build()
}
