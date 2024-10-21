package searchers

import (
	"context"
	"errors"
	pb "fjnkt98/atcodersearch/grpc/atcodersearch/v1"
	"fmt"
	"log/slog"
	"maps"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/go-viper/mapstructure/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/meilisearch/meilisearch-go"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrInvalidRequest = errors.New("invalid request")

func FieldList(doc any) []string {
	ty := reflect.TypeOf(doc)
	if ty.Kind() != reflect.Pointer {
		return nil
	}

	ty = ty.Elem()
	if ty.Kind() != reflect.Struct {
		return nil
	}

	fl := make([]string, 0, ty.NumField())
	for i := 0; i < ty.NumField(); i++ {
		f := ty.Field(i)

		var name string
		if tag, ok := f.Tag.Lookup("mapstructure"); ok {
			if tag == "-" {
				continue
			}
			n, _, _ := strings.Cut(tag, ",")
			name = n
		} else {
			name = f.Name
		}
		fl = append(fl, name)
	}
	return fl
}

func ParseFacetDistribution(facetDistribution any) map[string]map[string]int64 {
	facet, ok := facetDistribution.(map[string]any)
	if !ok {
		return nil
	}

	result := make(map[string]map[string]int64)
	for field, v := range facet {
		counts, ok := v.(map[string]any)
		if !ok {
			continue
		}

		fieldCounts := make(map[string]int64)
		for key, count := range counts {
			if count, ok := count.(float64); ok {
				fieldCounts[key] = int64(count)
			}
		}
		result[field] = fieldCounts
	}

	return result
}

type Searcher struct {
	pb.UnimplementedSearchServiceServer

	client meilisearch.ServiceManager
	pool   *pgxpool.Pool
}

func NewSearcher(client meilisearch.ServiceManager, pool *pgxpool.Pool) *Searcher {
	return &Searcher{
		client: client,
		pool:   pool,
	}
}

type Problem struct {
	ProblemID      string  `mapstructure:"problemId" bun:"problem_id"`
	ProblemTitle   string  `mapstructure:"problemTitle" bun:"problem_title"`
	ProblemURL     string  `mapstructure:"problemUrl" bun:"problem_url"`
	ContestID      string  `mapstructure:"contestId" bun:"contest_id"`
	ContestTitle   string  `mapstructure:"contestTitle" bun:"contest_title"`
	ContestURL     string  `mapstructure:"contestUrl" bun:"contest_url"`
	Difficulty     *int64  `mapstructure:"difficulty" bun:"difficulty"`
	Color          *string `mapstructure:"color" bun:"color"`
	StartAt        int64   `mapstructure:"startAt" bun:"start_at"`
	Duration       int64   `mapstructure:"duration" bun:"duration"`
	RateChange     string  `mapstructure:"rateChange" bun:"rate_change"`
	Category       string  `mapstructure:"category" bun:"category"`
	IsExperimental bool    `mapstructure:"isExperimental" bun:"is_experimental"`
}

func (p *Problem) Into() *pb.Problem {
	return &pb.Problem{
		ProblemId:      p.ProblemID,
		ProblemTitle:   p.ProblemTitle,
		ProblemUrl:     p.ProblemURL,
		ContestId:      p.ContestID,
		ContestTitle:   p.ContestTitle,
		ContestUrl:     p.ContestURL,
		Difficulty:     p.Difficulty,
		StartAt:        p.StartAt,
		Duration:       p.Duration,
		RateChange:     p.RateChange,
		Category:       p.Category,
		IsExperimental: p.IsExperimental,
	}
}

func (s *Searcher) SearchProblem(ctx context.Context, req *pb.SearchProblemRequest) (*pb.SearchProblemResponse, error) {
	start := time.Now()

	q, err := createSearchProblemQuery(s.pool, req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err)
	}

	items := make([]*pb.Problem, 0, req.GetLimit())
	var rows []Problem
	if err := q.Scan(ctx, &rows); err != nil {
		return nil, status.Errorf(codes.Unknown, "scan rows: %s", err)
	}
	for _, r := range rows {
		items = append(items, r.Into())
	}

	return &pb.SearchProblemResponse{
		Time:  int64(time.Since(start) / time.Millisecond),
		Total: 0,
		Index: 0,
		Pages: 0,
		Items: items,
	}, nil
}

func createSearchProblemQuery(pool *pgxpool.Pool, req *pb.SearchProblemRequest) (*bun.SelectQuery, error) {
	db := bun.NewDB(stdlib.OpenDBFromPool(pool), pgdialect.New())

	limit := int(req.GetLimit())
	if limit > 200 {
		return nil, fmt.Errorf("%w: too large limitation", ErrInvalidRequest)
	}

	var offset int
	if page := int(req.GetPage()); page == 0 {
		offset = 0
	} else {
		offset = (page - 1) * limit
	}

	q := db.NewSelect().
		ColumnExpr("p.problem_id").
		ColumnExpr("p.title AS problem_title").
		ColumnExpr("p.url AS problem_url").
		ColumnExpr("c.contest_id").
		ColumnExpr("c.title AS contest_title").
		ColumnExpr("CONCAT('https://atcoder.jp/contests/', c.contest_id) AS contest_url").
		ColumnExpr("d.difficulty").
		ColumnExpr("c.start_epoch_second AS start_at").
		ColumnExpr("c.duration_second AS duration").
		ColumnExpr("c.rate_change AS rate_change").
		ColumnExpr("c.category AS category").
		ColumnExpr("COALESCE(d.is_experimental, FALSE) AS is_experimental").
		TableExpr("problems as p").
		Join("INNER JOIN contests as c ON p.contest_id = c.contest_id").
		Join("LEFT JOIN difficulties as d ON p.problem_id = d.problem_id").
		Limit(limit).
		Offset(offset)

	fields := map[string]string{
		"startAt":    "c.start_epoch_second",
		"difficulty": "d.difficulty",
		"problemId":  "p.problem_id",
		"contestId":  "c.contest_id",
	}

	sort := make([]string, 0, 4)
	if sorts := req.GetSorts(); len(sorts) > 0 {
		for _, s := range sorts {
			field, direction, ok := strings.Cut(s, ":")
			if !ok {
				return nil, fmt.Errorf("%w: sort direction needed", ErrInvalidRequest)
			}
			if !slices.Contains([]string{"asc", "desc"}, direction) {
				return nil, fmt.Errorf("%w: invalid sort direction `%s`", ErrInvalidRequest, direction)
			}

			column, ok := fields[field]
			if !ok {
				return nil, fmt.Errorf("%w: invalid sort field `%s`", ErrInvalidRequest, field)
			}

			sort = append(sort, fmt.Sprintf("%s %s", column, direction))
		}
	}
	sort = append(sort, "p.problem_id asc")
	q = q.Order(sort...)

	if categories := req.GetCategories(); len(categories) > 0 {
		q = q.Where("c.category IN (?)", bun.In(categories))
	}

	if difficulty := req.GetDifficulty(); difficulty != nil {
		if from := difficulty.From; from != nil {
			q = q.Where("d.difficulty >= ?", *from)
		}
		if to := difficulty.To; to != nil {
			q = q.Where("d.difficulty < ?", *to)
		}
	}

	if experimental := req.Experimental; experimental != nil {
		q = q.Where("COALESCE(d.is_experimental, FALSE) = ?", *experimental)
	}

	if userID := req.GetUserId(); userID != "" {
		sub := db.NewSelect().
			Distinct().
			ColumnExpr("problem_id").
			TableExpr("submissions").
			Where("user_id = ?", userID).
			Where("result = ?", "AC")

		q = q.With("s", sub).
			Join("INNER JOIN s ON p.problem_id = s.problem_id")
	}

	return q, nil
}

func (s *Searcher) SearchProblemByKeyword(ctx context.Context, req *pb.SearchProblemByKeywordRequest) (*pb.SearchProblemByKeywordResponse, error) {
	index := s.client.Index("problems")

	q, err := createSearchProblemByKeywordQuery(req)
	if err != nil {
		if errors.Is(err, ErrInvalidRequest) {
			return nil, status.Errorf(codes.InvalidArgument, "%s", err)
		} else {
			return nil, status.Errorf(codes.Unknown, "parse error: %s", err)
		}
	}

	res, err := index.SearchWithContext(ctx, req.GetQ(), q)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "search failed", slog.Any("error", err))
		return nil, status.Errorf(codes.Unknown, "search: %s", err)
	}

	items := make([]*pb.Problem, len(res.Hits))
	for i, hit := range res.Hits {
		item, ok := hit.(map[string]any)
		if !ok {
			return nil, status.Errorf(codes.Unknown, "item conversion: item of res.Hits isn't an map[string]any")
		}

		var problem Problem
		if err := mapstructure.Decode(item, &problem); err != nil {
			return nil, status.Errorf(codes.Unknown, "item conversion: %s", err)
		}

		items[i] = problem.Into()
	}

	categories := make([]*pb.Count, 0, 16)
	difficulties := make([]*pb.Count, 0, 16)
	for field, counts := range ParseFacetDistribution(res.FacetDistribution) {
		switch field {
		case "category":
			for _, k := range slices.Sorted(maps.Keys(counts)) {
				categories = append(categories, &pb.Count{Label: k, Count: counts[k]})
			}
		case "difficultyFacet":
			for _, k := range slices.Sorted(maps.Keys(counts)) {
				difficulties = append(difficulties, &pb.Count{Label: k, Count: counts[k]})
			}
		}
	}

	return &pb.SearchProblemByKeywordResponse{
		Time:  res.ProcessingTimeMs,
		Total: res.TotalHits,
		Index: res.Page,
		Pages: res.TotalPages,
		Items: items,
		Facet: &pb.ProblemFacet{
			Categories:   categories,
			Difficulties: difficulties,
		},
	}, nil
}

func createSearchProblemByKeywordQuery(req *pb.SearchProblemByKeywordRequest) (*meilisearch.SearchRequest, error) {
	q := &meilisearch.SearchRequest{
		AttributesToRetrieve: FieldList(new(Problem)),
	}

	limit := req.GetLimit()
	if limit > 200 {
		return nil, fmt.Errorf("%w: too large limitation", ErrInvalidRequest)
	}
	q.HitsPerPage = int64(limit)

	if page := req.GetPage(); page == 0 {
		q.Page = 1
	} else {
		q.Page = int64(page)
	}

	if sorts := req.GetSorts(); len(sorts) > 0 {
		sort := make([]string, 0, len(sorts))
		for _, s := range sorts {
			field, direction, ok := strings.Cut(s, ":")
			if !ok {
				return nil, fmt.Errorf("%w: sort direction needed", ErrInvalidRequest)
			}
			if !slices.Contains([]string{"asc", "desc"}, direction) {
				return nil, fmt.Errorf("%w: invalid sort direction `%s`", ErrInvalidRequest, direction)
			}
			if !slices.Contains([]string{"startAt", "difficulty", "problemId", "contestId"}, field) {
				return nil, fmt.Errorf("%w: invalid sort field `%s`", ErrInvalidRequest, field)
			}

			sort = append(sort, fmt.Sprintf("%s:%s", field, direction))
		}
		sort = append(sort, "problemId:asc")

		q.Sort = sort
	}

	if facets := req.GetFacets(); len(facets) > 0 {
		facet := make([]string, 0, len(facets))
		for _, f := range facets {
			if !slices.Contains([]string{"category", "difficulty"}, f) {
				return nil, fmt.Errorf("%w: invalid facet field `%s`", ErrInvalidRequest, f)
			}

			switch f {
			case "category":
				facet = append(facet, f)
			case "difficulty":
				facet = append(facet, "difficultyFacet")
			}
		}
		q.Facets = facet
	}

	filters := make([][]string, 0, 3)
	if categories := req.GetCategories(); len(categories) != 0 {
		categoryFilter := make([]string, 0, len(categories))
		for _, c := range categories {
			categoryFilter = append(categoryFilter, fmt.Sprintf("category = '%s'", c))
		}
		filters = append(filters, categoryFilter)
	}
	if difficulty := req.GetDifficulty(); difficulty != nil {
		if from := difficulty.From; from != nil {
			filters = append(filters, []string{fmt.Sprintf("difficulty >= %d", *from)})
		}
		if to := difficulty.To; to != nil {
			filters = append(filters, []string{fmt.Sprintf("difficulty < %d", *to)})
		}
	}
	if experimental := req.Experimental; experimental != nil {
		filters = append(filters, []string{fmt.Sprintf("isExperimental = %t", *experimental)})
	}

	if len(filters) > 0 {
		q.Filter = filters
	}

	return q, nil
}

type User struct {
	UserID        string  `mapstructure:"userId" bun:"user_id"`
	Rating        int64   `mapstructure:"rating" bun:"rating"`
	HighestRating int64   `mapstructure:"highestRating" bun:"highest_rating"`
	Affiliation   *string `mapstructure:"affiliation" bun:"affiliation"`
	BirthYear     *int64  `mapstructure:"birthYear" bun:"birth_year"`
	Country       *string `mapstructure:"country" bun:"country"`
	Crown         *string `mapstructure:"crown" bun:"crown"`
	JoinCount     int64   `mapstructure:"joinCount" bun:"join_count"`
	Rank          int64   `mapstructure:"rank" bun:"rank"`
	ActiveRank    *int64  `mapstructure:"activeRank" bun:"active_rank"`
	Wins          int64   `mapstructure:"wins" bun:"wins"`
	UserURL       string  `mapstructure:"userUrl" bun:"user_url"`
}

func (u *User) Into() *pb.User {
	return &pb.User{
		UserId:        u.UserID,
		Rating:        u.Rating,
		HighestRating: u.HighestRating,
		Affiliation:   u.Affiliation,
		BirthYear:     u.BirthYear,
		Country:       u.Country,
		Crown:         u.Crown,
		JoinCount:     u.JoinCount,
		Rank:          u.Rank,
		ActiveRank:    u.ActiveRank,
		Wins:          u.Wins,
		UserUrl:       u.UserURL,
	}
}

func (s *Searcher) SearchUser(ctx context.Context, req *pb.SearchUserRequest) (*pb.SearchUserResponse, error) {
	index := s.client.Index("users")

	q, err := createSearchUserQuery(req)
	if err != nil {
		if errors.Is(err, ErrInvalidRequest) {
			return nil, status.Errorf(codes.InvalidArgument, "%s", err)
		} else {
			return nil, status.Errorf(codes.Unknown, "parse error: %s", err)
		}
	}

	res, err := index.SearchWithContext(ctx, req.GetQ(), q)
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "search failed", slog.Any("error", err))
		return nil, status.Errorf(codes.Unknown, "search: %s", err)
	}

	items := make([]*pb.User, len(res.Hits))
	for i, hit := range res.Hits {
		item, ok := hit.(map[string]any)
		if !ok {
			return nil, status.Errorf(codes.Unknown, "item conversion: item of res.Hits isn't an map[string]any")
		}

		var user User
		if err := mapstructure.Decode(item, &user); err != nil {
			return nil, status.Errorf(codes.Unknown, "item conversion: %s", err)
		}

		items[i] = user.Into()
	}

	countries := make([]*pb.Count, 0, 16)
	ratings := make([]*pb.Count, 0, 16)
	birthYears := make([]*pb.Count, 0, 16)
	joinCounts := make([]*pb.Count, 0, 16)
	for field, counts := range ParseFacetDistribution(res.FacetDistribution) {
		switch field {
		case "countries":
			for _, k := range slices.Sorted(maps.Keys(counts)) {
				countries = append(countries, &pb.Count{Label: k, Count: counts[k]})
			}
		case "ratingFacet":
			for _, k := range slices.Sorted(maps.Keys(counts)) {
				ratings = append(ratings, &pb.Count{Label: k, Count: counts[k]})
			}
		case "birthYearFacet":
			for _, k := range slices.Sorted(maps.Keys(counts)) {
				birthYears = append(birthYears, &pb.Count{Label: k, Count: counts[k]})
			}
		case "joinCountFacet":
			for _, k := range slices.Sorted(maps.Keys(counts)) {
				joinCounts = append(joinCounts, &pb.Count{Label: k, Count: counts[k]})
			}
		}
	}

	return &pb.SearchUserResponse{
		Time:  res.ProcessingTimeMs,
		Total: res.TotalHits,
		Index: res.Page,
		Pages: res.TotalPages,
		Items: items,
		Facet: &pb.UserFacet{
			Countries:  countries,
			Ratings:    ratings,
			BirthYears: birthYears,
			JoinCounts: joinCounts,
		},
	}, nil
}

func createSearchUserQuery(req *pb.SearchUserRequest) (*meilisearch.SearchRequest, error) {
	q := &meilisearch.SearchRequest{
		AttributesToRetrieve: FieldList(new(User)),
	}

	limit := int(req.GetLimit())
	if limit > 200 {
		return nil, fmt.Errorf("%w: too large limitation", ErrInvalidRequest)
	}
	q.HitsPerPage = int64(limit)

	if page := req.GetPage(); page == 0 {
		q.Page = 1
	} else {
		q.Page = int64(page)
	}

	if sorts := req.GetSorts(); len(sorts) > 0 {
		sort := make([]string, 0, len(sorts))
		for _, s := range sorts {
			field, direction, ok := strings.Cut(s, ":")
			if !ok {
				return nil, fmt.Errorf("%w: sort direction needed", ErrInvalidRequest)
			}
			if !slices.Contains([]string{"asc", "desc"}, direction) {
				return nil, fmt.Errorf("%w: invalid sort direction `%s`", ErrInvalidRequest, direction)
			}
			if !slices.Contains([]string{"rating", "birthYear", "userId"}, field) {
				return nil, fmt.Errorf("%w: invalid sort field `%s`", ErrInvalidRequest, field)
			}

			sort = append(sort, fmt.Sprintf("%s:%s", field, direction))
		}
		sort = append(sort, "userId:asc")

		q.Sort = sort
	}

	if facets := req.GetFacets(); len(facets) > 0 {
		facet := make([]string, 0, len(facets))
		for _, f := range facets {
			if !slices.Contains([]string{"country", "rating", "birthYear", "joinCount"}, f) {
				return nil, fmt.Errorf("%w: invalid facet field `%s`", ErrInvalidRequest, f)
			}

			switch f {
			case "country":
				facet = append(facet, f)
			case "rating", "birthYear", "joinCount":
				facet = append(facet, fmt.Sprintf("%sFacet", f))
			}
		}
		q.Facets = facet
	}

	filters := make([][]string, 0, 5)
	if userIDs := req.GetUserIds(); len(userIDs) > 0 {
		userIDFilter := make([]string, 0, len(userIDs))
		for _, u := range userIDs {
			userIDFilter = append(userIDFilter, fmt.Sprintf("userId = '%s'", u))
		}
		filters = append(filters, userIDFilter)
	}
	if rating := req.GetRating(); rating != nil {
		if from := rating.From; from != nil {
			filters = append(filters, []string{fmt.Sprintf("rating >= %d", *from)})
		}
		if to := rating.To; to != nil {
			filters = append(filters, []string{fmt.Sprintf("rating < %d", *to)})
		}
	}
	if birthYear := req.GetBirthYear(); birthYear != nil {
		if from := birthYear.From; from != nil {
			filters = append(filters, []string{fmt.Sprintf("birthYear >= %d", *from)})
		}
		if to := birthYear.To; to != nil {
			filters = append(filters, []string{fmt.Sprintf("birthYear < %d", *to)})
		}
	}
	if joinCount := req.GetJoinCount(); joinCount != nil {
		if from := joinCount.From; from != nil {
			filters = append(filters, []string{fmt.Sprintf("joinCount >= %d", *from)})
		}
		if to := joinCount.To; to != nil {
			filters = append(filters, []string{fmt.Sprintf("joinCount < %d", *to)})
		}
	}
	if countries := req.GetCountries(); len(countries) > 0 {
		countriesFilter := make([]string, 0, len(countries))
		for _, u := range countries {
			countriesFilter = append(countriesFilter, fmt.Sprintf("country = '%s'", u))
		}
		filters = append(filters, countriesFilter)
	}

	if len(filters) > 0 {
		q.Filter = filters
	}

	return q, nil
}

func (s *Searcher) SearchSubmission(ctx context.Context, req *pb.SearchSubmissionRequest) (*pb.SearchSubmissionResponse, error) {
	panic("")
}

func (s *Searcher) GetCategory(ctx context.Context, req *pb.GetCategoryRequest) (*pb.GetCategoryResponse, error) {
	panic("")
}

func (s *Searcher) GetContest(ctx context.Context, req *pb.GetContestRequest) (*pb.GetContestResponse, error) {
	panic("")
}

func (s *Searcher) GetLanguage(ctx context.Context, req *pb.GetLanguageRequest) (*pb.GetLanguageResponse, error) {
	panic("")
}

func (s *Searcher) GetProblem(ctx context.Context, req *pb.GetProblemRequest) (*pb.GetProblemResponse, error) {
	panic("")
}
