//go:build test_repository

package repository

import (
	"context"
	"testing"
)

func TestSaveLanguages(t *testing.T) {
	db := getTestDB()
	repository := NewLanguageRepository(db)

	languages := []Language{
		{
			Language: "C++ 20 (gcc 12.2)",
			Group:    "C++",
		},
		{
			Language: "Python (CPython 3.11.4)",
			Group:    "Python",
		},
	}
	ctx := context.Background()
	if err := repository.Save(ctx, languages); err != nil {
		t.Fatalf("failed to save languages: %s", err.Error())
	}
}
