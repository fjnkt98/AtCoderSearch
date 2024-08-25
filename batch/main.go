package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/fjnkt98/atcodersearch-batch/cmd"
	"github.com/urfave/cli/v2"
)

func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "atcodersearch"
	app.Commands = []*cli.Command{
		cmd.NewCrawlCmd(),
		cmd.NewUpdateCmd(),
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "database-url",
			EnvVars: []string{"DATABASE_URL"},
		},
		&cli.StringFlag{
			Name:    "engine-url",
			EnvVars: []string{"ENGINE_URL"},
		},
	}
	return app
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	app := NewApp()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	if err := app.RunContext(ctx, os.Args); err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "command failed", slog.Any("error", err))
		stop()
		os.Exit(1)
	}
	stop()
}
