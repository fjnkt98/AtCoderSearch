package cmd

import (
	"fjnkt98/atcodersearch/server"

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

		db := GetDB(GetEngine())

		server.RegisterProblemRoute(r)
		server.RegisterUserRoute(r)
		server.RegisterSubmissionRoute(r)
		server.RegisterRecommendRoute(r, db)

		r.Run("localhost:8000")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
