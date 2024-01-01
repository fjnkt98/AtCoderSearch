package language

import (
	"fjnkt98/atcodersearch/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LanguageController interface {
	HandleGET(*gin.Context)
}

type languageController struct {
	repo repository.LanguageRepository
}

func NewLanguageController(
	repo repository.LanguageRepository,
) LanguageController {
	return &languageController{
		repo: repo,
	}
}

func (c *languageController) HandleGET(ctx *gin.Context) {
	languages, err := c.repo.FetchLanguages(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, make([]string, 0))
		return
	}
	ctx.JSON(http.StatusOK, languages)
}

type LanguageGroupController interface {
	HandleGET(*gin.Context)
}

type languageGroupController struct {
	repo repository.LanguageRepository
}

func NewLanguageGroupController(
	repo repository.LanguageRepository,
) LanguageController {
	return &languageGroupController{
		repo: repo,
	}
}

func (c *languageGroupController) HandleGET(ctx *gin.Context) {
	groups, err := c.repo.FetchLanguageGroups(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, make([]string, 0))
		return
	}
	ctx.JSON(http.StatusOK, groups)
}
