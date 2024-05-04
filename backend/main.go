package main

import (
	"context"
	"fjnkt98/atcodersearch/cmd"
	"log/slog"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := godotenv.Load(); err != nil {
		slog.Warn("failed to load .env file", slog.String("error", err.Error()))
	}

	app := cmd.NewApp()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err := app.RunContext(ctx, os.Args); err != nil {
		slog.Error("command failed", slog.Any("error", err))
		os.Exit(1)
	}
}
