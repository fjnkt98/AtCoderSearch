package api

import (
	"fjnkt98/atcodersearch/pkg/solr"
	"reflect"
	"testing"
)

func TestRows(t *testing.T) {
	cases := []struct {
		name     string
		value    ParameterBase
		expected int
	}{
		{name: "nil", value: ParameterBase{Limit: nil}, expected: 20},
		{name: "zero", value: ParameterBase{Limit: ptr(0)}, expected: 0},
		{name: "some", value: ParameterBase{Limit: ptr(100)}, expected: 100},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.value.Rows()
			if tt.expected != actual {
				t.Errorf("expected %+v, but got %+v", tt.expected, actual)
			}
		})
	}
}

func TestStart(t *testing.T) {
	cases := []struct {
		name     string
		value    ParameterBase
		expected int
	}{
		{name: "nil-zero", value: ParameterBase{Limit: nil, Page: 0}, expected: 0},
		{name: "nil-one", value: ParameterBase{Limit: nil, Page: 1}, expected: 0},
		{name: "nil-two", value: ParameterBase{Limit: nil, Page: 2}, expected: 20},
		{name: "zero-zero", value: ParameterBase{Limit: ptr(0), Page: 0}, expected: 0},
		{name: "zero-one", value: ParameterBase{Limit: ptr(0), Page: 1}, expected: 0},
		{name: "zero-two", value: ParameterBase{Limit: ptr(0), Page: 2}, expected: 0},
		{name: "some-zero", value: ParameterBase{Limit: ptr(50), Page: 0}, expected: 0},
		{name: "some-one", value: ParameterBase{Limit: ptr(50), Page: 1}, expected: 0},
		{name: "some-two", value: ParameterBase{Limit: ptr(50), Page: 2}, expected: 50},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.value.Start()
			if tt.expected != actual {
				t.Errorf("expected %+v, but got %+v", tt.expected, actual)
			}
		})
	}
}

func TestParseSort(t *testing.T) {
	cases := []struct {
		name     string
		value    []string
		expected []string
	}{
		{name: "nil", value: nil, expected: []string{}},
		{name: "empty", value: []string{}, expected: []string{}},
		{name: "single-asc", value: []string{"score"}, expected: []string{"score asc"}},
		{name: "single-desc", value: []string{"-score"}, expected: []string{"score desc"}},
		{name: "multiple-asc", value: []string{"score", "id"}, expected: []string{"score asc", "id asc"}},
		{name: "multiple-desc", value: []string{"-score", "-id"}, expected: []string{"score desc", "id desc"}},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := ParseSort(tt.value)
			if !reflect.DeepEqual(tt.expected, actual) {
				t.Errorf("expected %+v, but got %+v", tt.expected, actual)
			}
		})
	}
}

func TestFacetCountsFromStringBucket(t *testing.T) {
	cases := []struct {
		name     string
		value    []solr.StringBucket
		expected []FacetCount
	}{
		{name: "nil", value: nil, expected: []FacetCount{}},
		{name: "empty", value: []solr.StringBucket{}, expected: []FacetCount{}},
		{name: "some", value: []solr.StringBucket{{Val: "foo", Count: 1}, {Val: "bar", Count: 2}}, expected: []FacetCount{{Label: "foo", Count: 1}, {Label: "bar", Count: 2}}},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := FacetCountsFromStringBucket(tt.value)
			if !reflect.DeepEqual(tt.expected, actual) {
				t.Errorf("expected %+v, but got %+v", tt.expected, actual)
			}
		})
	}
}

func TestFacetCountsFromRangeBucket(t *testing.T) {
	cases := []struct {
		name     string
		value    []solr.RangeBucket
		expected []FacetCount
	}{
		{name: "nil", value: nil, expected: []FacetCount{}},
		{name: "empty", value: []solr.RangeBucket{}, expected: []FacetCount{}},
		{name: "some", value: []solr.RangeBucket{
			{Begin: nil, End: ptr(0), Count: 0},
			{Begin: ptr(0), End: ptr(100), Count: 1},
			{Begin: ptr(100), End: nil, Count: 2},
		}, expected: []FacetCount{
			{Label: " ~ 0", Count: 0},
			{Label: "0 ~ 100", Count: 1},
			{Label: "100 ~ ", Count: 2},
		}},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := FacetCountsFromRangeBucket(tt.value)
			if !reflect.DeepEqual(tt.expected, actual) {
				t.Errorf("expected \n%+v\n, but got \n%+v\n", tt.expected, actual)
			}
		})
	}
}

func TestNewFacetCount(t *testing.T) {
	cases := []struct {
		name     string
		value    *solr.JSONFacetResponse
		expected map[string][]FacetCount
	}{
		{name: "nil", value: nil, expected: nil},
		{name: "terms-is-nil", value: &solr.JSONFacetResponse{Range: map[string]solr.RangeFacetCount{"difficulty": {Buckets: nil}}}, expected: map[string][]FacetCount{"difficulty": {}}},
		{name: "range-is-nil", value: &solr.JSONFacetResponse{Terms: map[string]solr.TermsFacetCount{"category": {Buckets: nil}}}, expected: map[string][]FacetCount{"category": {}}},
		{name: "both-is-nil", value: &solr.JSONFacetResponse{Terms: nil, Range: nil}, expected: map[string][]FacetCount{}},
		{name: "normal", value: &solr.JSONFacetResponse{
			Terms: map[string]solr.TermsFacetCount{"category": {Buckets: nil}},
			Range: map[string]solr.RangeFacetCount{"difficulty": {Buckets: nil}},
		}, expected: map[string][]FacetCount{"category": {}, "difficulty": {}}},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := NewFacetCount(tt.value)
			if !reflect.DeepEqual(tt.expected, actual) {
				t.Errorf("expected %+v, but got %+v", tt.expected, actual)
			}
		})
	}
}
