package submission

import (
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/server/utility"
	"log/slog"
)

type SearchResultResponse utility.SearchResultResponse[Submission]

type SubmissionPresenter interface {
	Format(req SearchParams, res solr.SelectResponse[Submission, SolrFacetCounts], t int) SearchResultResponse
	Error(msg string) SearchResultResponse
}

type submissionPresenter struct {
}

func NewSubmissionPresenter() SubmissionPresenter {
	return &submissionPresenter{}
}

func (p *submissionPresenter) Format(req SearchParams, res solr.SelectResponse[Submission, SolrFacetCounts], t int) SearchResultResponse {
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
		slog.String("domain", "submission"),
		slog.Int("elapsed_time", result.Stats.Time),
		slog.Int("hits", res.Response.NumFound),
		slog.Any("params", req),
	)
	return result
}

func (p *submissionPresenter) Error(msg string) SearchResultResponse {
	return SearchResultResponse(utility.NewErrorResponse[Submission](msg))
}
