//go:build test_repository

package repository

import (
	"context"
	"testing"
)

func TestSaveSubmission(t *testing.T) {
	db := getTestDB()
	repository := NewSubmissionRepository(db)

	submissions := []Submission{
		{
			ID:            48864123,
			EpochSecond:   1703593405,
			ProblemID:     "abc171_c",
			ContestID:     "abc171",
			UserID:        "OmameBeans",
			Language:      "C++ 20 (gcc 12.2)",
			Point:         300,
			Length:        2615,
			Result:        "AC",
			ExecutionTime: ptr(1),
		},
	}

	ctx := context.Background()
	if err := repository.Save(ctx, submissions); err != nil {
		t.Fatalf("failed to save submissions: %s", err.Error())
	}
}

func TestFetchLanguagesFromSubmission(t *testing.T) {
	db := getTestDB()
	repository := NewSubmissionRepository(db)

	ctx := context.Background()
	if _, err := repository.FetchLanguages(ctx); err != nil {
		t.Fatalf("failed to fetch languages: %s", err.Error())
	}
}
