package cmd

import (
	"fjnkt98/atcodersearch/atcodersearch/problem"
	"fjnkt98/atcodersearch/atcodersearch/user"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Launch API server",
	Long:  "Launch API server",
	Run: func(cmd *cobra.Command, args []string) {
		// Solr common configuration
		solrHost := os.Getenv("SOLR_HOST")
		if solrHost == "" {
			log.Fatalln("SOLR_HOST must be set.")
		}
		parsedSolrURL, err := url.Parse(solrHost)
		if err != nil {
			log.Fatalf("invalid Solr URL was given: %s", solrHost)
		}
		solrBaseURL := url.URL{Scheme: parsedSolrURL.Scheme, Host: parsedSolrURL.Host}

		// Problem searcher configuration
		problemCoreName := os.Getenv("PROBLEM_CORE_NAME")
		if problemCoreName == "" {
			log.Fatalln("PROBLEM_CORE_NAME must be set.")
		}
		problemSearcher, err := problem.NewSearcher(solrBaseURL.String(), problemCoreName)
		if err != nil {
			log.Fatalf("failed to instantiate problem searcher: %s", err.Error())
		}

		// User searcher configuration
		userCoreName := os.Getenv("USER_CORE_NAME")
		if userCoreName == "" {
			log.Fatalln("USER_CORE_NAME must be set.")
		}
		userSearcher, err := user.NewSearcher(solrBaseURL.String(), userCoreName)
		if err != nil {
			log.Fatalf("failed to instantiate user searcher: %s", err.Error())
		}

		// API handler registration
		http.HandleFunc("/api/search/problem", problemSearcher.HandleSearch)
		http.HandleFunc("/api/search/user", userSearcher.HandleSearch)
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
			log.Println("Environment variable API_SERVER_LISTEN_PORT was not given. Default port number 8000 will be used.")
		} else {
			log.Printf("API server will listen on %s", port)
		}
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
