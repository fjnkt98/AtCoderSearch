package cmd

import (
	"fjnkt98/atcodersearch/batch"
	"fjnkt98/atcodersearch/batch/generate"
	"fjnkt98/atcodersearch/config"
	"fjnkt98/atcodersearch/repository"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate document JSON files",
	Long:  "Generate document JSON files",
}

var generateProblemCmd = &cobra.Command{
	Use:   "problem",
	Short: "Generate problem document JSON files",
	Long:  "Generate problem document JSON files",
	Run: func(cmd *cobra.Command, args []string) {
		db := GetDB(GetEngine())

		generator := generate.NewProblemGenerator(
			config.Config.Problem.Generate,
			generate.NewProblemRowReader(db),
		)

		batch.RunBatch(generator)
	},
}

var generateUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Generate user document JSON files",
	Long:  "Generate user document JSON files",
	Run: func(cmd *cobra.Command, args []string) {
		db := GetDB(GetEngine())

		generator := generate.NewUserGenerator(
			config.Config.User.Generate,
			generate.NewUserRowReader(db),
		)

		batch.RunBatch(generator)
	},
}

var generateSubmissionCmd = &cobra.Command{
	Use:   "submission",
	Short: "Generate submission document JSON files",
	Long:  "Generate submission document JSON files",
	Run: func(cmd *cobra.Command, args []string) {
		db := GetDB(GetEngine())

		generator := generate.NewSubmissionGenerator(
			config.Config.Submission.Generate,
			generate.NewSubmissionRowReader(
				db,
				repository.NewUpdateHistoryRepository(db),
				config.Config.Submission.Read,
			),
		)

		batch.RunBatch(generator)
	},
}

var generateRecommendCmd = &cobra.Command{
	Use:   "recommend",
	Short: "Generate recommend document JSON files",
	Long:  "Generate recommend document JSON files",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	generateCmd.PersistentFlags().String("save-dir", "", "Directory path at which generated documents will be saved")
	viper.BindPFlag("problem.generate.save_dir", generateCmd.PersistentFlags().Lookup("save-dir"))
	viper.BindPFlag("user.generate.save_dir", generateCmd.PersistentFlags().Lookup("save-dir"))
	viper.BindPFlag("submission.generate.save_dir", generateCmd.PersistentFlags().Lookup("save-dir"))

	generateCmd.PersistentFlags().Int("concurrent", 10, "Concurrent number of document generation processes")
	viper.BindPFlag("problem.generate.concurrent", generateCmd.PersistentFlags().Lookup("concurrent"))
	viper.BindPFlag("user.generate.concurrent", generateCmd.PersistentFlags().Lookup("concurrent"))
	viper.BindPFlag("submission.generate.concurrent", generateCmd.PersistentFlags().Lookup("concurrent"))

	generateCmd.PersistentFlags().Int("chunk-size", 1000, "Number of documents to write in 1 file.")
	viper.BindPFlag("problem.generate.chunk_size", generateCmd.PersistentFlags().Lookup("chunk-size"))
	viper.BindPFlag("user.generate.chunk_size", generateCmd.PersistentFlags().Lookup("chunk-size"))
	viper.BindPFlag("submission.generate.chunk_size", generateCmd.PersistentFlags().Lookup("chunk-size"))

	generateCmd.AddCommand(generateProblemCmd)
	generateCmd.AddCommand(generateUserCmd)
	generateCmd.AddCommand(generateSubmissionCmd)
	generateCmd.AddCommand(generateRecommendCmd)

	rootCmd.AddCommand(generateCmd)
}
