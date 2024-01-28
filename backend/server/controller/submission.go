package controller

import (
	"fjnkt98/atcodersearch/server/domain"
	"fjnkt98/atcodersearch/server/presenter"
	"fjnkt98/atcodersearch/server/usecase"
	"fjnkt98/atcodersearch/server/utility"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	govalidator "github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
)

type SearchSubmissionController interface {
	HandleGET(*gin.Context)
	HandlePOST(*gin.Context)
}

type searchSubmissionController struct {
	uc        usecase.SearchSubmissionUsecase
	pr        presenter.SearchSubmissionPresenter
	validator *govalidator.Validate
	decoder   *schema.Decoder
}

func NewSearchSubmissionController(
	uc usecase.SearchSubmissionUsecase,
	pr presenter.SearchSubmissionPresenter,
) SearchSubmissionController {
	decoder := utility.NewSchemaDecoder()

	validator := govalidator.New()

	return &searchSubmissionController{
		uc:        uc,
		pr:        pr,
		decoder:   decoder,
		validator: validator,
	}
}

func (c *searchSubmissionController) HandleGET(ctx *gin.Context) {
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

	var params domain.SearchSubmissionParam
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
		slog.String("domain", "search/submission"),
		slog.Int("elapsed_time", result.Stats.Time),
		slog.Int("hits", res.Response.NumFound),
		slog.Any("params", params),
	)
	ctx.JSON(http.StatusOK, result)
}

func (c *searchSubmissionController) HandlePOST(ctx *gin.Context) {
	startTime := time.Now()

	var params domain.SearchSubmissionParam
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
		slog.String("domain", "search/submission"),
		slog.Int("elapsed_time", result.Stats.Time),
		slog.Int("hits", res.Response.NumFound),
		slog.Any("params", params),
	)
	ctx.JSON(http.StatusOK, result)
}
