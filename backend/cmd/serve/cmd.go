package serve

import (
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/repository"
	"fjnkt98/atcodersearch/server"
	"fjnkt98/atcodersearch/server/api/list"
	"fjnkt98/atcodersearch/server/api/recommend"
	"fjnkt98/atcodersearch/server/api/search"
	"fjnkt98/atcodersearch/settings"
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
				core, err := solr.NewSolrCore(host, settings.PROBLEM_CORE_NAME)
				if err != nil {
					return errs.Wrap(err)
				}
				search.NewSearchProblemHandler(core, pool).Register(e)
			}
			{
				core, err := solr.NewSolrCore(host, settings.USER_CORE_NAME)
				if err != nil {
					return errs.Wrap(err)
				}
				search.NewSearchUserHandler(core, pool).Register(e)
			}
			{
				core, err := solr.NewSolrCore(host, settings.PROBLEM_CORE_NAME)
				if err != nil {
					return errs.Wrap(err)
				}
				recommend.NewRecommendProblemHandler(core).Register(e)
			}
			{
				core, err := solr.NewSolrCore(host, settings.SUBMISSION_CORE_NAME)
				if err != nil {
					return errs.Wrap(err)
				}
				search.NewSearchSubmissionHandler(core).Register(e)
			}
			list.NewListHandler(pool).Register(e)

			port := ctx.Int("port")
			go func() {
				slog.Info("start server", slog.Int("port", port))
				if err := e.Start(fmt.Sprintf(":%d", port)); err != nil && err != http.ErrServerClosed {
					slog.Error("failed to start server", slog.Int("port", port), slog.Any("error", err))
					panic(fmt.Sprintf("failed to start server: %s", err.Error()))
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
