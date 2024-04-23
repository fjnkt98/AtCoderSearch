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
	var a int32 = 26
	var b int32 = 1
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
	ID            int64
	EpochSecond   int64
	ProblemID     string
	UserID        string
	Language      string
	Point         float64
	Length        int32
	Result        string
	ExecutionTime int32
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

func ptr[T any](v T) *T {
	return &v
}

func TestScrapeUsers(t *testing.T) {
	file, err := os.Open("./testdata/users.html")
	if err != nil {
		t.Fatalf("failed to open file `users.html`: %s", err.Error())
	}
	defer file.Close()

	result, err := scrapeUsers(file)
	if err != nil {
		t.Fatalf("failed to scrape users: %s", err.Error())
	}

	want := []User{
		{UserName: "tourist", Rating: 3863, HighestRating: 4229, Affiliation: ptr("ITMO University"), BirthYear: ptr(int32(1994)), Country: ptr("BY"), Crown: ptr("crown_champion"), JoinCount: 59, Rank: 1, ActiveRank: ptr(int32(1)), Wins: 22},
		{UserName: "w4yneb0t", Rating: 3710, HighestRating: 3802, Affiliation: ptr("ETH Zurich"), BirthYear: nil, Country: ptr("CH"), Crown: nil, JoinCount: 21, Rank: 2, ActiveRank: nil, Wins: 2},
		{UserName: "ksun48", Rating: 3681, HighestRating: 3802, Affiliation: ptr("MIT"), BirthYear: ptr(int32(1998)), Country: ptr("CA"), Crown: ptr("crown_gold"), JoinCount: 58, Rank: 3, ActiveRank: ptr(int32(2)), Wins: 5},
		{UserName: "ecnerwala", Rating: 3663, HighestRating: 3814, Affiliation: ptr("MIT"), BirthYear: ptr(int32(1997)), Country: ptr("US"), Crown: ptr("crown_gold"), JoinCount: 36, Rank: 4, ActiveRank: ptr(int32(3)), Wins: 2},
		{UserName: "Benq", Rating: 3658, HighestRating: 3683, Affiliation: ptr("MIT"), BirthYear: ptr(int32(2001)), Country: ptr("US"), Crown: nil, JoinCount: 48, Rank: 5, ActiveRank: nil, Wins: 0},
		{UserName: "cospleermusora", Rating: 3606, HighestRating: 3783, Affiliation: nil, BirthYear: nil, Country: ptr("RU"), Crown: nil, JoinCount: 25, Rank: 5, ActiveRank: nil, Wins: 3},
		{UserName: "apiad", Rating: 3600, HighestRating: 3852, Affiliation: nil, BirthYear: ptr(int32(1997)), Country: ptr("CN"), Crown: ptr("crown_gold"), JoinCount: 51, Rank: 7, ActiveRank: ptr(int32(4)), Wins: 6},
		{UserName: "Um_nik", Rating: 3571, HighestRating: 3948, Affiliation: nil, BirthYear: ptr(int32(1996)), Country: ptr("UA"), Crown: ptr("crown_gold"), JoinCount: 60, Rank: 8, ActiveRank: ptr(int32(5)), Wins: 7},
		{UserName: "mnbvmar", Rating: 3555, HighestRating: 3736, Affiliation: ptr("University of Warsaw"), BirthYear: ptr(int32(1996)), Country: ptr("PL"), Crown: ptr("crown_gold"), JoinCount: 22, Rank: 9, ActiveRank: ptr(int32(6)), Wins: 1},
		{UserName: "Stonefeang", Rating: 3554, HighestRating: 3658, Affiliation: ptr("University of Warsaw"), BirthYear: ptr(int32(1997)), Country: ptr("PL"), Crown: ptr("crown_gold"), JoinCount: 37, Rank: 10, ActiveRank: ptr(int32(7)), Wins: 2},
	}

	if !reflect.DeepEqual(result, want) {
		t.Errorf("expected %+v , but got %+v", want, result)
	}
}
