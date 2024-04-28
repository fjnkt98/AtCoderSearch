package serve

import (
	"fjnkt98/atcodersearch/server"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/goark/errs"
	"github.com/urfave/cli/v2"
)

func NewServeCmd() *cli.Command {
	return &cli.Command{
		Name: "serve",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "port",
				Value:   8000,
				EnvVars: []string{"PORT"},
			},
			&cli.StringSliceFlag{
				Name:    "allow-origin",
				EnvVars: []string{"ALLOW_ORIGIN"},
			},
		},
		Action: func(ctx *cli.Context) error {
			e, err := server.NewServer(server.ServerConfig{DatabaseURL: ctx.String("database-url"), SolrHost: ctx.String("solr-host")})
			if err != nil {
				return errs.Wrap(err)
			}

			port := ctx.Int("port")
			go func() {
				slog.Info("start server")
				if err := e.Start(fmt.Sprintf(":%d", port)); err != nil && err != http.ErrServerClosed {
					return
				}
			}()

			<-ctx.Done()
			slog.Info("shutdown server")
			if err := e.Shutdown(ctx.Context); err != nil {
				return errs.Wrap(err)
			}
			return nil
		},
	}
}
