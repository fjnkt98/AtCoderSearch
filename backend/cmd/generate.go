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

func newGenerateProblemCmd(args []string, runFunc func(cmd *cobra.Command, args []string)) *cobra.Command {
	generateProblemCmd := &cobra.Command{
		Use:   "problem",
		Short: "Generate problem document JSON files",
		Long:  "Generate problem document JSON files",
		PreRun: func(cmd *cobra.Command, args []string) {
			cmd.Flags().String("save-dir", "", "Directory path at which generated documents will be saved.")
			viper.BindPFlag("generate.problem.save_dir", cmd.Flags().Lookup("save-dir"))
			cmd.Flags().Int("chunk-size", 1000, "Number of documents to write in 1 file.")
			viper.BindPFlag("generate.problem.chunk_size", cmd.Flags().Lookup("chunk-size"))
			cmd.Flags().Int("concurrent", 10, "Concurrent number of document generation processes.")
			viper.BindPFlag("generate.problem.concurrent", cmd.Flags().Lookup("concurrent"))
		},
		Run: func(cmd *cobra.Command, args []string) {
			db := repository.MustGetDB(Config.DataBaseURL)

			generator := generate.NewProblemGenerator(
				generate.NewProblemRowReader(db),
				Config.Generate.Problem.SaveDir,
				Config.Generate.Problem.ChunkSize,
				Config.Generate.Problem.Concurrent,
			)

			batch.RunBatch(generator)
		},
	}

	generateProblemCmd.SetArgs(args)
	if runFunc != nil {
		generateProblemCmd.Run = runFunc
	}

	return generateProblemCmd

}

func newGenerateUserCmd(args []string, runFunc func(cmd *cobra.Command, args []string)) *cobra.Command {
	generateUserCmd := &cobra.Command{
		Use:   "user",
		Short: "Generate user document JSON files",
		Long:  "Generate user document JSON files",
		PreRun: func(cmd *cobra.Command, args []string) {
			cmd.Flags().String("save-dir", "", "Directory path at which generated documents will be saved.")
			viper.BindPFlag("generate.user.save_dir", cmd.Flags().Lookup("save-dir"))
			cmd.Flags().Int("chunk-size", 1000, "Number of documents to write in 1 file.")
			viper.BindPFlag("generate.user.chunk_size", cmd.Flags().Lookup("chunk-size"))
			cmd.Flags().Int("concurrent", 10, "Concurrent number of document generation processes.")
			viper.BindPFlag("generate.user.concurrent", cmd.Flags().Lookup("concurrent"))
		},
		Run: func(cmd *cobra.Command, args []string) {
			db := repository.MustGetDB(Config.DataBaseURL)

			generator := generate.NewUserGenerator(
				generate.NewUserRowReader(db),
				Config.Generate.User.SaveDir,
				Config.Generate.User.ChunkSize,
				Config.Generate.User.Concurrent,
			)

			batch.RunBatch(generator)
		},
	}

	generateUserCmd.SetArgs(args)
	if runFunc != nil {
		generateUserCmd.Run = runFunc
	}

	return generateUserCmd
}

func newGenerateSubmissionCmd(args []string, runFunc func(cmd *cobra.Command, args []string)) *cobra.Command {
	generateSubmissionCmd := &cobra.Command{
		Use:   "submission",
		Short: "Generate submission document JSON files",
		Long:  "Generate submission document JSON files",
		PreRun: func(cmd *cobra.Command, args []string) {
			cmd.Flags().String("save-dir", "", "Directory path at which generated documents will be saved.")
			viper.BindPFlag("generate.submission.save_dir", cmd.Flags().Lookup("save-dir"))
			cmd.Flags().Int("chunk-size", 1000, "Number of documents to write in 1 file.")
			viper.BindPFlag("generate.submission.chunk_size", cmd.Flags().Lookup("chunk-size"))
			cmd.Flags().Int("concurrent", 10, "Concurrent number of document generation processes.")
			viper.BindPFlag("generate.submission.concurrent", cmd.Flags().Lookup("concurrent"))
			cmd.Flags().Int("interval", 10, "The latest N days' submissions shall be considered.")
			viper.BindPFlag("generate.submission.interval", cmd.Flags().Lookup("interval"))
			cmd.Flags().BoolP("all", "a", false, "When false, crawl only the submissions which doesn't have been crawled (in the interval).")
			viper.BindPFlag("generate.submission.all", cmd.Flags().Lookup("all"))
		},
		Run: func(cmd *cobra.Command, args []string) {
			db := repository.MustGetDB(Config.DataBaseURL)

			generator := generate.NewSubmissionGenerator(
				generate.NewSubmissionRowReader(
					db,
					Config.Generate.Submission.Interval,
					Config.Generate.Submission.All,
				),
				Config.Generate.Submission.SaveDir,
				Config.Generate.Submission.ChunkSize,
				Config.Generate.Submission.Concurrent,
			)

			batch.RunBatch(generator)
		},
	}
	generateSubmissionCmd.SetArgs(args)
	if runFunc != nil {
		generateSubmissionCmd.Run = runFunc
	}

	return generateSubmissionCmd
}
