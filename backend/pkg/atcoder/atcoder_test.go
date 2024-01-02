package atcoder

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/goark/errs"
)

type DummySuccessAtCoderClient struct{}

func NewDummySuccessAtCoderClient() AtCoderClient {
	return &DummySuccessAtCoderClient{}
}

func (c *DummySuccessAtCoderClient) FetchSubmissions(ctx context.Context, contestID string, page int) ([]Submission, error) {
	a := 26
	b := 1
	return []Submission{
		{
			ID:            2054945,
			EpochSecond:   1517670136,
			ProblemID:     "apc001_b",
			ContestID:     "apc001",
			UserID:        "dko_n",
			Language:      "Python3 (3.4.3)",
			Point:         0,
			Length:        326,
			Result:        "WA",
			ExecutionTime: &a,
		},
		{
			ID:            3600716,
			EpochSecond:   1542210655,
			ProblemID:     "abc101_a",
			ContestID:     "abc101",
			UserID:        "kichi2004",
			Language:      "C++14 (GCC 5.4.1)",
			Point:         100,
			Length:        1458,
			Result:        "AC",
			ExecutionTime: &b,
		},
	}, nil
}

func (c *DummySuccessAtCoderClient) FetchSubmissionResult(ctx context.Context, contestID string, submissionID int64) (string, error) {
	return "AC", nil
}

func (c *DummySuccessAtCoderClient) FetchProblem(ctx context.Context, contestID string, problemID string) (string, error) {
	return "", nil
}

func (c *DummySuccessAtCoderClient) FetchUsers(ctx context.Context, page int) ([]User, error) {
	return nil, nil
}

func (c *DummySuccessAtCoderClient) Login(ctx context.Context, username, password string) error {
	return nil
}

type DummyFailAtCoderClient struct{}

func NewDummyFailAtCoderClient() AtCoderClient {
	return &DummyFailAtCoderClient{}
}

func (c *DummyFailAtCoderClient) FetchSubmissions(ctx context.Context, contestID string, page int) ([]Submission, error) {
	return nil, errs.New("request failed")
}

func (c *DummyFailAtCoderClient) FetchSubmissionResult(ctx context.Context, contestID string, submissionID int64) (string, error) {
	return "", errs.New("request failed")
}

func (c *DummyFailAtCoderClient) FetchProblem(ctx context.Context, contestID string, problemID string) (string, error) {
	return "", errs.New("request failed")
}

func (c *DummyFailAtCoderClient) FetchUsers(ctx context.Context, page int) ([]User, error) {
	return nil, errs.New("request failed")
}

func (c *DummyFailAtCoderClient) Login(ctx context.Context, username, password string) error {
	return errs.New("login failed")
}

type SubmissionPiece struct {
	ID            int
	EpochSecond   int64
	ProblemID     string
	UserID        string
	Language      string
	Point         float64
	Length        int
	Result        string
	ExecutionTime int
}

func TestScrapeSubmissions(t *testing.T) {
	file, err := os.Open("./testdata/submissions.html")
	if err != nil {
		t.Fatalf("failed to open file `submissions.html`: %s", err.Error())
	}
	defer file.Close()

	result, err := scrapeSubmissions(file)
	if err != nil {
		t.Fatalf("failed to scrape submissions: %s", err.Error())
	}

	want := []SubmissionPiece{
		{ID: 48852107, EpochSecond: 1703553569, ProblemID: "abc300_a", UserID: "Orkhon2010", Language: "C++ 20 (gcc 12.2)", Point: 100, Length: 259, Result: "AC", ExecutionTime: 1},
		{ID: 48852073, EpochSecond: 1703553403, ProblemID: "abc300_f", UserID: "ecsmtlir", Language: "C++ 20 (gcc 12.2)", Point: 500, Length: 14721, Result: "AC", ExecutionTime: 11},
	}
	for i := 0; i < 2; i++ {
		res := SubmissionPiece{
			ID:            result[i].ID,
			EpochSecond:   result[i].EpochSecond,
			ProblemID:     result[i].ProblemID,
			UserID:        result[i].UserID,
			Language:      result[i].Language,
			Point:         result[i].Point,
			Length:        result[i].Length,
			Result:        result[i].Result,
			ExecutionTime: *result[i].ExecutionTime,
		}
		if !reflect.DeepEqual(res, want[i]) {
			t.Errorf("scrape result %d is different from expected result, result: %+v, expected: %+v", i, res, want[i])
		}
	}
}
