package main

import (
	"context"
	"fjnkt98/atcodersearch/cmd"
	"log/slog"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	app := cli.NewApp()
	app.Name = "atcodersearch"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "database-url",
			Hidden:  true,
			EnvVars: []string{"DATABASE_URL"},
		},
		&cli.StringFlag{
			Name:    "engine-url",
			Hidden:  true,
			EnvVars: []string{"ENGINE_URL"},
		},
		&cli.StringFlag{
			Name:    "engine-master-key",
			Hidden:  true,
			EnvVars: []string{"ENGINE_MASTER_KEY"},
		},
	}
	app.Commands = []*cli.Command{
		cmd.NewCrawlCmd(),
		cmd.NewUpdateCmd(),
		cmd.NewServeCmd(),
	}

	if err := app.RunContext(context.Background(), os.Args); err != nil {
		slog.Error("command failed", slog.Any("error", err))
		os.Exit(1)
	}
}
