package usecase

import (
	"context"
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/server/domain"

	"github.com/goark/errs"
)

type SearchSubmissionUsecase interface {
	Search(ctx context.Context, params domain.SearchSubmissionParam) (solr.SelectResponse[domain.Submission, domain.SearchSubmissionFacetCounts], error)
}

type searchSubmissionUsecase struct {
	core solr.SolrCore
}

func NewSearchSubmissionUsecase(core solr.SolrCore) SearchSubmissionUsecase {
	return &searchSubmissionUsecase{
		core: core,
	}
}

func (u *searchSubmissionUsecase) Search(ctx context.Context, params domain.SearchSubmissionParam) (solr.SelectResponse[domain.Submission, domain.SearchSubmissionFacetCounts], error) {
	q := params.Query()
	res, err := solr.SelectWithContext[domain.Submission, domain.SearchSubmissionFacetCounts](ctx, u.core, q)
	if err != nil {
		return res, errs.New(
			"failed to execute select query",
			errs.WithCause(err),
		)
	}

	return res, nil
}
