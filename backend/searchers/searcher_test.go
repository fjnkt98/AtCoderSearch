package searchers

import (
	"context"
	"errors"
	pb "fjnkt98/atcodersearch/grpc/atcodersearch/v1"
	"fjnkt98/atcodersearch/internal/testutil"
	"fjnkt98/atcodersearch/repository"
	"fjnkt98/atcodersearch/update"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/meilisearch/meilisearch-go"
	"k8s.io/utils/ptr"
)

type doc struct {
	ID     string `mapstructure:"id"`
	Name   string `mapstructure:"name"`
	Grade  int
	Secret string `mapstructure:"-"`
}

func TestFieldList(t *testing.T) {
	expected := []string{"id", "name", "Grade"}
	actual := FieldList(new(doc))

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected \n%+v\n , but got \n%+v\n", expected, actual)
	}
}

func TestParseFacetDistribution(t *testing.T) {
	facet := map[string]interface{}{"category": map[string]interface{}{"ABC": 2237.0, "ABC-Like": 184.0, "AGC": 395.0, "AGC-Like": 77.0, "AHC": 38.0, "ARC": 810.0, "ARC-Like": 10.0, "JOI": 801.0, "Marathon": 22.0, "Other Contests": 1712.0, "Other Sponsored": 489.0, "PAST": 255.0}, "difficultyFacet": map[string]interface{}{"     ~    0": 3970.0, "   0 ~  400": 313.0, " 400 ~  800": 368.0, " 800 ~ 1200": 363.0, "1200 ~ 1600": 387.0, "1600 ~ 2000": 394.0, "2000 ~ 2400": 369.0, "2400 ~ 2800": 302.0, "2800 ~ 3200": 244.0, "3200 ~ 3600": 170.0, "3600 ~     ": 150.0}}

	res := ParseFacetDistribution(facet)

	want := map[string]map[string]int64{
		"category": {
			"ABC":             2237.0,
			"ABC-Like":        184.0,
			"AGC":             395.0,
			"AGC-Like":        77.0,
			"AHC":             38.0,
			"ARC":             810.0,
			"ARC-Like":        10.0,
			"JOI":             801.0,
			"Marathon":        22.0,
			"Other Contests":  1712.0,
			"Other Sponsored": 489.0,
			"PAST":            255.0,
		},
		"difficultyFacet": {
			"     ~    0": 3970.0,
			"   0 ~  400": 313.0,
			" 400 ~  800": 368.0,
			" 800 ~ 1200": 363.0,
			"1200 ~ 1600": 387.0,
			"1600 ~ 2000": 394.0,
			"2000 ~ 2400": 369.0,
			"2400 ~ 2800": 302.0,
			"2800 ~ 3200": 244.0,
			"3200 ~ 3600": 170.0,
			"3600 ~     ": 150.0,
		},
	}

	if !reflect.DeepEqual(want, res) {
		t.Errorf("expected %+v, but got %+v", want, res)
	}
}

func TestCreateSearchProblemByKeywordRequest(t *testing.T) {
	fields := []string{"problemId", "problemTitle", "problemUrl", "contestId", "contestTitle", "contestUrl", "difficulty", "color", "startAt", "duration", "rateChange", "category", "isExperimental"}

	{
		cases := []struct {
			name string
			req  *pb.SearchProblemByKeywordRequest
			want *meilisearch.SearchRequest
		}{
			{name: "empty", req: &pb.SearchProblemByKeywordRequest{}, want: &meilisearch.SearchRequest{Sort: []string{"problemId:asc"}, HitsPerPage: 0, Page: 1, AttributesToRetrieve: fields}},
			{name: "pagination", req: &pb.SearchProblemByKeywordRequest{Limit: ptr.To[uint64](20), Page: ptr.To[uint64](2)}, want: &meilisearch.SearchRequest{Sort: []string{"problemId:asc"}, HitsPerPage: 20, Page: 2, AttributesToRetrieve: fields}},
			{name: "sort(valid)", req: &pb.SearchProblemByKeywordRequest{Sorts: []string{"startAt:desc", "difficulty:asc"}}, want: &meilisearch.SearchRequest{Sort: []string{"startAt:desc", "difficulty:asc", "problemId:asc"}, Page: 1, AttributesToRetrieve: fields}},
			{name: "facet(valid)", req: &pb.SearchProblemByKeywordRequest{Facets: []string{"category", "difficulty"}}, want: &meilisearch.SearchRequest{Facets: []string{"category", "difficultyFacet"}, Sort: []string{"problemId:asc"}, Page: 1, AttributesToRetrieve: fields}},
			{name: "filter by category", req: &pb.SearchProblemByKeywordRequest{Categories: []string{"ABC", "ARC"}}, want: &meilisearch.SearchRequest{Filter: [][]string{{"category = 'ABC'", "category = 'ARC'"}}, Sort: []string{"problemId:asc"}, Page: 1, AttributesToRetrieve: fields}},
			{name: "filter by difficulty(from only)", req: &pb.SearchProblemByKeywordRequest{Difficulty: &pb.IntRange{From: ptr.To[int64](800)}}, want: &meilisearch.SearchRequest{Filter: [][]string{{"difficulty >= 800"}}, Sort: []string{"problemId:asc"}, Page: 1, AttributesToRetrieve: fields}},
			{name: "filter by difficulty(to only)", req: &pb.SearchProblemByKeywordRequest{Difficulty: &pb.IntRange{To: ptr.To[int64](1200)}}, want: &meilisearch.SearchRequest{Filter: [][]string{{"difficulty < 1200"}}, Sort: []string{"problemId:asc"}, Page: 1, AttributesToRetrieve: fields}},
			{name: "filter by difficulty(both)", req: &pb.SearchProblemByKeywordRequest{Difficulty: &pb.IntRange{From: ptr.To[int64](800), To: ptr.To[int64](1200)}}, want: &meilisearch.SearchRequest{Filter: [][]string{{"difficulty >= 800"}, {"difficulty < 1200"}}, Sort: []string{"problemId:asc"}, Page: 1, AttributesToRetrieve: fields}},
			{name: "filter for experimental", req: &pb.SearchProblemByKeywordRequest{Experimental: ptr.To(true)}, want: &meilisearch.SearchRequest{Filter: [][]string{{"isExperimental = true"}}, Sort: []string{"problemId:asc"}, Page: 1, AttributesToRetrieve: fields}},
			{name: "filter for not experimental", req: &pb.SearchProblemByKeywordRequest{Experimental: ptr.To(false)}, want: &meilisearch.SearchRequest{Filter: [][]string{{"isExperimental = false"}}, Sort: []string{"problemId:asc"}, Page: 1, AttributesToRetrieve: fields}},
		}

		for _, tt := range cases {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				actual, err := createSearchProblemByKeywordQuery(tt.req)
				if err != nil {
					t.Error(err)
				}
				if !reflect.DeepEqual(tt.want, actual) {
					t.Errorf("expected %+v, but got %+v", tt.want, actual)
				}
			})
		}
	}

	{
		cases := []struct {
			name string
			req  *pb.SearchProblemByKeywordRequest
		}{
			{name: "sort(without direction)", req: &pb.SearchProblemByKeywordRequest{Sorts: []string{"startAt"}}},
			{name: "sort(not allowed field)", req: &pb.SearchProblemByKeywordRequest{Sorts: []string{"isExperimental:asc"}}},
			{name: "facet(not allowed field)", req: &pb.SearchProblemByKeywordRequest{Facets: []string{"difficultyFacet"}}},
		}

		for _, tt := range cases {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				_, err := createSearchProblemByKeywordQuery(tt.req)
				if err == nil {
					t.Fatal("should error")
				}

				if !errors.Is(err, ErrInvalidRequest) {
					t.Errorf("expected %+v, but got %+v", ErrInvalidRequest, err)
				}
			})
		}
	}
}

func TestCreateSearchUserRequest(t *testing.T) {
	fields := []string{"userId", "rating", "highestRating", "affiliation", "birthYear", "country", "crown", "joinCount", "rank", "activeRank", "wins", "userUrl"}

	{
		cases := []struct {
			name string
			req  *pb.SearchUserRequest
			want *meilisearch.SearchRequest
		}{
			{name: "empty", req: &pb.SearchUserRequest{}, want: &meilisearch.SearchRequest{Sort: []string{"userId:asc"}, Page: 1, AttributesToRetrieve: fields}},
			{name: "pagination", req: &pb.SearchUserRequest{Limit: ptr.To[int64](20)}, want: &meilisearch.SearchRequest{HitsPerPage: 20, Sort: []string{"userId:asc"}, Page: 1, AttributesToRetrieve: fields}},
			{name: "sort(valid)", req: &pb.SearchUserRequest{Sorts: []string{"rating:desc", "birthYear:asc"}}, want: &meilisearch.SearchRequest{Sort: []string{"rating:desc", "birthYear:asc", "userId:asc"}, Page: 1, AttributesToRetrieve: fields}},
			{name: "facet(valid)", req: &pb.SearchUserRequest{Facets: []string{"country", "rating", "birthYear", "joinCount"}}, want: &meilisearch.SearchRequest{Facets: []string{"country", "ratingFacet", "birthYearFacet", "joinCountFacet"}, Sort: []string{"userId:asc"}, Page: 1, AttributesToRetrieve: fields}},
			{name: "filter by user id", req: &pb.SearchUserRequest{UserIds: []string{"user1", "user2"}}, want: &meilisearch.SearchRequest{Filter: [][]string{{"userId = 'user1'", "userId = 'user2'"}}, Sort: []string{"userId:asc"}, Page: 1, AttributesToRetrieve: fields}},
			{name: "filter by rating", req: &pb.SearchUserRequest{Rating: &pb.IntRange{From: ptr.To[int64](800), To: ptr.To[int64](1200)}}, want: &meilisearch.SearchRequest{Filter: [][]string{{"rating >= 800"}, {"rating < 1200"}}, Sort: []string{"userId:asc"}, Page: 1, AttributesToRetrieve: fields}},
			{name: "filter by birth year", req: &pb.SearchUserRequest{BirthYear: &pb.IntRange{From: ptr.To[int64](1998), To: ptr.To[int64](2000)}}, want: &meilisearch.SearchRequest{Filter: [][]string{{"birthYear >= 1998"}, {"birthYear < 2000"}}, Sort: []string{"userId:asc"}, Page: 1, AttributesToRetrieve: fields}},
			{name: "filter by join count", req: &pb.SearchUserRequest{JoinCount: &pb.IntRange{From: ptr.To[int64](5), To: ptr.To[int64](10)}}, want: &meilisearch.SearchRequest{Filter: [][]string{{"joinCount >= 5"}, {"joinCount < 10"}}, Sort: []string{"userId:asc"}, Page: 1, AttributesToRetrieve: fields}},
			{name: "filter by country", req: &pb.SearchUserRequest{Countries: []string{"JP"}}, want: &meilisearch.SearchRequest{Filter: [][]string{{"country = 'JP'"}}, Sort: []string{"userId:asc"}, Page: 1, AttributesToRetrieve: fields}},
		}

		for _, tt := range cases {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				actual, err := createSearchUserQuery(tt.req)
				if err != nil {
					t.Error(err)
				}
				if !reflect.DeepEqual(tt.want, actual) {
					t.Errorf("expected %+v, but got %+v", tt.want, actual)
				}
			})
		}
	}
}

var ABC300A = &pb.Problem{
	ProblemId:      "abc300_a",
	ProblemTitle:   "A. N-choice question",
	ProblemUrl:     "https://atcoder.jp/contests/abc300/tasks/abc300_a",
	ContestId:      "abc300",
	ContestTitle:   "ユニークビジョンプログラミングコンテスト2023 春 (AtCoder Beginner Contest 300)",
	ContestUrl:     "https://atcoder.jp/contests/abc300",
	Difficulty:     ptr.To[int64](-1147),
	StartAt:        1682769600,
	Duration:       6000,
	RateChange:     " ~ 1999",
	Category:       "ABC",
	IsExperimental: false,
}

var ABC300B = &pb.Problem{
	ProblemId:      "abc300_b",
	ProblemTitle:   "B. Same Map in the RPG World",
	ProblemUrl:     "https://atcoder.jp/contests/abc300/tasks/abc300_b",
	ContestId:      "abc300",
	ContestTitle:   "ユニークビジョンプログラミングコンテスト2023 春 (AtCoder Beginner Contest 300)",
	ContestUrl:     "https://atcoder.jp/contests/abc300",
	Difficulty:     ptr.To[int64](350),
	StartAt:        1682769600,
	Duration:       6000,
	RateChange:     " ~ 1999",
	Category:       "ABC",
	IsExperimental: false,
}

var ARC184A = &pb.Problem{
	ProblemId:      "arc184_a",
	ProblemTitle:   "A. Appraiser",
	ProblemUrl:     "https://atcoder.jp/contests/arc184/tasks/arc184_a",
	ContestId:      "arc184",
	ContestTitle:   "AtCoder Regular Contest 184",
	ContestUrl:     "https://atcoder.jp/contests/arc184",
	Difficulty:     ptr.To[int64](1383),
	StartAt:        1727006400,
	Duration:       7200,
	RateChange:     "1200 ~ 2799",
	Category:       "ARC",
	IsExperimental: false,
}

var ARC184B = &pb.Problem{
	ProblemId:      "arc184_b",
	ProblemTitle:   "B. 123 Set",
	ProblemUrl:     "https://atcoder.jp/contests/arc184/tasks/arc184_b",
	ContestId:      "arc184",
	ContestTitle:   "AtCoder Regular Contest 184",
	ContestUrl:     "https://atcoder.jp/contests/arc184",
	Difficulty:     ptr.To[int64](2867),
	StartAt:        1727006400,
	Duration:       7200,
	RateChange:     "1200 ~ 2799",
	Category:       "ARC",
	IsExperimental: false,
}

var User1 = &pb.User{
	UserId:        "user1",
	Rating:        2563,
	HighestRating: 2563,
	Affiliation:   nil,
	BirthYear:     nil,
	Country:       ptr.To("CN"),
	Crown:         ptr.To("user-orange-2"),
	JoinCount:     23,
	Rank:          470,
	ActiveRank:    nil,
	Wins:          0,
	UserUrl:       "https://atcoder.jp/users/user1",
}

var User2 = &pb.User{
	UserId:        "user2",
	Rating:        3710,
	HighestRating: 3802,
	Affiliation:   nil,
	BirthYear:     nil,
	Country:       ptr.To("CH"),
	Crown:         nil,
	JoinCount:     21,
	Rank:          3,
	ActiveRank:    nil,
	Wins:          2,
	UserUrl:       "https://atcoder.jp/users/user2",
}

var User3 = &pb.User{
	UserId:        "user3",
	Rating:        3658,
	HighestRating: 3683,
	Affiliation:   ptr.To("MIT"),
	BirthYear:     ptr.To[int64](2001),
	Country:       ptr.To("US"),
	Crown:         nil,
	JoinCount:     48,
	Rank:          3,
	ActiveRank:    nil,
	Wins:          0,
	UserUrl:       "https://atcoder.jp/users/user3",
}

var User4 = &pb.User{
	UserId:        "user4",
	Rating:        3604,
	HighestRating: 3814,
	Affiliation:   nil,
	BirthYear:     ptr.To[int64](1997),
	Country:       ptr.To("JP"),
	Crown:         ptr.To("crown_gold"),
	JoinCount:     38,
	Rank:          8,
	ActiveRank:    ptr.To[int64](5),
	Wins:          2,
	UserUrl:       "https://atcoder.jp/users/user4",
}

var Submission1 = &pb.Submission{
	SubmissionId:  1,
	SubmittedAt:   1729434074,
	SubmissionUrl: "https://atcoder.jp/contests/abc300/submissions/1",
	ProblemId:     "abc300_a",
	ProblemTitle:  "A. N-choice question",
	ProblemUrl:    "https://atcoder.jp/contests/abc300/tasks/abc300_a",
	ContestId:     "abc300",
	ContestTitle:  "ユニークビジョンプログラミングコンテスト2023 春 (AtCoder Beginner Contest 300)",
	ContestUrl:    "https://atcoder.jp/contests/abc300",
	Category:      "ABC",
	Difficulty:    ptr.To[int64](-1147),
	UserId:        "fjnkt98",
	Language:      "Python (CPython 3.11.4)",
	LanguageGroup: "Python",
	Point:         100.0,
	Length:        1024,
	Result:        "AC",
	ExecutionTime: ptr.To[int64](22),
}

var Submission2 = &pb.Submission{
	SubmissionId:  2,
	SubmittedAt:   1729434074,
	SubmissionUrl: "https://atcoder.jp/contests/abc300/submissions/2",
	ProblemId:     "abc300_b",
	ProblemTitle:  "B. Same Map in the RPG World",
	ProblemUrl:    "https://atcoder.jp/contests/abc300/tasks/abc300_b",
	ContestId:     "abc300",
	ContestTitle:  "ユニークビジョンプログラミングコンテスト2023 春 (AtCoder Beginner Contest 300)",
	ContestUrl:    "https://atcoder.jp/contests/abc300",
	Category:      "ABC",
	Difficulty:    ptr.To[int64](350),
	UserId:        "fjnkt98",
	Language:      "Python (CPython 3.11.4)",
	LanguageGroup: "Python",
	Point:         200.0,
	Length:        1024,
	Result:        "WA",
	ExecutionTime: nil,
}

func TestSearcher(t *testing.T) {
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

	_, url, key, stopEngine, err := testutil.CreateEngineContainer()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		stopDB()
		stopEngine()
	})

	client := meilisearch.New(url, meilisearch.WithAPIKey(key))

	if err := update.UpdateIndex(
		ctx,
		update.NewProblemRowReader(pool),
		update.NewProblemIndexer(client),
		1000,
		1,
	); err != nil {
		t.Fatal(err)
	}

	if err := update.UpdateIndex(
		ctx,
		update.NewUserRowReader(pool),
		update.NewUserIndexer(client),
		1000,
		1,
	); err != nil {
		t.Fatal(err)
	}

	searcher := NewSearcher(client, pool)

	t.Run("SearchProblem: empty", func(t *testing.T) {
		res, err := searcher.SearchProblem(ctx, &pb.SearchProblemRequest{})
		if err != nil {
			t.Error(err)
		}

		want := []*pb.Problem{
			ABC300A,
			ABC300B,
			ARC184A,
			ARC184B,
		}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	{
		cases := []struct {
			name string
			req  *pb.SearchProblemRequest
			want []*pb.Problem
		}{
			{name: "SearchProblem: sort by start_at asc", req: &pb.SearchProblemRequest{Sorts: []string{"startAt:asc"}}, want: []*pb.Problem{ABC300A, ABC300B, ARC184A, ARC184B}},
			{name: "SearchProblem: sort by start_at desc", req: &pb.SearchProblemRequest{Sorts: []string{"startAt:desc"}}, want: []*pb.Problem{ARC184A, ARC184B, ABC300A, ABC300B}},
			{name: "SearchProblem: sort by difficulty asc", req: &pb.SearchProblemRequest{Sorts: []string{"difficulty:asc"}}, want: []*pb.Problem{ABC300A, ABC300B, ARC184A, ARC184B}},
			{name: "SearchProblem: sort by difficulty desc", req: &pb.SearchProblemRequest{Sorts: []string{"difficulty:desc"}}, want: []*pb.Problem{ARC184B, ARC184A, ABC300B, ABC300A}},
		}

		for _, tt := range cases {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				res, err := searcher.SearchProblem(ctx, tt.req)
				if err != nil {
					t.Error(err)
				}

				if !reflect.DeepEqual(tt.want, res.Items) {
					t.Errorf("expect %+v, but got %+v", tt.want, res.Items)
				}
			})
		}
	}

	t.Run("SearchProblem: filter by category", func(t *testing.T) {
		res, err := searcher.SearchProblem(ctx, &pb.SearchProblemRequest{
			Categories: []string{"ABC"},
		})

		if err != nil {
			t.Error(err)
		}

		want := []*pb.Problem{ABC300A, ABC300B}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("SearchProblem: filter by difficulty", func(t *testing.T) {
		res, err := searcher.SearchProblem(ctx, &pb.SearchProblemRequest{
			Difficulty: &pb.IntRange{
				From: ptr.To[int64](1000),
				To:   ptr.To[int64](1400),
			},
		})

		if err != nil {
			t.Error(err)
		}

		want := []*pb.Problem{ARC184A}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("SearchProblem: filter by user id", func(t *testing.T) {
		res, err := searcher.SearchProblem(ctx, &pb.SearchProblemRequest{
			UserId: ptr.To("fjnkt98"),
		})

		if err != nil {
			t.Error(err)
		}

		want := []*pb.Problem{ABC300A}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("SearchProblemByKeyword: keyword", func(t *testing.T) {
		res, err := searcher.SearchProblemByKeyword(ctx, &pb.SearchProblemByKeywordRequest{
			Q: "ABC300",
		})

		if err != nil {
			t.Error(err)
		}

		if len(res.Items) != 2 {
			t.Errorf("expect length of items is 2, but got %d", len(res.Items))
		}
	})

	t.Run("SearchProblemByKeyword: filter by category", func(t *testing.T) {
		res, err := searcher.SearchProblemByKeyword(ctx, &pb.SearchProblemByKeywordRequest{
			Categories: []string{"ARC"},
		})

		if err != nil {
			t.Error(err)
		}

		want := []*pb.Problem{
			ARC184A,
			ARC184B,
		}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("SearchProblemByKeyword: facet", func(t *testing.T) {
		res, err := searcher.SearchProblemByKeyword(ctx, &pb.SearchProblemByKeywordRequest{
			Facets: []string{"category", "difficulty"},
		})

		if err != nil {
			t.Error(err)
		}

		want := &pb.ProblemFacet{
			Categories: []*pb.Count{
				{
					Label: "ABC",
					Count: 2,
				},
				{
					Label: "ARC",
					Count: 2,
				},
			},
			Difficulties: []*pb.Count{
				{
					Label: "     ~    0",
					Count: 1,
				},
				{
					Label: "   0 ~  400",
					Count: 1,
				},
				{
					Label: "1200 ~ 1600",
					Count: 1,
				},
				{
					Label: "2800 ~ 3200",
					Count: 1,
				},
			},
		}

		if !reflect.DeepEqual(want, res.Facet) {
			t.Errorf("expect %+v, but got %+v", want, res.Facet)
		}
	})

	t.Run("SearchUser: empty", func(t *testing.T) {
		res, err := searcher.SearchUser(ctx, &pb.SearchUserRequest{})

		if err != nil {
			t.Error(err)
		}

		want := []*pb.User{User1, User2, User3, User4}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("SearchUser: keyword", func(t *testing.T) {
		res, err := searcher.SearchUser(ctx, &pb.SearchUserRequest{
			Q: "MIT",
		})

		if err != nil {
			t.Error(err)
		}

		want := []*pb.User{User3}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	{
		cases := []struct {
			name string
			req  *pb.SearchUserRequest
			want []*pb.User
		}{
			{name: "SearchUser: sort by rating asc", req: &pb.SearchUserRequest{Sorts: []string{"rating:asc"}}, want: []*pb.User{User1, User4, User3, User2}},
			{name: "SearchUser: sort by rating desc", req: &pb.SearchUserRequest{Sorts: []string{"rating:desc"}}, want: []*pb.User{User2, User3, User4, User1}},
			{name: "SearchUser: sort by birthYear asc", req: &pb.SearchUserRequest{Sorts: []string{"birthYear:asc"}}, want: []*pb.User{User4, User3, User1, User2}},
			{name: "SearchUser: sort by birthYear desc", req: &pb.SearchUserRequest{Sorts: []string{"birthYear:desc"}}, want: []*pb.User{User3, User4, User1, User2}},
		}

		for _, tt := range cases {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				res, err := searcher.SearchUser(ctx, tt.req)

				if err != nil {
					t.Error(err)
				}

				if !reflect.DeepEqual(tt.want, res.Items) {
					t.Errorf("expect %+v, but got %+v", tt.want, res.Items)
				}
			})
		}
	}

	t.Run("SearchUser: filter by user id", func(t *testing.T) {
		res, err := searcher.SearchUser(ctx, &pb.SearchUserRequest{
			UserIds: []string{"user1", "user4"},
		})

		if err != nil {
			t.Error(err)
		}

		want := []*pb.User{User1, User4}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("SearchUser: filter by rating", func(t *testing.T) {
		res, err := searcher.SearchUser(ctx, &pb.SearchUserRequest{
			Rating: &pb.IntRange{From: ptr.To[int64](2563), To: ptr.To[int64](2564)},
		})

		if err != nil {
			t.Error(err)
		}

		want := []*pb.User{User1}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("SearchUser: filter by birth year", func(t *testing.T) {
		res, err := searcher.SearchUser(ctx, &pb.SearchUserRequest{
			BirthYear: &pb.IntRange{From: ptr.To[int64](1997), To: ptr.To[int64](1998)},
		})

		if err != nil {
			t.Error(err)
		}

		want := []*pb.User{User4}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("SearchUser: filter by join count", func(t *testing.T) {
		res, err := searcher.SearchUser(ctx, &pb.SearchUserRequest{
			JoinCount: &pb.IntRange{From: ptr.To[int64](45), To: ptr.To[int64](50)},
		})

		if err != nil {
			t.Error(err)
		}

		want := []*pb.User{User3}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("SearchUser: filter by country", func(t *testing.T) {
		res, err := searcher.SearchUser(ctx, &pb.SearchUserRequest{
			Countries: []string{"JP"},
		})

		if err != nil {
			t.Error(err)
		}

		want := []*pb.User{User4}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("SearchUser: facet", func(t *testing.T) {
		res, err := searcher.SearchUser(ctx, &pb.SearchUserRequest{
			Facets: []string{"country", "rating", "birthYear", "joinCount"},
		})

		if err != nil {
			t.Error(err)
		}

		want := &pb.UserFacet{
			BirthYears: []*pb.Count{
				{
					Label: "1990 ~ 2000",
					Count: 1,
				},
				{
					Label: "2000 ~ 2010",
					Count: 1,
				},
			},
			Countries: []*pb.Count{
				{
					Label: "CH",
					Count: 1,
				},
				{
					Label: "CN",
					Count: 1,
				},
				{
					Label: "JP",
					Count: 1,
				},
				{
					Label: "US",
					Count: 1,
				},
			},
			JoinCounts: []*pb.Count{
				{
					Label: "  20 ~   40",
					Count: 3,
				},
				{
					Label: "  40 ~   60",
					Count: 1,
				},
			},
			Ratings: []*pb.Count{
				{
					Label: "2400 ~ 2800",
					Count: 1,
				},
				{
					Label: "3600 ~ 4000",
					Count: 3,
				},
			},
		}

		if !reflect.DeepEqual(want, res.Facet) {
			t.Errorf("expect %+v, but got %+v", want, res.Facet)
		}
	})

	t.Run("SearchSubmission: empty", func(t *testing.T) {
		res, err := searcher.SearchSubmission(ctx, &pb.SearchSubmissionRequest{})
		if err != nil {
			t.Error(err)
		}

		want := []*pb.Submission{Submission2, Submission1}
		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	{
		cases := []struct {
			name string
			req  *pb.SearchSubmissionRequest
			want []*pb.Submission
		}{
			{name: "SearchSubmission: sort by execution time", req: &pb.SearchSubmissionRequest{Sorts: []string{"executionTime:asc"}}, want: []*pb.Submission{Submission1, Submission2}},
			{name: "SearchSubmission: sort by epoch second", req: &pb.SearchSubmissionRequest{Sorts: []string{"epochSecond:asc"}}, want: []*pb.Submission{Submission2, Submission1}},
			{name: "SearchSubmission: sort by point", req: &pb.SearchSubmissionRequest{Sorts: []string{"point:asc"}}, want: []*pb.Submission{Submission1, Submission2}},
			{name: "SearchSubmission: sort by length", req: &pb.SearchSubmissionRequest{Sorts: []string{"length:asc"}}, want: []*pb.Submission{Submission2, Submission1}},
		}

		for _, tt := range cases {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				res, err := searcher.SearchSubmission(ctx, tt.req)
				if err != nil {
					t.Error(err)
				}

				if !reflect.DeepEqual(tt.want, res.Items) {
					t.Errorf("expect %+v, but got %+v", tt.want, res.Items)
				}
			})
		}
	}

	t.Run("SearchSubmission: filter by epoch second", func(t *testing.T) {
		res, err := searcher.SearchSubmission(ctx, &pb.SearchSubmissionRequest{
			EpochSecond: &pb.IntRange{
				From: ptr.To[int64](1729434073),
				To:   ptr.To[int64](1729434075),
			},
		})

		if err != nil {
			t.Error(err)
		}

		want := []*pb.Submission{Submission2, Submission1}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("SearchSubmission: filter by problem id", func(t *testing.T) {
		res, err := searcher.SearchSubmission(ctx, &pb.SearchSubmissionRequest{
			ProblemIds: []string{"abc300_a"},
		})

		if err != nil {
			t.Error(err)
		}

		want := []*pb.Submission{Submission1}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("SearchSubmission: filter by contest id", func(t *testing.T) {
		res, err := searcher.SearchSubmission(ctx, &pb.SearchSubmissionRequest{
			ContestIds: []string{"abc300"},
		})

		if err != nil {
			t.Error(err)
		}

		want := []*pb.Submission{Submission2, Submission1}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("SearchSubmission: filter by category", func(t *testing.T) {
		res, err := searcher.SearchSubmission(ctx, &pb.SearchSubmissionRequest{
			Categories: []string{"ABC"},
		})

		if err != nil {
			t.Error(err)
		}

		want := []*pb.Submission{Submission2, Submission1}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("SearchSubmission: filter by user id", func(t *testing.T) {
		res, err := searcher.SearchSubmission(ctx, &pb.SearchSubmissionRequest{
			UserIds: []string{"fjnkt98"},
		})

		if err != nil {
			t.Error(err)
		}

		want := []*pb.Submission{Submission2, Submission1}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("SearchSubmission: filter by language", func(t *testing.T) {
		res, err := searcher.SearchSubmission(ctx, &pb.SearchSubmissionRequest{
			Languages: []string{"Python (CPython 3.11.4)"},
		})

		if err != nil {
			t.Error(err)
		}

		want := []*pb.Submission{Submission2, Submission1}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("SearchSubmission: filter by language group", func(t *testing.T) {
		res, err := searcher.SearchSubmission(ctx, &pb.SearchSubmissionRequest{
			LanguageGroups: []string{"Python"},
		})

		if err != nil {
			t.Error(err)
		}

		want := []*pb.Submission{Submission2, Submission1}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("SearchSubmission: filter by point", func(t *testing.T) {
		res, err := searcher.SearchSubmission(ctx, &pb.SearchSubmissionRequest{
			Point: &pb.FloatRange{
				From: ptr.To(100.0),
				To:   ptr.To(150.0),
			},
		})

		if err != nil {
			t.Error(err)
		}

		want := []*pb.Submission{Submission1}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("SearchSubmission: filter by length", func(t *testing.T) {
		res, err := searcher.SearchSubmission(ctx, &pb.SearchSubmissionRequest{
			Length: &pb.IntRange{
				From: ptr.To[int64](1024),
				To:   ptr.To[int64](2048),
			},
		})

		if err != nil {
			t.Error(err)
		}

		want := []*pb.Submission{Submission2, Submission1}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("SearchSubmission: filter by result", func(t *testing.T) {
		res, err := searcher.SearchSubmission(ctx, &pb.SearchSubmissionRequest{
			Results: []string{"WA"},
		})

		if err != nil {
			t.Error(err)
		}

		want := []*pb.Submission{Submission2}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("SearchSubmission: filter by execution time", func(t *testing.T) {
		res, err := searcher.SearchSubmission(ctx, &pb.SearchSubmissionRequest{
			ExecutionTime: &pb.IntRange{
				From: ptr.To[int64](20),
				To:   ptr.To[int64](1000),
			},
		})

		if err != nil {
			t.Error(err)
		}

		want := []*pb.Submission{Submission1}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("GetCategory", func(t *testing.T) {
		res, err := searcher.GetCategory(ctx, &pb.GetCategoryRequest{})

		if err != nil {
			t.Error(err)
		}

		want := []string{"ABC", "ARC"}

		if !reflect.DeepEqual(want, res.Categories) {
			t.Errorf("expect %+v, but got %+v", want, res.Categories)
		}
	})

	t.Run("GetContest: empty", func(t *testing.T) {
		res, err := searcher.GetContest(ctx, &pb.GetContestRequest{})

		if err != nil {
			t.Error(err)
		}

		want := []string{"arc184", "abc300"}

		if !reflect.DeepEqual(want, res.Contests) {
			t.Errorf("expect %+v, but got %+v", want, res.Contests)
		}
	})

	t.Run("GetContest: with category", func(t *testing.T) {
		res, err := searcher.GetContest(ctx, &pb.GetContestRequest{
			Categories: []string{"ABC"},
		})

		if err != nil {
			t.Error(err)
		}

		want := []string{"abc300"}

		if !reflect.DeepEqual(want, res.Contests) {
			t.Errorf("expect %+v, but got %+v", want, res.Contests)
		}
	})

	t.Run("GetLanguage: empty", func(t *testing.T) {
		res, err := searcher.GetLanguage(ctx, &pb.GetLanguageRequest{})

		if err != nil {
			t.Error(err)
		}

		want := []*pb.Language{
			{
				Group:     "Python",
				Languages: []string{"Python (CPython 3.11.4)"},
			},
		}

		if !reflect.DeepEqual(want, res.Languages) {
			t.Errorf("expect %+v, but got %+v", want, res.Languages)
		}
	})

	t.Run("GetLanguage: with group", func(t *testing.T) {
		res, err := searcher.GetLanguage(ctx, &pb.GetLanguageRequest{
			Groups: []string{"Python"},
		})

		if err != nil {
			t.Error(err)
		}

		want := []*pb.Language{
			{
				Group:     "Python",
				Languages: []string{"Python (CPython 3.11.4)"},
			},
		}

		if !reflect.DeepEqual(want, res.Languages) {
			t.Errorf("expect %+v, but got %+v", want, res.Languages)
		}
	})

	t.Run("GetProblem: empty", func(t *testing.T) {
		res, err := searcher.GetProblem(ctx, &pb.GetProblemRequest{})

		if err != nil {
			t.Error(err)
		}

		want := []string{"abc300_a", "abc300_b", "arc184_a", "arc184_b"}

		if !reflect.DeepEqual(want, res.Problems) {
			t.Errorf("expect %+v, but got %+v", want, res.Problems)
		}
	})

	t.Run("GetProblem: category", func(t *testing.T) {
		res, err := searcher.GetProblem(ctx, &pb.GetProblemRequest{
			Categories: []string{"ABC"},
		})

		if err != nil {
			t.Error(err)
		}

		want := []string{"abc300_a", "abc300_b"}

		if !reflect.DeepEqual(want, res.Problems) {
			t.Errorf("expect %+v, but got %+v", want, res.Problems)
		}
	})

	t.Run("GetProblem: contest", func(t *testing.T) {
		res, err := searcher.GetProblem(ctx, &pb.GetProblemRequest{
			Contests: []string{"arc184"},
		})

		if err != nil {
			t.Error(err)
		}

		want := []string{"arc184_a", "arc184_b"}

		if !reflect.DeepEqual(want, res.Problems) {
			t.Errorf("expect %+v, but got %+v", want, res.Problems)
		}
	})
}
