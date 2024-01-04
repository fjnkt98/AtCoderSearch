package domain

import (
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/server/utility"
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/text/unicode/norm"
)

type SearchUserParam struct {
	Keyword string   `json:"keyword" schema:"keyword" validate:"lte=200"`
	Sort    []string `json:"sort" schema:"sort" validate:"dive,oneof=-score rating -rating birth_year -birth_year"`
	utility.SearchParam[SearchUserFilterParam, SearchUserFacetParam]
}

func (p *SearchUserParam) GetSort() string {
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

type SearchUserFilterParam struct {
	UserID    []string             `json:"user_id" schema:"user_id"`
	Rating    utility.IntegerRange `json:"rating" schema:"rating"`
	BirthYear utility.IntegerRange `json:"birth_year" schema:"birth_year"`
	JoinCount utility.IntegerRange `json:"join_count" schema:"join_count"`
	Country   []string             `json:"country" schema:"country"`
	Color     []string             `json:"color" schema:"color"`
}

type SearchUserFacetParam struct {
	Term      []string                `json:"term" schema:"term" validate:"dive,oneof=country" facet:"country:country"`
	Rating    utility.RangeFacetParam `json:"rating" schema:"rating" facet:"rating:rating"`
	BirthYear utility.RangeFacetParam `json:"birth_year" schema:"birth_year" facet:"birth_year:birth_year"`
	JoinCount utility.RangeFacetParam `json:"join_count" schema:"join_count" facet:"join_count:join_count"`
}

func (p *SearchUserParam) Query() url.Values {
	return solr.NewEDisMaxQueryBuilder().
		Facet(p.GetFacet()).
		Fl(strings.Join(utility.FieldList(new(User)), ",")).
		Fq(p.GetFilter()).
		Op("AND").
		Q(solr.Sanitize(norm.NFKC.String(p.Keyword))).
		QAlt("*:*").
		Qf("text_unigram").
		Rows(p.GetRows()).
		Sort(p.GetSort()).
		Sow(true).
		Start(p.GetStart()).
		Build()
}

type User struct {
	UserName      string  `json:"user_name"`
	Rating        int     `json:"rating"`
	HighestRating int     `json:"highest_rating"`
	Affiliation   *string `json:"affiliation"`
	BirthYear     *int    `json:"birth_year"`
	Country       *string `json:"country"`
	Crown         *string `json:"crown"`
	JoinCount     int     `json:"join_count"`
	Rank          int     `json:"rank"`
	ActiveRank    *int    `json:"active_rank"`
	Wins          int     `json:"wins" `
	Color         string  `json:"color"`
	UserURL       string  `json:"user_url"`
}

type SearchUserFacetCounts struct {
	Rating    *solr.RangeFacetCount[int] `json:"rating"`
	BirthYear *solr.RangeFacetCount[int] `json:"birth_year"`
	JoinCount *solr.RangeFacetCount[int] `json:"join_count"`
	Country   *solr.TermFacetCount       `json:"country"`
}

func (f *SearchUserFacetCounts) Into(p SearchUserFacetParam) map[string][]utility.FacetPart {
	res := make(map[string][]utility.FacetPart)
	if f.Rating != nil {
		res["rating"] = utility.ConvertRangeBucket(*f.Rating, p.Rating)
	}
	if f.BirthYear != nil {
		res["birth_year"] = utility.ConvertRangeBucket(*f.BirthYear, p.BirthYear)
	}
	if f.JoinCount != nil {
		res["join_count"] = utility.ConvertRangeBucket(*f.JoinCount, p.JoinCount)
	}
	if f.Country != nil {
		res["country"] = utility.ConvertBucket(f.Country.Buckets)
	}

	return res
}
