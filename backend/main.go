package main

import (
	"fjnkt98/atcodersearch/cmd"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cmd.Execute()
}
