package presenter

import (
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/server/domain"
	"fjnkt98/atcodersearch/server/utility"
)

type SearchProblemResponse utility.SearchResultResponse[domain.Problem]

type SearchProblemPresenter interface {
	Format(req domain.SearchProblemParam, res solr.SelectResponse[domain.Problem, domain.SearchProblemFacetCounts], t int) SearchProblemResponse
	Error(msg string) SearchProblemResponse
}

type searchProblemPresenter struct{}

func NewSearchProblemPresenter() SearchProblemPresenter {
	return &searchProblemPresenter{}
}

func (p *searchProblemPresenter) Format(req domain.SearchProblemParam, res solr.SelectResponse[domain.Problem, domain.SearchProblemFacetCounts], t int) SearchProblemResponse {
	rows := req.GetRows()
	pages := 0
	index := 0
	if rows != 0 {
		pages = (res.Response.NumFound + rows) / rows
		index = (res.Response.Start / rows) + 1
	}

	result := SearchProblemResponse{
		Stats: utility.SearchResultStats{
			Time:   t,
			Total:  res.Response.NumFound,
			Index:  index,
			Count:  len(res.Response.Docs),
			Pages:  pages,
			Params: req,
			Facet:  res.FacetCounts.ToMap(req.Facet),
		},
		Items: res.Response.Docs,
	}
	return result
}

func (p *searchProblemPresenter) Error(msg string) SearchProblemResponse {
	return SearchProblemResponse(utility.NewErrorResponse[domain.Problem](msg))
}

type RecommendProblemResponse utility.SearchResultResponse[domain.RecommendedProblem]

type RecommendProblemPresenter interface {
	Format(req domain.RecommendProblemParam, res solr.SelectResponse[domain.RecommendedProblem, any], t int) RecommendProblemResponse
	Error(msg string) RecommendProblemResponse
}

type recommendProblemPresenter struct {
}

func NewRecommendProblemPresenter() RecommendProblemPresenter {
	return &recommendProblemPresenter{}
}

func (p *recommendProblemPresenter) Format(req domain.RecommendProblemParam, res solr.SelectResponse[domain.RecommendedProblem, any], t int) RecommendProblemResponse {
	rows := req.GetRows()
	pages := 0
	index := 0
	if rows != 0 {
		pages = (res.Response.NumFound + rows) / rows
		index = (res.Response.Start / rows) + 1
	}

	result := RecommendProblemResponse{
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
	return result
}

func (p *recommendProblemPresenter) Error(msg string) RecommendProblemResponse {
	return RecommendProblemResponse(utility.NewErrorResponse[domain.RecommendedProblem](msg))
}
