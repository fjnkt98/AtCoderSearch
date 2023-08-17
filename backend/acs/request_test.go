package acs

import (
	"reflect"
	"testing"
	"time"
)

func TestIntegerRange(t *testing.T) {
	a := 0
	b := 100
	cases := map[string]struct {
		r    IntegerRange[int]
		want string
	}{
		"0-100": {IntegerRange[int]{From: &a, To: &b}, "[0 TO 100}"},
		"0-*":   {IntegerRange[int]{From: &a, To: nil}, "[0 TO *}"},
		"*-100": {IntegerRange[int]{From: nil, To: &b}, "[* TO 100}"},
		"*-*":   {IntegerRange[int]{From: nil, To: nil}, ""},
	}

	for name, tt := range cases {
		if actual := tt.r.ToRange(); actual != tt.want {
			t.Errorf("%s: expected: %s, actual: %s", name, tt.want, actual)
		}
	}
}

func TestFloatRange(t *testing.T) {
	a := 0.0
	b := 3.14159265
	cases := map[string]struct {
		r    FloatRange[float64]
		want string
	}{
		"0.0-3.141593": {FloatRange[float64]{From: &a, To: &b}, "[0.000000 TO 3.141593}"},
		"0.0-*":        {FloatRange[float64]{From: &a, To: nil}, "[0.000000 TO *}"},
		"*-3.141593":   {FloatRange[float64]{From: nil, To: &b}, "[* TO 3.141593}"},
		"*-*":          {FloatRange[float64]{From: nil, To: nil}, ""},
	}

	for name, tt := range cases {
		if actual := tt.r.ToRange(); actual != tt.want {
			t.Errorf("%s: expected: %s, actual: %s", name, tt.want, actual)
		}
	}
}

func TestDateRange(t *testing.T) {
	a := time.Date(2023, 7, 15, 10, 30, 0, 0, time.Local)
	b := time.Date(2023, 8, 13, 20, 10, 0, 0, time.Local)
	cases := map[string]struct {
		r    DateRange
		want string
	}{
		"7/15-8/13": {DateRange{From: &a, To: &b}, "[2023-07-15T01:30:00Z TO 2023-08-13T11:10:00Z}"},
		"7/15-*":    {DateRange{From: &a, To: nil}, "[2023-07-15T01:30:00Z TO *}"},
		"*-8/13":    {DateRange{From: nil, To: &b}, "[* TO 2023-08-13T11:10:00Z}"},
		"*-*":       {DateRange{From: nil, To: nil}, ""},
	}

	for name, tt := range cases {
		if actual := tt.r.ToRange(); actual != tt.want {
			t.Errorf("%s: expected: %s, actual: %s", name, tt.want, actual)
		}
	}
}

func TestSanitizeStrings(t *testing.T) {
	cases := map[string]struct {
		s []string
		want []string
	}{
		"normal": {[]string{"a", "b", "c"}, []string{"a", "b", "c"}},
		"contains whitespace": {[]string{"a", "   ", "c"}, []string{"a", "c"}},
		"solr special characters": {[]string{"AND", "OR", "C++"}, []string{"\\AND", "\\OR", "C\\+\\+"}},
	}

	for name, tt := range cases {
		if actual := SanitizeStrings(tt.s); !reflect.DeepEqual(actual, tt.want) {
			t.Errorf("%s: expected: %s, actual: %s", name, tt.want, actual)
		}
	}
}