package utility

import (
	"net/url"
	"reflect"
	"testing"
	"time"
)

func ptr[T any](v T) *T {
	return &v
}

func TestIntegerRange(t *testing.T) {
	cases := []struct {
		name string
		from *int
		to   *int
		want string
	}{
		{name: "fulfill", from: ptr(100), to: ptr(200), want: "[100 TO 200}"},
		{name: "from_is_nil", from: nil, to: ptr(200), want: "[* TO 200}"},
		{name: "to_is_nil", from: ptr(100), to: nil, want: "[100 TO *}"},
		{name: "both_are_nil", from: nil, to: nil, want: ""},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := IntegerRange{From: tt.from, To: tt.to}
			result := r.ToRange()
			if result != tt.want {
				t.Errorf("expected %s, but got %s", tt.want, result)
			}
		})
	}
}

func TestFloatRange(t *testing.T) {
	cases := []struct {
		name string
		from *float64
		to   *float64
		want string
	}{
		{name: "fulfill", from: ptr(100.0), to: ptr(200.0), want: "[100.000000 TO 200.000000}"},
		{name: "from_is_nil", from: nil, to: ptr(200.0), want: "[* TO 200.000000}"},
		{name: "to_is_nil", from: ptr(100.0), to: nil, want: "[100.000000 TO *}"},
		{name: "both_are_nil", from: nil, to: nil, want: ""},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := FloatRange{From: tt.from, To: tt.to}
			result := r.ToRange()
			if result != tt.want {
				t.Errorf("expected %s, but got %s", tt.want, result)
			}
		})
	}
}

func TestDateRange(t *testing.T) {
	tz, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatalf("failed to load timezone: %s", err.Error())
	}

	cases := []struct {
		name string
		from *time.Time
		to   *time.Time
		want string
	}{
		{name: "fulfill", from: ptr(time.Date(1998, 7, 15, 0, 0, 0, 0, tz)), to: ptr(time.Date(1998, 8, 22, 0, 0, 0, 0, tz)), want: "[1998-07-14T15:00:00Z TO 1998-08-21T15:00:00Z}"},
		{name: "from_is_nil", from: nil, to: ptr(time.Date(1998, 8, 22, 0, 0, 0, 0, tz)), want: "[* TO 1998-08-21T15:00:00Z}"},
		{name: "to_is_nil", from: ptr(time.Date(1998, 7, 15, 0, 0, 0, 0, tz)), to: nil, want: "[1998-07-14T15:00:00Z TO *}"},
		{name: "both_are_nil", from: nil, to: nil, want: ""},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := DateRange{From: tt.from, To: tt.to}
			result := r.ToRange()
			if result != tt.want {
				t.Errorf("expected %s, but got %s", tt.want, result)
			}
		})
	}
}

func TestSanitizeStrings(t *testing.T) {
	cases := []struct {
		name string
		s    []string
		want []string
	}{
		{name: "normal", s: []string{"C++", "Python"}, want: []string{`C\+\+`, "Python"}},
		{name: "one", s: []string{"C++"}, want: []string{`C\+\+`}},
		{name: "empty", s: []string{}, want: []string{}},
		{name: "nil", s: nil, want: []string{}},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeStrings(tt.s)
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("expected %+v, but got %+v", tt.want, result)
			}
		})
	}
}

func TestQuoteStrings(t *testing.T) {
	cases := []struct {
		name string
		s    []string
		want []string
	}{
		{name: "normal", s: []string{"ABC", "Other Contests"}, want: []string{`"ABC"`, `"Other Contests"`}},
		{name: "one", s: []string{"ABC"}, want: []string{`"ABC"`}},
		{name: "empty", s: []string{}, want: []string{}},
		{name: "nil", s: nil, want: []string{}},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := QuoteStrings(tt.s)
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("expected %+v, but got %+v", tt.want, result)
			}
		})
	}
}

func TestRangeFacet(t *testing.T) {
	p := RangeFacetParam{From: ptr(0), To: ptr(2000), Gap: ptr(400)}
	result := p.ToFacet("difficulty", "difficulty")

	want := map[string]map[string]any{
		"difficulty": {
			"type":  "range",
			"field": "difficulty",
			"start": 0,
			"end":   2000,
			"gap":   400,
			"other": "all",
			"domain": map[string]any{
				"excludeTags": []string{"difficulty"},
			},
		},
	}

	if !reflect.DeepEqual(result, want) {
		t.Errorf("expected %+v, but got %+v", want, result)
	}
}

func TestLackedRangeFacet(t *testing.T) {
	cases := []struct {
		name  string
		param RangeFacetParam
	}{
		{name: "from_is_nil", param: RangeFacetParam{From: nil, To: ptr(2000), Gap: ptr(400)}},
		{name: "to_is_nil", param: RangeFacetParam{From: ptr(0), To: nil, Gap: ptr(400)}},
		{name: "gap_is_nil", param: RangeFacetParam{From: ptr(0), To: ptr(2000), Gap: nil}},
		{name: "from_is_not_nil", param: RangeFacetParam{From: ptr(0), To: nil, Gap: nil}},
		{name: "to_is_not_nil", param: RangeFacetParam{From: nil, To: ptr(2000), Gap: nil}},
		{name: "gap_is_not_nil", param: RangeFacetParam{From: nil, To: nil, Gap: ptr(400)}},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := tt.param.ToFacet("foo", "foo")

			if result != nil {
				t.Errorf("expected nil, but got %+v", result)
			}
		})
	}
}

func TestTermFacet(t *testing.T) {
	p := TermFacetParam{"category", "difficulty", "color"}
	result := p.ToFacet(map[string]string{"category": "contest_category", "difficulty": "difficulty"})

	want := map[string]map[string]any{
		"category": {
			"type":     "terms",
			"field":    "contest_category",
			"limit":    -1,
			"mincount": 0,
			"sort":     "index",
			"domain": map[string]any{
				"excludeTags": []string{"contest_category"},
			},
		},
		"difficulty": {
			"type":     "terms",
			"field":    "difficulty",
			"limit":    -1,
			"mincount": 0,
			"sort":     "index",
			"domain": map[string]any{
				"excludeTags": []string{"difficulty"},
			},
		},
		"color": {
			"type":     "terms",
			"field":    "color",
			"limit":    -1,
			"mincount": 0,
			"sort":     "index",
			"domain": map[string]any{
				"excludeTags": []string{"color"},
			},
		},
	}

	if !reflect.DeepEqual(result, want) {
		t.Errorf("expected %+v, but got %+v", want, result)
	}
}

func TestFacet(t *testing.T) {
	type param struct {
		Term       TermFacetParam  `facet:"category:category,color:color_facet"`
		Difficulty RangeFacetParam `facet:"difficulty:difficulty"`
	}

	p := param{
		Term:       TermFacetParam{"category", "color"},
		Difficulty: RangeFacetParam{From: ptr(200), To: ptr(1000), Gap: ptr(100)},
	}

	facet := Facet(p)
	want := `{"category":{"domain":{"excludeTags":["category"]},"field":"category","limit":-1,"mincount":0,"sort":"index","type":"terms"},"color":{"domain":{"excludeTags":["color_facet"]},"field":"color_facet","limit":-1,"mincount":0,"sort":"index","type":"terms"},"difficulty":{"domain":{"excludeTags":["difficulty"]},"end":1000,"field":"difficulty","gap":100,"other":"all","start":200,"type":"range"}}`

	if facet != want {
		t.Errorf("expected\n%s\n, but got \n%s\n", want, facet)
	}
}

func TestFilter(t *testing.T) {
	type param struct {
		Category   []string     `filter:"category"`
		Color      []string     `filter:"color,unquote"`
		Difficulty IntegerRange `filter:"difficulty"`
		User       []string     `filter:"-"`
		Country    []string
	}

	p := param{
		Category:   []string{"ABC", "Other Contests"},
		Color:      []string{"blue"},
		Difficulty: IntegerRange{From: ptr(0), To: ptr(2000)},
		User:       []string{"fjnkt98"},
		Country:    []string{"JP"},
	}

	fq := Filter(p)
	want := []string{
		`{!tag=category}category:("ABC" OR "Other Contests")`,
		`{!tag=color}color:(blue)`,
		`{!tag=difficulty}difficulty:[0 TO 2000}`,
		`{!tag=Country}Country:("JP")`,
	}

	if !reflect.DeepEqual(fq, want) {
		t.Errorf("expected %+v, but got %+v", want, fq)
	}
}

func TestSearchParamGetRows(t *testing.T) {
	cases := []struct {
		name  string
		param SearchParam[any, any]
		want  int
	}{
		{name: "normal", param: SearchParam[any, any]{Limit: ptr(50)}, want: 50},
		{name: "empty", param: SearchParam[any, any]{Limit: nil}, want: 20},
		{name: "zero", param: SearchParam[any, any]{Limit: ptr(0)}, want: 0},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := tt.param.GetRows()
			if result != tt.want {
				t.Errorf("expected %d, but got %d", tt.want, result)
			}
		})
	}
}

func TestSearchParamGetStart(t *testing.T) {
	cases := []struct {
		name  string
		param SearchParam[any, any]
		want  int
	}{
		{name: "normal", param: SearchParam[any, any]{Limit: ptr(20), Page: 5}, want: 80},
		{name: "one", param: SearchParam[any, any]{Limit: ptr(20), Page: 1}, want: 0},
		{name: "zero", param: SearchParam[any, any]{Limit: ptr(20), Page: 0}, want: 0},
		{name: "page_when_limit_zero", param: SearchParam[any, any]{Limit: ptr(0), Page: 1}, want: 0},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := tt.param.GetStart()
			if result != tt.want {
				t.Errorf("expected %d, but got %d", tt.want, result)
			}
		})
	}
}

func TestSearchParamGetSort(t *testing.T) {
	cases := []struct {
		name  string
		param SearchParam[any, any]
		want  string
	}{
		{name: "normal", param: SearchParam[any, any]{Sort: []string{"score"}}, want: "score asc"},
		{name: "multiple", param: SearchParam[any, any]{Sort: []string{"-score", "rating"}}, want: "score desc,rating asc"},
		{name: "empty", param: SearchParam[any, any]{Sort: []string{}}, want: ""},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := tt.param.GetSort()
			if result != tt.want {
				t.Errorf("expected %s, but got %s", tt.want, result)
			}
		})
	}
}

func TestSearchParamFilterAndFacet(t *testing.T) {
	type filterParams struct {
		Category []string `filter:"category"`
	}

	type facetParams struct {
		Term TermFacetParam `facet:"category:category"`
	}
	type params struct {
		Keyword string `json:"keyword" schema:"keyword"`
		SearchParam[filterParams, facetParams]
	}

	decoder := NewSchemaDecoder()

	v := url.Values{}
	v.Set("keyword", "foo")
	v.Set("limit", "100")
	v.Set("page", "1")
	v.Set("sort", "score,-start_at")
	v.Set("filter.category", "ABC")
	v.Set("facet.term", "category")
	v.Set("facet.term", "category")

	var p params
	if err := decoder.Decode(&p, v); err != nil {
		t.Fatalf("failed to decode query parameter: %s", err.Error())
	}

	{
		want := params{
			Keyword: "foo",
			SearchParam: SearchParam[filterParams, facetParams]{
				Limit: ptr(100),
				Page:  1,
				Filter: filterParams{
					Category: []string{"ABC"},
				},
				Sort: []string{"score", "-start_at"},
				Facet: facetParams{
					Term: TermFacetParam{"category"},
				},
			},
		}

		if !reflect.DeepEqual(p, want) {
			t.Errorf("expected %+v, but got %+v", want, p)
		}
	}

	{
		want := []string{
			`{!tag=category}category:("ABC")`,
		}
		fq := p.GetFilter()

		if !reflect.DeepEqual(fq, want) {
			t.Errorf("expected %+v, but got %+v", want, fq)
		}
	}

	{
		want := `{"category":{"domain":{"excludeTags":["category"]},"field":"category","limit":-1,"mincount":0,"sort":"index","type":"terms"}}`
		facet := p.GetFacet()

		if !reflect.DeepEqual(facet, want) {
			t.Errorf("expected %+v, but got %+v", want, facet)
		}
	}
}

func TestFieldList(t *testing.T) {
	type doc struct {
		ID    string
		Name  string `json:"name"`
		Grade string `json:"grade,omitempty"`
		Class string `json:"-"`
	}

	fl := FieldList(new(doc))
	want := []string{"ID", "name", "grade"}

	if !reflect.DeepEqual(fl, want) {
		t.Errorf("expected %+v, but got %+v", want, fl)
	}
}

func TestDecodeTermFacetParamType(t *testing.T) {
	type param struct {
		Term TermFacetParam `schema:"term"`
	}

	decoder := NewSchemaDecoder()

	cases := []struct {
		name   string
		raw    string
		want   TermFacetParam
		assert func(expected, actual TermFacetParam) bool
	}{
		{name: "empty", raw: "term=", want: TermFacetParam([]string{}), assert: func(expected, actual TermFacetParam) bool { return len(actual) == 0 }},
		{name: "one", raw: "term=category", want: TermFacetParam([]string{"category"}), assert: func(expected, actual TermFacetParam) bool { return reflect.DeepEqual(expected, actual) }},
		{name: "many", raw: "term=category,difficulty", want: TermFacetParam([]string{"category", "difficulty"}), assert: func(expected, actual TermFacetParam) bool { return reflect.DeepEqual(expected, actual) }},
		{name: "include_spaces", raw: "term=  category  ,     difficulty ", want: TermFacetParam([]string{"category", "difficulty"}), assert: func(expected, actual TermFacetParam) bool { return reflect.DeepEqual(expected, actual) }},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			query, err := url.ParseQuery(tt.raw)
			if err != nil {
				t.Fatalf("failed to parse query: %s", err.Error())
			}

			var result param
			err = decoder.Decode(&result, query)
			if err != nil {
				t.Fatalf("failed to decode query: %s", err.Error())
			}

			if !tt.assert(tt.want, result.Term) {
				t.Errorf("expected %+v, but got %+v", tt.want, result.Term)
			}
		})
	}

}
