package server

import (
	"fjnkt98/atcodersearch/config"
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/repository"
	"fjnkt98/atcodersearch/server/controller"
	"fjnkt98/atcodersearch/server/domain"
	"fjnkt98/atcodersearch/server/presenter"
	"fjnkt98/atcodersearch/server/usecase"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator"
	"github.com/uptrace/bun"
)

func NewRouter(db *bun.DB) *gin.Engine {
	r := gin.New()
	r.Use(
		gin.Recovery(),
	)

	RegisterSearchProblemRoute(r)
	RegisterSearchUserRoute(r)
	RegisterSearchSubmissionRoute(r)
	RegisterRecommendProblemRoute(r, db)
	RegisterCategoryListRoute(r, db)
	RegisterContestListRoute(r, db)
	RegisterLanguageListRoute(r, db)
	RegisterLanguageGroupListRoute(r, db)
	RegisterProblemListRoute(r, db)

	return r
}

func RegisterSearchProblemRoute(r *gin.Engine) {
	core, err := solr.NewSolrCore(config.Config.CommonConfig.SolrHost, "problem")
	if err != nil {
		slog.Error(
			"failed to initialize problem core",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	c := controller.NewSearchProblemController(
		usecase.NewSearchProblemUsecase(core),
		presenter.NewSearchProblemPresenter(),
	)
	r.GET("/api/search/problem", c.HandleGET)
	r.POST("/api/search/problem", c.HandlePOST)
}

func RegisterSearchUserRoute(r *gin.Engine) {
	core, err := solr.NewSolrCore(config.Config.CommonConfig.SolrHost, "user")
	if err != nil {
		slog.Error(
			"failed to initialize user core",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	c := controller.NewSearchUserController(
		usecase.NewSearchUserUsecase(core),
		presenter.NewSearchUserPresenter(),
	)
	r.GET("/api/search/user", c.HandleGET)
	r.POST("/api/search/user", c.HandlePOST)
}

func RegisterSearchSubmissionRoute(r *gin.Engine) {
	core, err := solr.NewSolrCore(config.Config.CommonConfig.SolrHost, "submission")
	if err != nil {
		slog.Error(
			"failed to initialize submission core",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	c := controller.NewSearchSubmissionController(
		usecase.NewSearchSubmissionUsecase(core),
		presenter.NewSubmissionPresenter(),
	)
	r.GET("/api/search/submission", c.HandleGET)
	r.POST("/api/search/submission", c.HandlePOST)
}

func RegisterRecommendProblemRoute(r *gin.Engine, db *bun.DB) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("model", domain.ValidateModel)
		v.RegisterValidation("option", domain.ValidateOption)
	}

	core, err := solr.NewSolrCore(config.Config.CommonConfig.SolrHost, "problem")
	if err != nil {
		slog.Error(
			"failed to initialize problem core",
			slog.Any("error", err),
		)
		os.Exit(1)
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
