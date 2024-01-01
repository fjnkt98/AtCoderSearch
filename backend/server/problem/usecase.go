package problem

import (
	"context"
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/server/utility"

	"github.com/goark/errs"
)

type ProblemUsecase interface {
	Search(ctx context.Context, params SearchParams) (solr.SelectResponse[Problem, SolrFacetCounts], error)
}

type problemUsecase struct {
	core solr.SolrCore
}

func NewProblemUsecase(core solr.SolrCore) ProblemUsecase {
	return &problemUsecase{
		core: core,
	}
}

func (u *problemUsecase) Search(ctx context.Context, params SearchParams) (solr.SelectResponse[Problem, SolrFacetCounts], error) {
	q := params.Query()
	res, err := solr.SelectWithContext[Problem, SolrFacetCounts](ctx, u.core, q)
	if err != nil {
		return res, errs.New(
			"failed to execute select query",
			errs.WithCause(err),
		)
	}

	return res, nil
}

type Problem struct {
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
}

type SolrFacetCounts struct {
	Category   *solr.TermFacetCount       `json:"category,omitempty"`
	Color      *solr.TermFacetCount       `json:"color,omitempty"`
	Difficulty *solr.RangeFacetCount[int] `json:"difficulty,omitempty"`
}

func (f *SolrFacetCounts) Into(p FacetParams) map[string][]utility.FacetPart {
	res := make(map[string][]utility.FacetPart)
	if f.Category != nil {
		res["category"] = utility.ConvertBucket(f.Category.Buckets)
	}
	if f.Color != nil {
		res["color"] = utility.ConvertBucket[string](f.Color.Buckets)
	}
	if f.Difficulty != nil {
		res["difficulty"] = utility.ConvertRangeBucket(*f.Difficulty, p.Difficulty)
	}

	return res
}
