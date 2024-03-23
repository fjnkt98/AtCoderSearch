package serve

import (
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/server/api/search/problem"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/goark/errs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/urfave/cli/v2"
)

type Validator struct{}

func (v *Validator) Validate(i any) error {
	if c, ok := i.(validation.Validatable); ok {
		return c.Validate()
	}
	return nil
}

func NewServeCmd() *cli.Command {
	return &cli.Command{
		Name: "serve",
		Action: func(c *cli.Context) error {
			e := echo.New()
			e.Use(middleware.Recover())
			e.Validator = new(Validator)

			core, err := solr.NewSolrCore("http://localhost:8983", "example")
			if err != nil {
				return errs.Wrap(err)
			}
			controller := problem.NewSearchProblemController(core)
			e.GET("/api/search/problem", controller.SearchProblem)

			go func() {
				if err := e.Start(":8000"); err != nil && err != http.ErrServerClosed {
					return
				}
			}()

			<-c.Done()
			if err := e.Shutdown(c.Context); err != nil {
				return errs.Wrap(err)
			}
			return nil
		},
	}
}
