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
