package main

import (
	"fjnkt98/atcodersearch/cmd"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/exp/slog"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := godotenv.Load(); err != nil {
		slog.Warn("failed to load .env file", slog.String("error", err.Error()))
	}
	cmd.Execute()
}
