package server

import (
	"fjnkt98/atcodersearch/batch/repository"
	"fjnkt98/atcodersearch/config"
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/server/problem"
	"fjnkt98/atcodersearch/server/recommend"
	"fjnkt98/atcodersearch/server/submission"
	"fjnkt98/atcodersearch/server/user"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator"
	"github.com/uptrace/bun"
)

func RegisterProblemRoute(r *gin.Engine) {
	core, err := solr.NewSolrCore(config.Config.CommonConfig.SolrHost, "problem")
	if err != nil {
		slog.Error(
			"failed to initialize problem core",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	c := problem.NewProblemController(
		problem.NewProblemUsecase(core),
		problem.NewProblemPresenter(),
	)
	r.GET("/api/search/problem", c.HandleGET)
	r.POST("/api/search/problem", c.HandlePOST)
}

func RegisterUserRoute(r *gin.Engine) {
	core, err := solr.NewSolrCore(config.Config.CommonConfig.SolrHost, "user")
	if err != nil {
		slog.Error(
			"failed to initialize user core",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	c := user.NewUserController(
		user.NewUserUsecase(core),
		user.NewUserPresenter(),
	)
	r.GET("/api/search/user", c.HandleGET)
	r.POST("/api/search/user", c.HandlePOST)
}

func RegisterSubmissionRoute(r *gin.Engine) {
	core, err := solr.NewSolrCore(config.Config.CommonConfig.SolrHost, "submission")
	if err != nil {
		slog.Error(
			"failed to initialize submission core",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	c := submission.NewSubmissionController(
		submission.NewSubmissionUsecase(core),
		submission.NewSubmissionPresenter(),
	)
	r.GET("/api/search/submission", c.HandleGET)
	r.POST("/api/search/submission", c.HandlePOST)
}

func RegisterRecommendRoute(r *gin.Engine, db *bun.DB) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("model", recommend.ValidateModel)
		v.RegisterValidation("option", recommend.ValidateOption)
	}

	core, err := solr.NewSolrCore(config.Config.CommonConfig.SolrHost, "problem")
	if err != nil {
		slog.Error(
			"failed to initialize recommend core",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	c := recommend.NewRecommendController(
		recommend.NewRecommendUsecase(
			core,
			repository.NewUserRepository(db),
		),
		recommend.NewRecommendPresenter(),
	)
	r.GET("/api/recommend/problem", c.HandleGET)
	r.POST("/api/recommend/problem", c.HandlePOST)
}
