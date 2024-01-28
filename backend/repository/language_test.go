//go:build test_repository

package repository

import (
	"context"
	"testing"
)

func TestSaveLanguages(t *testing.T) {
	db, err := getTestDB()
	if err != nil {
		t.Fatalf("%+v", err)
	}
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

func TestFetchLanguages(t *testing.T) {
	db, err := getTestDB()
	if err != nil {
		t.Fatalf("%+v", err)
	}

	repository := NewLanguageRepository(db)
	ctx := context.Background()
	if _, err := repository.FetchLanguagesByGroup(ctx, nil); err != nil {
		t.Errorf("failed to fetch languages: %s", err.Error())
	}

	if _, err := repository.FetchLanguagesByGroup(ctx, []string{"C++", "Python"}); err != nil {
		t.Errorf("failed to fetch languages: %s", err.Error())
	}
}

func TestFetchLanguageGroups(t *testing.T) {
	db, err := getTestDB()
	if err != nil {
		t.Fatalf("%+v", err)
	}

	repository := NewLanguageRepository(db)
	ctx := context.Background()
	if _, err := repository.FetchLanguageGroups(ctx); err != nil {
		t.Fatalf("failed to fetch language groups: %s", err.Error())
	}
}
