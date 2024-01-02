package usecase

import (
	"context"
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/repository"
	"fjnkt98/atcodersearch/server/domain"
	"log/slog"

	"github.com/goark/errs"
)

type SearchProblemUsecase interface {
	Search(ctx context.Context, params domain.SearchProblemParam) (solr.SelectResponse[domain.Problem, domain.SearchProblemFacetCounts], error)
}

type searchProblemUsecase struct {
	core solr.SolrCore
}

func NewSearchProblemUsecase(core solr.SolrCore) SearchProblemUsecase {
	return &searchProblemUsecase{
		core: core,
	}
}

func (u *searchProblemUsecase) Search(ctx context.Context, params domain.SearchProblemParam) (solr.SelectResponse[domain.Problem, domain.SearchProblemFacetCounts], error) {
	q := params.Query()
	res, err := solr.SelectWithContext[domain.Problem, domain.SearchProblemFacetCounts](ctx, u.core, q)
	if err != nil {
		return res, errs.New(
			"failed to execute select query",
			errs.WithCause(err),
		)
	}

	return res, nil
}

type RecommendProblemUsecase interface {
	Recommend(ctx context.Context, params domain.RecommendProblemParam) (solr.SelectResponse[domain.RecommendedProblem, any], error)
}

type recommendProblemUsecase struct {
	core solr.SolrCore
	repo repository.UserRepository
}

func NewRecommendProblemUsecase(core solr.SolrCore, repo repository.UserRepository) RecommendProblemUsecase {
	return &recommendProblemUsecase{
		core: core,
		repo: repo,
	}
}

func (u *recommendProblemUsecase) Recommend(ctx context.Context, params domain.RecommendProblemParam) (solr.SelectResponse[domain.RecommendedProblem, any], error) {
	if params.Rating == 0 && params.UserID != "" {
		if rating, err := u.repo.FetchRatingByUserName(ctx, params.UserID); err != nil {
			slog.Warn("invalid user id", slog.Any("error", err))
		} else {
			params.Rating = rating
		}
	}

	q := params.Query()
	res, err := solr.SelectWithContext[domain.RecommendedProblem, any](ctx, u.core, q)
	if err != nil {
		return res, errs.New(
			"failed to execute select query",
			errs.WithCause(err),
		)
	}

	return res, nil
}
