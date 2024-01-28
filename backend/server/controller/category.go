package controller

import (
	"fjnkt98/atcodersearch/repository"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryListController interface {
	HandleGET(*gin.Context)
}

type categoryListController struct {
	repo repository.ContestRepository
}

func NewCategoryListController(repo repository.ContestRepository) CategoryListController {
	return &categoryListController{
		repo: repo,
	}
}

func (c *categoryListController) HandleGET(ctx *gin.Context) {
	categories, err := c.repo.FetchCategories(ctx)
	if err != nil {
		slog.Error("failed to fetch categories", slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, make([]string, 0))
		return
	}

	ctx.JSON(http.StatusOK, categories)
}
