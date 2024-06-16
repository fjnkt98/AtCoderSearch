package solr

import (
	"fjnkt98/atcodersearch/settings"
	"net/url"
	"reflect"
	"testing"
)

func TestMoreLikeThisQuery(t *testing.T) {
	core, err := NewSolrCore("http://localhost:18983", settings.PROBLEM_CORE_NAME)
	if err != nil {
		t.Fatalf("failed to create core: %s", err.Error())
	}

	actual := core.NewMoreLikeThis().
		Start(0).
		Rows(3).
		Fl("problemId").
		Q("abc300_a").
		Qf("text_ja").
		MinTF(2).
		MinDF(2).
		MaxDF(5).
		MinWL(3).
		MaxWL(10).
		Raw()

	expected := url.Values{
		"start": []string{"0"},
		"rows":  []string{"3"},
		"fl":    []string{"problemId"},
		"q":     []string{"{!mlt qf=text_ja mintf=2 mindf=2 maxdf=5 minwl=3 maxwl=10}abc300_a"},
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %+v, but got %+v", expected, actual)
	}
}
