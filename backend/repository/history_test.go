//go:build test_repository

package repository

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	db := getTestDB()
	ctx := context.Background()
	if _, err := db.NewDelete().Model(new(UpdateHistory)).Where("0 = 0").Exec(ctx); err != nil {
		fmt.Printf("failed to delete records from `update_history`: %s", err.Error())
		os.Exit(1)
	}
	if _, err := db.NewDelete().Model(new(SubmissionCrawlHistory)).Where("0 = 0").Exec(ctx); err != nil {
		fmt.Printf("failed to delete records from `submission_crawl_history`: %s", err.Error())
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func TestSaveAndGetSubmissionCrawlHistory(t *testing.T) {
	db := getTestDB()
	repository := NewSubmissionCrawlHistoryRepository(db)

	// Get the latest history
	// expect an empty history to be returned
	ctx := context.Background()
	got, err := repository.GetLatestHistory(ctx, "abc300")
	if err != nil {
		t.Fatalf("failed to get latest history: %s", err.Error())
	}
	expected := SubmissionCrawlHistory{ContestID: "abc300"}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("expected empty struct %+v, but got %+v", expected, got)
	}

	// Save history
	history := SubmissionCrawlHistory{
		StartedAt: 10000000,
		ContestID: "abc300",
	}

	if err := repository.Save(ctx, history); err != nil {
		t.Fatalf("failed to save submission crawl history: %s", err.Error())
	}

	// Get the latest history again
	// expect the history which has been saved some time ago to be returned
	got, err = repository.GetLatestHistory(ctx, "abc300")
	if err != nil {
		t.Fatalf("failed to get latest history: %s", err.Error())
	}
	if !reflect.DeepEqual(history, got) {
		t.Fatalf("expected empty struct %+v, but got %+v", history, got)
	}
}

func TestSaveAndGetUpdateHistory(t *testing.T) {
	db := getTestDB()
	repository := NewUpdateHistoryRepository(db)

	ctx := context.Background()
	got, err := repository.GetLatest(ctx, "problem")
	if err != nil {
		t.Fatalf("failed to get update history: %s", err.Error())
	}
	expected := UpdateHistory{Domain: "problem", Status: "finished"}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("expected default history %+v, but got %+v", expected, got)
	}

	history := NewUpdateHistory("problem", "{}")
	if err := repository.Finish(ctx, &history); err != nil {
		t.Fatalf("failed to save finished update history: %s", err.Error())
	}

	got, err = repository.GetLatest(ctx, "problem")
	if err != nil {
		t.Fatalf("failed to get update history: %s", err.Error())
	}

	if !history.wasSaved {
		t.Fatalf("`wasSaved` of the saved history is still false")
	}
}
