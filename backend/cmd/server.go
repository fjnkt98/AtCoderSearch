package cmd

import (
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/repository"
	"fjnkt98/atcodersearch/server"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

func newServerCmd(args []string, config *RootConfig, runFunc func(cmd *cobra.Command, args []string)) *cobra.Command {
	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "Launch API server",
		Long:  "Launch API server",
		PreRun: func(cmd *cobra.Command, args []string) {
			MustLoadConfig(cmd.Flags(), config)
		},
		Run: func(cmd *cobra.Command, args []string) {
			db := repository.MustGetDB(config.DataBaseURL)

			r := gin.New()
			r.Use(
				gin.Recovery(),
			)
			r.SetTrustedProxies(config.TrustedProxies)

			problemCore := solr.MustNewSolrCore(config.SolrHost, config.ProblemCoreName)
			userCore := solr.MustNewSolrCore(config.SolrHost, config.UserCoreName)
			submissionCore := solr.MustNewSolrCore(config.SolrHost, config.SubmissionCoreName)

			server.RegisterSearchProblemRoute(r, problemCore)
			server.RegisterSearchUserRoute(r, userCore)
			server.RegisterSearchSubmissionRoute(r, submissionCore)
			server.RegisterRecommendProblemRoute(r, problemCore, db)
			server.RegisterCategoryListRoute(r, db)
			server.RegisterContestListRoute(r, db)
			server.RegisterLanguageListRoute(r, db)
			server.RegisterLanguageGroupListRoute(r, db)
			server.RegisterProblemListRoute(r, db)
			server.RegisterLivenessRoute(r, []solr.SolrCore{problemCore, userCore, submissionCore})

			r.Run("localhost:8000")
		},
	}

	serverCmd.SetArgs(args)
	if runFunc != nil {
		serverCmd.Run = runFunc
	}

	return serverCmd
}
