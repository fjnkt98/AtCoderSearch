package cmd

import (
	"fjnkt98/atcodersearch/problem"
	"fjnkt98/atcodersearch/recommend"
	"fjnkt98/atcodersearch/submission"
	"fjnkt98/atcodersearch/user"
	"fmt"
	"net/http"
	"net/url"
	"os"

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
			slog.Error("invalid Solr URL was given", slog.String("host", solrHost))
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

		// Recommend searcher configuration
		recommendSearcher, err := recommend.NewSearcher(solrBaseURL.String(), "problem")
		if err != nil {
			slog.Error("failed to instantiate recommend searcher", slog.String("error", fmt.Sprintf("%+v", err)))
		}

		// API handler registration
		http.HandleFunc("/api/search/problem", problemSearcher.HandleSearch)
		http.HandleFunc("/api/search/user", userSearcher.HandleSearch)
		http.HandleFunc("/api/search/submission", submissionSearcher.HandleSearch)
		http.HandleFunc("/api/recommend/problem", recommendSearcher.HandleSearch)
		http.HandleFunc("/api/liveness", func(w http.ResponseWriter, r *http.Request) {
			if problemSearcher.Liveness() && userSearcher.Liveness() {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		})
		http.HandleFunc("/api/readiness", func(w http.ResponseWriter, r *http.Request) {
			if problemSearcher.Readiness() && userSearcher.Readiness() {
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
		if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
			slog.Error("failed to listen and serve api server", slog.String("error", err.Error()))
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
