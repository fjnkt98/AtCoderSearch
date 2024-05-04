package api

import "testing"

func ptr[T any](v T) *T {
	return &v
}

func TestLocalParamStringer(t *testing.T) {
	p := LocalParam("tag", "category")
	want := `tag=category`
	actual := p.String()

	if want != actual {
		t.Errorf("expected \n%s\n, but got \n%s\n", want, actual)
	}
}

func TestTermsFilter(t *testing.T) {
	cases := []struct {
		name   string
		values []string
		field  string
		params []localParam
		want   string
	}{
		{name: "empty value", values: []string{""}, field: "category", want: ""},
		{name: "single value", values: []string{"ABC"}, field: "category", want: `category:("ABC")`},
		{name: "multiple values", values: []string{"ABC", "Other Contests"}, field: "category", want: `category:("ABC" OR "Other Contests")`},
		{name: "with local param", values: []string{"ABC"}, field: "category", params: []localParam{LocalParam("tag", "category")}, want: `{!tag=category}category:("ABC")`},
		{name: "with local params", values: []string{"ABC"}, field: "category", params: []localParam{LocalParam("tag", "category"), LocalParam("v", "'foo bar'")}, want: `{!tag=category v='foo bar'}category:("ABC")`},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := TermsFilter(tt.values, tt.field, tt.params...)

			if tt.want != actual {
				t.Errorf("expected \n%s\n, but got \n%s\n", tt.want, actual)
			}
		})
	}
}

func TestIntegerRangeFilter(t *testing.T) {
	cases := []struct {
		name   string
		from   *int
		to     *int
		field  string
		params []localParam
		want   string
	}{
		{name: "empty", from: nil, to: nil, field: "difficulty", want: ""},
		{name: "from only", from: ptr(0), to: nil, field: "difficulty", want: "difficulty:[0 TO *]"},
		{name: "to only", from: nil, to: ptr(100), field: "difficulty", want: "difficulty:[* TO 100]"},
		{name: "from and to", from: ptr(-100), to: ptr(100), field: "difficulty", want: "difficulty:[-100 TO 100]"},
		{name: "with local param", from: ptr(-100), to: ptr(100), field: "difficulty", params: []localParam{LocalParam("tag", "difficulty")}, want: "{!tag=difficulty}difficulty:[-100 TO 100]"},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := IntegerRangeFilter(tt.from, tt.to, tt.field, tt.params...)

			if tt.want != actual {
				t.Errorf("expected \n%s\n, but got \n%s\n", tt.want, actual)
			}
		})
	}
}

func TestFloatRangeFilter(t *testing.T) {
	cases := []struct {
		name   string
		from   *float64
		to     *float64
		field  string
		params []localParam
		want   string
	}{
		{name: "empty", from: nil, to: nil, field: "difficulty", want: ""},
		{name: "from only", from: ptr(0.0), to: nil, field: "difficulty", want: "difficulty:[0.00 TO *]"},
		{name: "to only", from: nil, to: ptr(100.0), field: "difficulty", want: "difficulty:[* TO 100.00]"},
		{name: "from and to", from: ptr(-100.0), to: ptr(100.0), field: "difficulty", want: "difficulty:[-100.00 TO 100.00]"},
		{name: "with local param", from: ptr(-100.0), to: ptr(100.0), field: "difficulty", params: []localParam{LocalParam("tag", "difficulty")}, want: "{!tag=difficulty}difficulty:[-100.00 TO 100.00]"},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := FloatRangeFilter(tt.from, tt.to, tt.field, tt.params...)

			if tt.want != actual {
				t.Errorf("expected \n%s\n, but got \n%s\n", tt.want, actual)
			}
		})
	}
}
