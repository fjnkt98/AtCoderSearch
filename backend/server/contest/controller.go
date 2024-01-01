package contest

import (
	"fjnkt98/atcodersearch/repository"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ContestController interface {
	HandleGET(*gin.Context)
}

type contestController struct {
	repo repository.ContestRepository
}

func NewContestController(repo repository.ContestRepository) ContestController {
	return &contestController{
		repo: repo,
	}
}

func (c *contestController) HandleGET(ctx *gin.Context) {
	categories := strings.Split(ctx.Query("category"), ",")

	ids, err := c.repo.FetchContestIDs(ctx, categories)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, make([]string, 0))
		return
	}

	ctx.JSON(http.StatusOK, ids)
}
