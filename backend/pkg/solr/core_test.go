//go:build test_solr

// Run this tests with the Docker container started with the following command.
//
// ```
// docker run --rm -d -p 18983:8983 --name solr_example solr:9.1.0 solr-precreate example
// ```

package solr

import (
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"
)

type SampleDocument struct {
	ID string `json:"id"`
}

func TestStatus(t *testing.T) {
	core, err := NewSolrCore("http://localhost:18983", "example")
	if err != nil {
		t.Fatalf("failed to create solr core: %s", err.Error())
	}
	status, err := Status(core)
	if err != nil {
		t.Errorf("failed to get core status: %s", err.Error())
	}

	if status.Name != "example" {
		t.Errorf("different core status name: expected `example` but got `%s`", status.Name)
	}
}

func TestReload(t *testing.T) {
	core, err := NewSolrCore("http://localhost:18983", "example")
	if err != nil {
		t.Fatalf("failed to create solr core: %s", err.Error())
	}

	before := time.Now()
	Reload(core)
	status, _ := Status(core)
	after, _ := time.Parse(time.RFC3339, status.StartTime)

	if before.After(after) {
		t.Errorf("invalid reloaded time: expected before(%s) < after(%s)", before, after)
	}

	duration := after.Sub(before)
	if !(duration.Abs().Milliseconds() < 1000) {
		t.Errorf("expected that duration time %s is less than 1000[ms]", duration)
	}
}

func TestPing(t *testing.T) {
	core, err := NewSolrCore("http://localhost:18983", "example")
	if err != nil {
		t.Fatalf("failed to create solr core: %s", err.Error())
		return
	}

	res, _ := Ping(core)
	if res.Status != "OK" {
		t.Errorf("ping returns non-ok message: expected `OK` but got `%s`", res.Status)
	}
}

func TestScenario(t *testing.T) {
	core, err := NewSolrCore("http://localhost:18983", "example")
	if err != nil {
		t.Fatalf("failed to create solr core: %s", err.Error())
	}

	if _, err := Truncate(core); err != nil {
		t.Fatalf("failed to truncate core: %s", err.Error())
	}
	if _, err := Commit(core); err != nil {
		t.Fatalf("failed to commit core: %s", err.Error())
	}

	params := url.Values{}
	params.Set("q", "*:*")

	res, err := Select[SampleDocument, any](core, params)
	if err != nil {
		t.Fatalf("failed to select document")
	}
	if res.Response.NumFound != 0 {
		t.Fatalf("unmatched number of document: expected 0, but got %d", res.Response.NumFound)
	}

	document := strings.NewReader(`[{"id":"001"}]`)
	if _, err := Post(core, document, "application/json"); err != nil {
		t.Fatalf("failed to post document: %s", err.Error())
	}

	if _, err := Commit(core); err != nil {
		t.Fatalf("failed to commit document")
	}

	res, err = Select[SampleDocument, any](core, params)
	if err != nil {
		t.Fatalf("failed to select document")
	}
	if res.Response.NumFound != 1 {
		t.Fatalf("unmatched number of document")
	}

	want := []SampleDocument{{ID: "001"}}
	if !reflect.DeepEqual(res.Response.Docs, want) {
		t.Fatalf("collection doesn't match the expected: %s", err.Error())
	}

	if _, err := Truncate(core); err != nil {
		t.Fatalf("failed to truncate core: %s", err.Error())
	}
	if _, err := Commit(core); err != nil {
		t.Fatalf("failed to commit document")
	}

	res, err = Select[SampleDocument, any](core, params)
	if err != nil {
		t.Fatalf("failed to select document")
	}
	if res.Response.NumFound != 0 {
		t.Fatalf("unmatched number of document: expected 0, but got %d", res.Response.NumFound)
	}
}
