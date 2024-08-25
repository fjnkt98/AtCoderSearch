package repository

import (
	"reflect"
	"testing"

	"github.com/fjnkt98/atcodersearch-batch/pkg/atcoder"
)

func ptr[T any](v T) *T {
	return &v
}

func TestNewContest(t *testing.T) {
	c := atcoder.Contest{
		ID:               "abc300",
		StartEpochSecond: 1682769600,
		DurationSecond:   6000,
		Title:            "ユニークビジョンプログラミングコンテスト2023 春 (AtCoder Beginner Contest 300)",
		RateChange:       " ~ 1999",
	}

	expected := Contest{
		ContestID:        "abc300",
		StartEpochSecond: 1682769600,
		DurationSecond:   6000,
		Title:            "ユニークビジョンプログラミングコンテスト2023 春 (AtCoder Beginner Contest 300)",
		RateChange:       " ~ 1999",
		Category:         "ABC",
	}
	actual := NewContest(c)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected \n%+v\n, but got \n%+v\n", expected, actual)
	}
}

func TestNewContests(t *testing.T) {
	contests := []atcoder.Contest{
		{
			ID:               "abc300",
			StartEpochSecond: 1682769600,
			DurationSecond:   6000,
			Title:            "ユニークビジョンプログラミングコンテスト2023 春 (AtCoder Beginner Contest 300)",
			RateChange:       " ~ 1999",
		},
		{
			ID:               "abc301",
			StartEpochSecond: 1683979200,
			DurationSecond:   6300,
			Title:            "パナソニックグループプログラミングコンテスト2023（AtCoder Beginner Contest 301）",
			RateChange:       " ~ 1999",
		},
	}

	expected := []Contest{
		{
			ContestID:        "abc300",
			StartEpochSecond: 1682769600,
			DurationSecond:   6000,
			Title:            "ユニークビジョンプログラミングコンテスト2023 春 (AtCoder Beginner Contest 300)",
			RateChange:       " ~ 1999",
			Category:         "ABC",
		},
		{
			ContestID:        "abc301",
			StartEpochSecond: 1683979200,
			DurationSecond:   6300,
			Title:            "パナソニックグループプログラミングコンテスト2023（AtCoder Beginner Contest 301）",
			RateChange:       " ~ 1999",
			Category:         "ABC",
		},
	}
	actual := NewContests(contests)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected \n%+v\n, but got \n%+v\n", expected, actual)
	}
}

func TestNewDifficulty(t *testing.T) {
	d := atcoder.Difficulty{
		Slope:            ptr(-0.0011830536916685426),
		Intercept:        ptr(7.8282534949123725),
		Variance:         ptr(0.39991055818601307),
		Difficulty:       ptr[int64](534),
		Discrimination:   ptr(0.004479398673070138),
		IrtLogLikelihood: ptr(-91.12539763628888),
		IrtUsers:         ptr(281.0),
		IsExperimental:   ptr(true),
	}

	expected := Difficulty{
		ProblemID:        "abc019_2",
		Slope:            ptr(-0.0011830536916685426),
		Intercept:        ptr(7.8282534949123725),
		Variance:         ptr(0.39991055818601307),
		Difficulty:       ptr[int64](534),
		Discrimination:   ptr(0.004479398673070138),
		IrtLoglikelihood: ptr(-91.12539763628888),
		IrtUsers:         ptr(281.0),
		IsExperimental:   ptr(true),
	}
	actual := NewDifficulty("abc019_2", d)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected \n%+v\n, but got \n%+v\n", expected, actual)
	}
}

func TestNewDifficulties(t *testing.T) {
	difficulties := map[string]atcoder.Difficulty{
		"abc019_2": {
			Slope:            ptr(-0.0011830536916685426),
			Intercept:        ptr(7.8282534949123725),
			Variance:         ptr(0.39991055818601307),
			Difficulty:       ptr[int64](534),
			Discrimination:   ptr(0.004479398673070138),
			IrtLogLikelihood: ptr(-91.12539763628888),
			IrtUsers:         ptr(281.0),
			IsExperimental:   ptr(true),
		},
		"abc182_e": {
			Slope:            ptr(-0.0008170719400144961),
			Intercept:        ptr(8.456092242494435),
			Variance:         ptr(0.27113303135125566),
			Difficulty:       ptr[int64](1098),
			Discrimination:   ptr(0.004479398673070138),
			IrtLogLikelihood: ptr(-1838.1064873009134),
			IrtUsers:         ptr(6840.0),
			IsExperimental:   ptr(false),
		},
	}

	expected := []Difficulty{
		{
			ProblemID:        "abc019_2",
			Slope:            ptr(-0.0011830536916685426),
			Intercept:        ptr(7.8282534949123725),
			Variance:         ptr(0.39991055818601307),
			Difficulty:       ptr[int64](534),
			Discrimination:   ptr(0.004479398673070138),
			IrtLoglikelihood: ptr(-91.12539763628888),
			IrtUsers:         ptr(281.0),
			IsExperimental:   ptr(true),
		},
		{
			ProblemID:        "abc182_e",
			Slope:            ptr(-0.0008170719400144961),
			Intercept:        ptr(8.456092242494435),
			Variance:         ptr(0.27113303135125566),
			Difficulty:       ptr[int64](1098),
			Discrimination:   ptr(0.004479398673070138),
			IrtLoglikelihood: ptr(-1838.1064873009134),
			IrtUsers:         ptr(6840.0),
			IsExperimental:   ptr(false),
		},
	}
	actual := NewDifficulties(difficulties)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected \n%+v\n, but got \n%+v\n", expected, actual)
	}
}

func TestNewSubmission(t *testing.T) {
	s := atcoder.Submission{
		ID:          4031,
		EpochSecond: 1305359825,
		ProblemID:   "utpc2011_9",
		ContestID:   "utpc2011",
		UserID:      "old_170",
		Language:    "C++ (GCC 4.4.7)",
		Point:       0,
		Length:      2566,
		Result:      "WR",
	}

	expected := Submission{
		ID:          4031,
		EpochSecond: 1305359825,
		ProblemID:   "utpc2011_9",
		ContestID:   ptr("utpc2011"),
		UserID:      ptr("old_170"),
		Language:    ptr("C++ (GCC 4.4.7)"),
		Point:       ptr(0.0),
		Length:      ptr[int32](2566),
		Result:      ptr("WR"),
	}
	actual := NewSubmission(s)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected \n%+v\n, but got \n%+v\n", expected, actual)
	}
}

func TestNewSubmissions(t *testing.T) {

	submissions := []atcoder.Submission{
		{
			ID:          4031,
			EpochSecond: 1305359825,
			ProblemID:   "utpc2011_9",
			ContestID:   "utpc2011",
			UserID:      "old_170",
			Language:    "C++ (GCC 4.4.7)",
			Point:       0,
			Length:      0,
			Result:      "WJ",
		},
		{
			ID:            2119,
			EpochSecond:   1305342271,
			ProblemID:     "utpc2011_1",
			ContestID:     "utpc2011",
			UserID:        "old_160",
			Language:      "C++ (GCC 4.4.7)",
			Point:         100,
			Length:        259,
			Result:        "AC",
			ExecutionTime: ptr[int32](20),
		},
		{
			ID:          4031,
			EpochSecond: 1305359825,
			ProblemID:   "utpc2011_9",
			ContestID:   "utpc2011",
			UserID:      "old_170",
			Language:    "C++ (GCC 4.4.7)",
			Point:       0,
			Length:      2566,
			Result:      "WR",
		},
	}

	expected := []Submission{
		{
			ID:            2119,
			EpochSecond:   1305342271,
			ProblemID:     "utpc2011_1",
			ContestID:     ptr("utpc2011"),
			UserID:        ptr("old_160"),
			Language:      ptr("C++ (GCC 4.4.7)"),
			Point:         ptr(100.0),
			Length:        ptr[int32](259),
			Result:        ptr("AC"),
			ExecutionTime: ptr[int32](20),
		},
		{
			ID:          4031,
			EpochSecond: 1305359825,
			ProblemID:   "utpc2011_9",
			ContestID:   ptr("utpc2011"),
			UserID:      ptr("old_170"),
			Language:    ptr("C++ (GCC 4.4.7)"),
			Point:       ptr(0.0),
			Length:      ptr[int32](2566),
			Result:      ptr("WR"),
		},
	}
	actual := NewSubmissions(submissions)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected \n%+v\n, but got \n%+v\n", expected, actual)
	}
}

func TestNewUser(t *testing.T) {
	u := atcoder.User{
		UserID:        "tourist",
		Rating:        3798,
		HighestRating: 4229,
		Affiliation:   ptr("ITMO University"),
		BirthYear:     ptr[int32](1994),
		Country:       ptr("BY"),
		Crown:         ptr("crown_champion"),
		JoinCount:     60,
		Rank:          1,
		ActiveRank:    ptr[int32](1),
		Wins:          22,
	}

	expected := User{
		UserID:        "tourist",
		Rating:        3798,
		HighestRating: 4229,
		Affiliation:   ptr("ITMO University"),
		BirthYear:     ptr[int32](1994),
		Country:       ptr("BY"),
		Crown:         ptr("crown_champion"),
		JoinCount:     60,
		Rank:          1,
		ActiveRank:    ptr[int32](1),
		Wins:          22,
	}
	actual := NewUser(u)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected \n%+v\n, but got \n%+v\n", expected, actual)
	}
}

func TestNewUsers(t *testing.T) {
	users := []atcoder.User{
		{
			UserID:        "tourist",
			Rating:        3798,
			HighestRating: 4229,
			Affiliation:   ptr("ITMO University"),
			BirthYear:     ptr[int32](1994),
			Country:       ptr("BY"),
			Crown:         ptr("crown_champion"),
			JoinCount:     60,
			Rank:          1,
			ActiveRank:    ptr[int32](1),
			Wins:          22,
		},
		{
			UserID:        "ksun48",
			Rating:        3724,
			HighestRating: 3802,
			Affiliation:   ptr("MIT"),
			BirthYear:     ptr[int32](1998),
			Country:       ptr("CA"),
			Crown:         ptr("crown_gold"),
			JoinCount:     59,
			Rank:          2,
			ActiveRank:    ptr[int32](2),
			Wins:          5,
		},
	}
	expected := []User{
		{
			UserID:        "tourist",
			Rating:        3798,
			HighestRating: 4229,
			Affiliation:   ptr("ITMO University"),
			BirthYear:     ptr[int32](1994),
			Country:       ptr("BY"),
			Crown:         ptr("crown_champion"),
			JoinCount:     60,
			Rank:          1,
			ActiveRank:    ptr[int32](1),
			Wins:          22,
		},
		{
			UserID:        "ksun48",
			Rating:        3724,
			HighestRating: 3802,
			Affiliation:   ptr("MIT"),
			BirthYear:     ptr[int32](1998),
			Country:       ptr("CA"),
			Crown:         ptr("crown_gold"),
			JoinCount:     59,
			Rank:          2,
			ActiveRank:    ptr[int32](2),
			Wins:          5,
		},
	}
	actual := NewUsers(users)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected \n%+v\n, but got \n%+v\n", expected, actual)
	}
}
