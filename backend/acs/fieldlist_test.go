package acs

import "testing"

type SampleDocument struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Grade string `json:"grade,omitempty"`
	Class string `json:"-"`
}

func TestFieldList(t *testing.T) {
	doc := SampleDocument{}

	fieldList := FieldList(doc)
	expected := "id,name,grade"

	if fieldList != expected {
		t.Errorf("expected `%s` but got `%s`", expected, fieldList)
	}
}
