package domain

import (
	"fjnkt98/atcodersearch/server/utility"
	"reflect"
	"testing"
)

func TestSearchSubmissionFilterParam(t *testing.T) {
	cases := []struct {
		name string
		p    SearchSubmissionFilterParam
		want []string
	}{
		{name: "empty", p: SearchSubmissionFilterParam{}, want: []string{}},
		{name: "epoch_second", p: SearchSubmissionFilterParam{EpochSecond: utility.IntegerRange{From: ptr(1000), To: ptr(2000)}}, want: []string{"{!tag=epoch_second}epoch_second:[1000 TO 2000}"}},
		{name: "problem_id", p: SearchSubmissionFilterParam{ProblemID: []string{"abc123_a", "abc123_b"}}, want: []string{`{!tag=problem_id}problem_id:("abc123_a" OR "abc123_b")`}},
		{name: "contest_id", p: SearchSubmissionFilterParam{ContestID: []string{"abc123", "abc234"}}, want: []string{`{!tag=contest_id}contest_id:("abc123" OR "abc234")`}},
		{name: "user_id", p: SearchSubmissionFilterParam{UserID: []string{"fjnkt98", "fjnkt99"}}, want: []string{`{!tag=user_id}user_id:("fjnkt98" OR "fjnkt99")`}},
		{name: "language", p: SearchSubmissionFilterParam{Language: []string{"C++ (GCC 9.6.2)", "CPython 3.11"}}, want: []string{`{!tag=language}language:("C\+\+ \(GCC 9.6.2\)" OR "CPython 3.11")`}},
		{name: "language_group", p: SearchSubmissionFilterParam{LanguageGroup: []string{"Python", "Go"}}, want: []string{`{!tag=language_group}language_group:("Python" OR "Go")`}},
		{name: "point", p: SearchSubmissionFilterParam{Point: utility.FloatRange{From: ptr(0.0), To: ptr(10.0)}}, want: []string{"{!tag=point}point:[0.000000 TO 10.000000}"}},
		{name: "length", p: SearchSubmissionFilterParam{Length: utility.IntegerRange{From: ptr(1000), To: ptr(2000)}}, want: []string{"{!tag=length}length:[1000 TO 2000}"}},
		{name: "result", p: SearchSubmissionFilterParam{Result: []string{"AC"}}, want: []string{`{!tag=result}result:("AC")`}},
		{name: "execution_time", p: SearchSubmissionFilterParam{ExecutionTime: utility.IntegerRange{From: ptr(1000), To: ptr(2000)}}, want: []string{"{!tag=execution_time}execution_time:[1000 TO 2000}"}},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := utility.Filter(tt.p)

			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("expected %+v, but got %+v", tt.want, result)
			}
		})
	}
}
