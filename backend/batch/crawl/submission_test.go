package crawl

import (
	"context"
	"fjnkt98/atcodersearch/internal/testutil"
	"fjnkt98/atcodersearch/pkg/atcoder"
	"fjnkt98/atcodersearch/pkg/ptr"
	"fjnkt98/atcodersearch/repository"
	"reflect"
	"testing"
	"time"
)

func TestFetchContestIDs(t *testing.T) {
	_, dsn, stop, err := testutil.CreateDBContainer()
	t.Cleanup(func() { stop() })

	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	pool, err := repository.NewPool(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}

	sql := `
INSERT INTO "contests" ("contest_id", "start_epoch_second", "duration_second", "title", "rate_change", "category") VALUES
('abc001', 0, 0, '', '-', 'ABC'),
('abc002', 0, 0, '', '-', 'ABC'),
('arc001', 0, 0, '', '-', 'ARC');`
	if _, err := pool.Exec(ctx, sql); err != nil {
		t.Fatal(err)
	}

	t.Run("all", func(t *testing.T) {
		result, err := FetchContestIDs(ctx, pool, []string{})
		if err != nil {
			t.Fatal(err)
		}

		want := []string{"abc001", "abc002", "arc001"}
		if !reflect.DeepEqual(result, want) {
			t.Errorf("result = %v, want %v", result, want)
		}
	})

	t.Run("ABC", func(t *testing.T) {
		result, err := FetchContestIDs(ctx, pool, []string{"ABC"})
		if err != nil {
			t.Fatal(err)
		}

		want := []string{"abc001", "abc002"}
		if !reflect.DeepEqual(result, want) {
			t.Errorf("result = %v, want %v", result, want)
		}
	})

	t.Run("ARC", func(t *testing.T) {
		result, err := FetchContestIDs(ctx, pool, []string{"ARC"})
		if err != nil {
			t.Fatal(err)
		}

		want := []string{"arc001"}
		if !reflect.DeepEqual(result, want) {
			t.Errorf("result = %v, want %v", result, want)
		}
	})
}

func TestSaveSubmissions(t *testing.T) {
	_, dsn, stop, err := testutil.CreateDBContainer()
	t.Cleanup(func() { stop() })

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
		submissions := make([]atcoder.Submission, 0)
		count, err := SaveSubmissions(ctx, pool, submissions, now)
		if err != nil {
			t.Fatal(err)
		}

		if count != 0 {
			t.Errorf("count = %d, want 0", count)
		}
	})

	t.Run("signle", func(t *testing.T) {
		submissions := []atcoder.Submission{
			{
				ID:            48852107,
				EpochSecond:   1703553569,
				ProblemID:     "abc300_a",
				UserID:        "Orkhon2010",
				ContestID:     "abc300",
				Language:      "C++ 20 (gcc 12.2)",
				Point:         100.0,
				Length:        259,
				Result:        "AC",
				ExecutionTime: ptr.To(int32(1)),
			},
		}
		count, err := SaveSubmissions(ctx, pool, submissions, now)
		if err != nil {
			t.Fatal(err)
		}

		if count != 1 {
			t.Errorf("count = %d, want 1", count)
		}
	})

	t.Run("multiple", func(t *testing.T) {
		submissions := []atcoder.Submission{
			{
				ID:            48852107,
				EpochSecond:   1703553569,
				ProblemID:     "abc300_a",
				UserID:        "Orkhon2010",
				ContestID:     "abc300",
				Language:      "C++ 20 (gcc 12.2)",
				Point:         100.0,
				Length:        259,
				Result:        "AC",
				ExecutionTime: ptr.To(int32(1)),
			},
			{
				ID:            48852073,
				EpochSecond:   1703553403,
				ProblemID:     "abc300_f",
				UserID:        "ecsmtlir",
				ContestID:     "abc300",
				Language:      "C++ 20 (gcc 12.2)",
				Point:         500.0,
				Length:        14721,
				Result:        "AC",
				ExecutionTime: ptr.To(int32(11)),
			},
		}
		count, err := SaveSubmissions(ctx, pool, submissions, now)
		if err != nil {
			t.Fatal(err)
		}

		if count != 2 {
			t.Errorf("count = %d, want 2", count)
		}
	})
}
