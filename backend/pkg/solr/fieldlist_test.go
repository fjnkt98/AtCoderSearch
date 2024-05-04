package solr

import (
	"reflect"
	"testing"
)

type doc struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Grade  int
	Secret string `json:"-"`
}

func TestFieldList(t *testing.T) {
	expected := []string{"id", "name", "Grade"}
	actual := FieldList(new(doc))

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected \n%+v\n , but got \n%+v\n", expected, actual)
	}
}
