package usecase

import (
	"context"
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/server/domain"

	"github.com/goark/errs"
)

type SearchUserUsecase interface {
	Search(ctx context.Context, params domain.SearchUserParam) (solr.SelectResponse[domain.User, domain.SearchUserFacetCounts], error)
}

type searchUserUsecase struct {
	core solr.SolrCore
}

func NewSearchUserUsecase(core solr.SolrCore) SearchUserUsecase {
	return &searchUserUsecase{
		core: core,
	}
}

func (u *searchUserUsecase) Search(ctx context.Context, params domain.SearchUserParam) (solr.SelectResponse[domain.User, domain.SearchUserFacetCounts], error) {
	q := params.Query()
	res, err := solr.SelectWithContext[domain.User, domain.SearchUserFacetCounts](ctx, u.core, q)
	if err != nil {
		return res, errs.New(
			"failed to execute select query",
			errs.WithCause(err),
		)
	}

	return res, nil
}
