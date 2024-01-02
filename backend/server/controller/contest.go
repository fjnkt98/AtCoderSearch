package controller

import (
	"fjnkt98/atcodersearch/repository"
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
	categories := strings.Split(ctx.Query("category"), ",")

	ids, err := c.repo.FetchContestIDs(ctx, categories)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, make([]string, 0))
		return
	}

	ctx.JSON(http.StatusOK, ids)
}
