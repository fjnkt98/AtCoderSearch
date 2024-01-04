//go:build test_repository

package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/goark/errs"
)

func TestSaveUser(t *testing.T) {
	db := getTestDB()
	repository := NewUserRepository(db)

	users := []User{
		{
			UserName:      "tourist",
			Rating:        3863,
			HighestRating: 4229,
			Affiliation:   ptr("TMO University"),
			BirthYear:     ptr(1994),
			Country:       ptr("BY"),
			Crown:         ptr("crown_champion"),
			JoinCount:     59,
			Rank:          1,
			ActiveRank:    ptr(1),
			Wins:          22,
		},
	}

	ctx := context.Background()
	if err := repository.Save(ctx, users); err != nil {
		t.Fatalf("failed to save user: %s", err.Error())
	}
}

func TestFetchRatingByUserName(t *testing.T) {
	db := getTestDB()
	repository := NewUserRepository(db)

	ctx := context.Background()
	if _, err := repository.FetchRatingByUserName(ctx, "fjnkt98"); err != nil && !errs.Is(err, sql.ErrNoRows) {
		t.Fatalf("failed to fetch rating by user name: %s", err.Error())
	}
}
