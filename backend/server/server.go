package server

import (
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/server/api/search"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/goark/errs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ServerConfig struct {
	DatabaseURL  string
	SolrHost     string
	AllowOrigins []string
}

type Validator struct{}

func (v *Validator) Validate(i any) error {
	if c, ok := i.(validation.Validatable); ok {
		return c.Validate()
	}
	return nil
}

func NewServer(c ServerConfig) (*echo.Echo, error) {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
		},
		AllowOrigins: c.AllowOrigins,
	}))
	e.HideBanner = true
	e.HidePort = true
	e.Validator = new(Validator)

	{
		core, err := solr.NewSolrCore(c.SolrHost, "problem")
		if err != nil {
			return nil, errs.Wrap(err)
		}
		handler := search.NewSearchProblemHandler(core)
		e.GET("/api/search/problem", handler.SearchProblem)
		e.POST("/api/search/problem", handler.SearchProblem)
	}

	return e, nil
}
