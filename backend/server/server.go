package server

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ServerConfig struct {
	AllowOrigins []string
}

type option func(*ServerConfig)

func WithAllowOrigins(origins []string) option {
	return func(opt *ServerConfig) {
		opt.AllowOrigins = origins
	}
}

type Validator struct{}

func (v *Validator) Validate(i any) error {
	if c, ok := i.(validation.Validatable); ok {
		return c.Validate()
	}
	return nil
}

func NewServer(options ...option) *echo.Echo {
	config := &ServerConfig{
		AllowOrigins: nil,
	}
	for _, opt := range options {
		opt(config)
	}

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
		AllowOrigins: config.AllowOrigins,
	}))
	e.HideBanner = true
	e.HidePort = true
	e.Validator = new(Validator)

	return e
}
