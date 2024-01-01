package recommend

import (
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/server/utility"
	"log/slog"
)

type SearchResultResponse utility.SearchResultResponse[Recommend]

type RecommendPresenter interface {
	Format(req SearchParams, res solr.SelectResponse[Recommend, any], t int) SearchResultResponse
	Error(msg string) SearchResultResponse
}

type recommendPresenter struct {
}

func NewRecommendPresenter() RecommendPresenter {
	return &recommendPresenter{}
}

func (p *recommendPresenter) Format(req SearchParams, res solr.SelectResponse[Recommend, any], t int) SearchResultResponse {
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
			Facet:  nil,
		},
		Items: res.Response.Docs,
	}
	slog.Info(
		"querylog",
		slog.String("domain", "recommend"),
		slog.Int("elapsed_time", result.Stats.Time),
		slog.Int("hits", res.Response.NumFound),
		slog.Any("params", req),
	)
	return result
}

func (p *recommendPresenter) Error(msg string) SearchResultResponse {
	return SearchResultResponse(utility.NewErrorResponse[Recommend](msg))
}
