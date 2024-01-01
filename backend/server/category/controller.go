package category

import (
	"fjnkt98/atcodersearch/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryController interface {
	HandleGET(*gin.Context)
}

type contestController struct {
	repo repository.ContestRepository
}

func NewCategoryController(repo repository.ContestRepository) CategoryController {
	return &contestController{
		repo: repo,
	}
}

func (c *contestController) HandleGET(ctx *gin.Context) {
	categories, err := c.repo.FetchCategories(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, make([]string, 0))
		return
	}

	ctx.JSON(http.StatusOK, categories)
}
