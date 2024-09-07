package atcoder

import (
	"fjnkt98/atcodersearch/pkg/ptr"
	"reflect"
	"testing"
)

// func TestFetchContests(t *testing.T) {
// 	c := NewAtCoderProblemsClient()
// 	_, err := c.FetchContests(context.Background())
// 	if err != nil {
// 		t.Fatalf("failed to fetch contests: %v", err)
// 	}
// }

// func TestFetchProblems(t *testing.T) {
// 	c := NewAtCoderProblemsClient()
// 	_, err := c.FetchProblems(context.Background())
// 	if err != nil {
// 		t.Fatalf("failed to fetch problems: %v", err)
// 	}
// }

// func TestFetchDifficulties(t *testing.T) {
// 	c := NewAtCoderProblemsClient()
// 	_, err := c.FetchDifficulties(context.Background())
// 	if err != nil {
// 		t.Fatalf("failed to fetch difficulties: %v", err)
// 	}
// }

func TestContestRatedTarget(t *testing.T) {
	cases := []struct {
		name    string
		contest Contest
		want    RatedTarget
	}{
		{name: "before agc001", contest: Contest{StartEpochSecond: 1468670399}, want: RatedTarget{Kind: UNRATED}},
		{name: "unrated", contest: Contest{StartEpochSecond: 1468670401, RateChange: "-"}, want: RatedTarget{Kind: UNRATED}},
		{name: "all", contest: Contest{StartEpochSecond: 1468670401, RateChange: "All"}, want: RatedTarget{Kind: ALL}},
		{name: "upper bound", contest: Contest{StartEpochSecond: 1468670401, RateChange: " ~ 1199"}, want: RatedTarget{Kind: UPPER_BOUND, To: ptr.To(1199)}},
		{name: "upper bound2", contest: Contest{StartEpochSecond: 1468670401, RateChange: " ~ 2799"}, want: RatedTarget{Kind: UPPER_BOUND, To: ptr.To(2799)}},
		{name: "lower bound", contest: Contest{StartEpochSecond: 1468670401, RateChange: "1200 ~ "}, want: RatedTarget{Kind: LOWER_BOUND, From: ptr.To(1200)}},
		{name: "range", contest: Contest{StartEpochSecond: 1468670401, RateChange: "1200 ~ 2799"}, want: RatedTarget{Kind: RANGE, From: ptr.To(1200), To: ptr.To(2799)}},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := tt.contest.RatedTarget()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("expected %v, but got %v", tt.want, got)
			}
		})
	}
}

func TestContestCategorize(t *testing.T) {
	cases := []struct {
		name    string
		contest Contest
		want    string
	}{
		{name: "abc", contest: Contest{ID: "abc042", Title: "AtCoder Beginner Contest 042", StartEpochSecond: 1469275200, RateChange: " ~ 1199"}, want: "ABC"},
		{name: "abc-like", contest: Contest{ID: "zone2021", Title: "ZONeエナジー プログラミングコンテスト  “HELLO SPACE”", StartEpochSecond: 1619870400, RateChange: " ~ 1999"}, want: "ABC-Like"},
		{name: "other sponsored", contest: Contest{ID: "jsc2019-final", Title: "第一回日本最強プログラマー学生選手権決勝", StartEpochSecond: 1569728700, RateChange: "-"}, want: "Other Sponsored"},
		{name: "ttpc2019", contest: Contest{ID: "ttpc2019", Title: "東京工業大学プログラミングコンテスト2019", StartEpochSecond: 1567224300, RateChange: "-"}, want: "Other Contests"},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := tt.contest.Categorize()
			if got != tt.want {
				t.Errorf("expected %s, but got %s", tt.want, got)
			}
		})
	}
}

func TestProblemURL(t *testing.T) {
	p := Problem{
		ID:        "abc001_a",
		ContestID: "abc001",
	}
	want := "https://atcoder.jp/contests/abc001/tasks/abc001_a"
	result := p.URL()
	if result != want {
		t.Errorf("expected %s, but got %s", want, result)
	}
}
