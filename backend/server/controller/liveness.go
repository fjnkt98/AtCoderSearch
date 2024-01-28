package controller

import (
	"fjnkt98/atcodersearch/pkg/solr"

	"github.com/gin-gonic/gin"
)

type LivenessController interface {
	HandleGET(*gin.Context)
}

type livenessController struct {
	cores []solr.SolrCore
}

func (c *livenessController) HandleGET(ctx *gin.Context) {
	isAlive := true

	for _, core := range c.cores {
		res, err := solr.PingWithContext(ctx, core)
		if err != nil {
			isAlive = false
		} else if res.Status != "OK" {
			isAlive = false
		}
	}

	if isAlive {
		ctx.String(200, "alive")
	} else {
		ctx.String(500, "not alive")
	}
}

func NewLivenessController(cores []solr.SolrCore) LivenessController {
	return &livenessController{
		cores: cores,
	}
}
