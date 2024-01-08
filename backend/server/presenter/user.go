package presenter

import (
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/server/domain"
	"fjnkt98/atcodersearch/server/utility"
)

type SearchUserResponse utility.SearchResultResponse[domain.User]

type SearchUserPresenter interface {
	Format(req domain.SearchUserParam, res solr.SelectResponse[domain.User, domain.SearchUserFacetCounts], t int) SearchUserResponse
	Error(msg string) SearchUserResponse
}

type searchYserPresenter struct {
}

func NewSearchUserPresenter() SearchUserPresenter {
	return &searchYserPresenter{}
}

func (p *searchYserPresenter) Format(req domain.SearchUserParam, res solr.SelectResponse[domain.User, domain.SearchUserFacetCounts], t int) SearchUserResponse {
	rows := req.GetRows()
	pages := 0
	index := 0
	if rows != 0 {
		pages = (res.Response.NumFound + rows) / rows
		index = (res.Response.Start / rows) + 1
	}

	result := SearchUserResponse{
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
	return result
}

func (p *searchYserPresenter) Error(msg string) SearchUserResponse {
	return SearchUserResponse(utility.NewErrorResponse[domain.User](msg))
}
