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

func TestGetCoreStatus(t *testing.T) {
	core, _ := NewSolrCore("example", "http://localhost:8983")
	status, err := Status(core)
	if err != nil {
		t.Errorf("failed to get core status: %s", err.Error())
	}

	if status.Name != "example" {
		t.Errorf("different core status name: expected `example` but got `%s`", status.Name)
	}
}

func TestReloadCore(t *testing.T) {
	core, _ := NewSolrCore("example", "http://localhost:8983")

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
	core, _ := NewSolrCore("example", "http://localhost:8983")

	res, _ := Ping(core)
	if res.Status != "OK" {
		t.Errorf("ping returns non-ok message: expected `OK` but got `%s`", res.Status)
	}
}

func TestSelect(t *testing.T) {
	core, _ := NewSolrCore("example", "http://localhost:8983")

	params := url.Values{}
	params.Set("q", "*:*")

	res, err := Select[any, any](core, params)
	if err != nil {
		t.Errorf("failed to select test: %s", err.Error())
	}

	if res.Header.Status != 0 {
		t.Errorf("select request returns non-0 status: %d", res.Header.Status)
	}
}
