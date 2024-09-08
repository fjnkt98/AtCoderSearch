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

	t.Run("single", func(t *testing.T) {
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
			t.Errorf("count = %d, want 2", count)
		}
	})
}

func TestSaveDifficulties(t *testing.T) {
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
		difficulties := make(map[string]atcoder.Difficulty)
		count, err := SaveDifficulties(ctx, pool, difficulties, now)
		if err != nil {
			t.Fatal(err)
		}

		if count != 0 {
			t.Errorf("count = %d, want 0", count)
		}
	})

	t.Run("single", func(t *testing.T) {
		difficulties := map[string]atcoder.Difficulty{
			"abc118_d": {
				Slope:            ptr.To(-0.0006619775680720775),
				Intercept:        ptr.To(8.881759153638702),
				Variance:         ptr.To(0.30752713797776526),
				Difficulty:       ptr.To(int64(1657)),
				Discrimination:   ptr.To(0.004479398673070138),
				IrtLoglikelihood: ptr.To(-491.8630322466751),
				IrtUsers:         ptr.To(2442.0),
				IsExperimental:   ptr.To(false),
			},
		}

		count, err := SaveDifficulties(ctx, pool, difficulties, now)
		if err != nil {
			t.Fatal(err)
		}
		if count != 1 {
			t.Errorf("count = %d, want 1", count)
		}
	})

	t.Run("multiple", func(t *testing.T) {
		difficulties := map[string]atcoder.Difficulty{
			"abc118_d": {
				Slope:            ptr.To(-0.0006619775680720775),
				Intercept:        ptr.To(8.881759153638702),
				Variance:         ptr.To(0.30752713797776526),
				Difficulty:       ptr.To(int64(1657)),
				Discrimination:   ptr.To(0.004479398673070138),
				IrtLoglikelihood: ptr.To(-491.8630322466751),
				IrtUsers:         ptr.To(2442.0),
				IsExperimental:   ptr.To(false),
			},
			"agc026_d": {
				Slope:            ptr.To(-0.0004027506918277324),
				Intercept:        ptr.To(9.274529080920633),
				Variance:         ptr.To(0.12135365008788429),
				Difficulty:       ptr.To(int64(2746)),
				Discrimination:   ptr.To(0.004479398673070138),
				IrtLoglikelihood: ptr.To(-145.66848869773756),
				IrtUsers:         ptr.To(1799.0),
				IsExperimental:   ptr.To(false),
			},
		}

		count, err := SaveDifficulties(ctx, pool, difficulties, now)
		if err != nil {
			t.Fatal(err)
		}
		if count != 2 {
			t.Errorf("count = %d, want 2", count)
		}
	})
}

func TestDetectDiff(t *testing.T) {
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

	sql := `
INSERT INTO "problems" ("problem_id", "contest_id", "problem_index", "name", "title", "url", "html")
VALUES
    ('abc001_a', 'abc001', 'A', 'test problem 1', 'A. test problem 1', 'url', 'html');`

	if _, err := pool.Exec(ctx, sql); err != nil {
		t.Fatal(err)
	}

	problems := []atcoder.Problem{
		{
			ID:           "abc001_a",
			ContestID:    "abc001",
			ProblemIndex: "A",
			Name:         "test problem 1",
			Title:        "A. test problem 1",
		},
		{
			ID:           "abc001_b",
			ContestID:    "abc001",
			ProblemIndex: "B",
			Name:         "test problem 2",
			Title:        "B. test problem 2",
		},
	}

	diff, err := DetectDiff(ctx, pool, problems)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(problems[1:], diff) {
		t.Errorf("diff = %v, want %v", diff, problems[1:])
	}
}
