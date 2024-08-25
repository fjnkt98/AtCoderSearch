package atcoder

import (
	"reflect"
	"testing"
)

func TestRatedTarget(t *testing.T) {
	cases := []struct {
		name    string
		contest Contest
		want    RatedTarget
	}{
		{name: "before agc001", contest: Contest{StartEpochSecond: 1468670399}, want: RatedTarget{UNRATED, nil, nil}},
		{name: "unrated", contest: Contest{StartEpochSecond: 1468670401, RateChange: "-"}, want: RatedTarget{UNRATED, nil, nil}},
		{name: "all", contest: Contest{StartEpochSecond: 1468670401, RateChange: "All"}, want: RatedTarget{ALL, nil, nil}},
		{name: "upper bound 1", contest: Contest{StartEpochSecond: 1468670401, RateChange: " ~ 1199"}, want: RatedTarget{UPPER_BOUND, nil, ptr(1199)}},
		{name: "upper bound 2", contest: Contest{StartEpochSecond: 1468670401, RateChange: " ~ 2799"}, want: RatedTarget{UPPER_BOUND, nil, ptr(2799)}},
		{name: "lower bound 1", contest: Contest{StartEpochSecond: 1468670401, RateChange: "1200 ~ "}, want: RatedTarget{LOWER_BOUND, ptr(1200), nil}},
		{name: "range", contest: Contest{StartEpochSecond: 1468670401, RateChange: "1200 ~ 2799"}, want: RatedTarget{RANGE, ptr(1200), ptr(2799)}},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.contest.RatedTarget()

			if !reflect.DeepEqual(tt.want, actual) {
				t.Errorf("expected %+v, but got %+v", tt.want, actual)
			}
		})
	}
}

func TestCategorize(t *testing.T) {
	cases := []struct {
		name    string
		contest Contest
		want    string
	}{
		{
			name: "abc",
			contest: Contest{
				DurationSecond:   6000,
				ID:               "abc042",
				RateChange:       " ~ 1199",
				StartEpochSecond: 1469275200,
				Title:            "AtCoder Beginner Contest 042",
			},
			want: "ABC",
		},
		{
			name: "abc-like",
			contest: Contest{
				DurationSecond:   6000,
				ID:               "zone2021",
				RateChange:       " ~ 1999",
				StartEpochSecond: 1619870400,
				Title:            "ZONeエナジー プログラミングコンテスト  “HELLO SPACE”",
			},
			want: "ABC-Like",
		},
		{
			name: "other sponsored",
			contest: Contest{
				DurationSecond:   10800,
				ID:               "jsc2019-final",
				RateChange:       "-",
				StartEpochSecond: 1569728700,
				Title:            "第一回日本最強プログラマー学生選手権決勝",
			},
			want: "Other Sponsored",
		},
		{
			name: "other contests",
			contest: Contest{
				DurationSecond:   18000,
				ID:               "ttpc2019",
				RateChange:       "-",
				StartEpochSecond: 1567224300,
				Title:            "東京工業大学プログラミングコンテスト2019",
			},
			want: "Other Contests",
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.contest.Categorize()

			if tt.want != actual {
				t.Errorf("expected %s, but got %s", tt.want, actual)
			}
		})
	}
}
