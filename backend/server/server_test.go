//go:build test_server

// docker run --rm -d -p 5432:5432 --name postgres -e POSTGRES_DB=test_atcodersearch -e POSTGRES_USER=test_atcodersearch -e POSTGRES_PASSWORD=test_atcodersearch --mount type=bind,src=./schema.sql,dst=/docker-entrypoint-initdb.d/schema.sql postgres:15

package server

import (
	"database/sql"
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/server/controller"
	"fjnkt98/atcodersearch/server/presenter"
	"fjnkt98/atcodersearch/server/usecase"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

func getTestDB() *bun.DB {
	os.Setenv("PGSSLMODE", "disable")
	engine, err := sql.Open("postgres", "postgres://test_atcodersearch:test_atcodersearch@localhost/test_atcodersearch")
	if err != nil {
		slog.Error("failed to open database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if err := engine.Ping(); err != nil {
		slog.Error("failed to connect database", slog.String("error", err.Error()))
		os.Exit(1)
	}

	db := bun.NewDB(engine, pgdialect.New())
	db.AddQueryHook(
		bundebug.NewQueryHook(
			bundebug.WithVerbose(true),
			bundebug.FromEnv("BUNDEBUG"),
		),
	)

	return db
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
		{name: "page", query: "limit=3", code: 200},
		{name: "filter by category", query: "filter.category=ABC,ARC", code: 200},
		{name: "filter by difficulty", query: "filter.difficulty.from=1000&filter.difficulty.to=2000", code: 200},
		{name: "filter by color", query: "filter.color=blue", code: 200},
		{name: "term facet", query: "facet.term=category,color", code: 200},
		{name: "range facet", query: "facet.difficulty.from=0&facet.difficulty.to=2000&facet.difficulty.gap=500", code: 200},
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
