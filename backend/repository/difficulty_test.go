//go:build test_repository

package repository

import (
	"context"
	"testing"
)

func TestSaveDifficulty(t *testing.T) {
	db, err := getTestDB()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	repository := NewDifficultyRepository(db)

	difficulties := []Difficulty{
		{
			ProblemID:        "abc019_1",
			Slope:            ptr(-0.0013332305474178768),
			Intercept:        ptr(7.39230222481235),
			Variance:         ptr(1.4089923968242148),
			Difficulty:       ptr(int64(-22)),
			Discrimination:   ptr(0.004479398673070138),
			IrtLogLikelihood: ptr(-39.42457119583683),
			IrtUsers:         ptr(281.0),
			IsExperimental:   ptr(true),
		},
	}

	ctx := context.Background()
	if err := repository.Save(ctx, difficulties); err != nil {
		t.Fatalf("failed to save difficulties: %s", err.Error())
	}
}
