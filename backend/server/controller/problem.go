package controller

import (
	"fjnkt98/atcodersearch/repository"
	"fjnkt98/atcodersearch/server/domain"
	"fjnkt98/atcodersearch/server/presenter"
	"fjnkt98/atcodersearch/server/usecase"
	"fjnkt98/atcodersearch/server/utility"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	govalidator "github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
)

type SearchProblemController interface {
	HandleGET(*gin.Context)
	HandlePOST(*gin.Context)
}

type searchProblemController struct {
	uc        usecase.SearchProblemUsecase
	pr        presenter.SearchProblemPresenter
	decoder   *schema.Decoder
	validator *govalidator.Validate
}

func NewSearchProblemController(
	uc usecase.SearchProblemUsecase,
	pr presenter.SearchProblemPresenter,
) SearchProblemController {
	decoder := utility.NewSchemaDecoder()

	validator := govalidator.New()

	return &searchProblemController{
		uc:        uc,
		pr:        pr,
		decoder:   decoder,
		validator: validator,
	}
}

func (c *searchProblemController) HandleGET(ctx *gin.Context) {
	startTime := time.Now()
	raw := ctx.Request.URL.RawQuery
	query, err := url.ParseQuery(raw)
	if err != nil {
		slog.Error(
			"failed to parse query string",
			slog.String("uri", ctx.Request.RequestURI),
			slog.Any("error", err),
		)
		ctx.JSON(http.StatusBadRequest, c.pr.Error("failed to parse query string"))
		return
	}

	var params domain.SearchProblemParam
	if err := c.decoder.Decode(&params, query); err != nil {
		slog.Error(
			"failed to decode request parameter",
			slog.String("uri", ctx.Request.RequestURI),
			slog.Any("error", err),
		)
		ctx.JSON(http.StatusBadRequest, c.pr.Error("failed to decode request parameter"))
		return
	}

	if err := c.validator.Struct(params); err != nil {
		slog.Error(
			"validation error",
			slog.String("uri", ctx.Request.RequestURI),
			slog.Any("params", params),
			slog.Any("error", err),
		)
		ctx.JSON(http.StatusBadRequest, c.pr.Error(fmt.Sprintf("validation error: %s", err.Error())))
		return
	}

	res, err := c.uc.Search(ctx, params)
	if err != nil {
		slog.Error(
			"failed to execute search query",
			slog.String("uri", ctx.Request.RequestURI),
			slog.Any("params", params),
			slog.Any("error", err),
		)
		ctx.JSON(http.StatusInternalServerError, c.pr.Error("internal server error"))
		return
	}

	result := c.pr.Format(params, res, int(time.Since(startTime).Milliseconds()))
	slog.Info(
		"querylog",
		slog.String("domain", "search/problem"),
		slog.Int("elapsed_time", result.Stats.Time),
		slog.Int("hits", res.Response.NumFound),
		slog.Any("params", params),
	)
	ctx.JSON(http.StatusOK, result)
}

func (c *searchProblemController) HandlePOST(ctx *gin.Context) {
	startTime := time.Now()

	var params domain.SearchProblemParam
	if err := ctx.ShouldBindJSON(&params); err != nil {
		slog.Error(
			"failed to bind request body",
			slog.String("uri", ctx.Request.RequestURI),
			slog.Any("error", err),
		)
		ctx.JSON(http.StatusBadRequest, c.pr.Error("failed to bind request body"))
		return
	}

	if err := c.validator.Struct(params); err != nil {
		slog.Error(
			"validation error",
			slog.String("uri", ctx.Request.RequestURI),
			slog.Any("params", params),
			slog.Any("error", err),
		)
		ctx.JSON(http.StatusBadRequest, c.pr.Error(fmt.Sprintf("validation error: %s", err.Error())))
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

	result := c.pr.Format(params, res, int(time.Since(startTime).Milliseconds()))
	slog.Info(
		"querylog",
		slog.String("domain", "search/problem"),
		slog.Int("elapsed_time", result.Stats.Time),
		slog.Int("hits", res.Response.NumFound),
		slog.Any("params", params),
	)
	ctx.JSON(http.StatusOK, result)
}

type RecommendProblemController interface {
	HandleGET(*gin.Context)
	HandlePOST(*gin.Context)
}

type recommendProblemController struct {
	uc usecase.RecommendProblemUsecase
	pr presenter.RecommendProblemPresenter
}

func NewRecommendProblemController(
	uc usecase.RecommendProblemUsecase,
	pr presenter.RecommendProblemPresenter,
) RecommendProblemController {
	return &recommendProblemController{
		uc: uc,
		pr: pr,
	}
}

func (c *recommendProblemController) HandleGET(ctx *gin.Context) {
	startTime := time.Now()

	var params domain.RecommendProblemParam
	if err := ctx.ShouldBindQuery(&params); err != nil {
		slog.Error(
			"failed to decode request parameter",
			slog.Any("error", err),
		)
		ctx.JSON(http.StatusBadRequest, c.pr.Error("failed to decode request parameter"))
		return
	}

	res, err := c.uc.Recommend(ctx, params)
	if err != nil {
		slog.Error(
			"failed to execute search query",
			slog.Any("error", err),
		)
		ctx.JSON(http.StatusInternalServerError, c.pr.Error("internal server error"))
		return
	}

	result := c.pr.Format(params, res, int(time.Since(startTime).Milliseconds()))
	slog.Info(
		"querylog",
		slog.String("domain", "recommend/problem"),
		slog.Int("elapsed_time", result.Stats.Time),
		slog.Int("hits", res.Response.NumFound),
		slog.Any("params", params),
	)
	ctx.JSON(http.StatusOK, result)
}

func (c *recommendProblemController) HandlePOST(ctx *gin.Context) {
	startTime := time.Now()

	var params domain.RecommendProblemParam
	if err := ctx.ShouldBindJSON(&params); err != nil {
		slog.Error(
			"failed to decode request parameter",
			slog.Any("error", err),
		)
		ctx.JSON(http.StatusBadRequest, c.pr.Error("failed to decode request parameter"))
		return
	}

	res, err := c.uc.Recommend(ctx, params)
	if err != nil {
		slog.Error(
			"failed to execute search query",
			slog.Any("error", err),
		)
		ctx.JSON(http.StatusInternalServerError, c.pr.Error("internal server error"))
		return
	}

	result := c.pr.Format(params, res, int(time.Since(startTime).Milliseconds()))
	slog.Info(
		"querylog",
		slog.String("domain", "recommend/problem"),
		slog.Int("elapsed_time", result.Stats.Time),
		slog.Int("hits", res.Response.NumFound),
		slog.Any("params", params),
	)
	ctx.JSON(http.StatusOK, result)
}

type ProblemListController interface {
	HandleGET(*gin.Context)
}

type problemListController struct {
	repo repository.ProblemRepository
}

func NewProblemListController(repo repository.ProblemRepository) ProblemListController {
	return &problemListController{
		repo: repo,
	}
}

func (c *problemListController) HandleGET(ctx *gin.Context) {
	targets := strings.Split(ctx.Query("contest_id"), ",")

	ids, err := c.repo.FetchIDsByContestID(ctx, targets)
	if err != nil {
		slog.Error("failed to fetch problem ids", slog.String("contest id", ctx.Query("contest_id")), slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, make([]string, 0))
		return
	}

	ctx.JSON(http.StatusOK, ids)
}
