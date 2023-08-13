// Run this tests with the Docker container started with the following command.
//
// ```
// docker run --rm -d -p 8983:8983 --name solr_example solr:9.1.0 solr-precreate example
// ```

package solr

import (
	"net/url"
	"testing"
	"time"
)

type SampleDocument struct {
	ID string `json:"id"`
}

func TestCreateCore(t *testing.T) {
	core, err := NewSolrCore[SampleDocument, map[string]any]("example", "http://localhost:8983")

	if err != nil {
		t.Error("failed to create core")
	}

	if core.adminURL.String() != "http://localhost:8983/solr/admin/cores" {
		t.Errorf("admin URL doesn't match: expected: `http://localhost:8983/solr/admin/cores` but got `%s`", core.adminURL.String())
	}
	if core.pingURL.String() != "http://localhost:8983/solr/example/admin/ping" {
		t.Errorf("admin URL doesn't match: expected: `http://localhost:8983/solr/example/admin/ping` but got `%s`", core.adminURL.String())
	}
	if core.selectURL.String() != "http://localhost:8983/solr/example/select" {
		t.Errorf("admin URL doesn't match: expected: `http://localhost:8983/solr/example/select` but got `%s`", core.adminURL.String())
	}
	if core.postURL.String() != "http://localhost:8983/solr/example/update" {
		t.Errorf("admin URL doesn't match: expected: `http://localhost:8983/solr/example/update` but got `%s`", core.adminURL.String())
	}
}

func TestGetCoreStatus(t *testing.T) {
	core, _ := NewSolrCore[SampleDocument, map[string]any]("example", "http://localhost:8983")
	status, err := core.Status()
	if err != nil {
		t.Errorf("failed to get core status: %s", err.Error())
	}

	if status.Name != "example" {
		t.Errorf("different core status name: expected `example` but got `%s`", status.Name)
	}
}

func TestReloadCore(t *testing.T) {
	core, _ := NewSolrCore[SampleDocument, map[string]any]("example", "http://localhost:8983")

	before := time.Now()
	core.Reload()
	status, _ := core.Status()
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
	core, _ := NewSolrCore[SampleDocument, map[string]any]("example", "http://localhost:8983")

	res, _ := core.Ping()
	if res.Status != "OK" {
		t.Errorf("ping returns non-ok message: expected `OK` but got `%s`", res.Status)
	}
}

func TestSelect(t *testing.T) {
	core, _ := NewSolrCore[SampleDocument, map[string]any]("example", "http://localhost:8983")

	params := url.Values{}
	params.Set("q", "*:*")

	res, err := core.Select(params)
	if err != nil {
		t.Errorf("failed to select test: %s", err.Error())
	}

	if res.Header.Status != 0 {
		t.Errorf("select request returns non-0 status: %d", res.Header.Status)
	}
}
