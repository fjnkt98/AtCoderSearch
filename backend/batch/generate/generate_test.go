package generate

import (
	"encoding/json"
	"fjnkt98/atcodersearch/pkg/solr"
	"reflect"
	"testing"
	"time"
)

type document struct {
	ID          string
	Name        string                `solr:"name"`
	Description string                `solr:"description,text_ja,text_en"`
	Birthday    solr.IntoSolrDateTime `solr:"birthday"`
	Class       string                `solr:"-"`
}

func TestStructToMap(t *testing.T) {
	tz, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatalf("failed to load timezone: %s", err.Error())
	}

	doc := document{
		ID:          "001",
		Name:        "fjnkt98",
		Description: "junior software engineer",
		Birthday:    solr.IntoSolrDateTime(time.Date(1998, 7, 15, 0, 0, 0, 0, tz)),
		Class:       "A01",
	}

	expanded := StructToMap(doc)
	expected := map[string]any{
		"ID":                   "001",
		"name":                 "fjnkt98",
		"description":          "junior software engineer",
		"description__text_ja": "junior software engineer",
		"description__text_en": "junior software engineer",
		"birthday":             solr.IntoSolrDateTime(time.Date(1998, 7, 15, 0, 0, 0, 0, tz)),
	}

	if !reflect.DeepEqual(expanded, expected) {
		t.Errorf("expected `%+v`, but got `%+v`", expected, expanded)
	}

	var serialized string
	res, err := json.Marshal(expanded)
	if err != nil {
		t.Errorf("failed to marshal map")
		return
	}
	serialized = string(res)
	want := `{"ID":"001","birthday":"1998-07-14T15:00:00Z","description":"junior software engineer","description__text_en":"junior software engineer","description__text_ja":"junior software engineer","name":"fjnkt98"}`
	if serialized != want {
		t.Errorf("expected `%s`, but got `%s`", want, serialized)
	}
}
