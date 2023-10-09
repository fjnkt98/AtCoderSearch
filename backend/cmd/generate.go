package cmd

import (
	"context"
	"fjnkt98/atcodersearch/acs"
	"fjnkt98/atcodersearch/problem"
	"fjnkt98/atcodersearch/recommend"
	"fjnkt98/atcodersearch/submission"
	"fjnkt98/atcodersearch/user"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/morikuni/failure"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
	"golang.org/x/sync/errgroup"
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
		saveDir, err := GetSaveDir(cmd, "problem")
		if err != nil {
			slog.Error("failed to get save dir", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}
		db := GetDB()
		concurrent := GetInt(cmd, "concurrent")
		chunkSize := GetInt(cmd, "chunk-size")

		ctx, cancel := context.WithCancel(context.Background())
		eg, ctx := errgroup.WithContext(ctx)

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)

		done := make(chan Msg, 1)

		eg.Go(func() error {
			if err := problem.Generate(ctx, db, saveDir, chunkSize, concurrent); err != nil {
				return failure.Wrap(err)
			}
			done <- Msg{}
			return nil
		})

		eg.Go(func() error {
			select {
			case <-quit:
				defer cancel()
				return failure.New(acs.Interrupt, failure.Message("generating problem documents has been interrupted"))
			case <-ctx.Done():
				return nil
			case <-done:
				return nil
			}
		})

		if err := eg.Wait(); err != nil {
			if failure.Is(err, acs.Interrupt) {
				slog.Error("generating problem documents has been interrupted", slog.String("error", fmt.Sprintf("%+v", err)))
				return
			} else {
				slog.Error("failed to generate problem documents", slog.String("error", fmt.Sprintf("%+v", err)))
				os.Exit(1)
			}
		}
		slog.Info("finished generating problem documents successfully.")
	},
}

var generateUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Generate user document JSON files",
	Long:  "Generate user document JSON files",
	Run: func(cmd *cobra.Command, args []string) {
		saveDir, err := GetSaveDir(cmd, "user")
		if err != nil {
			slog.Error("failed to get save dir", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}
		db := GetDB()
		concurrent := GetInt(cmd, "concurrent")
		chunkSize := GetInt(cmd, "chunk-size")

		ctx, cancel := context.WithCancel(context.Background())
		eg, ctx := errgroup.WithContext(ctx)

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)

		done := make(chan Msg, 1)

		eg.Go(func() error {
			if err := user.Generate(ctx, db, saveDir, chunkSize, concurrent); err != nil {
				return failure.Wrap(err)
			}
			done <- Msg{}
			return nil
		})

		eg.Go(func() error {
			select {
			case <-quit:
				defer cancel()
				return failure.New(acs.Interrupt, failure.Message("generating user documents has been interrupted"))
			case <-ctx.Done():
				return nil
			case <-done:
				return nil
			}
		})

		if err := eg.Wait(); err != nil {
			if failure.Is(err, acs.Interrupt) {
				slog.Error("generating user documents has been interrupted", slog.String("error", fmt.Sprintf("%+v", err)))
				return
			} else {
				slog.Error("failed to generate user documents", slog.String("error", fmt.Sprintf("%+v", err)))
				os.Exit(1)
			}
		}
		slog.Info("finished generating user documents successfully.")
	},
}

var generateSubmissionCmd = &cobra.Command{
	Use:   "submission",
	Short: "Generate submission document JSON files",
	Long:  "Generate submission document JSON files",
	Run: func(cmd *cobra.Command, args []string) {
		saveDir, err := GetSaveDir(cmd, "submission")
		if err != nil {
			slog.Error("failed to get save dir", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}
		db := GetDB()
		concurrent := GetInt(cmd, "concurrent")
		chunkSize := GetInt(cmd, "chunk-size")

		ctx, cancel := context.WithCancel(context.Background())
		eg, ctx := errgroup.WithContext(ctx)

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)

		done := make(chan Msg, 1)

		eg.Go(func() error {
			if err := submission.Generate(ctx, db, saveDir, chunkSize, concurrent, time.Time{}, 90); err != nil {
				return failure.Wrap(err)
			}
			done <- Msg{}
			return nil
		})

		eg.Go(func() error {
			select {
			case <-quit:
				defer cancel()
				return failure.New(acs.Interrupt, failure.Message("generating problem documents has been interrupted"))
			case <-ctx.Done():
				return nil
			case <-done:
				return nil
			}
		})

		if err := eg.Wait(); err != nil {
			if failure.Is(err, acs.Interrupt) {
				slog.Error("generating submission documents has been interrupted", slog.String("error", fmt.Sprintf("%+v", err)))
				return
			} else {
				slog.Error("failed to generate submission documents", slog.String("error", fmt.Sprintf("%+v", err)))
				os.Exit(1)
			}
		}
		slog.Info("finished generating submission documents successfully.")
	},
}

var generateRecommendCmd = &cobra.Command{
	Use:   "recommend",
	Short: "Generate recommend document JSON files",
	Long:  "Generate recommend document JSON files",
	Run: func(cmd *cobra.Command, args []string) {
		saveDir, err := GetSaveDir(cmd, "recommend")
		if err != nil {
			slog.Error("failed to get save dir", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}
		db := GetDB()
		concurrent := GetInt(cmd, "concurrent")
		chunkSize := GetInt(cmd, "chunk-size")

		ctx, cancel := context.WithCancel(context.Background())
		eg, ctx := errgroup.WithContext(ctx)

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)

		done := make(chan Msg, 1)

		eg.Go(func() error {
			if err := recommend.Generate(ctx, db, saveDir, chunkSize, concurrent); err != nil {
				return failure.Wrap(err)
			}
			done <- Msg{}
			return nil
		})

		eg.Go(func() error {
			select {
			case <-quit:
				defer cancel()
				return failure.New(acs.Interrupt, failure.Message("generating recommend documents has been interrupted"))
			case <-ctx.Done():
				return nil
			case <-done:
				return nil
			}
		})

		if err := eg.Wait(); err != nil {
			if failure.Is(err, acs.Interrupt) {
				slog.Error("generating recommend documents has been interrupted", slog.String("error", fmt.Sprintf("%+v", err)))
				return
			} else {
				slog.Error("failed to generate recommend documents", slog.String("error", fmt.Sprintf("%+v", err)))
				os.Exit(1)
			}
		}
		slog.Info("finished generating recommend documents successfully.")
	},
}

func init() {
	generateCmd.PersistentFlags().String("save-dir", "", "Directory path at which generated documents will be saved")
	generateCmd.PersistentFlags().Int("concurrent", 10, "Concurrent number of document generation processes")
	generateCmd.PersistentFlags().Int("chunk-size", 1000, "Number of documents to write in 1 file.")
	generateCmd.AddCommand(generateProblemCmd)
	generateCmd.AddCommand(generateUserCmd)
	generateCmd.AddCommand(generateSubmissionCmd)
	generateCmd.AddCommand(generateRecommendCmd)

	rootCmd.AddCommand(generateCmd)
}
