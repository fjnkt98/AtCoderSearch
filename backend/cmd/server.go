package cmd

import (
	"fjnkt98/atcodersearch/atcodersearch/problem"
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
		solrURL := os.Getenv("SOLR_URL")
		parsedSolrURL, err := url.Parse(solrURL)
		if err != nil {
			log.Fatalf("invalid Solr URL was given: %s\n", solrURL)
		}
		solrBaseURL := url.URL{Scheme: parsedSolrURL.Scheme, Host: parsedSolrURL.Host}

		problemCoreName := os.Getenv("PROBLEMS_CORE_NAME")
		if problemCoreName == "" {
			log.Fatalf("PROBLEMS_CORE_NAME must be set.")
		}

		problemSearcher, err := problem.NewProblemSearcher(solrBaseURL.String(), problemCoreName)
		if err != nil {
			log.Fatalf("failed to instantiate problems searcher: %s", err.Error())
		}

		http.HandleFunc("/api/search", problemSearcher.HandleSearchProblem)

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
