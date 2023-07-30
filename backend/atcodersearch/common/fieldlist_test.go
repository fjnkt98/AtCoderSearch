package common

import "testing"

type SampleDocument struct {
	ID    string `json:"id" solr:"id"`
	Name  string `json:"name" solr:"name"`
	Grade string `json:"grade" solr:"grade"`
	Class string `json:"class" solr:"class"`
}

func TestFieldList(t *testing.T) {
	doc := SampleDocument{}

	fieldList := FieldList(doc)
	expected := "id,name,grade,class"

	if fieldList != expected {
		t.Errorf("expected `%s` but got `%s`", expected, fieldList)
	}
}
