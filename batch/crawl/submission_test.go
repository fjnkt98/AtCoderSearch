package crawl

import (
	"context"
	"errors"
	"fjnkt98/atcodersearch/internal/testutil"
	"fjnkt98/atcodersearch/pkg/atcoder"
	"fjnkt98/atcodersearch/pkg/ptr"
	"fjnkt98/atcodersearch/repository"
	"reflect"
	"testing"
	"time"
)

func TestSubmission(t *testing.T) {
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

	t.Run("fetch all contest id", func(t *testing.T) {
		t.Parallel()

		result, err := FetchContestIDs(ctx, pool, []string{})
		if err != nil {
			t.Fatal(err)
		}

		want := []string{"abc001", "abc002", "arc001"}
		if !reflect.DeepEqual(result, want) {
			t.Errorf("result = %v, want %v", result, want)
		}
	})

	t.Run("fetch ABC contest id", func(t *testing.T) {
		t.Parallel()

		result, err := FetchContestIDs(ctx, pool, []string{"ABC"})
		if err != nil {
			t.Fatal(err)
		}

		want := []string{"abc001", "abc002"}
		if !reflect.DeepEqual(result, want) {
			t.Errorf("result = %v, want %v", result, want)
		}
	})

	t.Run("fetch ARC contest id", func(t *testing.T) {
		t.Parallel()

		result, err := FetchContestIDs(ctx, pool, []string{"ARC"})
		if err != nil {
			t.Fatal(err)
		}

		want := []string{"arc001"}
		if !reflect.DeepEqual(result, want) {
			t.Errorf("result = %v, want %v", result, want)
		}
	})

	t.Run("save empty submissions", func(t *testing.T) {
		t.Parallel()

		submissions := make([]atcoder.Submission, 0)
		count, err := SaveSubmissions(ctx, pool, submissions)
		if err != nil {
			t.Fatal(err)
		}

		if count != 0 {
			t.Errorf("count = %d, want 0", count)
		}
	})

	t.Run("save single submission", func(t *testing.T) {
		t.Parallel()

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
		count, err := SaveSubmissions(ctx, pool, submissions)
		if err != nil {
			t.Fatal(err)
		}

		if count != 1 {
			t.Errorf("count = %d, want 1", count)
		}
	})

	t.Run("save multiple submissions", func(t *testing.T) {
		t.Parallel()

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
		count, err := SaveSubmissions(ctx, pool, submissions)
		if err != nil {
			t.Fatal(err)
		}

		if count != 2 {
			t.Errorf("count = %d, want 2", count)
		}
	})

	t.Run("crawl submissions(success)", func(t *testing.T) {
		t.Parallel()

		crawler := NewSubmissionCrawler(
			&DummyAtCoderClientS{},
			pool,
			time.Second,
			0,
			0,
			nil,
		)
		if err := crawler.Crawl(ctx); err != nil {
			t.Errorf("expected no error, but got %v", err)
		}

		crawler = NewSubmissionCrawler(
			&DummyAtCoderClientS{},
			pool,
			time.Second,
			0,
			0,
			[]string{"ABC", "ARC"},
		)
		if err := crawler.Crawl(ctx); err != nil {
			t.Errorf("expected no error, but got %v", err)
		}
	})

	t.Run("crawl submissions(fail)", func(t *testing.T) {
		t.Parallel()

		crawler := NewSubmissionCrawler(
			&DummyAtCoderClientF{},
			pool,
			time.Second,
			0,
			0,
			nil,
		)
		if err := crawler.Crawl(ctx); !errors.Is(err, ErrDummy) {
			t.Errorf("expected ErrDummy, but got %#v", err)
		}
	})

	t.Run("crawl submissions(retry success)", func(t *testing.T) {
		t.Parallel()

		crawler := NewSubmissionCrawler(
			&RecoverableDummyAtCoderClient{
				errCount: 1,
			},
			pool,
			time.Second,
			1,
			100*time.Millisecond,
			nil,
		)

		if err := crawler.Crawl(ctx); err != nil {
			t.Errorf("expected no error, but got %v", err)
		}
	})

	t.Run("crawl submissions(retry fail)", func(t *testing.T) {
		t.Parallel()

		crawler := NewSubmissionCrawler(
			&RecoverableDummyAtCoderClient{
				errCount: 2,
			},
			pool,
			time.Second,
			1,
			100*time.Millisecond,
			nil,
		)

		if err := crawler.Crawl(ctx); !errors.Is(err, ErrDummy) {
			t.Errorf("expected ErrDummy, but got %#v", err)
		}
	})
}

type RecoverableDummyAtCoderClient struct {
	count    int
	errCount int
}

var _ atcoder.AtCoderClient = new(RecoverableDummyAtCoderClient)

func (c *RecoverableDummyAtCoderClient) Login(ctx context.Context, username, password string) error {
	return nil
}
func (c *RecoverableDummyAtCoderClient) FetchProblemHTML(ctx context.Context, contestID, problemID string) (string, error) {
	return "", nil
}
func (c *RecoverableDummyAtCoderClient) FetchSubmissions(ctx context.Context, contestID string, page int) ([]atcoder.Submission, error) {
	if page >= 2 {
		return []atcoder.Submission{}, nil
	}

	if c.count < c.errCount {
		c.count += 1
		return nil, ErrDummy
	} else {
		return []atcoder.Submission{
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
		}, nil
	}
}
func (c *RecoverableDummyAtCoderClient) FetchUsers(ctx context.Context, page int) ([]atcoder.User, error) {
	return nil, nil
}
