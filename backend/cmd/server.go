package cmd

import (
	"fjnkt98/atcodersearch/list"
	"fjnkt98/atcodersearch/problem"
	"fjnkt98/atcodersearch/recommend"
	"fjnkt98/atcodersearch/submission"
	"fjnkt98/atcodersearch/user"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/coocood/freecache"
	cache "github.com/gitsight/go-echo-cache"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Launch API server",
	Long:  "Launch API server",
	Run: func(cmd *cobra.Command, args []string) {
		// Solr common configuration
		solrHost := os.Getenv("SOLR_HOST")
		if solrHost == "" {
			slog.Error("SOLR_HOST must be set")
			os.Exit(1)
		}
		parsedSolrURL, err := url.Parse(solrHost)
		if err != nil {
			slog.Error("invalid Solr URL was given", slog.String("host", solrHost), slog.String("error", err.Error()))
			os.Exit(1)
		}
		solrBaseURL := url.URL{Scheme: parsedSolrURL.Scheme, Host: parsedSolrURL.Host}

		// Problem searcher configuration
		problemSearcher, err := problem.NewSearcher(solrBaseURL.String(), "problem")
		if err != nil {
			slog.Error("failed to instantiate problem searcher", slog.String("error", fmt.Sprintf("%+v", err)))
		}

		// User searcher configuration
		userSearcher, err := user.NewSearcher(solrBaseURL.String(), "user")
		if err != nil {
			slog.Error("failed to instantiate user searcher", slog.String("error", fmt.Sprintf("%+v", err)))
		}

		// Submission searcher configuration
		submissionSearcher, err := submission.NewSearcher(solrBaseURL.String(), "submission")
		if err != nil {
			slog.Error("failed to instantiate submission searcher", slog.String("error", fmt.Sprintf("%+v", err)))
		}

		db := GetDB()

		// Recommend searcher configuration
		recommendSearcher, err := recommend.NewSearcher(solrBaseURL.String(), "problem", db)
		if err != nil {
			slog.Error("failed to instantiate recommend searcher", slog.String("error", fmt.Sprintf("%+v", err)))
		}

		listSearcher := list.NewSearcher(db)

		c := freecache.NewCache(100 * 1024 * 1024)

		e := echo.New()
		e.Use(middleware.Recover())
		e.Use(cache.New(&cache.Config{}, c))

		// API handler registration
		e.GET("/api/search/problem", problemSearcher.HandleGET)
		e.GET("/api/search/user", userSearcher.HandleGET)
		e.GET("/api/search/submission", submissionSearcher.HandleGET)
		e.GET("/api/recommend/problem", recommendSearcher.HandleGET)
		e.GET("/api/list/category", listSearcher.HandleCategory)
		e.GET("/api/list/language", listSearcher.HandleLanguage)
		e.GET("/api/list/contest", listSearcher.HandleContest)
		e.GET("/api/list/problem", listSearcher.HandleProblem)

		http.HandleFunc("/api/liveness", func(w http.ResponseWriter, r *http.Request) {
			if problemSearcher.Liveness() && userSearcher.Liveness() && submissionSearcher.Liveness() && recommendSearcher.Liveness() {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		})
		http.HandleFunc("/api/readiness", func(w http.ResponseWriter, r *http.Request) {
			if problemSearcher.Readiness() && userSearcher.Readiness() && submissionSearcher.Readiness() && recommendSearcher.Readiness() {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		})

		// Launch server
		port := os.Getenv("API_SERVER_LISTEN_PORT")
		if port == "" {
			port = "8000"
			slog.Warn("Environment variable API_SERVER_LISTEN_PORT was not given. Default port number 8000 will be used.")
		} else {
			slog.Info(fmt.Sprintf("API server will listen on %s", port))
		}

		e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
