package submission

import (
	"context"
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/server/utility"

	"github.com/goark/errs"
)

type SubmissionUsecase interface {
	Search(ctx context.Context, params SearchParams) (solr.SelectResponse[Submission, SolrFacetCounts], error)
}

type submissionUsecase struct {
	core solr.SolrCore
}

func NewSubmissionUsecase(core solr.SolrCore) SubmissionUsecase {
	return &submissionUsecase{
		core: core,
	}
}

func (u *submissionUsecase) Search(ctx context.Context, params SearchParams) (solr.SelectResponse[Submission, SolrFacetCounts], error) {
	q := params.Query()
	res, err := solr.SelectWithContext[Submission, SolrFacetCounts](ctx, u.core, q)
	if err != nil {
		return res, errs.New(
			"failed to execute select query",
			errs.WithCause(err),
		)
	}

	return res, nil
}

type Submission struct {
	SubmissionID  int64                 `json:"submission_id"`
	SubmittedAt   solr.FromSolrDateTime `json:"submitted_at"`
	SubmissionURL string                `json:"submission_url"`
	ProblemID     string                `json:"problem_id"`
	ProblemTitle  string                `json:"problem_title"`
	ContestID     string                `json:"contest_id"`
	ContestTitle  string                `json:"contest_title"`
	Category      string                `json:"category"`
	Difficulty    int                   `json:"difficulty"`
	Color         string                `json:"color"`
	UserID        string                `json:"user_id"`
	Language      string                `json:"language"`
	Point         float64               `json:"point"`
	Length        int64                 `json:"length"`
	Result        string                `json:"result"`
	ExecutionTime *int64                `json:"execution_time"`
}

type SolrFacetCounts struct {
	ContestID     *solr.TermFacetCount       `json:"contest_id,omitempty"`
	ProblemID     *solr.TermFacetCount       `json:"problem_id,omitempty"`
	UserID        *solr.TermFacetCount       `json:"user_id,omitempty"`
	Language      *solr.TermFacetCount       `json:"language,omitempty"`
	LanguageGroup *solr.TermFacetCount       `json:"language_group,omitempty"`
	Result        *solr.TermFacetCount       `json:"result,omitempty"`
	Length        *solr.RangeFacetCount[int] `json:"length,omitempty"`
	ExecutionTime *solr.RangeFacetCount[int] `json:"execution_time,omitempty"`
}

func (f *SolrFacetCounts) Into(p FacetParams) map[string][]utility.FacetPart {
	res := make(map[string][]utility.FacetPart)
	if f.ContestID != nil {
		res["contest_id"] = utility.ConvertBucket(f.ContestID.Buckets)
	}
	if f.ProblemID != nil {
		res["problem_id"] = utility.ConvertBucket(f.ProblemID.Buckets)
	}
	if f.UserID != nil {
		res["user_id"] = utility.ConvertBucket(f.UserID.Buckets)
	}
	if f.Language != nil {
		res["language"] = utility.ConvertBucket(f.Language.Buckets)
	}
	if f.LanguageGroup != nil {
		res["language_group"] = utility.ConvertBucket(f.LanguageGroup.Buckets)
	}
	if f.Result != nil {
		res["result"] = utility.ConvertBucket(f.Result.Buckets)
	}
	if f.Length != nil {
		res["length"] = utility.ConvertRangeBucket(*f.Length, p.Length)
	}
	if f.ExecutionTime != nil {
		res["execution_time"] = utility.ConvertRangeBucket(*f.ExecutionTime, p.ExecutionTime)
	}

	return res
}
