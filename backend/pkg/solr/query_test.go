package solr

import (
	"fjnkt98/atcodersearch/settings"
	"net/url"
	"reflect"
	"testing"
)

func TestSelectQuery(t *testing.T) {
	core, err := NewSolrCore("http://localhost:18983", settings.PROBLEM_CORE_NAME)
	if err != nil {
		t.Fatalf("failed to create core: %s", err.Error())
	}

	actual := core.NewSelect().
		Sort("score desc").
		Start(20).
		Rows(20).
		Fq("id:10", "name:bob").
		Fl("id,name,grade").
		Wt("json").
		Q("foo bar").
		Qf("text_ja text_en").
		Sow(true).
		Some(KeyValue{"qq", "bar"}, KeyValue{"qq2", "baz"}).
		Raw()

	expected := url.Values{
		"defType": []string{"edismax"},
		"sort":    []string{"score desc"},
		"start":   []string{"20"},
		"rows":    []string{"20"},
		"fq":      []string{"id:10", "name:bob"},
		"fl":      []string{"id,name,grade"},
		"wt":      []string{"json"},
		"q":       []string{"foo bar"},
		"qf":      []string{"text_ja text_en"},
		"sow":     []string{"true"},
		"qq":      []string{"bar"},
		"qq2":     []string{"baz"},
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected %+v, but got %+v", expected, actual)
	}

}
