package cmd

import (
	"fjnkt98/atcodersearch/config"
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/server/problem"
	"fjnkt98/atcodersearch/server/submission"
	"fjnkt98/atcodersearch/server/user"
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

		// Register problem search api handlers
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
		// Register user search api handlers
		{
			core, err := solr.NewSolrCore(config.Config.CommonConfig.SolrHost, "user")
			if err != nil {
				slog.Error(
					"failed to initialize user core",
					slog.Any("error", err),
				)
				os.Exit(1)
			}

			c := user.NewUserController(
				user.NewUserUsecase(core),
				user.NewUserPresenter(),
			)
			r.GET("/api/search/user", c.HandleGET)
			r.POST("/api/search/user", c.HandlePOST)
		}

		// Register user search api handlers
		{
			core, err := solr.NewSolrCore(config.Config.CommonConfig.SolrHost, "submission")
			if err != nil {
				slog.Error(
					"failed to initialize submission core",
					slog.Any("error", err),
				)
				os.Exit(1)
			}

			c := submission.NewSubmissionController(
				submission.NewSubmissionUsecase(core),
				submission.NewSubmissionPresenter(),
			)
			r.GET("/api/search/submission", c.HandleGET)
			r.POST("/api/search/submission", c.HandlePOST)
		}

		r.Run("localhost:8000")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
