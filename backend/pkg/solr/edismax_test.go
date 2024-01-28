package solr

import (
	"net/url"
	"reflect"
	"testing"
)

func TestWithNoParams(t *testing.T) {
	params := NewEDisMaxQueryBuilder()

	expected, _ := url.ParseQuery("defType=edismax")

	if !reflect.DeepEqual(params.inner, expected) {
		t.Errorf("expected `%s` but got `%s`", expected.Encode(), params.inner.Encode())
	}
}

func TestWithCommonParams(t *testing.T) {
	params := NewEDisMaxQueryBuilder().Start(10).Rows(20).Fq([]string{"name:alice", "{!collapse field=grade}"}).Fl("id,name,grade")

	expected, _ := url.ParseQuery("defType=edismax&start=10&rows=20&fq=name:alice&fq={!collapse field=grade}&fl=id,name,grade")

	if !reflect.DeepEqual(params.inner, expected) {
		t.Errorf("expected `%s` but got `%s`", expected.Encode(), params.inner.Encode())
	}
}
