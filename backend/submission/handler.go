package submission

import (
	"encoding/json"
	"fjnkt98/atcodersearch/acs"
	"fjnkt98/atcodersearch/solr"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"github.com/labstack/echo/v4"
	"github.com/morikuni/failure"
	"golang.org/x/exp/slog"
)

type SearchParams struct {
	Limit  *int         `json:"limit" schema:"limit" validate:"omitempty,lte=1000"`
	Page   int          `json:"page" schema:"page" validate:"lte=1000"`
	Filter FilterParams `json:"filter" schema:"filter"`
	Sort   string       `json:"sort" schema:"sort" validate:"omitempty,oneof=execution_time -execution_time submitted_at -submitted_at point -point length -length"`
	Facet  FacetParams  `json:"facet" schema:"facet"`
}

type FilterParams struct {
	EpochSecond   acs.IntegerRange `json:"epoch_second" schema:"epoch_second"`
	ProblemID     []string         `json:"problem_id" schema:"problem_id"`
	ContestID     []string         `json:"contest_id" schema:"contest_id"`
	Category      []string         `json:"category" schema:"category"`
	UserID        []string         `json:"user_id" schema:"user_id"`
	Language      []string         `json:"language" schema:"language"`
	LanguageGroup []string         `json:"language_group" schema:"language_group"`
	Point         acs.FloatRange   `json:"point" schema:"point"`
	Length        acs.IntegerRange `json:"length" schema:"length"`
	Result        []string         `json:"result" schema:"result"`
	ExecutionTime acs.IntegerRange `json:"execution_time" schema:"execution_time"`
}

type FacetParams struct {
	Term          []string            `json:"term" schema:"term" validate:"dive,oneof=problem_id user_id language language_group result contest_id"`
	Length        acs.RangeFacetParam `json:"length" schema:"length"`
	ExecutionTime acs.RangeFacetParam `json:"execution_time" schema:"execution_time"`
}

func (p *SearchParams) ToQuery() url.Values {
	return solr.NewLuceneQueryBuilder().
		Facet(p.facet()).
		Fl(acs.FieldList(Response{})).
		Fq(p.fq()).
		Op("AND").
		Q("*:*").
		Rows(p.rows()).
		Sort(p.sort()).
		Start(p.start()).
		Build()
}

func (p *SearchParams) rows() int {
	if p.Limit == nil {
		return 20
	}
	return *p.Limit
}

func (p *SearchParams) start() int {
	if p.Page == 0 || p.rows() == 0 {
		return 0
	}

	return (p.Page - 1) * p.rows()
}

func (p *SearchParams) sort() string {
	if p.Sort == "" {
		return "submitted_at desc"
	}
	if strings.HasPrefix(p.Sort, "-") {
		return fmt.Sprintf("%s desc", p.Sort[1:])
	} else {
		return fmt.Sprintf("%s asc", p.Sort)
	}
}

func (p *SearchParams) facet() string {
	facets := make(map[string]any)

	for _, f := range p.Facet.Term {
		facets[f] = map[string]any{
			"type":     "terms",
			"field":    f,
			"limit":    -1,
			"mincount": 0,
			"sort":     "index",
			"domain": map[string]any{
				"excludeTags": []string{f},
			},
		}
	}
	if f := p.Facet.Length.ToFacet("length"); f != nil {
		facets["length"] = f
	}
	if f := p.Facet.ExecutionTime.ToFacet("execution_time"); f != nil {
		facets["execution_time"] = f
	}

	facet, err := json.Marshal(facets)
	if err != nil {
		slog.Warn("failed to marshal json.facet parameter", slog.Any("facet", p.Facet))
		return ""
	}

	return string(facet)
}

func (p *SearchParams) fq() []string {
	fq := make([]string, 0)

	if r := p.Filter.EpochSecond.ToRange(); r != "" {
		fq = append(fq, fmt.Sprintf("{!tag=epoch_second}epoch_second:%s", r))
	}
	if r := p.Filter.Point.ToRange(); r != "" {
		fq = append(fq, fmt.Sprintf("{!tag=point}point:%s", r))
	}
	if r := p.Filter.Length.ToRange(); r != "" {
		fq = append(fq, fmt.Sprintf("{!tag=length}length:%s", r))
	}
	if r := p.Filter.ExecutionTime.ToRange(); r != "" {
		fq = append(fq, fmt.Sprintf("{!tag=execution_time}execution_time:%s", r))
	}

	if expr := strings.Join(acs.SanitizeStrings(p.Filter.ProblemID), " OR "); expr != "" {
		fq = append(fq, fmt.Sprintf("{!tag=problem_id}problem_id:(%s)", expr))
	}
	if expr := strings.Join(acs.SanitizeStrings(p.Filter.ContestID), " OR "); expr != "" {
		fq = append(fq, fmt.Sprintf("{!tag=contest_id}contest_id:(%s)", expr))
	}
	if expr := strings.Join(acs.QuoteStrings(acs.SanitizeStrings(p.Filter.Category)), " OR "); expr != "" {
		fq = append(fq, fmt.Sprintf("{!tag=category}category:(%s)", expr))
	}
	if expr := strings.Join(acs.SanitizeStrings(p.Filter.UserID), " OR "); expr != "" {
		fq = append(fq, fmt.Sprintf("{!tag=user_id}user_id:(%s)", expr))
	}
	if expr := strings.Join(acs.QuoteStrings(acs.SanitizeStrings(p.Filter.Language)), " OR "); expr != "" {
		fq = append(fq, fmt.Sprintf("{!tag=language}language:(%s)", expr))
	}
	if expr := strings.Join(acs.QuoteStrings(acs.SanitizeStrings(p.Filter.LanguageGroup)), " OR "); expr != "" {
		fq = append(fq, fmt.Sprintf("{!tag=language_group}language_group:(%s)", expr))
	}
	if expr := strings.Join(acs.SanitizeStrings(p.Filter.Result), " OR "); expr != "" {
		fq = append(fq, fmt.Sprintf("{!tag=result}result:(%s)", expr))
	}

	return fq
}

type Response struct {
	SubmissionID  int64                 `json:"submission_id"`
	SubmittedAt   solr.FromSolrDateTime `json:"submitted_at"`
	SubmissionURL string                `json:"submission_url"`
	ProblemID     string                `json:"problem_id"`
	ProblemTitle  string                `json:"problem_title"`
	ContestID     string                `json:"contest_id"`
	ContestTitle  string                `json:"contest_title"`
	Category      string                `json:"category"`
	Difficulty    int                   `json:"difficulty"`
	Color         string                `json:"color"`
	UserID        string                `json:"user_id"`
	Language      string                `json:"language"`
	Point         float64               `json:"point"`
	Length        int64                 `json:"length"`
	Result        string                `json:"result"`
	ExecutionTime *int64                `json:"execution_time"`
}

type FacetCounts struct {
	ContestID     *solr.TermFacetCount       `json:"contest_id,omitempty"`
	ProblemID     *solr.TermFacetCount       `json:"problem_id,omitempty"`
	UserID        *solr.TermFacetCount       `json:"user_id,omitempty"`
	Language      *solr.TermFacetCount       `json:"language,omitempty"`
	LanguageGroup *solr.TermFacetCount       `json:"language_group,omitempty"`
	Result        *solr.TermFacetCount       `json:"result,omitempty"`
	Length        *solr.RangeFacetCount[int] `json:"length,omitempty"`
	ExecutionTime *solr.RangeFacetCount[int] `json:"execution_time,omitempty"`
}

func (f *FacetCounts) Into(p FacetParams) FacetResponse {
	var contestID []acs.FacetPart
	if f.ContestID != nil {
		contestID = acs.ConvertBucket[string](f.ContestID.Buckets)
	}

	var problemID []acs.FacetPart
	if f.ProblemID != nil {
		problemID = acs.ConvertBucket[string](f.ProblemID.Buckets)
	}
	var userID []acs.FacetPart
	if f.UserID != nil {
		userID = acs.ConvertBucket[string](f.UserID.Buckets)
	}
	var language []acs.FacetPart
	if f.Language != nil {
		language = acs.ConvertBucket[string](f.Language.Buckets)
	}
	var languageGroup []acs.FacetPart
	if f.LanguageGroup != nil {
		languageGroup = acs.ConvertBucket[string](f.LanguageGroup.Buckets)
	}
	var result []acs.FacetPart
	if f.Result != nil {
		result = acs.ConvertBucket[string](f.Result.Buckets)
	}
	var length []acs.FacetPart
	if f.Length != nil {
		length = acs.ConvertRangeBucket(f.Length, p.Length)
	}
	var executionTime []acs.FacetPart
	if f.ExecutionTime != nil {
		executionTime = acs.ConvertRangeBucket(f.ExecutionTime, p.ExecutionTime)
	}

	return FacetResponse{
		ContestID:     contestID,
		ProblemID:     problemID,
		UserID:        userID,
		Language:      language,
		LanguageGroup: languageGroup,
		Result:        result,
		Length:        length,
		ExecutionTime: executionTime,
	}
}

type FacetResponse struct {
	ContestID     []acs.FacetPart `json:"contest_id,omitempty"`
	ProblemID     []acs.FacetPart `json:"problem_id,omitempty"`
	UserID        []acs.FacetPart `json:"user_id,omitempty"`
	Language      []acs.FacetPart `json:"language,omitempty"`
	LanguageGroup []acs.FacetPart `json:"language_group,omitempty"`
	Result        []acs.FacetPart `json:"result,omitempty"`
	Length        []acs.FacetPart `json:"length,omitempty"`
	ExecutionTime []acs.FacetPart `json:"execution_time,omitempty"`
}

type Searcher struct {
	core      *solr.Core
	validator *validator.Validate
	decoder   *schema.Decoder
}

func NewSearcher(baseURL string, coreName string) (Searcher, error) {
	core, err := solr.NewSolrCore(coreName, baseURL)
	if err != nil {
		return Searcher{}, failure.Translate(err, acs.SearcherInitializeError, failure.Context{"baseURL": baseURL, "coreName": coreName}, failure.Message("failed to create user searcher"))
	}

	validator := validator.New()
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	decoder.RegisterConverter([]string{}, func(input string) reflect.Value {
		return reflect.ValueOf(strings.Split(input, ","))
	})

	searcher := Searcher{
		core:      core,
		validator: validator,
		decoder:   decoder,
	}
	return searcher, nil
}

func NewErrorResponse(msg string, params any) acs.SearchResultResponse[Response] {
	return acs.NewErrorResponse[Response](msg, params)
}

func (s *Searcher) HandleGET(c echo.Context) error {
	raw := c.Request().URL.RawQuery
	query, err := url.ParseQuery(raw)
	if err != nil {
		slog.Error("failed to parse query string", slog.String("uri", c.Request().RequestURI), slog.String("error", fmt.Sprintf("%+v", err)))
		return c.JSON(http.StatusBadRequest, NewErrorResponse(fmt.Sprintf("failed to parse query string `%s`", raw), nil))
	}

	var params SearchParams
	if err := s.decoder.Decode(&params, query); err != nil {
		slog.Error("failed to decode request parameter", slog.String("uri", c.Request().RequestURI), slog.String("error", fmt.Sprintf("%+v", err)))
		return c.JSON(http.StatusBadRequest, NewErrorResponse(fmt.Sprintf("failed to decode request parameter `%s`", raw), nil))
	}

	if err := s.validator.Struct(params); err != nil {
		slog.Error("validation error", slog.String("uri", c.Request().RequestURI), slog.Any("params", params), slog.String("error", fmt.Sprintf("%+v", err)))
		return c.JSON(http.StatusBadRequest, NewErrorResponse(fmt.Sprintf("validation error `%s`: %s", raw, err.Error()), params))
	}
	code, res := s.search(c.Request(), params)
	return c.JSON(code, res)
}

func (s *Searcher) HandlePOST(c echo.Context) error {
	var params SearchParams
	if err := c.Bind(&params); err != nil {
		slog.Error("failed to decode request parameter", slog.String("uri", c.Request().RequestURI), slog.String("error", fmt.Sprintf("%+v", err)))
		return c.JSON(http.StatusBadRequest, NewErrorResponse("failed to decode request parameter", nil))
	}

	if err := s.validator.Struct(params); err != nil {
		slog.Error("validation error", slog.String("uri", c.Request().RequestURI), slog.Any("params", params), slog.String("error", fmt.Sprintf("%+v", err)))
		return c.JSON(http.StatusBadRequest, NewErrorResponse(fmt.Sprintf("validation error: %s", err.Error()), params))
	}

	code, res := s.search(c.Request(), params)
	return c.JSON(code, res)
}

func (s *Searcher) search(r *http.Request, params SearchParams) (int, acs.SearchResultResponse[Response]) {
	startTime := time.Now()

	query := params.ToQuery()
	res, err := solr.Select[Response, FacetCounts](s.core, query)
	if err != nil {
		slog.Error("failed to request to solr", slog.String("url", r.URL.String()), slog.Any("query", query), slog.Any("params", params), slog.String("error", fmt.Sprintf("%+v", err)))
		return 500, NewErrorResponse("internal error", params)
	}

	rows, _ := strconv.Atoi(query.Get("rows"))
	var pages = 0
	var index = 0
	if rows != 0 {
		pages = (res.Response.NumFound + rows) / rows
		index = (res.Response.Start / rows) + 1
	}

	result := acs.SearchResultResponse[Response]{
		Stats: acs.SearchResultStats{
			Time:   int(time.Since(startTime).Milliseconds()),
			Total:  res.Response.NumFound,
			Index:  index,
			Count:  len(res.Response.Docs),
			Pages:  pages,
			Params: &params,
			Facet:  res.FacetCounts.Into(params.Facet),
		},
		Items: res.Response.Docs,
	}
	slog.Info("querylog", slog.String("domain", "submission"), slog.Int("elapsed_time", result.Stats.Time), slog.Int("hits", res.Response.NumFound), slog.Any("params", params))

	return http.StatusOK, result
}

func (s *Searcher) Liveness() bool {
	ping, err := solr.Ping(s.core)
	if err != nil {
		slog.Error("submission core doesn't alive", slog.String("error", fmt.Sprintf("%+v", err)))
		return false
	}

	return ping.Status == "OK"
}

func (s *Searcher) Readiness() bool {
	status, err := solr.Status(s.core)
	if err != nil {
		slog.Error("submission core isn't ready", slog.String("error", fmt.Sprintf("%+v", err)))
		return false
	}

	return status.Index.NumDocs != 0
}
