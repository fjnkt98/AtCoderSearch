package atcoder

import (
	"bytes"
	"os"
	"reflect"
	"testing"

	"k8s.io/utils/ptr"
)

func TestNewAtCoderClient(t *testing.T) {
	_, err := NewAtCoderClient()
	if err != nil {
		t.Errorf("failed to create client: %s", err.Error())
	}
}

func TestExtractCSRFToken(t *testing.T) {
	file, err := os.Open("./testdata/login.html")
	if err != nil {
		t.Fatalf("failed to open file `login.html`: %s", err.Error())
	}
	defer file.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(file); err != nil {
		t.Fatalf("failed to read file: %s", err.Error())
	}
	token, err := extractCSRFToken(buf.String())
	if err != nil {
		t.Fatalf("failed to extract CSRF token: %s", err.Error())
	}
	want := "KrVShPadRMxPBKM9LmjWJHaQvjC7ALXz6DXgHOCL1LQ="

	if token != want {
		t.Errorf("expected %s, but got %s", want, token)
	}
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

	want := []Submission{
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
	if !reflect.DeepEqual(result[:2], want) {
		t.Errorf("expected \n%+v\n , but got \n%+v\n", want, result[:2])
	}
}

func TestScrapeSubmissionsARC027(t *testing.T) {
	file, err := os.Open("./testdata/submissions.arc027.html")
	if err != nil {
		t.Fatalf("failed to open file `submissions.arc027.html`: %s", err.Error())
	}
	defer file.Close()

	result, err := scrapeSubmissions(file)
	if err != nil {
		t.Fatalf("failed to scrape submissions: %s", err.Error())
	}

	want := []Submission{
		{
			ID:            208118,
			EpochSecond:   1407106143,
			ProblemID:     "arc027_2",
			UserID:        "iab",
			ContestID:     "arc027",
			Language:      "OCaml (3.12.1)",
			Point:         100.0,
			Length:        1972,
			Result:        "AC",
			ExecutionTime: ptr.To(int32(36)),
		},
		{
			ID:            208117,
			EpochSecond:   1407102628,
			ProblemID:     "arc027_3",
			UserID:        "ne240214",
			ContestID:     "arc027",
			Language:      "Java (OpenJDK 1.7.0)",
			Point:         0.0,
			Length:        1628,
			Result:        "WA",
			ExecutionTime: ptr.To(int32(2106)),
		},
	}
	if !reflect.DeepEqual(result[:2], want) {
		t.Errorf("expected \n%+v\n , but got \n%+v\n", want, result[:2])
	}
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
		{UserID: "tourist", Rating: 3863, HighestRating: 4229, Affiliation: ptr.To("ITMO University"), BirthYear: ptr.To(int32(1994)), Country: ptr.To("BY"), Crown: ptr.To("crown_champion"), JoinCount: 59, Rank: 1, ActiveRank: ptr.To(int32(1)), Wins: 22},
		{UserID: "w4yneb0t", Rating: 3710, HighestRating: 3802, Affiliation: ptr.To("ETH Zurich"), BirthYear: nil, Country: ptr.To("CH"), Crown: nil, JoinCount: 21, Rank: 2, ActiveRank: nil, Wins: 2},
		{UserID: "ksun48", Rating: 3681, HighestRating: 3802, Affiliation: ptr.To("MIT"), BirthYear: ptr.To(int32(1998)), Country: ptr.To("CA"), Crown: ptr.To("crown_gold"), JoinCount: 58, Rank: 3, ActiveRank: ptr.To(int32(2)), Wins: 5},
		{UserID: "ecnerwala", Rating: 3663, HighestRating: 3814, Affiliation: ptr.To("MIT"), BirthYear: ptr.To(int32(1997)), Country: ptr.To("US"), Crown: ptr.To("crown_gold"), JoinCount: 36, Rank: 4, ActiveRank: ptr.To(int32(3)), Wins: 2},
		{UserID: "Benq", Rating: 3658, HighestRating: 3683, Affiliation: ptr.To("MIT"), BirthYear: ptr.To(int32(2001)), Country: ptr.To("US"), Crown: nil, JoinCount: 48, Rank: 5, ActiveRank: nil, Wins: 0},
		{UserID: "cospleermusora", Rating: 3606, HighestRating: 3783, Affiliation: nil, BirthYear: nil, Country: ptr.To("RU"), Crown: nil, JoinCount: 25, Rank: 5, ActiveRank: nil, Wins: 3},
		{UserID: "apiad", Rating: 3600, HighestRating: 3852, Affiliation: nil, BirthYear: ptr.To(int32(1997)), Country: ptr.To("CN"), Crown: ptr.To("crown_gold"), JoinCount: 51, Rank: 7, ActiveRank: ptr.To(int32(4)), Wins: 6},
		{UserID: "Um_nik", Rating: 3571, HighestRating: 3948, Affiliation: nil, BirthYear: ptr.To(int32(1996)), Country: ptr.To("UA"), Crown: ptr.To("crown_gold"), JoinCount: 60, Rank: 8, ActiveRank: ptr.To(int32(5)), Wins: 7},
		{UserID: "mnbvmar", Rating: 3555, HighestRating: 3736, Affiliation: ptr.To("University of Warsaw"), BirthYear: ptr.To(int32(1996)), Country: ptr.To("PL"), Crown: ptr.To("crown_gold"), JoinCount: 22, Rank: 9, ActiveRank: ptr.To(int32(6)), Wins: 1},
		{UserID: "Stonefeang", Rating: 3554, HighestRating: 3658, Affiliation: ptr.To("University of Warsaw"), BirthYear: ptr.To(int32(1997)), Country: ptr.To("PL"), Crown: ptr.To("crown_gold"), JoinCount: 37, Rank: 10, ActiveRank: ptr.To(int32(7)), Wins: 2},
	}

	if !reflect.DeepEqual(result, want) {
		t.Errorf("expected \n%+v\n , but got \n%+v\n", want, result)
	}
}
