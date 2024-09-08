package repository

import (
	"fjnkt98/atcodersearch/pkg/atcoder"
	"fjnkt98/atcodersearch/pkg/ptr"
	"reflect"
	"slices"
	"testing"
	"time"
)

func TestNewContest(t *testing.T) {
	c := atcoder.Contest{
		ID:               "abc300",
		StartEpochSecond: 1682769600,
		DurationSecond:   6000,
		Title:            "ユニークビジョンプログラミングコンテスト2023 春 (AtCoder Beginner Contest 300)",
		RateChange:       " ~ 1999",
	}

	updatedAt := time.Now()
	want := Contest{
		ContestID:        "abc300",
		StartEpochSecond: 1682769600,
		DurationSecond:   6000,
		Title:            "ユニークビジョンプログラミングコンテスト2023 春 (AtCoder Beginner Contest 300)",
		RateChange:       " ~ 1999",
		Category:         "ABC",
		UpdatedAt:        updatedAt,
	}
	actual := NewContest(c, updatedAt)

	if !reflect.DeepEqual(want, actual) {
		t.Errorf("expected \n%+v\n, but got \n%+v\n", want, actual)
	}
}

func TestNewDifficulties(t *testing.T) {
	difficulties := map[string]atcoder.Difficulty{
		"abc019_2": {
			Slope:            ptr.To(-0.0011830536916685426),
			Intercept:        ptr.To(7.8282534949123725),
			Variance:         ptr.To(0.39991055818601307),
			Difficulty:       ptr.To[int64](534),
			Discrimination:   ptr.To(0.004479398673070138),
			IrtLogLikelihood: ptr.To(-91.12539763628888),
			IrtUsers:         ptr.To(281.0),
			IsExperimental:   ptr.To(true),
		},
		"abc182_e": {
			Slope:            ptr.To(-0.0008170719400144961),
			Intercept:        ptr.To(8.456092242494435),
			Variance:         ptr.To(0.27113303135125566),
			Difficulty:       ptr.To[int64](1098),
			Discrimination:   ptr.To(0.004479398673070138),
			IrtLogLikelihood: ptr.To(-1838.1064873009134),
			IrtUsers:         ptr.To(6840.0),
			IsExperimental:   ptr.To(false),
		},
	}

	updatedAt := time.Now()
	expected := []Difficulty{
		{
			ProblemID:        "abc019_2",
			Slope:            ptr.To(-0.0011830536916685426),
			Intercept:        ptr.To(7.8282534949123725),
			Variance:         ptr.To(0.39991055818601307),
			Difficulty:       ptr.To[int64](534),
			Discrimination:   ptr.To(0.004479398673070138),
			IrtLoglikelihood: ptr.To(-91.12539763628888),
			IrtUsers:         ptr.To(281.0),
			IsExperimental:   ptr.To(true),
			UpdatedAt:        updatedAt,
		},
		{
			ProblemID:        "abc182_e",
			Slope:            ptr.To(-0.0008170719400144961),
			Intercept:        ptr.To(8.456092242494435),
			Variance:         ptr.To(0.27113303135125566),
			Difficulty:       ptr.To[int64](1098),
			Discrimination:   ptr.To(0.004479398673070138),
			IrtLoglikelihood: ptr.To(-1838.1064873009134),
			IrtUsers:         ptr.To(6840.0),
			IsExperimental:   ptr.To(false),
			UpdatedAt:        updatedAt,
		},
	}
	actual := NewDifficulties(difficulties, updatedAt)

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

	updatedAt := time.Now()
	want := Submission{
		ID:          4031,
		EpochSecond: 1305359825,
		ProblemID:   "utpc2011_9",
		ContestID:   ptr.To("utpc2011"),
		UserID:      ptr.To("old_170"),
		Language:    ptr.To("C++ (GCC 4.4.7)"),
		Point:       ptr.To(0.0),
		Length:      ptr.To[int32](2566),
		Result:      ptr.To("WR"),
		UpdatedAt:   updatedAt,
	}
	actual := NewSubmission(s, updatedAt)

	if !reflect.DeepEqual(want, actual) {
		t.Errorf("expected \n%+v\n, but got \n%+v\n", want, actual)
	}
}

func TestNewUser(t *testing.T) {
	u := atcoder.User{
		UserID:        "tourist",
		Rating:        3798,
		HighestRating: 4229,
		Affiliation:   ptr.To("ITMO University"),
		BirthYear:     ptr.To[int32](1994),
		Country:       ptr.To("BY"),
		Crown:         ptr.To("crown_champion"),
		JoinCount:     60,
		Rank:          1,
		ActiveRank:    ptr.To[int32](1),
		Wins:          22,
	}

	updatedAt := time.Now()
	want := User{
		UserID:        "tourist",
		Rating:        3798,
		HighestRating: 4229,
		Affiliation:   ptr.To("ITMO University"),
		BirthYear:     ptr.To[int32](1994),
		Country:       ptr.To("BY"),
		Crown:         ptr.To("crown_champion"),
		JoinCount:     60,
		Rank:          1,
		ActiveRank:    ptr.To[int32](1),
		Wins:          22,
		UpdatedAt:     updatedAt,
	}
	actual := NewUser(u, updatedAt)

	if !reflect.DeepEqual(want, actual) {
		t.Errorf("expected \n%+v\n, but got \n%+v\n", want, actual)
	}
}

func TestMap(t *testing.T) {
	src := []atcoder.Contest{
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

	updatedAt := time.Now()
	want := []Contest{
		{
			ContestID:        "abc300",
			StartEpochSecond: 1682769600,
			DurationSecond:   6000,
			Title:            "ユニークビジョンプログラミングコンテスト2023 春 (AtCoder Beginner Contest 300)",
			RateChange:       " ~ 1999",
			Category:         "ABC",
			UpdatedAt:        updatedAt,
		},
		{
			ContestID:        "abc301",
			StartEpochSecond: 1683979200,
			DurationSecond:   6300,
			Title:            "パナソニックグループプログラミングコンテスト2023（AtCoder Beginner Contest 301）",
			RateChange:       " ~ 1999",
			Category:         "ABC",
			UpdatedAt:        updatedAt,
		},
	}

	actual := slices.Collect(Map(NewContest, slices.Values(src), updatedAt))

	if !reflect.DeepEqual(want, actual) {
		t.Errorf("expected \n%+v\n, but got \n%+v\n", want, actual)
	}
}
