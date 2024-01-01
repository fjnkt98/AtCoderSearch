package main

import (
	"fjnkt98/atcodersearch/cmd"
	"os"

	"golang.org/x/exp/slog"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cmd.Execute()
}
