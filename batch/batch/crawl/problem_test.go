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

func TestProblem(t *testing.T) {
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

	t.Run("save empty problems", func(t *testing.T) {
		t.Parallel()
		contests := make([]atcoder.Contest, 0)
		count, err := SaveContests(ctx, pool, contests)
		if err != nil {
			t.Fatal(err)
		}

		if count != 0 {
			t.Errorf("count = %d, want 0", count)
		}
	})

	t.Run("save single problem", func(t *testing.T) {
		t.Parallel()
		contests := []atcoder.Contest{
			{
				ID:               "abc001",
				StartEpochSecond: 1468670400,
				DurationSecond:   6000,
				Title:            "AtCoder Beginner Contest 001",
				RateChange:       "-",
			},
		}

		count, err := SaveContests(ctx, pool, contests)
		if err != nil {
			t.Fatal(err)
		}
		if count != 1 {
			t.Errorf("count = %d, want 1", count)
		}
	})

	t.Run("save multiple problems", func(t *testing.T) {
		t.Parallel()
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

		count, err := SaveContests(ctx, pool, contests)
		if err != nil {
			t.Fatal(err)
		}
		if count != 2 {
			t.Errorf("count = %d, want 2", count)
		}
	})

	t.Run("save empty difficulties", func(t *testing.T) {
		t.Parallel()
		difficulties := make(map[string]atcoder.Difficulty)
		count, err := SaveDifficulties(ctx, pool, difficulties)
		if err != nil {
			t.Fatal(err)
		}

		if count != 0 {
			t.Errorf("count = %d, want 0", count)
		}
	})

	t.Run("save single difficulty", func(t *testing.T) {
		t.Parallel()
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

		count, err := SaveDifficulties(ctx, pool, difficulties)
		if err != nil {
			t.Fatal(err)
		}
		if count != 1 {
			t.Errorf("count = %d, want 1", count)
		}
	})

	t.Run("save multiple difficulties", func(t *testing.T) {
		t.Parallel()
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

		count, err := SaveDifficulties(ctx, pool, difficulties)
		if err != nil {
			t.Fatal(err)
		}
		if count != 2 {
			t.Errorf("count = %d, want 2", count)
		}
	})

	t.Run("detect diff", func(t *testing.T) {
		t.Parallel()
		sql := `
INSERT INTO "problems" ("problem_id", "contest_id", "problem_index", "name", "title", "url", "html")
VALUES
    ('diff_1', 'diff', 'A', 'detect diff 1', 'A. detect diff 1', 'url', 'html');`

		if _, err := pool.Exec(ctx, sql); err != nil {
			t.Fatal(err)
		}

		problems := []atcoder.Problem{
			{
				ID:           "diff_1",
				ContestID:    "diff",
				ProblemIndex: "A",
				Name:         "detect diff 1",
				Title:        "A. detect diff 1",
			},
			{
				ID:           "diff_2",
				ContestID:    "diff",
				ProblemIndex: "A",
				Name:         "detect diff 2",
				Title:        "A. detect diff 2",
			},
		}

		diff, err := DetectDiff(ctx, pool, problems)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(problems[1:], diff) {
			t.Errorf("diff = %v, want %v", diff, problems[1:])
		}
	})

	t.Run("crawl problem(success, diff)", func(t *testing.T) {
		t.Parallel()
		crawler := NewProblemCrawler(
			&DummyAtCoderClientS{},
			&DummyAtCoderProblemsClientS{},
			pool,
			time.Second,
			false,
		)

		if err := crawler.CrawlContests(ctx); err != nil {
			t.Errorf("expected no error, but got %v", err)
		}

		if err := crawler.CrawlDifficulties(ctx); err != nil {
			t.Errorf("expected no error, but got %v", err)
		}

		if err := crawler.CrawlProblems(ctx); err != nil {
			t.Errorf("expected no error, but got %v", err)
		}
	})

	t.Run("crawl problem(success, all)", func(t *testing.T) {
		t.Parallel()
		crawler := NewProblemCrawler(
			&DummyAtCoderClientS{},
			&DummyAtCoderProblemsClientS{},
			pool,
			time.Second,
			true,
		)
		if err := crawler.CrawlProblems(ctx); err != nil {
			t.Errorf("expected no error, but got %v", err)
		}
	})

	t.Run("crawl problem(fail atcoder)", func(t *testing.T) {
		t.Parallel()
		crawler := NewProblemCrawler(
			&DummyAtCoderClientF{},
			&DummyAtCoderProblemsClientS{},
			pool,
			time.Second,
			true,
		)

		if err := crawler.CrawlProblems(ctx); !errors.Is(err, ErrDummy) {
			t.Errorf("expected ErrDummy, but got %#v", err)
		}
	})

	t.Run("crawl problem(fail atcoder problems)", func(t *testing.T) {
		t.Parallel()
		crawler := NewProblemCrawler(
			&DummyAtCoderClientS{},
			&DummyAtCoderProblemsClientF{},
			pool,
			time.Second,
			false,
		)

		if err := crawler.CrawlContests(ctx); !errors.Is(err, ErrDummy) {
			t.Errorf("expected ErrDummy, but got %#v", err)
		}

		if err := crawler.CrawlDifficulties(ctx); !errors.Is(err, ErrDummy) {
			t.Errorf("expected ErrDummy, but got %#v", err)
		}

		if err := crawler.CrawlProblems(ctx); !errors.Is(err, ErrDummy) {
			t.Errorf("expected ErrDummy, but got %#v", err)
		}
	})

	t.Run("crawl problem(fail)", func(t *testing.T) {
		t.Parallel()
		crawler := NewProblemCrawler(
			&DummyAtCoderClientF{},
			&DummyAtCoderProblemsClientF{},
			pool,
			time.Second,
			false,
		)

		if err := crawler.CrawlContests(ctx); !errors.Is(err, ErrDummy) {
			t.Errorf("expected ErrDummy, but got %#v", err)
		}

		if err := crawler.CrawlDifficulties(ctx); !errors.Is(err, ErrDummy) {
			t.Errorf("expected ErrDummy, but got %#v", err)
		}

		if err := crawler.CrawlProblems(ctx); !errors.Is(err, ErrDummy) {
			t.Errorf("expected ErrDummy, but got %#v", err)
		}
	})
}

type DummyAtCoderClientS struct{}

var _ atcoder.AtCoderClient = &DummyAtCoderClientS{}

func (c *DummyAtCoderClientS) Login(ctx context.Context, username, password string) error {
	return nil
}
func (c *DummyAtCoderClientS) FetchProblemHTML(ctx context.Context, contestID, problemID string) (string, error) {
	return "", nil
}
func (c *DummyAtCoderClientS) FetchSubmissions(ctx context.Context, contestID string, page int) ([]atcoder.Submission, error) {
	if page >= 2 {
		return []atcoder.Submission{}, nil
	}
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
func (c *DummyAtCoderClientS) FetchUsers(ctx context.Context, page int) ([]atcoder.User, error) {
	if page >= 2 {
		return []atcoder.User{}, nil
	}
	return []atcoder.User{
		{
			UserID:        "tourist",
			Rating:        3863,
			HighestRating: 4229,
			Affiliation:   ptr.To("ITMO University"),
			BirthYear:     ptr.To(int32(1994)),
			Country:       ptr.To("BY"),
			Crown:         ptr.To("crown_champion"),
			JoinCount:     59,
			Rank:          1,
			ActiveRank:    ptr.To(int32(1)),
			Wins:          22,
		},
		{
			UserID:        "w4yneb0t",
			Rating:        3710,
			HighestRating: 3802,
			Affiliation:   ptr.To("ETH Zurich"),
			BirthYear:     nil,
			Country:       ptr.To("CH"),
			Crown:         nil,
			JoinCount:     21,
			Rank:          2,
			ActiveRank:    nil,
			Wins:          2,
		},
	}, nil
}

type DummyAtCoderClientF struct{}

var ErrDummy = errors.New("dummy error")

func (c *DummyAtCoderClientF) Login(ctx context.Context, username, password string) error {
	return ErrDummy
}
func (c *DummyAtCoderClientF) FetchProblemHTML(ctx context.Context, contestID, problemID string) (string, error) {
	return "", ErrDummy
}
func (c *DummyAtCoderClientF) FetchSubmissions(ctx context.Context, contestID string, page int) ([]atcoder.Submission, error) {
	return nil, ErrDummy
}
func (c *DummyAtCoderClientF) FetchUsers(ctx context.Context, page int) ([]atcoder.User, error) {
	return nil, ErrDummy
}

type DummyAtCoderProblemsClientS struct{}

var _ atcoder.AtCoderProblemsClient = &DummyAtCoderProblemsClientS{}

func (c *DummyAtCoderProblemsClientS) FetchProblems(ctx context.Context) ([]atcoder.Problem, error) {
	return []atcoder.Problem{
		{
			ID:           "test001_a",
			ContestID:    "test001",
			ProblemIndex: "A",
			Name:         "test",
			Title:        "A. test",
		},
	}, nil
}
func (c *DummyAtCoderProblemsClientS) FetchDifficulties(ctx context.Context) (map[string]atcoder.Difficulty, error) {
	return map[string]atcoder.Difficulty{
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
	}, nil
}
func (c *DummyAtCoderProblemsClientS) FetchContests(ctx context.Context) ([]atcoder.Contest, error) {
	return []atcoder.Contest{
		{
			ID:               "abc001",
			StartEpochSecond: 1468670400,
			DurationSecond:   6000,
			Title:            "AtCoder Beginner Contest 001",
			RateChange:       "-",
		},
	}, nil
}

type DummyAtCoderProblemsClientF struct{}

var _ atcoder.AtCoderProblemsClient = &DummyAtCoderProblemsClientF{}

func (c *DummyAtCoderProblemsClientF) FetchProblems(ctx context.Context) ([]atcoder.Problem, error) {
	return nil, ErrDummy
}
func (c *DummyAtCoderProblemsClientF) FetchDifficulties(ctx context.Context) (map[string]atcoder.Difficulty, error) {
	return nil, ErrDummy
}
func (c *DummyAtCoderProblemsClientF) FetchContests(ctx context.Context) ([]atcoder.Contest, error) {
	return nil, ErrDummy
}
