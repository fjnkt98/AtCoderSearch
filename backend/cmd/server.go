package cmd

import (
	"fjnkt98/atcodersearch/server"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Launch API server",
	Long:  "Launch API server",
	Run: func(cmd *cobra.Command, args []string) {
		db := GetDB(GetEngine())
		r := server.NewRouter(db)

		r.Run("localhost:8000")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
