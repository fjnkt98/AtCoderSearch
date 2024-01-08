package controller

import (
	"fjnkt98/atcodersearch/repository"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LanguageListController interface {
	HandleGET(*gin.Context)
}

type languageListController struct {
	repo repository.LanguageRepository
}

func NewLanguageListController(
	repo repository.LanguageRepository,
) LanguageListController {
	return &languageListController{
		repo: repo,
	}
}

func (c *languageListController) HandleGET(ctx *gin.Context) {
	languages, err := c.repo.FetchLanguages(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, make([]string, 0))
		slog.Error("failed to fetch languages", slog.Any("error", err))
		return
	}
	ctx.JSON(http.StatusOK, languages)
}

type LanguageGroupListController interface {
	HandleGET(*gin.Context)
}

type languageGroupListController struct {
	repo repository.LanguageRepository
}

func NewLanguageGroupListController(
	repo repository.LanguageRepository,
) LanguageListController {
	return &languageGroupListController{
		repo: repo,
	}
}

func (c *languageGroupListController) HandleGET(ctx *gin.Context) {
	groups, err := c.repo.FetchLanguageGroups(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, make([]string, 0))
		slog.Error("failed to fetch language groups", slog.Any("error", err))
		return
	}
	ctx.JSON(http.StatusOK, groups)
}
