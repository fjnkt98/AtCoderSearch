//go:build test_repository

package repository

import (
	"context"
	"testing"
)

func TestSaveProblem(t *testing.T) {
	db := getTestDB()
	repository := NewProblemRepository(db)

	problems := []Problem{
		{
			ProblemID:    "abc300_a",
			ContestID:    "abc300",
			Name:         "A. N-choice question",
			Title:        "N-choice question",
			ProblemIndex: "A",
			URL:          "https://atcoder.jp/contests/abc300/tasks/abc300_a",
			HTML:         "",
		},
	}

	ctx := context.Background()
	if err := repository.Save(ctx, problems); err != nil {
		t.Fatalf("failed to save problems: %s", err.Error())
	}
}

func TestFetchIDs(t *testing.T) {
	db := getTestDB()
	repository := NewProblemRepository(db)

	ctx := context.Background()
	_, err := repository.FetchIDs(ctx)
	if err != nil {
		t.Fatalf("failed to fetch problem ids: %s", err.Error())
	}
}

func TestFetchIDsByContestID(t *testing.T) {
	db := getTestDB()
	repository := NewProblemRepository(db)

	cases := []struct {
		name string
		args []string
		want error
	}{
		{name: "one", args: []string{"abc300"}, want: nil},
		{name: "many", args: []string{"abc300", "abc301"}, want: nil},
		{name: "empty", args: []string{}, want: nil},
		{name: "nil", args: nil, want: nil},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			_, err := repository.FetchIDsByContestID(ctx, tt.args)
			if err != tt.want {
				t.Errorf("failed to fetch problem ids by contest id: %s", err.Error())
			}
		})
	}

}
