//go:build test_repository

package repository

import (
	"context"
	"testing"
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
		t.Fatalf("failed to save submissions: %s", err.Error())
	}

}
