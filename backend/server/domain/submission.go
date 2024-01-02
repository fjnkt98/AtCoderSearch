package domain

import (
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/server/utility"
	"fmt"
	"net/url"
	"strings"
)

type SearchSubmissionParam struct {
	Sort []string `json:"sort" schema:"sort" validate:"dive,oneof=-score rating -rating birth_year -birth_year"`
	utility.SearchParams[SearchSubmissionFilterParam, SearchSubmissionFacetParam]
}

func (p *SearchSubmissionParam) GetSort() string {
	orders := make([]string, 0, len(p.Sort))
	for _, s := range p.Sort {
		if strings.HasPrefix(s, "-") {
			orders = append(orders, fmt.Sprintf("%s desc", s[1:]))
		} else {
			orders = append(orders, fmt.Sprintf("%s asc", s))
		}
	}

	return strings.Join(orders, ",")
}

type SearchSubmissionFilterParam struct {
	EpochSecond   utility.IntegerRange `json:"epoch_second" schema:"epoch_second"`
	ProblemID     []string             `json:"problem_id" schema:"problem_id"`
	ContestID     []string             `json:"contest_id" schema:"contest_id"`
	Category      []string             `json:"category" schema:"category"`
	UserID        []string             `json:"user_id" schema:"user_id"`
	Language      []string             `json:"language" schema:"language"`
	LanguageGroup []string             `json:"language_group" schema:"language_group"`
	Point         utility.FloatRange   `json:"point" schema:"point"`
	Length        utility.IntegerRange `json:"length" schema:"length"`
	Result        []string             `json:"result" schema:"result"`
	ExecutionTime utility.IntegerRange `json:"execution_time" schema:"execution_time"`
}

type SearchSubmissionFacetParam struct {
	Term          []string                `json:"term" schema:"term" validate:"dive,oneof=problem_id user_id language language_group result contest_id" facet:"problem_id:problem_id,user_id:user_id,language:language,language_group:language_group,result:result,contest_id:contest_id"`
	Length        utility.RangeFacetParam `json:"length" schema:"length" facet:"length:length"`
	ExecutionTime utility.RangeFacetParam `json:"execution_time" schema:"execution_time" facet:"execution_time:execution_time"`
}

func (p *SearchSubmissionParam) Query() url.Values {
	return solr.NewLuceneQueryBuilder().
		Facet(p.GetFacet()).
		Fl(strings.Join(utility.FieldList(new(Submission)), ",")).
		Fq(p.GetFilter()).
		Op("AND").
		Q("*:*").
		Rows(p.GetRows()).
		Sort(p.GetSort()).
		Start(p.GetStart()).
		Build()
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

type SearchSubmissionFacetCounts struct {
	ContestID     *solr.TermFacetCount       `json:"contest_id,omitempty"`
	ProblemID     *solr.TermFacetCount       `json:"problem_id,omitempty"`
	UserID        *solr.TermFacetCount       `json:"user_id,omitempty"`
	Language      *solr.TermFacetCount       `json:"language,omitempty"`
	LanguageGroup *solr.TermFacetCount       `json:"language_group,omitempty"`
	Result        *solr.TermFacetCount       `json:"result,omitempty"`
	Length        *solr.RangeFacetCount[int] `json:"length,omitempty"`
	ExecutionTime *solr.RangeFacetCount[int] `json:"execution_time,omitempty"`
}

func (f *SearchSubmissionFacetCounts) Into(p SearchSubmissionFacetParam) map[string][]utility.FacetPart {
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
