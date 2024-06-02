package search

import (
	"fjnkt98/atcodersearch/pkg/solr"
	"fjnkt98/atcodersearch/settings"
	"reflect"
	"testing"
)

var core = solr.MustNewSolrCore("http://localhost:18983", settings.PROBLEM_CORE_NAME)

func ptr[T any](v T) *T {
	return &v
}

func TestSearchProblemFiltering(t *testing.T) {
	cases := []struct {
		name     string
		param    ProblemParameter
		expected []string
	}{
		{name: "default", param: ProblemParameter{}, expected: ([]string)(nil)},
		{name: "category 1", param: ProblemParameter{Category: []string{"ABC"}}, expected: []string{`{!tag=category}category:("ABC")`}},
		{name: "category 2", param: ProblemParameter{Category: []string{"ABC", "Other Contests"}}, expected: []string{`{!tag=category}category:("ABC" OR "Other Contests")`}},
		{name: "color 1", param: ProblemParameter{Color: []string{"green"}}, expected: []string{`{!tag=color}color:("green")`}},
		{name: "color 2", param: ProblemParameter{Color: []string{"green", "cyan"}}, expected: []string{`{!tag=color}color:("green" OR "cyan")`}},
		{name: "difficulty from", param: ProblemParameter{DifficultyFrom: ptr(400)}, expected: []string{`{!tag=difficulty}difficulty:[400 TO *]`}},
		{name: "difficulty to", param: ProblemParameter{DifficultyTo: ptr(800)}, expected: []string{`{!tag=difficulty}difficulty:[* TO 800]`}},
		{name: "difficulty from to", param: ProblemParameter{DifficultyFrom: ptr(400), DifficultyTo: ptr(800)}, expected: []string{`{!tag=difficulty}difficulty:[400 TO 800]`}},
		{name: "is experimental", param: ProblemParameter{Experimental: ptr(true)}, expected: []string{`isExperimental:true`}},
		{name: "is not experimental", param: ProblemParameter{Experimental: ptr(false)}, expected: []string{`isExperimental:false`}},
		{name: "excludeSolved without userId", param: ProblemParameter{ExcludeSolved: true}, expected: ([]string)(nil)},
		{name: "excludeSolved with userId", param: ProblemParameter{ExcludeSolved: true, UserID: "fjnkt98"}, expected: []string{`-{!join fromIndex=submissions from=problemId to=problemId v='+userId:"fjnkt98" +result:AC'}`}},
		{name: "excludeSolved with userId contains special character", param: ProblemParameter{ExcludeSolved: true, UserID: "C++"}, expected: []string{`-{!join fromIndex=submissions from=problemId to=problemId v='+userId:"C\+\+" +result:AC'}`}},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := map[string][]string(tt.param.Query(core).Raw())["fq"]

			if !reflect.DeepEqual(tt.expected, actual) {
				t.Errorf("expected \n%#v\n, but got \n%#v\n", tt.expected, actual)
			}
		})
	}
}

func TestSearchProblemBoosting(t *testing.T) {
	cases := []struct {
		name     string
		param    ProblemParameter
		expected []string
	}{
		{name: "default", param: ProblemParameter{}, expected: ([]string)(nil)},
		{name: "prioritize recent", param: ProblemParameter{PrioritizeRecent: true}, expected: []string{`{!boost b=7}{!func}pow(2,mul(-1,div(ms(NOW,startAt),2592000000)))`}},
		{name: "rating", param: ProblemParameter{Difficulty: ptr(1000)}, expected: []string{`{!boost b=10}{!func}pow(2.71828182846,mul(-1,div(pow(sub(1000,difficulty),2),20000)))`}},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := map[string][]string(tt.param.Query(core).Raw())["bq"]

			if !reflect.DeepEqual(tt.expected, actual) {
				t.Errorf("expected \n%#v\n, but got \n%#v\n", tt.expected, actual)
			}
		})
	}
}

func TestSearchProblemFacet(t *testing.T) {
	cases := []struct {
		name     string
		param    ProblemParameter
		expected []string
	}{
		{name: "default", param: ProblemParameter{}, expected: []string{"{}"}},
		{name: "category", param: ProblemParameter{Facet: []string{"category"}}, expected: []string{`{"category":{"domain":{"excludeTags":["category"]},"field":"category","limit":-1,"mincount":0,"sort":"index","type":"terms"}}`}},
		{name: "difficulty", param: ProblemParameter{Facet: []string{"difficulty"}}, expected: []string{`{"difficulty":{"domain":{"excludeTags":["difficulty"]},"end":4000,"field":"difficulty","gap":400,"other":"all","start":0,"type":"range"}}`}},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := map[string][]string(tt.param.Query(core).Raw())["json.facet"]

			if !reflect.DeepEqual(tt.expected, actual) {
				t.Errorf("expected \n%#v\n, but got \n%#v\n", tt.expected, actual)
			}
		})
	}
}

func TestSearchProblemSort(t *testing.T) {
	cases := []struct {
		name     string
		param    ProblemParameter
		expected []string
	}{
		{name: "empty", param: ProblemParameter{}, expected: []string{"problemId asc"}},
		{name: "single", param: ProblemParameter{Sort: []string{"-score"}}, expected: []string{"score desc,problemId asc"}},
		{name: "multiple", param: ProblemParameter{Sort: []string{"startedAt", "-score"}}, expected: []string{"startedAt asc,score desc,problemId asc"}},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := map[string][]string(tt.param.Query(core).Raw())["sort"]

			if !reflect.DeepEqual(tt.expected, actual) {
				t.Errorf("expected \n%#v\n, but got \n%#v\n", tt.expected, actual)
			}
		})
	}
}
