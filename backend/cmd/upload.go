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

func newUploadProblemCmd(args []string, runFunc func(cmd *cobra.Command, args []string)) *cobra.Command {
	uploadProblemCmd := &cobra.Command{
		Use:   "problem",
		Short: "Upload document JSON files into problem core",
		Long:  "Upload document JSON files into problem core",
		PreRun: func(cmd *cobra.Command, args []string) {
			cmd.Flags().BoolP("optimize", "o", false, "When true, send optimize request to Solr")
			viper.BindPFlag("upload.problem.optimize", cmd.Flags().Lookup("optimize"))

			cmd.Flags().BoolP("truncate", "t", false, "When true, truncate index before upload")
			viper.BindPFlag("upload.problem.truncate", cmd.Flags().Lookup("truncate"))

			cmd.Flags().String("save-dir", "", "Directory path at which generated documents will be saved")
			viper.BindPFlag("upload.problem.save_dir", cmd.Flags().Lookup("save-dir"))

			cmd.Flags().Int("concurrent", 3, "Concurrent number of document upload processes")
			viper.BindPFlag("upload.problem.concurrent", cmd.Flags().Lookup("concurrent"))
		},
		Run: func(cmd *cobra.Command, args []string) {
			core, err := solr.NewSolrCore(Config.SolrHost, Config.ProblemCoreName)
			if err != nil {
				slog.Error("failed to create core", slog.Any("error", err))
				panic("failed to create core")
			}

			uploader := upload.NewDocumentUploader(
				core,
				Config.Upload.Problem.SaveDir,
				Config.Upload.Problem.Concurrent,
				Config.Upload.Problem.Optimize,
				Config.Upload.Problem.Truncate,
			)

			batch.RunBatch(uploader)
		},
	}

	uploadProblemCmd.SetArgs(args)
	if runFunc != nil {
		uploadProblemCmd.Run = runFunc
	}

	return uploadProblemCmd
}

func newUploadUserCmd(args []string, runFunc func(cmd *cobra.Command, args []string)) *cobra.Command {
	var uploadUserCmd = &cobra.Command{
		Use:   "user",
		Short: "Upload document JSON files into user core",
		Long:  "Upload document JSON files into user core",
		PreRun: func(cmd *cobra.Command, args []string) {
			cmd.Flags().BoolP("optimize", "o", false, "When true, send optimize request to Solr")
			viper.BindPFlag("upload.user.optimize", cmd.Flags().Lookup("optimize"))

			cmd.Flags().BoolP("truncate", "t", false, "When true, truncate index before upload")
			viper.BindPFlag("upload.user.truncate", cmd.Flags().Lookup("truncate"))

			cmd.Flags().String("save-dir", "", "Directory path at which generated documents will be saved")
			viper.BindPFlag("upload.user.save_dir", cmd.Flags().Lookup("save-dir"))

			cmd.Flags().Int("concurrent", 3, "Concurrent number of document upload processes")
			viper.BindPFlag("upload.user.concurrent", cmd.Flags().Lookup("concurrent"))
		},
		Run: func(cmd *cobra.Command, args []string) {
			core, err := solr.NewSolrCore(Config.SolrHost, Config.UserCoreName)
			if err != nil {
				slog.Error("failed to create core", slog.Any("error", err))
				panic("failed to create core")
			}

			uploader := upload.NewDocumentUploader(
				core,
				Config.Upload.User.SaveDir,
				Config.Upload.User.Concurrent,
				Config.Upload.User.Optimize,
				Config.Upload.User.Truncate,
			)

			batch.RunBatch(uploader)
		},
	}

	uploadUserCmd.SetArgs(args)
	if runFunc != nil {
		uploadUserCmd.Run = runFunc
	}

	return uploadUserCmd
}

func newUploadSubmissionCmd(args []string, runFunc func(cmd *cobra.Command, args []string)) *cobra.Command {
	uploadSubmissionCmd := &cobra.Command{
		Use:   "submission",
		Short: "Upload document JSON files into submission core",
		Long:  "Upload document JSON files into submission core",
		PreRun: func(cmd *cobra.Command, args []string) {
			cmd.Flags().BoolP("optimize", "o", false, "When true, send optimize request to Solr")
			viper.BindPFlag("upload.submission.optimize", cmd.Flags().Lookup("optimize"))

			cmd.Flags().BoolP("truncate", "t", false, "When true, truncate index before upload")
			viper.BindPFlag("upload.submission.truncate", cmd.Flags().Lookup("truncate"))

			cmd.Flags().String("save-dir", "", "Directory path at which generated documents will be saved")
			viper.BindPFlag("upload.submission.save_dir", cmd.Flags().Lookup("save-dir"))

			cmd.Flags().Int("concurrent", 3, "Concurrent number of document upload processes")
			viper.BindPFlag("upload.submission.concurrent", cmd.Flags().Lookup("concurrent"))
		},
		Run: func(cmd *cobra.Command, args []string) {
			core, err := solr.NewSolrCore(Config.SolrHost, Config.SubmissionCoreName)
			if err != nil {
				slog.Error("failed to create core", slog.Any("error", err))
				panic("failed to create core")
			}

			uploader := upload.NewDocumentUploader(
				core,
				Config.Upload.Submission.SaveDir,
				Config.Upload.Submission.Concurrent,
				Config.Upload.Submission.Optimize,
				Config.Upload.Submission.Truncate,
			)

			batch.RunBatch(uploader)
		},
	}

	uploadSubmissionCmd.SetArgs(args)
	if runFunc != nil {
		uploadSubmissionCmd.Run = runFunc
	}

	return uploadSubmissionCmd
}
