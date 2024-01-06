package domain

import (
	"fjnkt98/atcodersearch/server/utility"
	"reflect"
	"testing"

	govalidator "github.com/go-playground/validator/v10"
)

func ptr[T any](v T) *T {
	return &v
}

func TestSearchUserParamValidation(t *testing.T) {
	cases := []struct {
		name  string
		param SearchUserParam
		isnil bool
	}{
		{name: "length_of_keyword_equals_to_200", param: SearchUserParam{Keyword: "1oeK4YdWnR4AQ3bu/vG+SH7OF08pbo6GspuS9MFso/GiDN/OaBOvp8PzNfUMu7lnjRwA3rEZNUBhV9xR+q7VbRSz1OK5YKGH7yEvsfq1hyNDgEOVwLBEYeXjodpMnaXYvB9sv/G2EdsQeF3hHZjLS5GX/25MA7jBEg3FSGAhy/cFG/GP3HFfSWlAUUUBdCeWjyfovEne"}, isnil: true},
		{name: "length_of_keyword_is_greater_than_200", param: SearchUserParam{Keyword: "kRCcgkuNLuQyBQiYV45wBGly/3/cb+0Yltg+KSsPhQNx24q+oTvdHAJjXFGaDQPc3UdLaQrzjF5JX1lrAzNCFoE4pPFoK66H4AROEu5VH2WwNxklENyiUYLv41PUyXAsq0s+aS08OenQvPLQOWkJVqEaWrRr0pvN9Z8SbqVE09odtKbC5W2/PGqT/mROS193usFJVtvSg"}, isnil: false},
		{name: "sort_valid", param: SearchUserParam{Sort: []string{"-score", "rating", "-rating", "birth_year", "-birth_year"}}, isnil: true},
		{name: "sort_invalid", param: SearchUserParam{Sort: []string{"score"}}, isnil: false},
		{name: "facet_valid", param: SearchUserParam{SearchParam: utility.SearchParam[SearchUserFilterParam, SearchUserFacetParam]{Facet: SearchUserFacetParam{Term: []string{"country"}}}}, isnil: true},
		{name: "facet_invalid", param: SearchUserParam{SearchParam: utility.SearchParam[SearchUserFilterParam, SearchUserFacetParam]{Facet: SearchUserFacetParam{Term: []string{"user_name"}}}}, isnil: false},
	}

	validator := govalidator.New()
	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Struct(tt.param)
			isnil := (err == nil)
			if isnil != tt.isnil {
				t.Errorf("expected %t, but got %t", tt.isnil, isnil)
			}
		})
	}
}

func TestSearchUserQuery(t *testing.T) {
	cases := []struct {
		name      string
		param     SearchUserParam
		wantField string
		wantValue []string
	}{
		{name: "keyword", param: SearchUserParam{Keyword: "fjnkt98"}, wantField: "q", wantValue: []string{"fjnkt98"}},
		{name: "sort", param: SearchUserParam{Sort: []string{"-rating", "birth_year"}}, wantField: "sort", wantValue: []string{"rating desc,birth_year asc"}},
		{name: "limit", param: SearchUserParam{SearchParam: utility.SearchParam[SearchUserFilterParam, SearchUserFacetParam]{Limit: ptr(50)}}, wantField: "rows", wantValue: []string{"50"}},
		{name: "empty_limit", param: SearchUserParam{}, wantField: "rows", wantValue: []string{"20"}},
		{name: "page_with_default_limit", param: SearchUserParam{SearchParam: utility.SearchParam[SearchUserFilterParam, SearchUserFacetParam]{Page: 2}}, wantField: "start", wantValue: []string{"20"}},
		{name: "page_with_specified_limit", param: SearchUserParam{SearchParam: utility.SearchParam[SearchUserFilterParam, SearchUserFacetParam]{Limit: ptr(50), Page: 3}}, wantField: "start", wantValue: []string{"100"}},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := tt.param.Query()

			if !reflect.DeepEqual(result[tt.wantField], tt.wantValue) {
				t.Errorf("expected %+v in field `%s`, but got %+v", tt.wantValue, tt.wantField, result[tt.wantField])
			}
		})
	}
}
