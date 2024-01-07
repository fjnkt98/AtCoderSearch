package presenter

import (
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/server/domain"
	"fjnkt98/atcodersearch/server/utility"
	"log/slog"
)

type SearchSubmissionResponse utility.SearchResultResponse[domain.Submission]

type SearchSubmissionPresenter interface {
	Format(req domain.SearchSubmissionParam, res solr.SelectResponse[domain.Submission, domain.SearchSubmissionFacetCounts], t int) SearchSubmissionResponse
	Error(msg string) SearchSubmissionResponse
}

type searchSubmissionPresenter struct {
}

func NewSearchSubmissionPresenter() SearchSubmissionPresenter {
	return &searchSubmissionPresenter{}
}

func (p *searchSubmissionPresenter) Format(req domain.SearchSubmissionParam, res solr.SelectResponse[domain.Submission, domain.SearchSubmissionFacetCounts], t int) SearchSubmissionResponse {
	rows := req.GetRows()
	pages := 0
	index := 0
	if rows != 0 {
		pages = (res.Response.NumFound + rows) / rows
		index = (res.Response.Start / rows) + 1
	}

	result := SearchSubmissionResponse{
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
		slog.String("domain", "submission"),
		slog.Int("elapsed_time", result.Stats.Time),
		slog.Int("hits", res.Response.NumFound),
		slog.Any("params", req),
	)
	return result
}

func (p *searchSubmissionPresenter) Error(msg string) SearchSubmissionResponse {
	return SearchSubmissionResponse(utility.NewErrorResponse[domain.Submission](msg))
}
