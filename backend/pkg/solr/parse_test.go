package solr

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []SearchWord
	}{
		{name: "single word", input: "foo", want: []SearchWord{{Word: "foo"}}},
		{name: "multiple word", input: "foo bar", want: []SearchWord{{Word: "foo"}, {Word: "bar"}}},
		{name: "leading whitespace included", input: "     foo  bar ", want: []SearchWord{{Word: "foo"}, {Word: "bar"}}},
		{name: "negative search", input: "-foo", want: []SearchWord{{Word: "foo", Negative: true}}},
		{name: "negative mixed search", input: "bar -foo", want: []SearchWord{{Word: "bar"}, {Word: "foo", Negative: true}}},
		{name: "empty negative search word", input: "-", want: []SearchWord{}},
		{name: "consequence negative search", input: "----foo", want: []SearchWord{{Word: "---foo", Negative: true}}},
		{name: "phrase search", input: `"foo bar"`, want: []SearchWord{{Word: "foo bar", Phrase: true}}},
		{name: "single word phrase search", input: `"foo"`, want: []SearchWord{{Word: "foo", Phrase: true}}},
		{name: "mixed phrase word", input: `foo "bar baz"`, want: []SearchWord{{Word: "foo"}, {Word: "bar baz", Phrase: true}}},
		{name: "negative phrase search", input: `-"bar baz"`, want: []SearchWord{{Word: "bar baz", Negative: true, Phrase: true}}},
		{name: "empty phrase search word", input: `""`, want: []SearchWord{}},
		{name: "invalid phrase syntax", input: `"foo bar`, want: []SearchWord{{Word: "foo bar", Phrase: true}}},
		{name: "invalid negative phrase syntax", input: `-"foo bar`, want: []SearchWord{{Word: "foo bar", Negative: true, Phrase: true}}},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := Parse(tt.input)

			if !reflect.DeepEqual(tt.want, actual) {
				t.Errorf("expected %+v, but got %+v", tt.want, actual)
			}
		})
	}
}
