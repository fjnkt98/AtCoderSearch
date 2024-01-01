package recommend

import (
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/server/utility"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type RecommendController interface {
	HandleGET(*gin.Context)
	HandlePOST(*gin.Context)
}

type recommendController struct {
	uc RecommendUsecase
	pr RecommendPresenter
}

func NewRecommendController(
	uc RecommendUsecase,
	pr RecommendPresenter,
) RecommendController {
	return &recommendController{
		uc: uc,
		pr: pr,
	}
}

func (c *recommendController) HandleGET(ctx *gin.Context) {
	startTime := time.Now()

	var params SearchParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		slog.Error(
			"failed to decode request parameter",
			slog.Any("error", err),
		)
		ctx.JSON(http.StatusBadRequest, c.pr.Error("failed to decode request parameter"))
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

func (c *recommendController) HandlePOST(ctx *gin.Context) {
	startTime := time.Now()

	var params SearchParams
	if err := ctx.ShouldBindJSON(&params); err != nil {
		slog.Error(
			"failed to decode request parameter",
			slog.Any("error", err),
		)
		ctx.JSON(http.StatusBadRequest, c.pr.Error("failed to decode request parameter"))
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
	Model    int    `json:"model" form:"model" validate:"required,model"`
	Option   string `json:"option" form:"option" validate:"omitempty,option"`
	UserID   string `json:"user_id" form:"user_id"`
	Rating   int    `json:"rating" form:"rating"`
	Limit    *int   `json:"limit" form:"limit" validate:"omitempty,lte=200"`
	Page     int    `json:"page" form:"page"`
	Unsolved bool   `json:"unsolved" form:"unsolved"`
}

func (p *SearchParams) Query() url.Values {
	return solr.NewEDisMaxQueryBuilder().
		Bq(p.GetBq()).
		Fq(p.GetFq()).
		Fl(strings.Join(utility.FieldList(new(Recommend)), ",")).
		QAlt("*:*^=0").
		Rows(p.GetRows()).
		Start(p.GetStart()).
		Sort("score desc,problem_id asc").
		Build()
}

func (p *SearchParams) GetRows() int {
	if p.Limit == nil {
		return 20
	}

	return *p.Limit
}

func (p *SearchParams) GetStart() int {
	if p.Page == 0 || p.GetRows() == 0 {
		return 0
	}

	return (p.Page - 1) * p.GetRows()
}

func (p *SearchParams) GetFq() []string {
	fq := make([]string, 1)
	if p.Unsolved {
		fq = append(fq, fmt.Sprintf(`-{!join fromIndex=submission from=problem_id to=problem_id v="+user_id:%s +result:AC"}`, solr.Sanitize(p.UserID)))
	}
	return fq
}

const (
	RECENT  = 1
	RATING  = 2
	HISTORY = 3
)

type Weights struct {
	Trend           int
	Difficulty      int
	ABC             int
	ARC             int
	AGC             int
	Other           int
	NotExperimental int
}

func (p *SearchParams) GetBq() []string {
	bq := make([]string, 0)

	var w Weights
	var rate int

	switch p.Model {
	case RECENT:
		w = Weights{Trend: 10}
	case RATING:
		w = Weights{Trend: 3, Difficulty: 10, ABC: 5, ARC: 5, AGC: 5, Other: 1, NotExperimental: 0}

		if p.Option != "" {
			opt := []rune(p.Option)

			switch opt[0] {
			case '0':
				rate = p.Rating - 200
			case '1':
				rate = p.Rating
			case '2':
				rate = p.Rating + 200
			}

			switch opt[1] {
			case '1':
				w.ABC = 16
				w.ARC = 4
				w.AGC = 2
			case '2':
				w.ABC = 2
				w.ARC = 16
				w.AGC = 4
			case '3':
				w.ABC = 2
				w.ARC = 4
				w.AGC = 16
			default:
			}

			switch opt[2] {
			case '0':
				w.Trend = 3
			case '1':
				w.Trend = 7
			}

			switch opt[3] {
			case '0':
				w.NotExperimental = 0
			case '1':
				w.NotExperimental = 10
			}
		}
	case HISTORY:
		// TODO
		w = Weights{Trend: 10}
	}

	bq = append(bq, fmt.Sprintf("{!boost b=%d}{!func}pow(2,mul(-1,div(ms(NOW,start_at),2592000000)))", w.Trend))
	bq = append(bq, fmt.Sprintf("{!boost b=%d}{!func}pow(2.71828182846,mul(-1,div(pow(sub(%d,difficulty),2),20000)))", w.Difficulty, rate))
	bq = append(bq, fmt.Sprintf(`{!boost b=%d}(category:"ABC" OR category:"ABC-Like"^0.5)`, w.ABC))
	bq = append(bq, fmt.Sprintf(`{!boost b=%d}(category:"ARC" OR category:"ARC-Like"^0.5)`, w.ARC))
	bq = append(bq, fmt.Sprintf(`{!boost b=%d}(category:"AGC" OR category:"AGC-Like"^0.5)`, w.AGC))
	bq = append(bq, fmt.Sprintf(`{!boost b=%d}category:("JOI" OR "Other Sponsored" OR "Other Contests" OR "PAST")`, w.Other))
	bq = append(bq, fmt.Sprintf(`{!boost b=%d}is_experimental:false`, w.NotExperimental))
	bq = append(bq, "{!boost b=0.2}{!join fromIndex=recommend from=problem_id to=problem_id score=max}{!func v=log(add(solved_count,1))}")

	return bq
}

func ValidateModel(fl validator.FieldLevel) bool {
	if !fl.Field().CanInt() {
		return false
	}

	if model := fl.Field().Int(); model == RECENT || model == RATING || model == HISTORY {
		return true
	}
	return false
}

func ValidateOption(fl validator.FieldLevel) bool {
	s := fl.Field().String()

	if _, err := strconv.Atoi(s); err != nil {
		return false
	}

	opt := []rune(s)
	if len(opt) != 4 {
		return false
	}

	if !('0' <= opt[0] && opt[0] <= '2') {
		return false
	}
	if !('0' <= opt[1] && opt[1] <= '3') {
		return false
	}
	if !('0' <= opt[2] && opt[2] <= '1') {
		return false
	}
	if !('0' <= opt[3] && opt[3] <= '1') {
		return false
	}

	return true
}
