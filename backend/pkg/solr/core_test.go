package solr

import (
	"context"
	"fjnkt98/atcodersearch/settings"
	"reflect"
	"strings"
	"testing"
	"time"
)

type problem struct {
	ProblemID string `json:"problemId"`
}

func TestStatus(t *testing.T) {
	core, err := NewSolrCore("http://localhost:18983", settings.PROBLEM_CORE_NAME)
	if err != nil {
		t.Fatalf("failed to create solr core: %s", err.Error())
	}
	ctx := context.Background()
	status, err := core.Status(ctx)
	if err != nil {
		t.Errorf("failed to get core status: %s", err.Error())
	}

	if status.Name != settings.PROBLEM_CORE_NAME {
		t.Errorf("different core status name: expected `example` but got `%s`", status.Name)
	}
}

func TestReload(t *testing.T) {
	core, err := NewSolrCore("http://localhost:18983", settings.PROBLEM_CORE_NAME)
	if err != nil {
		t.Fatalf("failed to create solr core: %s", err.Error())
	}

	before := time.Now()
	ctx := context.Background()
	if _, err := core.Reload(ctx); err != nil {
		t.Fatalf("failed to reload core: %s", err.Error())
	}
	status, err := core.Status(ctx)
	if err != nil {
		t.Fatalf("failed to get status of the core: %s", err.Error())
	}
	after, err := time.Parse(time.RFC3339, status.StartTime)
	if err != nil {
		t.Fatalf("failed to parse the start time of the core: %s", err.Error())
	}

	if before.After(after) {
		t.Errorf("invalid reloaded time: expected before(%s) < after(%s)", before, after)
	}

	duration := after.Sub(before)
	if !(duration.Abs().Milliseconds() < 1000) {
		t.Errorf("expected that duration time %s is less than 1000[ms]", duration)
	}
}

func TestPing(t *testing.T) {
	core, err := NewSolrCore("http://localhost:18983", settings.PROBLEM_CORE_NAME)
	if err != nil {
		t.Fatalf("failed to create solr core: %s", err.Error())
		return
	}

	ctx := context.Background()
	res, err := core.Ping(ctx)
	if err != nil {
		t.Fatalf("failed to ping: %s", err.Error())
	}
	if res.Status != "OK" {
		t.Errorf("ping returns non-ok message: expected `OK` but got `%s`", res.Status)
	}
}

func TestScenario(t *testing.T) {
	core, err := NewSolrCore("http://localhost:18983", settings.PROBLEM_CORE_NAME)
	if err != nil {
		t.Fatalf("failed to create solr core: %s", err.Error())
	}

	ctx := context.Background()
	if _, err := core.Delete(ctx); err != nil {
		t.Fatalf("failed to delete the documents of the core: %s", err.Error())
	}
	if _, err := core.Commit(ctx); err != nil {
		t.Fatalf("failed to commit core: %s", err.Error())
	}

	res, err := core.NewSelect().Q("*:*").Exec(ctx)
	if err != nil {
		t.Fatalf("failed to select document")
	}
	if res.Raw.Response.NumFound != 0 {
		t.Fatalf("unmatched number of document: expected 0, but got %d", res.Raw.Response.NumFound)
	}

	document := strings.NewReader(`[{"problemId":"abc300_a"}]`)
	if _, err := core.Post(ctx, document, "application/json"); err != nil {
		t.Fatalf("failed to post document: %s", err.Error())
	}

	if _, err := core.Commit(ctx); err != nil {
		t.Fatalf("failed to commit document")
	}

	res, err = core.NewSelect().Q("*:*").Exec(ctx)
	if err != nil {
		t.Fatalf("failed to select document")
	}
	if res.Raw.Response.NumFound != 1 {
		t.Fatalf("unmatched number of document")
	}

	want := []problem{{ProblemID: "abc300_a"}}
	var actual []problem
	if err := res.Scan(&actual); err != nil {
		t.Fatalf("failed to scan the documents: %s", err.Error())
	}
	if !reflect.DeepEqual(actual, want) {
		t.Fatalf("collection doesn't match the expected: %s", err.Error())
	}

	if _, err := core.Delete(ctx); err != nil {
		t.Fatalf("failed to truncate core: %s", err.Error())
	}
	if _, err := core.Commit(ctx); err != nil {
		t.Fatalf("failed to commit document")
	}

	res, err = core.NewSelect().Q("*:*").Exec(ctx)
	if err != nil {
		t.Fatalf("failed to select document")
	}
	if res.Raw.Response.NumFound != 0 {
		t.Fatalf("unmatched number of document: expected 0, but got %d", res.Raw.Response.NumFound)
	}
}
