package problem

import (
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/server/utility"
	"log/slog"
)

type SearchResultResponse utility.SearchResultResponse[Problem]

type ProblemPresenter interface {
	Format(req SearchParams, res solr.SelectResponse[Problem, SolrFacetCounts], t int) SearchResultResponse
	Error(msg string) SearchResultResponse
}

type problemPresenter struct {
}

func NewProblemPresenter() ProblemPresenter {
	return &problemPresenter{}
}

func (p *problemPresenter) Format(req SearchParams, res solr.SelectResponse[Problem, SolrFacetCounts], t int) SearchResultResponse {
	rows := req.GetRows()
	pages := 0
	index := 0
	if rows != 0 {
		pages = (res.Response.NumFound + rows) / rows
		index = (res.Response.Start / rows) + 1
	}

	result := SearchResultResponse{
		Stats: utility.SearchResultStats{
			Time:   t,
			Total:  res.Response.NumFound,
			Index:  index,
			Count:  len(res.Response.Docs),
			Pages:  pages,
			Params: req,
			Facet:  res.FacetCounts.Into(req.Facet),
		},
		Items: res.Response.Docs,
	}
	slog.Info(
		"querylog",
		slog.String("domain", "problem"),
		slog.Int("elapsed_time", result.Stats.Time),
		slog.Int("hits", res.Response.NumFound),
		slog.Any("params", req),
	)
	return result
}

func (p *problemPresenter) Error(msg string) SearchResultResponse {
	return SearchResultResponse(utility.NewErrorResponse[Problem](msg))
}
