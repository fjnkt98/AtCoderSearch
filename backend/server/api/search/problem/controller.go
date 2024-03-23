package problem

import (
	"fjnkt98/atcodersearch/pkg/solr"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SearchProblemController struct {
	core *solr.SolrCore
}

func NewSearchProblemController(core *solr.SolrCore) *SearchProblemController {
	return &SearchProblemController{
		core: core,
	}
}

func (c *SearchProblemController) SearchProblem(ctx echo.Context) error {
	var param Parameter
	if err := ctx.Bind(&param); err != nil {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	if err := ctx.Validate(param); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	return ctx.JSON(http.StatusOK, param)
}
