package main

import (
	"context"
	"fjnkt98/atcodersearch/cmd"
	"log/slog"
	"os"
	"os/signal"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	app := cmd.NewApp()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := app.RunContext(ctx, os.Args); err != nil {
		slog.Error("command failed", slog.Any("error", err))
		os.Exit(1)
	}
}
