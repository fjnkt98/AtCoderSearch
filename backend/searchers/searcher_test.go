package searchers

import (
	"context"
	"fjnkt98/atcodersearch/api"
	"fjnkt98/atcodersearch/internal/testutil"
	"fjnkt98/atcodersearch/repository"
	"fjnkt98/atcodersearch/update"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/meilisearch/meilisearch-go"
)

var ABC300A = api.Problem{
	ProblemId:      "abc300_a",
	ProblemTitle:   "A. N-choice question",
	ProblemUrl:     "https://atcoder.jp/contests/abc300/tasks/abc300_a",
	ContestId:      "abc300",
	ContestTitle:   "ユニークビジョンプログラミングコンテスト2023 春 (AtCoder Beginner Contest 300)",
	ContestUrl:     "https://atcoder.jp/contests/abc300",
	Difficulty:     api.NewOptInt(-1147),
	StartAt:        1682769600,
	Duration:       6000,
	RateChange:     " ~ 1999",
	Category:       "ABC",
	IsExperimental: false,
}

var ABC300B = api.Problem{
	ProblemId:      "abc300_b",
	ProblemTitle:   "B. Same Map in the RPG World",
	ProblemUrl:     "https://atcoder.jp/contests/abc300/tasks/abc300_b",
	ContestId:      "abc300",
	ContestTitle:   "ユニークビジョンプログラミングコンテスト2023 春 (AtCoder Beginner Contest 300)",
	ContestUrl:     "https://atcoder.jp/contests/abc300",
	Difficulty:     api.NewOptInt(350),
	StartAt:        1682769600,
	Duration:       6000,
	RateChange:     " ~ 1999",
	Category:       "ABC",
	IsExperimental: false,
}

var ARC184A = api.Problem{
	ProblemId:      "arc184_a",
	ProblemTitle:   "A. Appraiser",
	ProblemUrl:     "https://atcoder.jp/contests/arc184/tasks/arc184_a",
	ContestId:      "arc184",
	ContestTitle:   "AtCoder Regular Contest 184",
	ContestUrl:     "https://atcoder.jp/contests/arc184",
	Difficulty:     api.NewOptInt(1383),
	StartAt:        1727006400,
	Duration:       7200,
	RateChange:     "1200 ~ 2799",
	Category:       "ARC",
	IsExperimental: false,
}

var ARC184B = api.Problem{
	ProblemId:      "arc184_b",
	ProblemTitle:   "B. 123 Set",
	ProblemUrl:     "https://atcoder.jp/contests/arc184/tasks/arc184_b",
	ContestId:      "arc184",
	ContestTitle:   "AtCoder Regular Contest 184",
	ContestUrl:     "https://atcoder.jp/contests/arc184",
	Difficulty:     api.NewOptInt(2867),
	StartAt:        1727006400,
	Duration:       7200,
	RateChange:     "1200 ~ 2799",
	Category:       "ARC",
	IsExperimental: false,
}

var User1 = api.User{
	UserId:        "user1",
	Rating:        2563,
	HighestRating: 2563,
	Affiliation:   api.OptString{},
	BirthYear:     api.OptInt{},
	Country:       api.NewOptString("CN"),
	Crown:         api.NewOptString("user-orange-2"),
	JoinCount:     23,
	Rank:          470,
	ActiveRank:    api.OptInt{},
	Wins:          0,
	UserUrl:       "https://atcoder.jp/users/user1",
}

var User2 = api.User{
	UserId:        "user2",
	Rating:        3710,
	HighestRating: 3802,
	Affiliation:   api.OptString{},
	BirthYear:     api.OptInt{},
	Country:       api.NewOptString("CH"),
	Crown:         api.OptString{},
	JoinCount:     21,
	Rank:          3,
	ActiveRank:    api.OptInt{},
	Wins:          2,
	UserUrl:       "https://atcoder.jp/users/user2",
}

var User3 = api.User{
	UserId:        "user3",
	Rating:        3658,
	HighestRating: 3683,
	Affiliation:   api.NewOptString("MIT"),
	BirthYear:     api.NewOptInt(2001),
	Country:       api.NewOptString("US"),
	Crown:         api.OptString{},
	JoinCount:     48,
	Rank:          3,
	ActiveRank:    api.OptInt{},
	Wins:          0,
	UserUrl:       "https://atcoder.jp/users/user3",
}

var User4 = api.User{
	UserId:        "user4",
	Rating:        3604,
	HighestRating: 3814,
	Affiliation:   api.OptString{},
	BirthYear:     api.NewOptInt(1997),
	Country:       api.NewOptString("JP"),
	Crown:         api.NewOptString("crown_gold"),
	JoinCount:     38,
	Rank:          8,
	ActiveRank:    api.NewOptInt(5),
	Wins:          2,
	UserUrl:       "https://atcoder.jp/users/user4",
}

var Submission1 = api.Submission{
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
	Difficulty:    api.NewOptInt(-1147),
	UserId:        "fjnkt98",
	Language:      "Python (CPython 3.11.4)",
	LanguageGroup: "Python",
	Point:         100.0,
	Length:        1024,
	Result:        "AC",
	ExecutionTime: api.NewOptInt(22),
}

var Submission2 = api.Submission{
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
	Difficulty:    api.NewOptInt(350),
	UserId:        "fjnkt98",
	Language:      "Python (CPython 3.11.4)",
	LanguageGroup: "Python",
	Point:         200.0,
	Length:        1024,
	Result:        "WA",
	ExecutionTime: api.OptInt{},
}

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

	want := map[string]map[string]int{
		"category": {
			"ABC":             2237,
			"ABC-Like":        184,
			"AGC":             395,
			"AGC-Like":        77,
			"AHC":             38,
			"ARC":             810,
			"ARC-Like":        10,
			"JOI":             801,
			"Marathon":        22,
			"Other Contests":  1712,
			"Other Sponsored": 489,
			"PAST":            255,
		},
		"difficultyFacet": {
			"     ~    0": 3970,
			"   0 ~  400": 313,
			" 400 ~  800": 368,
			" 800 ~ 1200": 363,
			"1200 ~ 1600": 387,
			"1600 ~ 2000": 394,
			"2000 ~ 2400": 369,
			"2400 ~ 2800": 302,
			"2800 ~ 3200": 244,
			"3200 ~ 3600": 170,
			"3600 ~     ": 150,
		},
	}

	if !reflect.DeepEqual(want, res) {
		t.Errorf("expected %+v, but got %+v", want, res)
	}
}

func TestCreateSearchProblemRequest(t *testing.T) {
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
	searcher := NewSearcher(nil, pool)

	fields := []string{"problemId", "problemTitle", "problemUrl", "contestId", "contestTitle", "contestUrl", "difficulty", "color", "startAt", "duration", "rateChange", "category", "isExperimental"}

	cases := []struct {
		name string
		req  *api.APIProblemPostReq
		want *meilisearch.SearchRequest
	}{
		{name: "empty", req: &api.APIProblemPostReq{}, want: &meilisearch.SearchRequest{Sort: []string{"problemId:asc"}, HitsPerPage: 0, Page: 1, AttributesToRetrieve: fields}},
		{name: "pagination", req: &api.APIProblemPostReq{Limit: api.NewOptInt(20), Page: api.NewOptInt(2)}, want: &meilisearch.SearchRequest{Sort: []string{"problemId:asc"}, HitsPerPage: 20, Page: 2, AttributesToRetrieve: fields}},
		{name: "sort(valid)", req: &api.APIProblemPostReq{Sort: []api.APIProblemPostReqSortItem{"startAt:desc", "difficulty:asc"}}, want: &meilisearch.SearchRequest{Sort: []string{"startAt:desc", "difficulty:asc", "problemId:asc"}, Page: 1, AttributesToRetrieve: fields}},
		{name: "facet(valid)", req: &api.APIProblemPostReq{Facet: []api.APIProblemPostReqFacetItem{"category", "difficulty"}}, want: &meilisearch.SearchRequest{Facets: []string{"category", "difficultyFacet"}, Sort: []string{"problemId:asc"}, Page: 1, AttributesToRetrieve: fields}},
		{name: "filter by category", req: &api.APIProblemPostReq{Category: []string{"ABC", "ARC"}}, want: &meilisearch.SearchRequest{Filter: [][]string{{"category = 'ABC'", "category = 'ARC'"}}, Sort: []string{"problemId:asc"}, Page: 1, AttributesToRetrieve: fields}},
		{name: "filter by difficulty(from only)", req: &api.APIProblemPostReq{Difficulty: api.NewOptIntRange(api.IntRange{From: api.NewOptInt(800)})}, want: &meilisearch.SearchRequest{Filter: [][]string{{"difficulty >= 800"}}, Sort: []string{"problemId:asc"}, Page: 1, AttributesToRetrieve: fields}},
		{name: "filter by difficulty(to only)", req: &api.APIProblemPostReq{Difficulty: api.NewOptIntRange(api.IntRange{To: api.NewOptInt(1200)})}, want: &meilisearch.SearchRequest{Filter: [][]string{{"difficulty < 1200"}}, Sort: []string{"problemId:asc"}, Page: 1, AttributesToRetrieve: fields}},
		{name: "filter by difficulty(both)", req: &api.APIProblemPostReq{Difficulty: api.NewOptIntRange(api.IntRange{From: api.NewOptInt(800), To: api.NewOptInt(1200)})}, want: &meilisearch.SearchRequest{Filter: [][]string{{"difficulty >= 800"}, {"difficulty < 1200"}}, Sort: []string{"problemId:asc"}, Page: 1, AttributesToRetrieve: fields}},
		{name: "filter for experimental", req: &api.APIProblemPostReq{Experimental: api.NewOptBool(true)}, want: &meilisearch.SearchRequest{Filter: [][]string{{"isExperimental = true"}}, Sort: []string{"problemId:asc"}, Page: 1, AttributesToRetrieve: fields}},
		{name: "filter for not experimental", req: &api.APIProblemPostReq{Experimental: api.NewOptBool(false)}, want: &meilisearch.SearchRequest{Filter: [][]string{{"isExperimental = false"}}, Sort: []string{"problemId:asc"}, Page: 1, AttributesToRetrieve: fields}},
		{name: "filter for user id", req: &api.APIProblemPostReq{UserId: api.NewOptString("fjnkt98")}, want: &meilisearch.SearchRequest{Filter: [][]string{{"problemId NOT IN [abc300_a]"}}, Sort: []string{"problemId:asc"}, Page: 1, AttributesToRetrieve: fields}},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual, err := searcher.createSearchProblemQuery(ctx, tt.req)
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(tt.want, actual) {
				t.Errorf("expected %+v, but got %+v", tt.want, actual)
			}
		})
	}
}

func TestAPIProblem(t *testing.T) {
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

	searcher := NewSearcher(client, pool)

	t.Run("search problem empty", func(t *testing.T) {
		res, err := searcher.APIProblemPost(ctx, &api.APIProblemPostReq{})
		if err != nil {
			t.Error(err)
		}

		want := []api.Problem{
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
			req  *api.APIProblemPostReq
			want []api.Problem
		}{
			{name: "search problem with sort by start_at asc", req: &api.APIProblemPostReq{Sort: []api.APIProblemPostReqSortItem{"startAt:asc"}}, want: []api.Problem{ABC300A, ABC300B, ARC184A, ARC184B}},
			{name: "search problem with sort by start_at desc", req: &api.APIProblemPostReq{Sort: []api.APIProblemPostReqSortItem{"startAt:desc"}}, want: []api.Problem{ARC184A, ARC184B, ABC300A, ABC300B}},
			{name: "search problem with sort by difficulty asc", req: &api.APIProblemPostReq{Sort: []api.APIProblemPostReqSortItem{"difficulty:asc"}}, want: []api.Problem{ABC300A, ABC300B, ARC184A, ARC184B}},
			{name: "search problem with sort by difficulty desc", req: &api.APIProblemPostReq{Sort: []api.APIProblemPostReqSortItem{"difficulty:desc"}}, want: []api.Problem{ARC184B, ARC184A, ABC300B, ABC300A}},
		}

		for _, tt := range cases {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				res, err := searcher.APIProblemPost(ctx, tt.req)
				if err != nil {
					t.Error(err)
				}

				if !reflect.DeepEqual(tt.want, res.Items) {
					t.Errorf("expect %+v, but got %+v", tt.want, res.Items)
				}
			})
		}
	}

	t.Run("search problem with filter by category", func(t *testing.T) {
		res, err := searcher.APIProblemPost(ctx, &api.APIProblemPostReq{
			Category: []string{"ABC"},
		})

		if err != nil {
			t.Error(err)
		}

		want := []api.Problem{ABC300A, ABC300B}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("search problem filter by difficulty", func(t *testing.T) {
		res, err := searcher.APIProblemPost(ctx, &api.APIProblemPostReq{
			Difficulty: api.NewOptIntRange(api.IntRange{
				From: api.NewOptInt(1000),
				To:   api.NewOptInt(1400),
			}),
		})

		if err != nil {
			t.Error(err)
		}

		want := []api.Problem{ARC184A}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("search problem filter by user id", func(t *testing.T) {
		res, err := searcher.APIProblemPost(ctx, &api.APIProblemPostReq{
			UserId: api.NewOptString("fjnkt98"),
		})

		if err != nil {
			t.Error(err)
		}

		want := []api.Problem{ABC300B, ARC184A, ARC184B}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("search problem with keyword", func(t *testing.T) {
		res, err := searcher.APIProblemPost(ctx, &api.APIProblemPostReq{
			Q: api.NewOptString("ABC300"),
		})

		if err != nil {
			t.Error(err)
		}

		if len(res.Items) != 2 {
			t.Errorf("expect length of items is 2, but got %d", len(res.Items))
		}
	})

	t.Run("search problem with facet", func(t *testing.T) {
		res, err := searcher.APIProblemPost(ctx, &api.APIProblemPostReq{
			Facet: []api.APIProblemPostReqFacetItem{"category", "difficulty"},
		})

		if err != nil {
			t.Error(err)
		}

		want := api.APIProblemPostOKFacet{
			"category": []api.Count{
				{
					Label: "ABC",
					Count: 2,
				},
				{
					Label: "ARC",
					Count: 2,
				},
			},
			"difficulty": []api.Count{
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

	t.Run("get problem with empty", func(t *testing.T) {
		res, err := searcher.APIProblemGet(ctx, api.APIProblemGetParams{})

		if err != nil {
			t.Error(err)
		}

		want := []string{"abc300_a", "abc300_b", "arc184_a", "arc184_b"}

		if !reflect.DeepEqual(want, res.Problems) {
			t.Errorf("expect %+v, but got %+v", want, res.Problems)
		}
	})

	t.Run("get problem with category", func(t *testing.T) {
		res, err := searcher.APIProblemGet(ctx, api.APIProblemGetParams{
			Category: []string{"ABC"},
		})

		if err != nil {
			t.Error(err)
		}

		want := []string{"abc300_a", "abc300_b"}

		if !reflect.DeepEqual(want, res.Problems) {
			t.Errorf("expect %+v, but got %+v", want, res.Problems)
		}
	})

	t.Run("get problem with contest id", func(t *testing.T) {
		res, err := searcher.APIProblemGet(ctx, api.APIProblemGetParams{
			ContestId: []string{"arc184"},
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

func TestCreateSearchUserRequest(t *testing.T) {
	fields := []string{"userId", "rating", "highestRating", "affiliation", "birthYear", "country", "crown", "joinCount", "rank", "activeRank", "wins", "userUrl", "accepted", "submissionCount"}

	{
		cases := []struct {
			name string
			req  *api.APIUserPostReq
			want *meilisearch.SearchRequest
		}{
			{name: "empty", req: &api.APIUserPostReq{}, want: &meilisearch.SearchRequest{Sort: []string{"userId:asc"}, Page: 1, AttributesToRetrieve: fields}},
			{name: "pagination", req: &api.APIUserPostReq{Limit: api.NewOptInt(20)}, want: &meilisearch.SearchRequest{HitsPerPage: 20, Sort: []string{"userId:asc"}, Page: 1, AttributesToRetrieve: fields}},
			{name: "sort(valid)", req: &api.APIUserPostReq{Sort: []api.APIUserPostReqSortItem{"rating:desc", "birthYear:asc"}}, want: &meilisearch.SearchRequest{Sort: []string{"rating:desc", "birthYear:asc", "userId:asc"}, Page: 1, AttributesToRetrieve: fields}},
			{name: "facet(valid)", req: &api.APIUserPostReq{Facet: []api.APIUserPostReqFacetItem{"country", "rating", "birthYear", "joinCount"}}, want: &meilisearch.SearchRequest{Facets: []string{"country", "ratingFacet", "birthYearFacet", "joinCountFacet"}, Sort: []string{"userId:asc"}, Page: 1, AttributesToRetrieve: fields}},
			{name: "filter by user id", req: &api.APIUserPostReq{UserId: []string{"user1", "user2"}}, want: &meilisearch.SearchRequest{Filter: [][]string{{"userId = 'user1'", "userId = 'user2'"}}, Sort: []string{"userId:asc"}, Page: 1, AttributesToRetrieve: fields}},
			{name: "filter by rating", req: &api.APIUserPostReq{Rating: api.NewOptIntRange(api.IntRange{From: api.NewOptInt(800), To: api.NewOptInt(1200)})}, want: &meilisearch.SearchRequest{Filter: [][]string{{"rating >= 800"}, {"rating < 1200"}}, Sort: []string{"userId:asc"}, Page: 1, AttributesToRetrieve: fields}},
			{name: "filter by birth year", req: &api.APIUserPostReq{BirthYear: api.NewOptIntRange(api.IntRange{From: api.NewOptInt(1998), To: api.NewOptInt(2000)})}, want: &meilisearch.SearchRequest{Filter: [][]string{{"birthYear >= 1998"}, {"birthYear < 2000"}}, Sort: []string{"userId:asc"}, Page: 1, AttributesToRetrieve: fields}},
			{name: "filter by join count", req: &api.APIUserPostReq{JoinCount: api.NewOptIntRange(api.IntRange{From: api.NewOptInt(5), To: api.NewOptInt(10)})}, want: &meilisearch.SearchRequest{Filter: [][]string{{"joinCount >= 5"}, {"joinCount < 10"}}, Sort: []string{"userId:asc"}, Page: 1, AttributesToRetrieve: fields}},
			{name: "filter by country", req: &api.APIUserPostReq{Country: []string{"JP"}}, want: &meilisearch.SearchRequest{Filter: [][]string{{"country = 'JP'"}}, Sort: []string{"userId:asc"}, Page: 1, AttributesToRetrieve: fields}},
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

func TestAPIUser(t *testing.T) {
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
		update.NewUserRowReader(pool),
		update.NewUserIndexer(client),
		1000,
		1,
	); err != nil {
		t.Fatal(err)
	}

	searcher := NewSearcher(client, pool)

	t.Run("search user with empty", func(t *testing.T) {
		res, err := searcher.APIUserPost(ctx, &api.APIUserPostReq{})

		if err != nil {
			t.Error(err)
		}

		want := []api.User{User1, User2, User3, User4}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("search user with keyword", func(t *testing.T) {
		res, err := searcher.APIUserPost(ctx, &api.APIUserPostReq{
			Q: api.NewOptString("MIT"),
		})

		if err != nil {
			t.Error(err)
		}

		want := []api.User{User3}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	{
		cases := []struct {
			name string
			req  *api.APIUserPostReq
			want []api.User
		}{
			{name: "SearchUser: sort by rating asc", req: &api.APIUserPostReq{Sort: []api.APIUserPostReqSortItem{"rating:asc"}}, want: []api.User{User1, User4, User3, User2}},
			{name: "SearchUser: sort by rating desc", req: &api.APIUserPostReq{Sort: []api.APIUserPostReqSortItem{"rating:desc"}}, want: []api.User{User2, User3, User4, User1}},
			{name: "SearchUser: sort by birthYear asc", req: &api.APIUserPostReq{Sort: []api.APIUserPostReqSortItem{"birthYear:asc"}}, want: []api.User{User4, User3, User1, User2}},
			{name: "SearchUser: sort by birthYear desc", req: &api.APIUserPostReq{Sort: []api.APIUserPostReqSortItem{"birthYear:desc"}}, want: []api.User{User3, User4, User1, User2}},
		}

		for _, tt := range cases {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				res, err := searcher.APIUserPost(ctx, tt.req)

				if err != nil {
					t.Error(err)
				}

				if !reflect.DeepEqual(tt.want, res.Items) {
					t.Errorf("expect %+v, but got %+v", tt.want, res.Items)
				}
			})
		}
	}

	t.Run("search user with filter by user id", func(t *testing.T) {
		res, err := searcher.APIUserPost(ctx, &api.APIUserPostReq{
			UserId: []string{"user1", "user4"},
		})

		if err != nil {
			t.Error(err)
		}

		want := []api.User{User1, User4}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("search user with filter by rating", func(t *testing.T) {
		res, err := searcher.APIUserPost(ctx, &api.APIUserPostReq{
			Rating: api.NewOptIntRange(api.IntRange{From: api.NewOptInt(2563), To: api.NewOptInt(2564)}),
		})

		if err != nil {
			t.Error(err)
		}

		want := []api.User{User1}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("search user with filter by birth year", func(t *testing.T) {
		res, err := searcher.APIUserPost(ctx, &api.APIUserPostReq{
			BirthYear: api.NewOptIntRange(api.IntRange{From: api.NewOptInt(1997), To: api.NewOptInt(1998)}),
		})

		if err != nil {
			t.Error(err)
		}

		want := []api.User{User4}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("search user with filter by join count", func(t *testing.T) {
		res, err := searcher.APIUserPost(ctx, &api.APIUserPostReq{
			JoinCount: api.NewOptIntRange(api.IntRange{From: api.NewOptInt(45), To: api.NewOptInt(50)}),
		})

		if err != nil {
			t.Error(err)
		}

		want := []api.User{User3}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("search user with filter by country", func(t *testing.T) {
		res, err := searcher.APIUserPost(ctx, &api.APIUserPostReq{
			Country: []string{"JP"},
		})

		if err != nil {
			t.Error(err)
		}

		want := []api.User{User4}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("search user with facet", func(t *testing.T) {
		res, err := searcher.APIUserPost(ctx, &api.APIUserPostReq{
			Facet: []api.APIUserPostReqFacetItem{"country", "rating", "birthYear", "joinCount"},
		})

		if err != nil {
			t.Error(err)
		}

		want := api.APIUserPostOKFacet{
			"birthYear": []api.Count{
				{
					Label: "1990 ~ 2000",
					Count: 1,
				},
				{
					Label: "2000 ~ 2010",
					Count: 1,
				},
			},
			"country": []api.Count{
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
			"joinCount": []api.Count{
				{
					Label: "  20 ~   40",
					Count: 3,
				},
				{
					Label: "  40 ~   60",
					Count: 1,
				},
			},
			"rating": []api.Count{
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
}

func TestAPISubmission(t *testing.T) {
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

	searcher := NewSearcher(nil, pool)

	t.Run("search submission empty", func(t *testing.T) {
		res, err := searcher.APISubmissionPost(ctx, &api.APISubmissionPostReq{})
		if err != nil {
			t.Error(err)
		}

		want := []api.Submission{Submission2, Submission1}
		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	{
		cases := []struct {
			name string
			req  *api.APISubmissionPostReq
			want []api.Submission
		}{
			{name: "search submission with sort by execution time", req: &api.APISubmissionPostReq{Sort: []api.APISubmissionPostReqSortItem{"executionTime:asc"}}, want: []api.Submission{Submission1, Submission2}},
			{name: "search submission with sort by epoch second", req: &api.APISubmissionPostReq{Sort: []api.APISubmissionPostReqSortItem{"epochSecond:asc"}}, want: []api.Submission{Submission2, Submission1}},
			{name: "search submission with sort by point", req: &api.APISubmissionPostReq{Sort: []api.APISubmissionPostReqSortItem{"point:asc"}}, want: []api.Submission{Submission1, Submission2}},
			{name: "search submission with sort by length", req: &api.APISubmissionPostReq{Sort: []api.APISubmissionPostReqSortItem{"length:asc"}}, want: []api.Submission{Submission2, Submission1}},
		}

		for _, tt := range cases {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				res, err := searcher.APISubmissionPost(ctx, tt.req)
				if err != nil {
					t.Error(err)
				}

				if !reflect.DeepEqual(tt.want, res.Items) {
					t.Errorf("expect %+v, but got %+v", tt.want, res.Items)
				}
			})
		}
	}

	t.Run("search submission with filter by epoch second", func(t *testing.T) {
		res, err := searcher.APISubmissionPost(ctx, &api.APISubmissionPostReq{
			EpochSecond: api.NewOptIntRange(api.IntRange{
				From: api.NewOptInt(1729434073),
				To:   api.NewOptInt(1729434075),
			}),
		})

		if err != nil {
			t.Error(err)
		}

		want := []api.Submission{Submission2, Submission1}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("search submission filter by problem id", func(t *testing.T) {
		res, err := searcher.APISubmissionPost(ctx, &api.APISubmissionPostReq{
			ProblemId: []string{"abc300_a"},
		})

		if err != nil {
			t.Error(err)
		}

		want := []api.Submission{Submission1}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("search submission filter by contest id", func(t *testing.T) {
		res, err := searcher.APISubmissionPost(ctx, &api.APISubmissionPostReq{
			ContestId: []string{"abc300"},
		})

		if err != nil {
			t.Error(err)
		}

		want := []api.Submission{Submission2, Submission1}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("search submission filter by category", func(t *testing.T) {
		res, err := searcher.APISubmissionPost(ctx, &api.APISubmissionPostReq{
			Category: []string{"ABC"},
		})

		if err != nil {
			t.Error(err)
		}

		want := []api.Submission{Submission2, Submission1}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("search submission filter by user id", func(t *testing.T) {
		res, err := searcher.APISubmissionPost(ctx, &api.APISubmissionPostReq{
			UserId: []string{"fjnkt98"},
		})

		if err != nil {
			t.Error(err)
		}

		want := []api.Submission{Submission2, Submission1}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("search submission filter by language", func(t *testing.T) {
		res, err := searcher.APISubmissionPost(ctx, &api.APISubmissionPostReq{
			Language: []string{"Python (CPython 3.11.4)"},
		})

		if err != nil {
			t.Error(err)
		}

		want := []api.Submission{Submission2, Submission1}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("search submission filter by language group", func(t *testing.T) {
		res, err := searcher.APISubmissionPost(ctx, &api.APISubmissionPostReq{
			LanguageGroup: []string{"Python"},
		})

		if err != nil {
			t.Error(err)
		}

		want := []api.Submission{Submission2, Submission1}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("search submission filter by point", func(t *testing.T) {
		res, err := searcher.APISubmissionPost(ctx, &api.APISubmissionPostReq{
			Point: api.NewOptFloatRange(api.FloatRange{
				From: api.NewOptFloat64(100.0),
				To:   api.NewOptFloat64(150.0),
			}),
		})

		if err != nil {
			t.Error(err)
		}

		want := []api.Submission{Submission1}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("search submission filter by length", func(t *testing.T) {
		res, err := searcher.APISubmissionPost(ctx, &api.APISubmissionPostReq{
			Length: api.NewOptIntRange(api.IntRange{
				From: api.NewOptInt(1024),
				To:   api.NewOptInt(2048),
			}),
		})

		if err != nil {
			t.Error(err)
		}

		want := []api.Submission{Submission2, Submission1}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("search submission filter by result", func(t *testing.T) {
		res, err := searcher.APISubmissionPost(ctx, &api.APISubmissionPostReq{
			Result: []string{"WA"},
		})

		if err != nil {
			t.Error(err)
		}

		want := []api.Submission{Submission2}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("search submission filter by execution time", func(t *testing.T) {
		res, err := searcher.APISubmissionPost(ctx, &api.APISubmissionPostReq{
			ExecutionTime: api.NewOptIntRange(api.IntRange{
				From: api.NewOptInt(20),
				To:   api.NewOptInt(1000),
			}),
		})

		if err != nil {
			t.Error(err)
		}

		want := []api.Submission{Submission1}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})
}

func TestAPIGet(t *testing.T) {
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

	searcher := NewSearcher(nil, pool)

	t.Run("get category", func(t *testing.T) {
		res, err := searcher.APICategoryGet(ctx)

		if err != nil {
			t.Error(err)
		}

		want := []string{"ABC", "ARC"}

		if !reflect.DeepEqual(want, res.Categories) {
			t.Errorf("expect %+v, but got %+v", want, res.Categories)
		}
	})

	t.Run("get contest", func(t *testing.T) {
		res, err := searcher.APIContestGet(ctx, api.APIContestGetParams{})

		if err != nil {
			t.Error(err)
		}

		want := []string{"arc184", "abc300"}

		if !reflect.DeepEqual(want, res.Contests) {
			t.Errorf("expect %+v, but got %+v", want, res.Contests)
		}
	})

	t.Run("get contest with category", func(t *testing.T) {
		res, err := searcher.APIContestGet(ctx, api.APIContestGetParams{
			Category: []string{"ABC"},
		})

		if err != nil {
			t.Error(err)
		}

		want := []string{"abc300"}

		if !reflect.DeepEqual(want, res.Contests) {
			t.Errorf("expect %+v, but got %+v", want, res.Contests)
		}
	})

	t.Run("get language", func(t *testing.T) {
		res, err := searcher.APILanguageGet(ctx, api.APILanguageGetParams{})

		if err != nil {
			t.Error(err)
		}

		want := []api.Language{
			{
				Group:     "Python",
				Languages: []string{"Python (CPython 3.11.4)"},
			},
		}

		if !reflect.DeepEqual(want, res.Languages) {
			t.Errorf("expect %+v, but got %+v", want, res.Languages)
		}
	})

	t.Run("get language with group", func(t *testing.T) {
		res, err := searcher.APILanguageGet(ctx, api.APILanguageGetParams{
			Group: []string{"Python"},
		})

		if err != nil {
			t.Error(err)
		}

		want := []api.Language{
			{
				Group:     "Python",
				Languages: []string{"Python (CPython 3.11.4)"},
			},
		}

		if !reflect.DeepEqual(want, res.Languages) {
			t.Errorf("expect %+v, but got %+v", want, res.Languages)
		}
	})
}
