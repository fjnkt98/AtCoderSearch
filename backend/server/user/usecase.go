package user

import (
	"context"
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/server/utility"

	"github.com/goark/errs"
)

type UserUsecase interface {
	Search(ctx context.Context, params SearchParams) (solr.SelectResponse[User, SolrFacetCounts], error)
}

type userUsecase struct {
	core solr.SolrCore
}

func NewUserUsecase(core solr.SolrCore) UserUsecase {
	return &userUsecase{
		core: core,
	}
}

func (u *userUsecase) Search(ctx context.Context, params SearchParams) (solr.SelectResponse[User, SolrFacetCounts], error) {
	q := params.Query()
	res, err := solr.SelectWithContext[User, SolrFacetCounts](ctx, u.core, q)
	if err != nil {
		return res, errs.New(
			"failed to execute select query",
			errs.WithCause(err),
		)
	}

	return res, nil
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

type SolrFacetCounts struct {
	Rating    *solr.RangeFacetCount[int] `json:"rating"`
	BirthYear *solr.RangeFacetCount[int] `json:"birth_year"`
	JoinCount *solr.RangeFacetCount[int] `json:"join_count"`
	Country   *solr.TermFacetCount       `json:"country"`
}

func (f *SolrFacetCounts) Into(p FacetParams) map[string][]utility.FacetPart {
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
