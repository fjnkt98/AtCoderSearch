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
			{name: "empty", req: &pb.SearchProblemByKeywordRequest{}, want: &meilisearch.SearchRequest{HitsPerPage: 0, Page: 1, AttributesToRetrieve: fields}},
			{name: "pagination", req: &pb.SearchProblemByKeywordRequest{Limit: ptr.To[uint64](20), Page: ptr.To[uint64](2)}, want: &meilisearch.SearchRequest{HitsPerPage: 20, Page: 2, AttributesToRetrieve: fields}},
			{name: "sort(valid)", req: &pb.SearchProblemByKeywordRequest{Sorts: []string{"startAt:desc", "difficulty:asc"}}, want: &meilisearch.SearchRequest{Sort: []string{"startAt:desc", "difficulty:asc", "problemId:asc"}, Page: 1, AttributesToRetrieve: fields}},
			{name: "facet(valid)", req: &pb.SearchProblemByKeywordRequest{Facets: []string{"category", "difficulty"}}, want: &meilisearch.SearchRequest{Facets: []string{"category", "difficultyFacet"}, Page: 1, AttributesToRetrieve: fields}},
			{name: "filter by category", req: &pb.SearchProblemByKeywordRequest{Categories: []string{"ABC", "ARC"}}, want: &meilisearch.SearchRequest{Filter: [][]string{{"category = 'ABC'", "category = 'ARC'"}}, Page: 1, AttributesToRetrieve: fields}},
			{name: "filter by difficulty(from only)", req: &pb.SearchProblemByKeywordRequest{Difficulty: &pb.IntRange{From: ptr.To[int64](800)}}, want: &meilisearch.SearchRequest{Filter: [][]string{{"difficulty >= 800"}}, Page: 1, AttributesToRetrieve: fields}},
			{name: "filter by difficulty(to only)", req: &pb.SearchProblemByKeywordRequest{Difficulty: &pb.IntRange{To: ptr.To[int64](1200)}}, want: &meilisearch.SearchRequest{Filter: [][]string{{"difficulty < 1200"}}, Page: 1, AttributesToRetrieve: fields}},
			{name: "filter by difficulty(both)", req: &pb.SearchProblemByKeywordRequest{Difficulty: &pb.IntRange{From: ptr.To[int64](800), To: ptr.To[int64](1200)}}, want: &meilisearch.SearchRequest{Filter: [][]string{{"difficulty >= 800"}, {"difficulty < 1200"}}, Page: 1, AttributesToRetrieve: fields}},
			{name: "filter for experimental", req: &pb.SearchProblemByKeywordRequest{Experimental: ptr.To(true)}, want: &meilisearch.SearchRequest{Filter: [][]string{{"isExperimental = true"}}, Page: 1, AttributesToRetrieve: fields}},
			{name: "filter for not experimental", req: &pb.SearchProblemByKeywordRequest{Experimental: ptr.To(false)}, want: &meilisearch.SearchRequest{Filter: [][]string{{"isExperimental = false"}}, Page: 1, AttributesToRetrieve: fields}},
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

func TestSearchProblemByKeyword(t *testing.T) {
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

	t.Run("keyword", func(t *testing.T) {
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

	t.Run("filter by category", func(t *testing.T) {
		res, err := searcher.SearchProblemByKeyword(ctx, &pb.SearchProblemByKeywordRequest{
			Categories: []string{"ARC"},
			Sorts:      []string{"startAt:desc"},
		})

		if err != nil {
			t.Error(err)
		}

		want := []*pb.Problem{
			{
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
			},
			{
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
			},
		}

		if !reflect.DeepEqual(want, res.Items) {
			t.Errorf("expect %+v, but got %+v", want, res.Items)
		}
	})

	t.Run("facet", func(t *testing.T) {
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
}
