package solr

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestMarshalJSONFacetQuery(t *testing.T) {
	q := NewJSONFacetQuery(
		NewTermsFacetQuery("category").Name("cate").Limit(-1).ExcludeTags("category"),
		NewRangeFacetQuery("difficulty", 0, 2000, 400).Name("diff").Other("all").ExcludeTags("difficulty"),
	)

	b, err := json.Marshal(q)
	if err != nil {
		t.Fatalf("failed to marshal JSONFacetQuery: %s", err.Error())
	}
	actual := string(b)
	expected := `{"cate":{"domain":{"excludeTags":["category"]},"field":"category","limit":-1,"type":"terms"},"diff":{"domain":{"excludeTags":["difficulty"]},"end":2000,"field":"difficulty","gap":400,"other":"all","start":0,"type":"range"}}`

	if expected != actual {
		t.Errorf("mismatch result: expected \n%s\n, but got \n%s\n", expected, actual)
	}
}

func TestTermsFacetQuery(t *testing.T) {
	q := NewTermsFacetQuery("category").
		Limit(-1).
		MinCount(0).
		ExcludeTags("category")

	expected := &termsFacetQuery{
		name: "category",
		params: map[string]any{
			"type":     "terms",
			"field":    "category",
			"limit":    -1,
			"mincount": 0,
			"domain": map[string]any{
				"excludeTags": []string{"category"},
			},
		},
	}

	if !reflect.DeepEqual(q, expected) {
		t.Errorf("mismatch result: expected %+v but got %+v", expected, q)
	}
}

func TestUnmarshalJSONFacetResponse(t *testing.T) {
	s := []byte(`
	{
		"count": 2481,
		"category":{"buckets":[{"val":"ABC","count":1213},{"val":"ARC","count":615},{"val":"AGC","count":320}]},
		"difficulty":{"buckets":[{"val":0,"count":255},{"val":400,"count":310},{"val":800,"count":283}],"before":{"count":0},"after":{"count":0},"between":{"count":2481}}
	}
	`)
	want := RawJSONFacetResponse(map[string]json.RawMessage{
		"category":   json.RawMessage([]byte(`{"buckets":[{"val":"ABC","count":1213},{"val":"ARC","count":615},{"val":"AGC","count":320}]}`)),
		"difficulty": json.RawMessage([]byte(`{"buckets":[{"val":0,"count":255},{"val":400,"count":310},{"val":800,"count":283}],"before":{"count":0},"after":{"count":0},"between":{"count":2481}}`)),
	})

	var actual RawJSONFacetResponse
	if err := json.Unmarshal(s, &actual); err != nil {
		t.Fatalf("failed to unmarshal: %s", err.Error())
	}

	if !reflect.DeepEqual(want, actual) {
		t.Errorf("expected \n%+v\n, but got \n%+v\n", want, actual)
	}
}

func ptr[T any](v T) *T {
	return &v
}

func TestParseJSONFacetResponse(t *testing.T) {
	s := []byte(`
	{
		"count": 2481,
		"category":{"buckets":[{"val":"ABC","count":1213},{"val":"ARC","count":615},{"val":"AGC","count":320}]},
		"difficulty":{"buckets":[{"val":0,"count":255},{"val":400,"count":310},{"val":800,"count":283}],"before":{"count":0},"after":{"count":0},"between":{"count":2481}}
	}
	`)

	q := NewJSONFacetQuery(
		NewTermsFacetQuery("category").MinCount(0),
		NewRangeFacetQuery("difficulty", 0, 1200, 400).MinCount(0),
	)

	var raw RawJSONFacetResponse
	if err := json.Unmarshal(s, &raw); err != nil {
		t.Fatalf("failed to unmarshal: %s", err.Error())
	}

	actual, err := raw.Parse(q)
	if err != nil {
		t.Fatalf("failed to parse json facet response: %s", err.Error())
	}

	want := &JSONFacetResponse{
		Terms: map[string]TermsFacetCount{
			"category": {
				Buckets: []StringBucket{{Val: "ABC", Count: 1213}, {Val: "ARC", Count: 615}, {Val: "AGC", Count: 320}},
			},
		},
		Range: map[string]RangeFacetCount{
			"difficulty": {
				Buckets: []RangeBucket{{End: ptr(0)}, {Begin: ptr(0), End: ptr(400), Count: 255}, {Begin: ptr(400), End: ptr(800), Count: 310}, {Begin: ptr(800), End: ptr(1200), Count: 283}, {Begin: ptr(1200)}},
			},
		},
	}

	if !reflect.DeepEqual(want, actual) {
		t.Errorf("expected \n%+v\n, but got \n%+v\n", want, actual)
	}
}
