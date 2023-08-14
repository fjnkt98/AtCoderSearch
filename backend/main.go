package main

import (
	"fjnkt98/atcodersearch/cmd"
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/exp/slog"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := godotenv.Load(); err != nil {
		log.Printf("failed to load .env file: %s.", err.Error())
	}
	cmd.Execute()
}
