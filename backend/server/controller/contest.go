package controller

import (
	"fjnkt98/atcodersearch/repository"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ContestListController interface {
	HandleGET(*gin.Context)
}

type contestListController struct {
	repo repository.ContestRepository
}

func NewContestListController(repo repository.ContestRepository) ContestListController {
	return &contestListController{
		repo: repo,
	}
}

func (c *contestListController) HandleGET(ctx *gin.Context) {
	categories := make([]string, 0)
	for _, category := range strings.Split(ctx.Query("category"), ",") {
		if category != "" {
			categories = append(categories, category)
		}
	}

	ids, err := c.repo.FetchContestIDs(ctx, categories)
	if err != nil {
		slog.Error("failed to fetch contests", slog.String("category", ctx.Query("category")), slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, make([]string, 0))
		return
	}

	ctx.JSON(http.StatusOK, ids)
}
