package searchers

import (
	"context"
	"encoding/json"
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
	"github.com/jackc/pgx/v5"
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

type Into[T any] interface {
	Into() T
}

func ScanItems[S any, T Into[S]](hits []any) ([]S, error) {
	items := make([]S, len(hits))
	for i, hit := range hits {
		hit, ok := hit.(map[string]any)
		if !ok {
			return nil, status.Errorf(codes.Unknown, "item conversion: item of res.Hits isn't an map[string]any")
		}

		var item T
		if err := mapstructure.Decode(hit, &item); err != nil {
			return nil, status.Errorf(codes.Unknown, "item conversion: %s", err)
		}

		items[i] = item.Into()
	}

	return items, nil
}

func GenSort(sorts []string, sep rune, allowed []string, defaultKey ...string) ([]string, error) {
	sort := make([]string, 0, len(sorts)+len(defaultKey))

	for _, s := range sorts {
		field, direction, ok := strings.Cut(s, ":")
		if !ok {
			return nil, fmt.Errorf("%w: sort direction needed", ErrInvalidRequest)
		}
		if !slices.Contains([]string{"asc", "desc"}, direction) {
			return nil, fmt.Errorf("%w: invalid sort direction `%s`", ErrInvalidRequest, direction)
		}
		if len(allowed) > 0 && !slices.Contains(allowed, field) {
			return nil, fmt.Errorf("%w: invalid sort field `%s`", ErrInvalidRequest, field)
		}

		sort = append(sort, fmt.Sprintf("%s%c%s", field, sep, direction))
	}
	sort = append(sort, defaultKey...)

	return sort, nil
}

func GenFacet(facets []string, mapping map[string]string) ([]string, error) {
	if len(facets) == 0 {
		return nil, nil
	}

	facet := make([]string, 0, len(facets))
	for _, f := range facets {
		if field, ok := mapping[f]; ok {
			facet = append(facet, field)
		} else {
			return nil, fmt.Errorf("%w: invalid facet field `%s`", ErrInvalidRequest, f)
		}
	}

	return facet, nil
}

func StringSliceFilter(field string, values []string) []string {
	if len(values) == 0 {
		return nil
	}

	result := make([]string, 0, len(values))
	for _, v := range values {
		result = append(result, fmt.Sprintf("%s = '%s'", field, v))
	}

	return result
}

func IntRangeFilter(field string, value *pb.IntRange) [][]string {
	if value == nil {
		return nil
	}

	result := make([][]string, 0, 2)
	if from := value.From; from != nil {
		result = append(result, []string{fmt.Sprintf("%s >= %d", field, *from)})
	}
	if to := value.To; to != nil {
		result = append(result, []string{fmt.Sprintf("%s < %d", field, *to)})
	}

	if len(result) == 0 {
		return nil
	}

	return result
}

func NullableBoolFilter(field string, value *bool) []string {
	if value == nil {
		return nil
	}

	return []string{fmt.Sprintf("%s = %t", field, *value)}
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

func (p Problem) Into() *pb.Problem {
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
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Errorf(codes.Unknown, "scan rows: %s", err)
		}
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
	if limit > 1000 {
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

	items, err := ScanItems[*pb.Problem, Problem](res.Hits)
	if err != nil {
		return nil, err
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
	if limit > 1000 {
		return nil, fmt.Errorf("%w: too large limitation", ErrInvalidRequest)
	}
	q.HitsPerPage = int64(limit)

	if page := req.GetPage(); page == 0 {
		q.Page = 1
	} else {
		q.Page = int64(page)
	}

	if sort, err := GenSort(req.GetSorts(), ':', []string{"startAt", "difficulty", "problemId", "contestId"}, "problemId:asc"); err != nil {
		return nil, err
	} else {
		q.Sort = sort
	}

	if facet, err := GenFacet(req.GetFacets(), map[string]string{"category": "category", "difficulty": "difficultyFacet"}); err != nil {
		return nil, err
	} else {
		q.Facets = facet
	}

	filters := make([][]string, 0, 3)
	if categories := StringSliceFilter("category", req.GetCategories()); categories != nil {
		filters = append(filters, categories)
	}
	if difficulty := IntRangeFilter("difficulty", req.GetDifficulty()); difficulty != nil {
		filters = append(filters, difficulty...)
	}
	if experimental := NullableBoolFilter("isExperimental", req.Experimental); experimental != nil {
		filters = append(filters, experimental)
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

func (u User) Into() *pb.User {
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

	items, err := ScanItems[*pb.User, User](res.Hits)
	if err != nil {
		return nil, err
	}

	countries := make([]*pb.Count, 0, 16)
	ratings := make([]*pb.Count, 0, 16)
	birthYears := make([]*pb.Count, 0, 16)
	joinCounts := make([]*pb.Count, 0, 16)
	for field, counts := range ParseFacetDistribution(res.FacetDistribution) {
		switch field {
		case "country":
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
	if limit > 1000 {
		return nil, fmt.Errorf("%w: too large limitation", ErrInvalidRequest)
	}
	q.HitsPerPage = int64(limit)

	if page := req.GetPage(); page == 0 {
		q.Page = 1
	} else {
		q.Page = int64(page)
	}

	if sort, err := GenSort(req.GetSorts(), ':', []string{"rating", "birthYear", "userId"}, "userId:asc"); err != nil {
		return nil, err
	} else {
		q.Sort = sort
	}

	if facet, err := GenFacet(req.GetFacets(), map[string]string{
		"country":   "country",
		"rating":    "ratingFacet",
		"birthYear": "birthYearFacet",
		"joinCount": "joinCountFacet",
	}); err != nil {
		return nil, err
	} else {
		q.Facets = facet
	}

	filters := make([][]string, 0, 5)
	if users := StringSliceFilter("userId", req.GetUserIds()); users != nil {
		filters = append(filters, users)
	}
	if rating := IntRangeFilter("rating", req.GetRating()); rating != nil {
		filters = append(filters, rating...)
	}
	if birthYear := IntRangeFilter("birthYear", req.GetBirthYear()); birthYear != nil {
		filters = append(filters, birthYear...)
	}
	if joinCount := IntRangeFilter("joinCount", req.GetJoinCount()); joinCount != nil {
		filters = append(filters, joinCount...)
	}
	if countries := StringSliceFilter("country", req.GetCountries()); countries != nil {
		filters = append(filters, countries)
	}

	if len(filters) > 0 {
		q.Filter = filters
	}

	return q, nil
}

type Submission struct {
	SubmissionID  int64   `bun:"submission_id"`
	SubmittedAt   int64   `bun:"submitted_at"`
	SubmissionURL string  `bun:"submission_url"`
	ProblemID     string  `bun:"problem_id"`
	ProblemTitle  string  `bun:"problem_title"`
	ProblemURL    string  `bun:"problem_url"`
	ContestID     string  `bun:"contest_id"`
	ContestTitle  string  `bun:"contest_title"`
	ContestURL    string  `bun:"contest_url"`
	Category      string  `bun:"category"`
	Difficulty    *int64  `bun:"difficulty"`
	UserID        string  `bun:"user_id"`
	Language      string  `bun:"language"`
	LanguageGroup string  `bun:"language_group"`
	Point         float64 `bun:"point"`
	Length        int64   `bun:"length"`
	Result        string  `bun:"result"`
	ExecutionTime *int64  `bun:"execution_time"`
}

func (s Submission) Into() *pb.Submission {
	return &pb.Submission{
		SubmissionId:  s.SubmissionID,
		SubmittedAt:   s.SubmittedAt,
		SubmissionUrl: s.SubmissionURL,
		ProblemId:     s.ProblemID,
		ProblemTitle:  s.ProblemTitle,
		ProblemUrl:    s.ProblemURL,
		ContestId:     s.ContestID,
		ContestTitle:  s.ContestTitle,
		ContestUrl:    s.ContestURL,
		Category:      s.Category,
		Difficulty:    s.Difficulty,
		UserId:        s.UserID,
		Language:      s.Language,
		LanguageGroup: s.LanguageGroup,
		Point:         s.Point,
		Length:        s.Length,
		Result:        s.Result,
		ExecutionTime: s.ExecutionTime,
	}
}

func (s *Searcher) SearchSubmission(ctx context.Context, req *pb.SearchSubmissionRequest) (*pb.SearchSubmissionResponse, error) {
	start := time.Now()

	q, err := createSearchSubmissionQuery(s.pool, req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "%s", err)
	}

	items := make([]*pb.Submission, 0, req.GetLimit())
	var rows []Submission
	if err := q.Scan(ctx, &rows); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Errorf(codes.Unknown, "scan rows: %s", err)
		}
	}
	for _, r := range rows {
		items = append(items, r.Into())
	}

	return &pb.SearchSubmissionResponse{
		Time:  int64(time.Since(start) / time.Millisecond),
		Index: 0,
		Items: items,
	}, nil
}

func createSearchSubmissionQuery(pool *pgxpool.Pool, req *pb.SearchSubmissionRequest) (*bun.SelectQuery, error) {
	db := bun.NewDB(stdlib.OpenDBFromPool(pool), pgdialect.New())

	limit := int(req.GetLimit())
	if limit > 1000 {
		return nil, fmt.Errorf("%w: too large limitation", ErrInvalidRequest)
	}

	var offset int
	if page := int(req.GetPage()); page == 0 {
		offset = 0
	} else {
		offset = (page - 1) * limit
	}

	q := db.NewSelect().
		ColumnExpr("s.id AS submission_id").
		ColumnExpr("s.epoch_second AS submitted_at").
		ColumnExpr("FORMAT('https://atcoder.jp/contests/%s/submissions/%s', s.contest_id, s.id) AS submission_url").
		ColumnExpr("s.problem_id").
		ColumnExpr("p.title AS problem_title").
		ColumnExpr("p.url AS problem_url").
		ColumnExpr("s.contest_id").
		ColumnExpr("c.title AS contest_title").
		ColumnExpr("CONCAT('https://atcoder.jp/contests/', c.contest_id) AS contest_url").
		ColumnExpr("c.category").
		ColumnExpr("d.difficulty").
		ColumnExpr("s.user_id").
		ColumnExpr("s.language").
		ColumnExpr("l.group AS language_group").
		ColumnExpr("s.point").
		ColumnExpr("s.length").
		ColumnExpr("s.result").
		ColumnExpr("s.execution_time").
		TableExpr("submissions AS s").
		Join("LEFT JOIN contests AS c ON s.contest_id = c.contest_id").
		Join("LEFT JOIN problems AS p ON s.problem_id = p.problem_id").
		Join("LEFT JOIN difficulties AS d ON s.problem_id = d.problem_id").
		Join("LEFT JOIN languages AS l ON s.language = l.language").
		Limit(limit).
		Offset(offset)

	fields := map[string]string{
		"executionTime": "execution_time",
		"epochSecond":   "epoch_second",
		"point":         "point",
		"length":        "length",
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
	sort = append(sort, "id desc")
	q = q.Order(sort...)

	if epochSecond := req.GetEpochSecond(); epochSecond != nil {
		if from := epochSecond.From; from != nil {
			q = q.Where("s.epoch_second >= ?", *from)
		}
		if to := epochSecond.To; to != nil {
			q = q.Where("s.epoch_second < ?", *to)
		}
	}
	if problems := req.GetProblemIds(); len(problems) > 0 {
		q = q.Where("s.problem_id IN (?)", bun.In(problems))
	}
	if contests := req.GetContestIds(); len(contests) > 0 {
		q = q.Where("s.contest_id IN (?)", bun.In(contests))
	}
	if categories := req.GetCategories(); len(categories) > 0 {
		q = q.Where("c.category IN (?)", bun.In(categories))
	}
	if users := req.GetUserIds(); len(users) > 0 {
		q = q.Where("s.user_id IN (?)", bun.In(users))
	}
	if languages := req.GetLanguages(); len(languages) > 0 {
		q = q.Where("s.language IN (?)", bun.In(languages))
	}
	if groups := req.GetLanguageGroups(); len(groups) > 0 {
		q = q.Where("l.group IN (?)", bun.In(groups))
	}
	if point := req.GetPoint(); point != nil {
		if from := point.From; from != nil {
			q = q.Where("s.point >= ?", *from)
		}
		if to := point.To; to != nil {
			q = q.Where("s.point < ?", *to)
		}
	}
	if length := req.GetLength(); length != nil {
		if from := length.From; from != nil {
			q = q.Where("s.length >= ?", *from)
		}
		if to := length.To; to != nil {
			q = q.Where("s.length < ?", *to)
		}
	}
	if results := req.GetResults(); len(results) > 0 {
		q = q.Where("s.result IN (?)", bun.In(results))
	}
	if exec := req.GetExecutionTime(); exec != nil {
		if from := exec.From; from != nil {
			q = q.Where("s.execution_time >= ?", *from)
		}
		if to := exec.To; to != nil {
			q = q.Where("s.execution_time < ?", *to)
		}
	}

	return q, nil
}

func (s *Searcher) GetCategory(ctx context.Context, req *pb.GetCategoryRequest) (*pb.GetCategoryResponse, error) {
	db := bun.NewDB(stdlib.OpenDBFromPool(s.pool), pgdialect.New())

	q := db.NewSelect().Distinct().
		Column("category").
		Table("contests").
		Order("category ASC")

	var categories []string
	if err := q.Scan(ctx, &categories); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Errorf(codes.Unknown, "scan rows: %s", err)
		}
	}

	return &pb.GetCategoryResponse{
		Categories: categories,
	}, nil
}

func (s *Searcher) GetContest(ctx context.Context, req *pb.GetContestRequest) (*pb.GetContestResponse, error) {
	db := bun.NewDB(stdlib.OpenDBFromPool(s.pool), pgdialect.New())

	q := db.NewSelect().
		Column("contest_id").
		Table("contests").
		Order("start_epoch_second DESC")

	if c := req.GetCategories(); len(c) > 0 {
		q = q.Where("category IN (?)", bun.In(c))
	}

	var contests []string
	if err := q.Scan(ctx, &contests); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Errorf(codes.Unknown, "scan rows: %s", err)
		}
	}

	return &pb.GetContestResponse{
		Contests: contests,
	}, nil
}

func (s *Searcher) GetLanguage(ctx context.Context, req *pb.GetLanguageRequest) (*pb.GetLanguageResponse, error) {
	db := bun.NewDB(stdlib.OpenDBFromPool(s.pool), pgdialect.New())

	cte := db.NewSelect().
		ColumnExpr("?", bun.Ident("group")).
		ColumnExpr("ARRAY_AGG(language ORDER BY language ASC) AS languages").
		Table("languages").
		GroupExpr("?", bun.Ident("group")).
		OrderExpr("? ASC", bun.Ident("group"))

	if g := req.GetGroups(); len(g) > 0 {
		cte = cte.Where("? IN (?)", bun.Ident("group"), bun.In(g))
	}

	q := db.NewSelect().
		With("l", cte).
		ColumnExpr("JSON_AGG(l)").
		Table("l")

	var body []byte
	if err := q.Scan(ctx, &body); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &pb.GetLanguageResponse{}, nil
		} else {
			return nil, status.Errorf(codes.Unknown, "scan rows: %s", err)
		}
	}

	languages := make([]*pb.Language, 0)
	if len(body) > 0 {
		if err := json.Unmarshal(body, &languages); err != nil {
			return nil, status.Errorf(codes.Unknown, "unmarshal body: %s", err)
		}
	}

	return &pb.GetLanguageResponse{
		Languages: languages,
	}, nil
}

func (s *Searcher) GetProblem(ctx context.Context, req *pb.GetProblemRequest) (*pb.GetProblemResponse, error) {
	db := bun.NewDB(stdlib.OpenDBFromPool(s.pool), pgdialect.New())

	q := db.NewSelect().
		Column("problem_id").
		Table("problems").
		Join("LEFT JOIN contests USING(contest_id)").
		Order("problem_id ASC")

	if c := req.GetCategories(); len(c) > 0 {
		q = q.Where("category IN (?)", bun.In(c))
	}

	if c := req.GetContests(); len(c) > 0 {
		q = q.Where("contest_id IN (?)", bun.In(c))
	}

	var problems []string
	if err := q.Scan(ctx, &problems); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Errorf(codes.Unknown, "scan rows: %s", err)
		}
	}

	return &pb.GetProblemResponse{
		Problems: problems,
	}, nil
}
