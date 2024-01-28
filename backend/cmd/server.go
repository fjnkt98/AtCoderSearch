package cmd

import (
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/repository"
	"fjnkt98/atcodersearch/server"
	"time"

	"github.com/gin-contrib/cache/persistence"
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

			store := persistence.NewInMemoryStore(10 * time.Second)

			server.RegisterSearchProblemRoute(r, problemCore)
			server.RegisterSearchUserRoute(r, userCore)
			server.RegisterSearchSubmissionRoute(r, submissionCore)
			server.RegisterRecommendProblemRoute(r, problemCore, db)
			server.RegisterCategoryListRoute(r, db, store)
			server.RegisterContestListRoute(r, db, store)
			server.RegisterLanguageListRoute(r, db, store)
			server.RegisterLanguageGroupListRoute(r, db, store)
			server.RegisterProblemListRoute(r, db, store)
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
