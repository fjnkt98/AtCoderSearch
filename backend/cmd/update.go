package cmd

import (
	"fjnkt98/atcodersearch/repository"
	"fjnkt98/atcodersearch/update"
	"fmt"

	"github.com/meilisearch/meilisearch-go"
	"github.com/urfave/cli/v2"
)

func NewUpdateCmd() *cli.Command {
	return &cli.Command{
		Name: "update",
		Subcommands: []*cli.Command{
			{
				Name: "problem",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "chunk-size",
						Value: 1000,
					},
					&cli.IntFlag{
						Name:  "concurrent",
						Value: 4,
					},
				},
				Action: func(ctx *cli.Context) error {
					pool, err := repository.NewPool(ctx.Context, ctx.String("database-url"))
					if err != nil {
						return fmt.Errorf("new pool: %w", err)
					}
					client := meilisearch.New(ctx.String("engine-url"), meilisearch.WithAPIKey(ctx.String("engine-master-key")))
					indexer := update.NewProblemIndexer(client)

					reader := update.NewProblemRowReader(pool)

					if err := update.UpdateIndex(
						ctx.Context,
						reader,
						indexer,
						ctx.Int("chunk-size"),
						ctx.Int("concurrent"),
					); err != nil {
						return fmt.Errorf("update index: %w", err)
					}

					return nil
				},
			},
			{
				Name: "user",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "chunk-size",
						Value: 10000,
					},
					&cli.IntFlag{
						Name:  "concurrent",
						Value: 4,
					},
				},
				Action: func(ctx *cli.Context) error {
					pool, err := repository.NewPool(ctx.Context, ctx.String("database-url"))
					if err != nil {
						return fmt.Errorf("new pool: %w", err)
					}
					client := meilisearch.New(ctx.String("engine-url"), meilisearch.WithAPIKey(ctx.String("engine-master-key")))
					indexer := update.NewUserIndexer(client)

					reader := update.NewUserRowReader(pool)

					if err := update.UpdateIndex(
						ctx.Context,
						reader,
						indexer,
						ctx.Int("chunk-size"),
						ctx.Int("concurrent"),
					); err != nil {
						return fmt.Errorf("update index: %w", err)
					}

					return nil
				},
			},
			{
				Name:  "language",
				Flags: []cli.Flag{},
				Action: func(ctx *cli.Context) error {
					pool, err := repository.NewPool(ctx.Context, ctx.String("database-url"))
					if err != nil {
						return fmt.Errorf("new pool: %w", err)
					}

					if err := update.UpdateLanguage(ctx.Context, pool); err != nil {
						return fmt.Errorf("update language: %w", err)
					}
					return nil
				},
			},
		},
		Flags: []cli.Flag{},
	}
}
