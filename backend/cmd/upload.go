package cmd

import (
	"fjnkt98/atcodersearch/batch"
	"fjnkt98/atcodersearch/batch/upload"
	"fjnkt98/atcodersearch/pkg/solr"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

func newUploadCmd(args []string, sub ...*cobra.Command) *cobra.Command {
	uploadCmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload document JSON files into Solr core",
		Long:  "Upload document JSON files into Solr core",
	}

	uploadCmd.SetArgs(args)
	uploadCmd.AddCommand(sub...)

	return uploadCmd
}

func newUploadProblemCmd(args []string, config *RootConfig, runFunc func(cmd *cobra.Command, args []string)) *cobra.Command {
	uploadProblemCmd := &cobra.Command{
		Use:   "problem",
		Short: "Upload document JSON files into problem core",
		Long:  "Upload document JSON files into problem core",
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlag("upload.problem.optimize", cmd.Flags().Lookup("optimize"))
			viper.BindPFlag("upload.problem.truncate", cmd.Flags().Lookup("truncate"))
			viper.BindPFlag("upload.problem.save_dir", cmd.Flags().Lookup("save-dir"))
			viper.BindPFlag("upload.problem.concurrent", cmd.Flags().Lookup("concurrent"))

			MustLoadConfigFromFlags(cmd.Flags(), config)
		},
		Run: func(cmd *cobra.Command, args []string) {
			core, err := solr.NewSolrCore(config.SolrHost, config.ProblemCoreName)
			if err != nil {
				slog.Error("failed to create core", slog.Any("error", err))
				panic("failed to create core")
			}

			uploader := upload.NewDocumentUploader(
				core,
				config.Upload.Problem.SaveDir,
				config.Upload.Problem.Concurrent,
				config.Upload.Problem.Optimize,
				config.Upload.Problem.Truncate,
			)

			batch.RunBatch(uploader)
		},
	}

	uploadProblemCmd.SetArgs(args)
	if runFunc != nil {
		uploadProblemCmd.Run = runFunc
	}
	uploadProblemCmd.Flags().BoolP("optimize", "o", false, "When true, send optimize request to Solr")
	uploadProblemCmd.Flags().BoolP("truncate", "t", false, "When true, truncate index before upload")
	uploadProblemCmd.Flags().String("save-dir", "", "Directory path at which generated documents will be saved")
	uploadProblemCmd.Flags().Int("concurrent", 3, "Concurrent number of document upload processes")

	return uploadProblemCmd
}

func newUploadUserCmd(args []string, config *RootConfig, runFunc func(cmd *cobra.Command, args []string)) *cobra.Command {
	var uploadUserCmd = &cobra.Command{
		Use:   "user",
		Short: "Upload document JSON files into user core",
		Long:  "Upload document JSON files into user core",
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlag("upload.user.optimize", cmd.Flags().Lookup("optimize"))
			viper.BindPFlag("upload.user.truncate", cmd.Flags().Lookup("truncate"))
			viper.BindPFlag("upload.user.save_dir", cmd.Flags().Lookup("save-dir"))
			viper.BindPFlag("upload.user.concurrent", cmd.Flags().Lookup("concurrent"))

			MustLoadConfigFromFlags(cmd.Flags(), config)
		},
		Run: func(cmd *cobra.Command, args []string) {
			core, err := solr.NewSolrCore(config.SolrHost, config.UserCoreName)
			if err != nil {
				slog.Error("failed to create core", slog.Any("error", err))
				panic("failed to create core")
			}

			uploader := upload.NewDocumentUploader(
				core,
				config.Upload.User.SaveDir,
				config.Upload.User.Concurrent,
				config.Upload.User.Optimize,
				config.Upload.User.Truncate,
			)

			batch.RunBatch(uploader)
		},
	}

	uploadUserCmd.SetArgs(args)
	if runFunc != nil {
		uploadUserCmd.Run = runFunc
	}
	uploadUserCmd.Flags().BoolP("optimize", "o", false, "When true, send optimize request to Solr")
	uploadUserCmd.Flags().BoolP("truncate", "t", false, "When true, truncate index before upload")
	uploadUserCmd.Flags().String("save-dir", "", "Directory path at which generated documents will be saved")
	uploadUserCmd.Flags().Int("concurrent", 3, "Concurrent number of document upload processes")

	return uploadUserCmd
}

func newUploadSubmissionCmd(args []string, config *RootConfig, runFunc func(cmd *cobra.Command, args []string)) *cobra.Command {
	uploadSubmissionCmd := &cobra.Command{
		Use:   "submission",
		Short: "Upload document JSON files into submission core",
		Long:  "Upload document JSON files into submission core",
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlag("upload.submission.optimize", cmd.Flags().Lookup("optimize"))
			viper.BindPFlag("upload.submission.truncate", cmd.Flags().Lookup("truncate"))
			viper.BindPFlag("upload.submission.save_dir", cmd.Flags().Lookup("save-dir"))
			viper.BindPFlag("upload.submission.concurrent", cmd.Flags().Lookup("concurrent"))

			MustLoadConfigFromFlags(cmd.Flags(), config)
		},
		Run: func(cmd *cobra.Command, args []string) {
			core, err := solr.NewSolrCore(config.SolrHost, config.SubmissionCoreName)
			if err != nil {
				slog.Error("failed to create core", slog.Any("error", err))
				panic("failed to create core")
			}

			uploader := upload.NewDocumentUploader(
				core,
				config.Upload.Submission.SaveDir,
				config.Upload.Submission.Concurrent,
				config.Upload.Submission.Optimize,
				config.Upload.Submission.Truncate,
			)

			batch.RunBatch(uploader)
		},
	}

	uploadSubmissionCmd.SetArgs(args)
	if runFunc != nil {
		uploadSubmissionCmd.Run = runFunc
	}
	uploadSubmissionCmd.Flags().BoolP("optimize", "o", false, "When true, send optimize request to Solr")
	uploadSubmissionCmd.Flags().BoolP("truncate", "t", false, "When true, truncate index before upload")
	uploadSubmissionCmd.Flags().String("save-dir", "", "Directory path at which generated documents will be saved")
	uploadSubmissionCmd.Flags().Int("concurrent", 3, "Concurrent number of document upload processes")

	return uploadSubmissionCmd
}
