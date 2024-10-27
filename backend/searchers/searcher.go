package searchers

import (
	"context"
	"encoding/json"
	"errors"
	"fjnkt98/atcodersearch/api"
	"fmt"
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

func ParseFacetDistribution(facetDistribution any) map[string]map[string]int {
	facet, ok := facetDistribution.(map[string]any)
	if !ok {
		return nil
	}

	result := make(map[string]map[string]int)
	for field, v := range facet {
		counts, ok := v.(map[string]any)
		if !ok {
			continue
		}

		fieldCounts := make(map[string]int)
		for key, count := range counts {
			if count, ok := count.(float64); ok {
				fieldCounts[key] = int(count)
			}
		}
		result[field] = fieldCounts
	}

	return result
}

type Into[T any] interface {
	Into() T
}

var ErrUnexpectedHitType = errors.New("item of res.Hits isn't an map[string]any")

func ScanItems[S any, T Into[S]](hits []any) ([]S, error) {
	items := make([]S, len(hits))
	for i, hit := range hits {
		hit, ok := hit.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("type assert: %w", ErrUnexpectedHitType)
		}

		var item T
		if err := mapstructure.Decode(hit, &item); err != nil {
			return nil, fmt.Errorf("decode items: %w", err)
		}

		items[i] = item.Into()
	}

	return items, nil
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

func IntRangeFilter(field string, value api.OptIntRange) [][]string {
	v, ok := value.Get()
	if !ok {
		return nil
	}
	result := make([][]string, 0, 2)

	if from, ok := v.From.Get(); ok {
		result = append(result, []string{fmt.Sprintf("%s >= %d", field, from)})
	}
	if to, ok := v.To.Get(); ok {
		result = append(result, []string{fmt.Sprintf("%s < %d", field, to)})
	}

	if len(result) == 0 {
		return nil
	}

	return result
}

type Searcher struct {
	client meilisearch.ServiceManager
	pool   *pgxpool.Pool
}

func NewSearcher(client meilisearch.ServiceManager, pool *pgxpool.Pool) *Searcher {
	return &Searcher{
		client: client,
		pool:   pool,
	}
}

func (s *Searcher) APICategoryGet(ctx context.Context) (*api.APICategoryGetOK, error) {
	db := bun.NewDB(stdlib.OpenDBFromPool(s.pool), pgdialect.New())

	q := db.NewSelect().Distinct().
		Column("category").
		Table("contests").
		Order("category ASC")

	var categories []string
	if err := q.Scan(ctx, &categories); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("scan rows: %w", err)
		}
	}

	return &api.APICategoryGetOK{
		Categories: categories,
	}, nil
}

func (s *Searcher) APIContestGet(ctx context.Context, params api.APIContestGetParams) (*api.APIContestGetOK, error) {
	db := bun.NewDB(stdlib.OpenDBFromPool(s.pool), pgdialect.New())

	q := db.NewSelect().
		Column("contest_id").
		Table("contests").
		Order("start_epoch_second DESC")

	if c := params.Category; len(c) > 0 {
		q = q.Where("category IN (?)", bun.In(c))
	}

	var contests []string
	if err := q.Scan(ctx, &contests); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("scan rows: %w", err)
		}
	}

	return &api.APIContestGetOK{
		Contests: contests,
	}, nil
}

func (s *Searcher) APILanguageGet(ctx context.Context, params api.APILanguageGetParams) (*api.APILanguageGetOK, error) {
	db := bun.NewDB(stdlib.OpenDBFromPool(s.pool), pgdialect.New())

	cte := db.NewSelect().
		ColumnExpr("?", bun.Ident("group")).
		ColumnExpr("ARRAY_AGG(language ORDER BY language ASC) AS languages").
		Table("languages").
		GroupExpr("?", bun.Ident("group")).
		OrderExpr("? ASC", bun.Ident("group"))

	if g := params.Group; len(g) > 0 {
		cte = cte.Where("? IN (?)", bun.Ident("group"), bun.In(g))
	}

	q := db.NewSelect().
		With("l", cte).
		ColumnExpr("JSON_AGG(l)").
		Table("l")

	var body []byte
	if err := q.Scan(ctx, &body); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &api.APILanguageGetOK{
				Languages: []api.Language{},
			}, nil
		} else {
			return nil, fmt.Errorf("scan rows: %w", err)
		}
	}

	languages := make([]api.Language, 0)
	if len(body) > 0 {
		if err := json.Unmarshal(body, &languages); err != nil {
			return nil, fmt.Errorf("unmarshal body: %w", err)
		}
	}

	return &api.APILanguageGetOK{
		Languages: languages,
	}, nil
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
	StartAt        int     `mapstructure:"startAt" bun:"start_at"`
	Duration       int     `mapstructure:"duration" bun:"duration"`
	RateChange     string  `mapstructure:"rateChange" bun:"rate_change"`
	Category       string  `mapstructure:"category" bun:"category"`
	IsExperimental bool    `mapstructure:"isExperimental" bun:"is_experimental"`
}

func (p Problem) Into() api.Problem {
	var difficulty api.OptInt
	if d := p.Difficulty; d != nil {
		difficulty = api.NewOptInt(int(*d))
	}

	return api.Problem{
		ProblemId:      p.ProblemID,
		ProblemTitle:   p.ProblemTitle,
		ProblemUrl:     p.ProblemURL,
		ContestId:      p.ContestID,
		ContestTitle:   p.ContestTitle,
		ContestUrl:     p.ContestURL,
		Difficulty:     difficulty,
		StartAt:        p.StartAt,
		Duration:       p.Duration,
		RateChange:     p.RateChange,
		Category:       p.Category,
		IsExperimental: p.IsExperimental,
	}
}

func (s *Searcher) APIProblemGet(ctx context.Context, params api.APIProblemGetParams) (*api.APIProblemGetOK, error) {
	db := bun.NewDB(stdlib.OpenDBFromPool(s.pool), pgdialect.New())

	q := db.NewSelect().
		Column("problem_id").
		Table("problems").
		Join("LEFT JOIN contests USING(contest_id)").
		Order("problem_id ASC")

	if c := params.Category; len(c) > 0 {
		q = q.Where("category IN (?)", bun.In(c))
	}

	if c := params.ContestId; len(c) > 0 {
		q = q.Where("contest_id IN (?)", bun.In(c))
	}

	var problems []string
	if err := q.Scan(ctx, &problems); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("scan rows: %w", err)
		}
	}

	return &api.APIProblemGetOK{
		Problems: problems,
	}, nil
}

func (s *Searcher) APIProblemPost(ctx context.Context, req *api.APIProblemPostReq) (*api.APIProblemPostOK, error) {
	index := s.client.Index("problems")

	q, err := s.createSearchProblemQuery(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("create query: %w", err)
	}

	res, err := index.SearchWithContext(ctx, req.Q.Value, q)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}

	items, err := ScanItems[api.Problem, Problem](res.Hits)
	if err != nil {
		return nil, fmt.Errorf("scan items: %w", err)
	}

	facet := make(map[string][]api.Count)
	mapping := map[string]string{
		"category":        "category",
		"difficultyFacet": "difficulty",
	}
	for field, countMap := range ParseFacetDistribution(res.FacetDistribution) {
		counts := make([]api.Count, 0, 16)
		for _, k := range slices.Sorted(maps.Keys(countMap)) {
			counts = append(counts, api.Count{Label: k, Count: countMap[k]})
		}
		if f, ok := mapping[field]; ok {
			facet[f] = counts
		}
	}

	return &api.APIProblemPostOK{
		Time:  int(res.ProcessingTimeMs),
		Total: int(res.TotalHits),
		Index: int(res.Page),
		Pages: int(res.TotalPages),
		Items: items,
		Facet: facet,
	}, nil
}

func (s *Searcher) createSearchProblemQuery(ctx context.Context, req *api.APIProblemPostReq) (*meilisearch.SearchRequest, error) {
	q := &meilisearch.SearchRequest{
		AttributesToRetrieve: FieldList(new(Problem)),
	}

	if limit, ok := req.Limit.Get(); ok {
		q.HitsPerPage = int64(limit)
	}

	if page, ok := req.Page.Get(); ok {
		if page <= 0 {
			q.Page = 1
		} else {
			q.Page = int64(page)
		}
	} else {
		q.Page = 1
	}

	sorts := make([]string, 0, len(req.Sort)+1)
	for _, s := range req.Sort {
		sorts = append(sorts, string(s))
	}
	sorts = append(sorts, "problemId:asc")
	q.Sort = sorts

	fields := make([]string, len(req.Facet))
	for i, f := range req.Facet {
		fields[i] = string(f)
	}

	if facet, err := GenFacet(fields, map[string]string{"category": "category", "difficulty": "difficultyFacet"}); err != nil {
		return nil, err
	} else {
		q.Facets = facet
	}

	filters := make([][]string, 0, 3)
	if categories := StringSliceFilter("category", req.Category); categories != nil {
		filters = append(filters, categories)
	}
	if difficulty := IntRangeFilter("difficulty", req.Difficulty); difficulty != nil {
		filters = append(filters, difficulty...)
	}
	if experimental, ok := req.Experimental.Get(); ok {
		filters = append(filters, []string{fmt.Sprintf("isExperimental = %t", experimental)})
	}

	if userID, ok := req.UserId.Get(); ok {
		db := bun.NewDB(stdlib.OpenDBFromPool(s.pool), pgdialect.New())

		var solvedProblems []string
		if err := db.NewSelect().
			Distinct().
			Column("problem_id").
			Table("submissions").
			Where("user_id = ?", userID).
			Where("result = ?", "AC").
			Scan(ctx, &solvedProblems); err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return nil, fmt.Errorf("fetch solved problems: %w", err)
			}
		}

		if len(solvedProblems) > 0 {
			filters = append(filters, []string{fmt.Sprintf("problemId NOT IN [%s]", strings.Join(solvedProblems, ", "))})
		}
	}

	if len(filters) > 0 {
		q.Filter = filters
	}

	return q, nil
}

type User struct {
	UserID        string  `mapstructure:"userId" bun:"user_id"`
	Rating        int     `mapstructure:"rating" bun:"rating"`
	HighestRating int     `mapstructure:"highestRating" bun:"highest_rating"`
	Affiliation   *string `mapstructure:"affiliation" bun:"affiliation"`
	BirthYear     *int    `mapstructure:"birthYear" bun:"birth_year"`
	Country       *string `mapstructure:"country" bun:"country"`
	Crown         *string `mapstructure:"crown" bun:"crown"`
	JoinCount     int     `mapstructure:"joinCount" bun:"join_count"`
	Rank          int     `mapstructure:"rank" bun:"rank"`
	ActiveRank    *int    `mapstructure:"activeRank" bun:"active_rank"`
	Wins          int     `mapstructure:"wins" bun:"wins"`
	UserURL       string  `mapstructure:"userUrl" bun:"user_url"`
}

func (u User) Into() api.User {
	var affiliation api.OptString
	if a := u.Affiliation; a != nil {
		affiliation = api.NewOptString(*a)
	}

	var birthYear api.OptInt
	if b := u.BirthYear; b != nil {
		birthYear = api.NewOptInt(*b)
	}

	var country api.OptString
	if a := u.Country; a != nil {
		country = api.NewOptString(*a)
	}

	var crown api.OptString
	if a := u.Crown; a != nil {
		crown = api.NewOptString(*a)
	}

	var activeRank api.OptInt
	if b := u.ActiveRank; b != nil {
		activeRank = api.NewOptInt(*b)
	}

	return api.User{
		UserId:        u.UserID,
		Rating:        u.Rating,
		HighestRating: u.HighestRating,
		Affiliation:   affiliation,
		BirthYear:     birthYear,
		Country:       country,
		Crown:         crown,
		JoinCount:     u.JoinCount,
		Rank:          u.Rank,
		ActiveRank:    activeRank,
		Wins:          u.Wins,
		UserUrl:       u.UserURL,
	}
}

func (s *Searcher) APIUserPost(ctx context.Context, req *api.APIUserPostReq) (*api.APIUserPostOK, error) {
	index := s.client.Index("users")

	q, err := createSearchUserQuery(req)
	if err != nil {
		return nil, fmt.Errorf("create query: %w", err)
	}

	res, err := index.SearchWithContext(ctx, req.Q.Value, q)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}

	items, err := ScanItems[api.User, User](res.Hits)
	if err != nil {
		return nil, err
	}

	facet := make(map[string][]api.Count)
	mapping := map[string]string{
		"country":        "country",
		"ratingFacet":    "rating",
		"birthYearFacet": "birthYear",
		"joinCountFacet": "joinCount",
	}
	for field, countMap := range ParseFacetDistribution(res.FacetDistribution) {
		counts := make([]api.Count, 0, 16)
		for _, k := range slices.Sorted(maps.Keys(countMap)) {
			counts = append(counts, api.Count{Label: k, Count: countMap[k]})
		}
		if f, ok := mapping[field]; ok {
			facet[f] = counts
		}
	}

	return &api.APIUserPostOK{
		Time:  int(res.ProcessingTimeMs),
		Total: int(res.TotalHits),
		Index: int(res.Page),
		Pages: int(res.TotalPages),
		Items: items,
		Facet: facet,
	}, nil
}

func createSearchUserQuery(req *api.APIUserPostReq) (*meilisearch.SearchRequest, error) {
	q := &meilisearch.SearchRequest{
		AttributesToRetrieve: FieldList(new(User)),
	}

	if limit, ok := req.Limit.Get(); ok {
		q.HitsPerPage = int64(limit)
	}

	if page, ok := req.Page.Get(); ok {
		if page <= 0 {
			q.Page = 1
		} else {
			q.Page = int64(page)
		}
	} else {
		q.Page = 1
	}

	sorts := make([]string, 0, len(req.Sort)+1)
	for _, s := range req.Sort {
		sorts = append(sorts, string(s))
	}
	sorts = append(sorts, "userId:asc")
	q.Sort = sorts

	fields := make([]string, len(req.Facet))
	for i, f := range req.Facet {
		fields[i] = string(f)
	}

	if facet, err := GenFacet(fields, map[string]string{
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
	if users := StringSliceFilter("userId", req.UserId); users != nil {
		filters = append(filters, users)
	}
	if rating := IntRangeFilter("rating", req.Rating); rating != nil {
		filters = append(filters, rating...)
	}
	if birthYear := IntRangeFilter("birthYear", req.BirthYear); birthYear != nil {
		filters = append(filters, birthYear...)
	}
	if joinCount := IntRangeFilter("joinCount", req.JoinCount); joinCount != nil {
		filters = append(filters, joinCount...)
	}
	if countries := StringSliceFilter("country", req.Country); countries != nil {
		filters = append(filters, countries)
	}

	if len(filters) > 0 {
		q.Filter = filters
	}

	return q, nil
}

type Submission struct {
	SubmissionID  int     `bun:"submission_id"`
	SubmittedAt   int     `bun:"submitted_at"`
	SubmissionURL string  `bun:"submission_url"`
	ProblemID     string  `bun:"problem_id"`
	ProblemTitle  string  `bun:"problem_title"`
	ProblemURL    string  `bun:"problem_url"`
	ContestID     string  `bun:"contest_id"`
	ContestTitle  string  `bun:"contest_title"`
	ContestURL    string  `bun:"contest_url"`
	Category      string  `bun:"category"`
	Difficulty    *int    `bun:"difficulty"`
	UserID        string  `bun:"user_id"`
	Language      string  `bun:"language"`
	LanguageGroup string  `bun:"language_group"`
	Point         float64 `bun:"point"`
	Length        int     `bun:"length"`
	Result        string  `bun:"result"`
	ExecutionTime *int    `bun:"execution_time"`
}

func (s Submission) Into() api.Submission {
	var difficulty api.OptInt
	if d := s.Difficulty; d != nil {
		difficulty = api.NewOptInt(*d)
	}

	var executionTime api.OptInt
	if e := s.ExecutionTime; e != nil {
		executionTime = api.NewOptInt(*e)
	}

	return api.Submission{
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
		Difficulty:    difficulty,
		UserId:        s.UserID,
		Language:      s.Language,
		LanguageGroup: s.LanguageGroup,
		Point:         s.Point,
		Length:        s.Length,
		Result:        s.Result,
		ExecutionTime: executionTime,
	}
}

func (s *Searcher) APISubmissionPost(ctx context.Context, req *api.APISubmissionPostReq) (*api.APISubmissionPostOK, error) {
	start := time.Now()

	q, err := createSearchSubmissionQuery(s.pool, req)
	if err != nil {
		return nil, fmt.Errorf("create query: %w", err)
	}

	items := make([]api.Submission, 0, req.Limit.Value)
	var rows []Submission
	if err := q.Scan(ctx, &rows); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("scan rows: %w", err)
		}
	}
	for _, r := range rows {
		items = append(items, r.Into())
	}

	return &api.APISubmissionPostOK{
		Time:  int(time.Since(start) / time.Millisecond),
		Items: items,
	}, nil
}

func createSearchSubmissionQuery(pool *pgxpool.Pool, req *api.APISubmissionPostReq) (*bun.SelectQuery, error) {
	db := bun.NewDB(stdlib.OpenDBFromPool(pool), pgdialect.New())

	limit := req.Limit.Value

	var offset int
	if page, ok := req.Page.Get(); ok {
		if page <= 0 {
			offset = 0
		} else {
			offset = (page - 1) * limit
		}
	} else {
		offset = 0
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
	if sorts := req.Sort; len(sorts) > 0 {
		for _, s := range sorts {
			field, direction, ok := strings.Cut(string(s), ":")
			if !ok {
				return nil, fmt.Errorf("%w: invalid sort key `%s`", ErrInvalidRequest, field)
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

	if epochSecond, ok := req.EpochSecond.Get(); ok {
		if from, ok := epochSecond.From.Get(); ok {
			q = q.Where("s.epoch_second >= ?", from)
		}
		if to, ok := epochSecond.To.Get(); ok {
			q = q.Where("s.epoch_second < ?", to)
		}
	}
	if problems := req.ProblemId; len(problems) > 0 {
		q = q.Where("s.problem_id IN (?)", bun.In(problems))
	}
	if contests := req.ContestId; len(contests) > 0 {
		q = q.Where("s.contest_id IN (?)", bun.In(contests))
	}
	if categories := req.Category; len(categories) > 0 {
		q = q.Where("c.category IN (?)", bun.In(categories))
	}
	if users := req.UserId; len(users) > 0 {
		q = q.Where("s.user_id IN (?)", bun.In(users))
	}
	if languages := req.Language; len(languages) > 0 {
		q = q.Where("s.language IN (?)", bun.In(languages))
	}
	if groups := req.LanguageGroup; len(groups) > 0 {
		q = q.Where("l.group IN (?)", bun.In(groups))
	}
	if point, ok := req.Point.Get(); ok {
		if from, ok := point.From.Get(); ok {
			q = q.Where("s.point >= ?", from)
		}
		if to, ok := point.To.Get(); ok {
			q = q.Where("s.point < ?", to)
		}
	}
	if length, ok := req.Length.Get(); ok {
		if from, ok := length.From.Get(); ok {
			q = q.Where("s.length >= ?", from)
		}
		if to, ok := length.To.Get(); ok {
			q = q.Where("s.length < ?", to)
		}
	}
	if results := req.Result; len(results) > 0 {
		q = q.Where("s.result IN (?)", bun.In(results))
	}
	if exec, ok := req.ExecutionTime.Get(); ok {
		if from, ok := exec.From.Get(); ok {
			q = q.Where("s.execution_time >= ?", from)
		}
		if to, ok := exec.To.Get(); ok {
			q = q.Where("s.execution_time < ?", to)
		}
	}

	return q, nil
}
