package cmd

import (
	"context"
	"fjnkt98/atcodersearch/acs"
	"fjnkt98/atcodersearch/solr"
	"fmt"
	"os"
	"os/signal"

	"github.com/morikuni/failure"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
	"golang.org/x/sync/errgroup"
)

func post(domain string, core *solr.Core, saveDir string, optimize bool, truncate bool, concurrent int) {
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	done := make(chan Msg, 1)

	eg.Go(func() error {
		if err := acs.PostDocument(ctx, core, saveDir, optimize, truncate, concurrent); err != nil {
			return failure.Wrap(err)
		}
		done <- Msg{}
		return nil
	})

	eg.Go(func() error {
		select {
		case <-quit:
			defer cancel()
			return failure.New(acs.Interrupt, failure.Messagef("post %s documents has been interrupted", domain))
		case <-ctx.Done():
			return nil
		case <-done:
			return nil
		}
	})

	if err := eg.Wait(); err != nil {
		if failure.Is(err, acs.Interrupt) {
			slog.Error(fmt.Sprintf("post %s documents has been interrupted", domain), slog.String("error", fmt.Sprintf("%+v", err)))
			return
		} else {
			slog.Error(fmt.Sprintf("failed to post %s documents", domain), slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}
	}

	slog.Info(fmt.Sprintf("finished post %s documents successfully.", domain))
}

var postCmd = &cobra.Command{
	Use:   "post",
	Short: "Post document JSON files into Solr core",
	Long:  "Post document JSON files into Solr core",
}

var postProblemCmd = &cobra.Command{
	Use:   "problem",
	Short: "Post document JSON files into problem core",
	Long:  "Post document JSON files into problem core",
	Run: func(cmd *cobra.Command, args []string) {
		saveDir, err := GetSaveDir(cmd, "problem")
		if err != nil {
			slog.Error("failed to get save dir", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}

		solrURL := os.Getenv("SOLR_HOST")
		if solrURL == "" {
			slog.Error("environment variable `SOLR_HOST` must be set.")
			os.Exit(1)
		}
		core, err := solr.NewSolrCore("problem", solrURL)
		if err != nil {
			slog.Error("failed to create `problem` core", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}

		optimize := GetBool(cmd, "optimize")
		truncate := GetBool(cmd, "truncate")
		concurrent := GetInt(cmd, "concurrent")

		post("problem", core, saveDir, optimize, truncate, concurrent)
	},
}

var postUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Post document JSON files into user core",
	Long:  "Post document JSON files into user core",
	Run: func(cmd *cobra.Command, args []string) {
		saveDir, err := GetSaveDir(cmd, "user")
		if err != nil {
			slog.Error("failed to get save dir", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}

		solrURL := os.Getenv("SOLR_HOST")
		if solrURL == "" {
			slog.Error("environment variable `SOLR_HOST` must be set.")
			os.Exit(1)
		}
		core, err := solr.NewSolrCore("user", solrURL)
		if err != nil {
			slog.Error("failed to create `user` core", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}

		optimize := GetBool(cmd, "optimize")
		truncate := GetBool(cmd, "truncate")
		concurrent := GetInt(cmd, "concurrent")

		post("user", core, saveDir, optimize, truncate, concurrent)
	},
}

var postSubmissionCmd = &cobra.Command{
	Use:   "submission",
	Short: "Post document JSON files into submission core",
	Long:  "Post document JSON files into submission core",
	Run: func(cmd *cobra.Command, args []string) {
		saveDir, err := GetSaveDir(cmd, "submission")
		if err != nil {
			slog.Error("failed to get save dir", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}

		solrURL := os.Getenv("SOLR_HOST")
		if solrURL == "" {
			slog.Error("environment variable `SOLR_HOST` must be set.")
			os.Exit(1)
		}
		core, err := solr.NewSolrCore("submission", solrURL)
		if err != nil {
			slog.Error("failed to create `user` core", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}

		optimize := GetBool(cmd, "optimize")
		truncate := GetBool(cmd, "truncate")
		concurrent := GetInt(cmd, "concurrent")

		post("submission", core, saveDir, optimize, truncate, concurrent)
	},
}

var postRecommendCmd = &cobra.Command{
	Use:   "recommend",
	Short: "Post document JSON files into recommend core",
	Long:  "Post document JSON files into recommend core",
	Run: func(cmd *cobra.Command, args []string) {
		saveDir, err := GetSaveDir(cmd, "recommend")
		if err != nil {
			slog.Error("failed to get save dir", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}

		solrURL := os.Getenv("SOLR_HOST")
		if solrURL == "" {
			slog.Error("environment variable `SOLR_HOST` must be set.")
			os.Exit(1)
		}
		core, err := solr.NewSolrCore("recommend", solrURL)
		if err != nil {
			slog.Error("failed to create `user` core", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}

		optimize := GetBool(cmd, "optimize")
		truncate := GetBool(cmd, "truncate")
		concurrent := GetInt(cmd, "concurrent")

		post("recommend", core, saveDir, optimize, truncate, concurrent)
	},
}

func init() {
	postCmd.PersistentFlags().BoolP("optimize", "o", false, "When true, send optimize request to Solr")
	postCmd.PersistentFlags().BoolP("truncate", "t", false, "When true, truncate index before post")
	postCmd.PersistentFlags().String("save-dir", "", "Directory path at which generated documents will be saved")
	postCmd.PersistentFlags().Int("concurrent", 3, "Concurrent number of document upload processes")

	postCmd.AddCommand(postProblemCmd)
	postCmd.AddCommand(postUserCmd)
	postCmd.AddCommand(postSubmissionCmd)
	postCmd.AddCommand(postRecommendCmd)

	rootCmd.AddCommand(postCmd)
}
