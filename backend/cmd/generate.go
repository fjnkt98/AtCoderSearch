package cmd

import (
	"fjnkt98/atcodersearch/batch"
	"fjnkt98/atcodersearch/batch/generate"
	"fjnkt98/atcodersearch/repository"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newGenerateCmd(args []string, sub ...*cobra.Command) *cobra.Command {
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate document JSON files",
		Long:  "Generate document JSON files",
	}

	generateCmd.SetArgs(args)
	generateCmd.AddCommand(sub...)

	return generateCmd
}

func newGenerateProblemCmd(args []string, config *RootConfig, runFunc func(cmd *cobra.Command, args []string)) *cobra.Command {
	generateProblemCmd := &cobra.Command{
		Use:   "problem",
		Short: "Generate problem document JSON files",
		Long:  "Generate problem document JSON files",
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlag("generate.problem.save_dir", cmd.Flags().Lookup("save-dir"))
			viper.BindPFlag("generate.problem.chunk_size", cmd.Flags().Lookup("chunk-size"))
			viper.BindPFlag("generate.problem.concurrent", cmd.Flags().Lookup("concurrent"))

			MustLoadConfigFromFlags(cmd.Flags(), config)
		},
		Run: func(cmd *cobra.Command, args []string) {
			db := repository.MustGetDB(config.DataBaseURL)

			generator := generate.NewProblemGenerator(
				generate.NewProblemRowReader(db),
				config.Generate.Problem.SaveDir,
				config.Generate.Problem.ChunkSize,
				config.Generate.Problem.Concurrent,
			)

			batch.RunBatch(generator)
		},
	}

	generateProblemCmd.SetArgs(args)
	if runFunc != nil {
		generateProblemCmd.Run = runFunc
	}
	generateProblemCmd.Flags().String("save-dir", "", "Directory path at which generated documents will be saved.")
	generateProblemCmd.Flags().Int("chunk-size", 1000, "Number of documents to write in 1 file.")
	generateProblemCmd.Flags().Int("concurrent", 10, "Concurrent number of document generation processes.")

	return generateProblemCmd

}

func newGenerateUserCmd(args []string, config *RootConfig, runFunc func(cmd *cobra.Command, args []string)) *cobra.Command {
	generateUserCmd := &cobra.Command{
		Use:   "user",
		Short: "Generate user document JSON files",
		Long:  "Generate user document JSON files",
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlag("generate.user.save_dir", cmd.Flags().Lookup("save-dir"))
			viper.BindPFlag("generate.user.chunk_size", cmd.Flags().Lookup("chunk-size"))
			viper.BindPFlag("generate.user.concurrent", cmd.Flags().Lookup("concurrent"))

			MustLoadConfigFromFlags(cmd.Flags(), config)
		},
		Run: func(cmd *cobra.Command, args []string) {
			db := repository.MustGetDB(config.DataBaseURL)

			generator := generate.NewUserGenerator(
				generate.NewUserRowReader(db),
				config.Generate.User.SaveDir,
				config.Generate.User.ChunkSize,
				config.Generate.User.Concurrent,
			)

			batch.RunBatch(generator)
		},
	}

	generateUserCmd.SetArgs(args)
	if runFunc != nil {
		generateUserCmd.Run = runFunc
	}
	generateUserCmd.Flags().String("save-dir", "", "Directory path at which generated documents will be saved.")
	generateUserCmd.Flags().Int("chunk-size", 1000, "Number of documents to write in 1 file.")
	generateUserCmd.Flags().Int("concurrent", 10, "Concurrent number of document generation processes.")

	return generateUserCmd
}

func newGenerateSubmissionCmd(args []string, config *RootConfig, runFunc func(cmd *cobra.Command, args []string)) *cobra.Command {
	generateSubmissionCmd := &cobra.Command{
		Use:   "submission",
		Short: "Generate submission document JSON files",
		Long:  "Generate submission document JSON files",
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlag("generate.submission.save_dir", cmd.Flags().Lookup("save-dir"))
			viper.BindPFlag("generate.submission.chunk_size", cmd.Flags().Lookup("chunk-size"))
			viper.BindPFlag("generate.submission.concurrent", cmd.Flags().Lookup("concurrent"))
			viper.BindPFlag("generate.submission.interval", cmd.Flags().Lookup("interval"))
			viper.BindPFlag("generate.submission.all", cmd.Flags().Lookup("all"))

			MustLoadConfigFromFlags(cmd.Flags(), config)
		},
		Run: func(cmd *cobra.Command, args []string) {
			db := repository.MustGetDB(config.DataBaseURL)

			generator := generate.NewSubmissionGenerator(
				generate.NewSubmissionRowReader(
					db,
					config.Generate.Submission.Interval,
					config.Generate.Submission.All,
				),
				config.Generate.Submission.SaveDir,
				config.Generate.Submission.ChunkSize,
				config.Generate.Submission.Concurrent,
			)

			batch.RunBatch(generator)
		},
	}
	generateSubmissionCmd.SetArgs(args)
	if runFunc != nil {
		generateSubmissionCmd.Run = runFunc
	}
	generateSubmissionCmd.Flags().String("save-dir", "", "Directory path at which generated documents will be saved.")
	generateSubmissionCmd.Flags().Int("chunk-size", 1000, "Number of documents to write in 1 file.")
	generateSubmissionCmd.Flags().Int("concurrent", 10, "Concurrent number of document generation processes.")
	generateSubmissionCmd.Flags().Int("interval", 10, "The latest N days' submissions shall be considered.")
	generateSubmissionCmd.Flags().BoolP("all", "a", false, "When false, crawl only the submissions which doesn't have been crawled (in the interval).")

	return generateSubmissionCmd
}
