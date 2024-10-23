package update

import (
	"context"
	"reflect"
	"testing"

	"k8s.io/utils/ptr"
)

func TestUserDocument(t *testing.T) {
	cases := []struct {
		name string
		user UserRow
		want UserDocument
	}{
		{name: "less than", user: UserRow{Rating: 0, BirthYear: ptr.To(1969), JoinCount: 0}, want: UserDocument{Rating: 0, RatingFacet: "   0 ~  400", BirthYear: ptr.To(1969), BirthYearFacet: "     ~ 1970", JoinCount: 0, JoinCountFacet: "   0 ~   20", UserURL: "https://atcoder.jp/users/"}},
		{name: "lower bound", user: UserRow{Rating: 0, BirthYear: ptr.To(1970), JoinCount: 0}, want: UserDocument{Rating: 0, RatingFacet: "   0 ~  400", BirthYear: ptr.To(1970), BirthYearFacet: "1970 ~ 1980", JoinCount: 0, JoinCountFacet: "   0 ~   20", UserURL: "https://atcoder.jp/users/"}},
		{name: "upper bound", user: UserRow{Rating: 3999, BirthYear: ptr.To(2019), JoinCount: 99}, want: UserDocument{Rating: 3999, RatingFacet: "3600 ~ 4000", BirthYear: ptr.To(2019), BirthYearFacet: "2010 ~ 2020", JoinCount: 99, JoinCountFacet: "  80 ~  100", UserURL: "https://atcoder.jp/users/"}},
		{name: "greater than", user: UserRow{Rating: 4000, BirthYear: ptr.To(2020), JoinCount: 100}, want: UserDocument{Rating: 4000, RatingFacet: "4000 ~     ", BirthYear: ptr.To(2020), BirthYearFacet: "2020 ~     ", JoinCount: 100, JoinCountFacet: " 100 ~     ", UserURL: "https://atcoder.jp/users/"}},
	}

	ctx := context.Background()
	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual, err := tt.user.Document(ctx)
			if err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(tt.want, actual) {
				t.Errorf("expect %+v, but got %+v", tt.want, actual)
			}
		})
	}
}
