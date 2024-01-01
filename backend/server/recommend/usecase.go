package recommend

import (
	"context"
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/repository"
	"log/slog"

	"github.com/goark/errs"
)

type RecommendUsecase interface {
	Search(ctx context.Context, params SearchParams) (solr.SelectResponse[Recommend, any], error)
}

type recommendUsecase struct {
	core solr.SolrCore
	repo repository.UserRepository
}

func NewRecommendUsecase(core solr.SolrCore, repo repository.UserRepository) RecommendUsecase {
	return &recommendUsecase{
		core: core,
		repo: repo,
	}
}

func (u *recommendUsecase) Search(ctx context.Context, params SearchParams) (solr.SelectResponse[Recommend, any], error) {
	if params.Rating == 0 && params.UserID != "" {
		if rating, err := u.repo.FetchRatingByUserName(ctx, params.UserID); err != nil {
			slog.Warn("invalid user id", slog.Any("error", err))
		} else {
			params.Rating = rating
		}
	}

	q := params.Query()
	res, err := solr.SelectWithContext[Recommend, any](ctx, u.core, q)
	if err != nil {
		return res, errs.New(
			"failed to execute select query",
			errs.WithCause(err),
		)
	}

	return res, nil
}

type Recommend struct {
	ProblemID    string                `json:"problem_id"`
	ProblemTitle string                `json:"problem_title"`
	ProblemURL   string                `json:"problem_url"`
	ContestID    string                `json:"contest_id"`
	ContestTitle string                `json:"contest_title"`
	ContestURL   string                `json:"contest_url"`
	Difficulty   *int                  `json:"difficulty"`
	Color        *string               `json:"color"`
	StartAt      solr.FromSolrDateTime `json:"start_at"`
	Duration     int                   `json:"duration"`
	RateChange   string                `json:"rate_change"`
	Category     string                `json:"category"`
	Score        float64               `json:"score"`
}
