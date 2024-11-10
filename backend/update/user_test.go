package update

import (
	"context"
	"fjnkt98/atcodersearch/internal/testutil"
	"fjnkt98/atcodersearch/repository"
	"path/filepath"
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
		{name: "counts", user: UserRow{UserID: "fjnkt98", Rating: 1123, HighestRating: 1123, Affiliation: nil, BirthYear: ptr.To(1998), Country: ptr.To("JP"), Crown: ptr.To("user-green-4"), JoinCount: 40, Rank: 18728, ActiveRank: ptr.To(10256), Wins: 0, SubmissionCount: 2, Accepted: 1}, want: UserDocument{UserID: "fjnkt98", Rating: 1123, RatingFacet: " 800 ~ 1200", HighestRating: 1123, Affiliation: nil, BirthYear: ptr.To(1998), BirthYearFacet: "1990 ~ 2000", Country: ptr.To("JP"), Crown: ptr.To("user-green-4"), JoinCount: 40, JoinCountFacet: "  40 ~   60", Rank: 18728, ActiveRank: ptr.To(10256), Wins: 0, UserURL: "https://atcoder.jp/users/fjnkt98", Accepted: 1, SubmissionCount: 2}},
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

func TestReadUserRows(t *testing.T) {
	matches, err := filepath.Glob("./testdata/*.sql")
	if err != nil {
		t.Fatal(err)
	}
	files := make([]string, len(matches))
	for i, m := range matches {
		file, err := filepath.Abs(m)
		if err != nil {
			t.Fatal(err)
		}
		files[i] = file
	}

	ctx := context.Background()

	_, dsn, stopDB, err := testutil.CreateDBContainer(files...)
	if err != nil {
		t.Fatal(err)
	}
	pool, err := repository.NewPool(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		stopDB()
	})

	reader := NewUserRowReader(pool)

	ch := make(chan UserRow, 10)
	if err := reader.ReadRows(ctx, ch); err != nil {
		t.Error(err)
	}
	close(ch)

	rows := make([]UserRow, 0)
	for row := range ch {
		rows = append(rows, row)
	}

	want := []UserRow{
		{
			UserID:          "fjnkt98",
			Rating:          1123,
			HighestRating:   1123,
			Affiliation:     nil,
			BirthYear:       ptr.To(1998),
			Country:         ptr.To("JP"),
			Crown:           ptr.To("user-green-4"),
			JoinCount:       40,
			Rank:            18728,
			ActiveRank:      ptr.To(10256),
			Wins:            0,
			SubmissionCount: 2,
			Accepted:        1,
		},
	}

	if !reflect.DeepEqual(want, rows) {
		t.Errorf("expected %+v, but got %+v", want, rows)
	}
}
