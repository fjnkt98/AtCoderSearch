package acs

import "testing"

type SampleDocument struct {
	ID    string `json:"id" solr:"id"`
	Name  string `json:"name" solr:"name"`
	Grade string `json:"grade" solr:"grade"`
	Class string `json:"class" solr:"class"`
}

func TestFieldList(t *testing.T) {
	lister := NewFieldLister()
	doc := SampleDocument{}

	fieldList, err := lister.FieldList(doc)
	if err != nil {
		t.Errorf("test failed: %s", err.Error())
	}
	expected := "id,name,grade,class"

	if fieldList != expected {
		t.Errorf("expected `%s` but got `%s`", expected, fieldList)
	}
}
