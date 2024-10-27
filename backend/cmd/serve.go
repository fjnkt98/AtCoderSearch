package cmd

import (
	"fjnkt98/atcodersearch/api"
	"fjnkt98/atcodersearch/repository"
	"fjnkt98/atcodersearch/searchers"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/meilisearch/meilisearch-go"
	"github.com/urfave/cli/v2"
)

type Problem struct {
}

func NewServeCmd() *cli.Command {
	return &cli.Command{
		Name: "serve",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "port",
				Value:   8000,
				EnvVars: []string{"BACKEND_PORT"},
			},
		},
		Action: func(c *cli.Context) error {
			port := c.Int("port")

			pool, err := repository.NewPool(c.Context, c.String("database-url"))
			if err != nil {
				return fmt.Errorf("new pool: %w", err)
			}

			client, err := meilisearch.Connect(c.String("engine-url"), meilisearch.WithAPIKey(c.String("engine-master-key")))
			if err != nil {
				return fmt.Errorf("new connection: %w", err)
			}

			ctx, stop := signal.NotifyContext(c.Context, os.Interrupt)
			defer stop()

			searcher := searchers.NewSearcher(client, pool)

			s, err := api.NewServer(searcher)
			if err != nil {
				return fmt.Errorf("create server: %w", err)
			}

			srv := &http.Server{
				Addr:              fmt.Sprintf(":%d", port),
				Handler:           s,
				ReadHeaderTimeout: 30 * time.Second,
			}

			go func() {
				slog.LogAttrs(ctx, slog.LevelInfo, "start server", slog.Int("port", port))
				srv.ListenAndServe()
			}()

			<-ctx.Done()
			if err := srv.Shutdown(ctx); err != nil {
				return fmt.Errorf("shutdown server: %w", err)
			} else {
				slog.LogAttrs(ctx, slog.LevelInfo, "shutdown server")
			}

			return nil
		},
	}
}
