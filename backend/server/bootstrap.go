package server

import (
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/repository"
	"fjnkt98/atcodersearch/server/controller"
	"fjnkt98/atcodersearch/server/domain"
	"fjnkt98/atcodersearch/server/presenter"
	"fjnkt98/atcodersearch/server/usecase"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator"
	"github.com/uptrace/bun"
)

func RegisterSearchProblemRoute(r *gin.Engine, core solr.SolrCore) {
	c := controller.NewSearchProblemController(
		usecase.NewSearchProblemUsecase(core),
		presenter.NewSearchProblemPresenter(),
	)
	r.GET("/api/search/problem", c.HandleGET)
	r.POST("/api/search/problem", c.HandlePOST)
}

func RegisterSearchUserRoute(r *gin.Engine, core solr.SolrCore) {
	c := controller.NewSearchUserController(
		usecase.NewSearchUserUsecase(core),
		presenter.NewSearchUserPresenter(),
	)
	r.GET("/api/search/user", c.HandleGET)
	r.POST("/api/search/user", c.HandlePOST)
}

func RegisterSearchSubmissionRoute(r *gin.Engine, core solr.SolrCore) {
	c := controller.NewSearchSubmissionController(
		usecase.NewSearchSubmissionUsecase(core),
		presenter.NewSearchSubmissionPresenter(),
	)
	r.GET("/api/search/submission", c.HandleGET)
	r.POST("/api/search/submission", c.HandlePOST)
}

func RegisterRecommendProblemRoute(r *gin.Engine, core solr.SolrCore, db *bun.DB) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("model", domain.ValidateModel)
		v.RegisterValidation("option", domain.ValidateOption)
	}

	c := controller.NewRecommendProblemController(
		usecase.NewRecommendProblemUsecase(
			core,
			repository.NewUserRepository(db),
		),
		presenter.NewRecommendProblemPresenter(),
	)
	r.GET("/api/recommend/problem", c.HandleGET)
	r.POST("/api/recommend/problem", c.HandlePOST)
}

func RegisterCategoryListRoute(r *gin.Engine, db *bun.DB) {
	c := controller.NewCategoryListController(
		repository.NewContestRepository(db),
	)

	r.GET("/api/list/category", c.HandleGET)
}

func RegisterContestListRoute(r *gin.Engine, db *bun.DB) {
	c := controller.NewContestListController(
		repository.NewContestRepository(db),
	)

	r.GET("/api/list/contest", c.HandleGET)
}

func RegisterLanguageListRoute(r *gin.Engine, db *bun.DB) {
	c := controller.NewLanguageListController(
		repository.NewLanguageRepository(db),
	)

	r.GET("/api/list/language", c.HandleGET)
}

func RegisterLanguageGroupListRoute(r *gin.Engine, db *bun.DB) {
	c := controller.NewLanguageGroupListController(
		repository.NewLanguageRepository(db),
	)

	r.GET("/api/list/language/group", c.HandleGET)
}

func RegisterProblemListRoute(r *gin.Engine, db *bun.DB) {
	c := controller.NewProblemListController(
		repository.NewProblemRepository(db),
	)

	r.GET("/api/list/problem", c.HandleGET)
}
