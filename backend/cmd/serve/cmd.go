package serve

import (
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/repository"
	"fjnkt98/atcodersearch/server"
	"fjnkt98/atcodersearch/server/api/search"
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
			e := server.NewServer(
				server.WithAllowOrigins(ctx.StringSlice("allow-origin")),
			)

			pool, err := repository.NewPool(ctx.Context, ctx.String("database-url"))
			if err != nil {
				return errs.Wrap(err)
			}

			host := ctx.String("solr-host")

			{
				core, err := solr.NewSolrCore(host, "problem")
				if err != nil {
					return errs.Wrap(err)
				}
				search.NewSearchProblemHandler(core, pool).Register(e)
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
