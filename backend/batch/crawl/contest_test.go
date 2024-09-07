package crawl

import (
	"context"
	"fjnkt98/atcodersearch/internal/testutil"
	"fjnkt98/atcodersearch/pkg/atcoder"
	"fjnkt98/atcodersearch/repository"
	"testing"
	"time"
)

func TestSaveContests(t *testing.T) {
	_, dsn, stop, err := testutil.CreateDBContainer()
	defer stop()

	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	pool, err := repository.NewPool(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now()

	t.Run("empty", func(t *testing.T) {
		contests := make([]atcoder.Contest, 0)
		count, err := SaveContests(ctx, pool, contests, now)
		if err != nil {
			t.Fatal(err)
		}

		if count != 0 {
			t.Errorf("count = %d, want 0", count)
		}
	})

	t.Run("insert", func(t *testing.T) {
		contests := []atcoder.Contest{
			{
				ID:               "abc001",
				StartEpochSecond: 1468670400,
				DurationSecond:   6000,
				Title:            "AtCoder Beginner Contest 001",
				RateChange:       "-",
			},
		}

		count, err := SaveContests(ctx, pool, contests, now)
		if err != nil {
			t.Fatal(err)
		}
		if count != 1 {
			t.Errorf("count = %d, want 1", count)
		}
	})

	t.Run("multiple", func(t *testing.T) {
		contests := []atcoder.Contest{
			{
				ID:               "abc001",
				StartEpochSecond: 1468670400,
				DurationSecond:   6000,
				Title:            "AtCoder Beginner Contest 001",
				RateChange:       "-",
			},
			{
				ID:               "abc002",
				StartEpochSecond: 1468670400,
				DurationSecond:   6000,
				Title:            "AtCoder Beginner Contest 002",
				RateChange:       "-",
			},
		}

		count, err := SaveContests(ctx, pool, contests, now)
		if err != nil {
			t.Fatal(err)
		}
		if count != 2 {
			t.Errorf("count = %d, want 1", count)
		}
	})
}
