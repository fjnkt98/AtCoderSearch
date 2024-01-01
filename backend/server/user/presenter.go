package user

import (
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/server/utility"
	"log/slog"
)

type SearchResultResponse utility.SearchResultResponse[User]

type UserPresenter interface {
	Format(req SearchParams, res solr.SelectResponse[User, SolrFacetCounts], t int) SearchResultResponse
	Error(msg string) SearchResultResponse
}

type userPresenter struct {
}

func NewUserPresenter() UserPresenter {
	return &userPresenter{}
}

func (p *userPresenter) Format(req SearchParams, res solr.SelectResponse[User, SolrFacetCounts], t int) SearchResultResponse {
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
		slog.String("domain", "user"),
		slog.Int("elapsed_time", result.Stats.Time),
		slog.Int("hits", res.Response.NumFound),
		slog.Any("params", req),
	)
	return result
}

func (p *userPresenter) Error(msg string) SearchResultResponse {
	return SearchResultResponse(utility.NewErrorResponse[User](msg))
}
