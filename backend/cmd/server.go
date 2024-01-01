package cmd

import (
	"fjnkt98/atcodersearch/config"
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/server/problem"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Launch API server",
	Long:  "Launch API server",
	Run: func(cmd *cobra.Command, args []string) {
		r := gin.New()
		r.Use(
			gin.Recovery(),
		)

		{
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

		r.Run("localhost:8000")

		// // Solr common configuration
		// solrHost := os.Getenv("SOLR_HOST")
		// if solrHost == "" {
		// 	slog.Error("SOLR_HOST must be set")
		// 	os.Exit(1)
		// }
		// parsedSolrURL, err := url.Parse(solrHost)
		// if err != nil {
		// 	slog.Error("invalid Solr URL was given", slog.String("host", solrHost), slog.String("error", err.Error()))
		// 	os.Exit(1)
		// }
		// solrBaseURL := url.URL{Scheme: parsedSolrURL.Scheme, Host: parsedSolrURL.Host}

		// // Problem searcher configuration
		// problemSearcher, err := problem.NewSearcher(solrBaseURL.String(), "problem")
		// if err != nil {
		// 	slog.Error("failed to instantiate problem searcher", slog.String("error", fmt.Sprintf("%+v", err)))
		// }

		// // User searcher configuration
		// userSearcher, err := user.NewSearcher(solrBaseURL.String(), "user")
		// if err != nil {
		// 	slog.Error("failed to instantiate user searcher", slog.String("error", fmt.Sprintf("%+v", err)))
		// }

		// // Submission searcher configuration
		// submissionSearcher, err := submission.NewSearcher(solrBaseURL.String(), "submission")
		// if err != nil {
		// 	slog.Error("failed to instantiate submission searcher", slog.String("error", fmt.Sprintf("%+v", err)))
		// }

		// db := GetDB()

		// // Recommend searcher configuration
		// recommendSearcher, err := recommend.NewSearcher(solrBaseURL.String(), "problem", db)
		// if err != nil {
		// 	slog.Error("failed to instantiate recommend searcher", slog.String("error", fmt.Sprintf("%+v", err)))
		// }

		// listSearcher := list.NewSearcher(db)

		// e := echo.New()
		// e.Use(middleware.Recover())
		// e.Use(middleware.Gzip())
		// e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// 	AllowMethods: []string{
		// 		http.MethodGet,
		// 		http.MethodPost,
		// 	},
		// 	AllowHeaders: []string{
		// 		echo.HeaderOrigin,
		// 	},
		// 	AllowOrigins: []string{
		// 		"https://atcoder-search.fjnkt98.com",
		// 	},
		// }))
		// e.HideBanner = true
		// e.HidePort = true

		// // API handler registration
		// e.GET("/api/search/problem", problemSearcher.HandleGET)
		// e.POST("/api/search/problem", problemSearcher.HandlePOST)
		// e.GET("/api/search/user", userSearcher.HandleGET)
		// e.POST("/api/search/user", userSearcher.HandlePOST)
		// e.GET("/api/search/submission", submissionSearcher.HandleGET)
		// e.POST("/api/search/submission", submissionSearcher.HandlePOST)
		// e.GET("/api/recommend/problem", recommendSearcher.HandleGET)
		// e.POST("/api/recommend/problem", recommendSearcher.HandlePOST)
		// e.GET("/api/list/category", listSearcher.HandleCategory)
		// e.GET("/api/list/language", listSearcher.HandleLanguage)
		// e.GET("/api/list/language/group", listSearcher.HandleLanguageGroup)
		// e.GET("/api/list/contest", listSearcher.HandleContest)
		// e.GET("/api/list/problem", listSearcher.HandleProblem)

		// e.GET("/api/liveness", func(c echo.Context) error {
		// 	if problemSearcher.Liveness() && userSearcher.Liveness() && submissionSearcher.Liveness() {
		// 		return c.String(http.StatusOK, "")
		// 	} else {
		// 		return c.String(http.StatusInternalServerError, "")
		// 	}
		// })

		// e.GET("/api/readiness", func(c echo.Context) error {
		// 	if problemSearcher.Readiness() && userSearcher.Readiness() && submissionSearcher.Readiness() {
		// 		return c.String(http.StatusOK, "")
		// 	} else {
		// 		return c.String(http.StatusInternalServerError, "")
		// 	}
		// })

		// // Launch server
		// port := os.Getenv("API_SERVER_LISTEN_PORT")
		// if port == "" {
		// 	port = "8000"
		// 	slog.Warn("Environment variable API_SERVER_LISTEN_PORT was not given. Default port number 8000 will be used.")
		// } else {
		// 	slog.Info(fmt.Sprintf("API server will listen on %s", port))
		// }

		// go func() {
		// 	if err := e.Start(fmt.Sprintf(":%s", port)); err != nil && err != http.ErrServerClosed {
		// 		e.Logger.Fatal("shutting down the server")
		// 	}
		// }()

		// quit := make(chan os.Signal, 1)
		// signal.Notify(quit, os.Interrupt)
		// <-quit
		// slog.Info("shutting down the server gracefully")
		// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		// defer cancel()
		// if err := e.Shutdown(ctx); err != nil {
		// 	e.Logger.Fatal(err)
		// }
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
