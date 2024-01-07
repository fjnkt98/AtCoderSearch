//go:build test_server

// docker run --rm -d -p 5432:5432 --name postgres -e POSTGRES_DB=test_atcodersearch -e POSTGRES_USER=test_atcodersearch -e POSTGRES_PASSWORD=test_atcodersearch --mount type=bind,src=./schema.sql,dst=/docker-entrypoint-initdb.d/schema.sql postgres:15

package server

import (
	"database/sql"
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/server/controller"
	"fjnkt98/atcodersearch/server/presenter"
	"fjnkt98/atcodersearch/server/usecase"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/goark/errs"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

func getTestDB() (*bun.DB, error) {
	os.Setenv("PGSSLMODE", "disable")
	engine, err := sql.Open("postgres", "postgres://test_atcodersearch:test_atcodersearch@localhost/test_atcodersearch")
	if err != nil {
		return nil, errs.New("failed to open database", errs.WithCause(err))
	}

	if err := engine.Ping(); err != nil {
		return nil, errs.New("failed to connect database", errs.WithCause(err))
	}

	db := bun.NewDB(engine, pgdialect.New())
	db.AddQueryHook(
		bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
			bundebug.FromEnv("BUNDEBUG"),
		),
	)

	return db, nil
}

func TestSearchProblemRequest(t *testing.T) {
	r := gin.New()
	r.Use(
		gin.Recovery(),
	)

	core, err := solr.NewSolrCore("http://localhost:8983", "problem")
	if err != nil {
		t.Fatalf("failed to initialize solr core: %s", err.Error())
	}

	c := controller.NewSearchProblemController(
		usecase.NewSearchProblemUsecase(core),
		presenter.NewSearchProblemPresenter(),
	)

	r.GET("/api/search/problem", c.HandleGET)
	r.POST("/api/search/problem", c.HandlePOST)

	cases := []struct {
		name  string
		query string
		code  int
	}{
		{name: "default", query: "", code: 200},
		{name: "keyword", query: "keyword=chokudai", code: 200},
		{name: "empty_keyword", query: "keyword=", code: 200},
		{name: "limit", query: "limit=50", code: 200},
		{name: "page", query: "page=2", code: 200},
		{name: "filter by category", query: "filter.category=ABC,ARC", code: 200},
		{name: "filter by difficulty", query: "filter.difficulty.from=1000&filter.difficulty.to=2000", code: 200},
		{name: "filter by difficulty (only from)", query: "filter.difficulty.from=1000", code: 200},
		{name: "filter by difficulty (only to)", query: "filter.difficulty.to=2000", code: 200},
		{name: "filter by color", query: "filter.color=blue", code: 200},
		{name: "term facet", query: "facet.term=category,color", code: 200},
		{name: "invalid term facet", query: "facet.term=difficulty", code: 400},
		{name: "range facet", query: "facet.difficulty.from=0&facet.difficulty.to=2000&facet.difficulty.gap=500", code: 200},
		{name: "valid sort", query: "sort=-score,start_at,-start_at,difficulty,-difficulty,problem_id,-problem_id", code: 200},
		{name: "invalid sort", query: "sort=score", code: 400},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			v, err := url.ParseQuery(tt.query)
			if err != nil {
				t.Fatalf("failed to parse query: %s", err.Error())
			}
			req, err := http.NewRequest("GET", "/api/search/problem?"+v.Encode(), nil)
			if err != nil {
				t.Fatalf("failed to create request: %s", err.Error())
			}

			r.ServeHTTP(w, req)

			if w.Code != tt.code {
				t.Errorf("expected code %d, but got code %d", tt.code, w.Code)
			}
		})
	}
}

func TestSearchUserRequest(t *testing.T) {
	r := gin.New()
	r.Use(
		gin.Recovery(),
	)

	core, err := solr.NewSolrCore("http://localhost:8983", "user")
	if err != nil {
		t.Fatalf("failed to initialize solr core: %s", err.Error())
	}

	c := controller.NewSearchUserController(
		usecase.NewSearchUserUsecase(core),
		presenter.NewSearchUserPresenter(),
	)

	r.GET("/api/search/user", c.HandleGET)
	r.POST("/api/search/user", c.HandlePOST)

	cases := []struct {
		name  string
		query string
		code  int
	}{
		{name: "default", query: "", code: 200},
		{name: "keyword", query: "keyword=chokudai", code: 200},
		{name: "empty_keyword", query: "keyword=", code: 200},
		{name: "limit", query: "limit=50", code: 200},
		{name: "page", query: "page=2", code: 200},
		{name: "filter by user id", query: "filter.user_id=tourist,fjnkt98", code: 200},
		{name: "filter by rating", query: "filter.rating.from=400&filter.rating.to=1000", code: 200},
		{name: "filter by birth year", query: "filter.birth_year.from=1990&filter.birth_year.to=2000", code: 200},
		{name: "filter by join count", query: "filter.join_count.from=50&filter.join_count.to=100", code: 200},
		{name: "filter by country", query: "filter.country=JP", code: 200},
		{name: "filter by color", query: "filter.color=blue", code: 200},
		{name: "term facet", query: "facet.term=country", code: 200},
		{name: "invalid term facet", query: "facet.term=color", code: 400},
		{name: "facet by rating", query: "facet.rating.from=0&facet.rating.to=2000&facet.rating.gap=500", code: 200},
		{name: "facet by birth year", query: "facet.birth_year.from=0&facet.birth_year.to=2000&facet.birth_year.gap=500", code: 200},
		{name: "facet by join count", query: "facet.join_count.from=0&facet.join_count.to=2000&facet.join_count.gap=500", code: 200},
		{name: "valid sort", query: "sort=-score,rating,-rating,birth_year,-birth_year", code: 200},
		{name: "invalid sort", query: "sort=score", code: 400},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			v, err := url.ParseQuery(tt.query)
			if err != nil {
				t.Fatalf("failed to parse query: %s", err.Error())
			}
			req, err := http.NewRequest("GET", "/api/search/user?"+v.Encode(), nil)
			if err != nil {
				t.Fatalf("failed to create request: %s", err.Error())
			}

			r.ServeHTTP(w, req)

			if w.Code != tt.code {
				t.Errorf("expected code %d, but got code %d", tt.code, w.Code)
			}
		})
	}
}
