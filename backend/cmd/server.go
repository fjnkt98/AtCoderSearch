package cmd

import (
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/repository"
	"fjnkt98/atcodersearch/server"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

func newServerCmd(args []string, runFunc func(cmd *cobra.Command, args []string)) *cobra.Command {
	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "Launch API server",
		Long:  "Launch API server",
		Run: func(cmd *cobra.Command, args []string) {
			db := repository.MustGetDB(Config.DataBaseURL)

			r := gin.New()
			r.Use(
				gin.Recovery(),
			)

			server.RegisterSearchProblemRoute(r, solr.MustNewSolrCore(Config.SolrHost, Config.ProblemCoreName))
			server.RegisterSearchUserRoute(r, solr.MustNewSolrCore(Config.SolrHost, Config.UserCoreName))
			server.RegisterSearchSubmissionRoute(r, solr.MustNewSolrCore(Config.SolrHost, Config.SubmissionCoreName))
			server.RegisterRecommendProblemRoute(r, solr.MustNewSolrCore(Config.SolrHost, Config.ProblemCoreName), db)
			server.RegisterCategoryListRoute(r, db)
			server.RegisterContestListRoute(r, db)
			server.RegisterLanguageListRoute(r, db)
			server.RegisterLanguageGroupListRoute(r, db)
			server.RegisterProblemListRoute(r, db)

			r.Run("localhost:8000")
		},
	}

	serverCmd.SetArgs(args)
	if runFunc != nil {
		serverCmd.Run = runFunc
	}

	return serverCmd
}
