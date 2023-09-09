package problem

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

	"golang.org/x/exp/slog"
	"golang.org/x/text/unicode/norm"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"github.com/labstack/echo/v4"
	"github.com/morikuni/failure"
)

type SearchParams struct {
	Keyword string       `json:"keyword" schema:"keyword" validate:"lte=200"`
	Limit   int          `json:"limit" schema:"limit" validate:"lte=200"`
	Page    int          `json:"page" schema:"page"`
	Filter  FilterParams `json:"filter" schema:"filter"`
	Sort    string       `json:"sort" schema:"sort" validate:"omitempty,oneof=-score start_at -start_at difficulty -difficulty"`
	Facet   FacetParams  `json:"facet" schema:"facet"`
}

type FilterParams struct {
	Category   []string         `json:"category" schema:"category"`
	Difficulty acs.IntegerRange `json:"difficulty" schema:"difficulty"`
	Color      []string         `json:"color" schema:"color"`
}

type FacetParams struct {
	Term       []string            `json:"term" schema:"term" validate:"dive,oneof=category color"`
	Difficulty acs.RangeFacetParam `json:"difficulty" schema:"difficulty"`
}

func (p *SearchParams) ToQuery() url.Values {
	return solr.NewEDisMaxQueryBuilder().
		Facet(p.facet()).
		Fl(acs.FieldList(Response{})).
		Fq(p.fq()).
		Op("AND").
		Q(solr.Sanitize(norm.NFKC.String(p.Keyword))).
		QAlt("*:*").
		Qf("text_ja text_en text_reading").
		Rows(p.rows()).
		Sort(p.sort()).
		Sow(true).
		Start(p.start()).
		Build()
}

func (p *SearchParams) rows() int {
	if p.Limit == 0 {
		return 20
	}
	return p.Limit
}

func (p *SearchParams) start() int {
	if p.Page == 0 {
		return 0
	}

	return int(int(p.Page)-1) * p.rows()
}

func (p *SearchParams) sort() string {
	if p.Sort == "" {
		return "start_at desc"
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

	if f := p.Facet.Difficulty.ToFacet("difficulty"); f != nil {
		facets["difficulty"] = f
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

	if c := acs.QuoteStrings(acs.SanitizeStrings(p.Filter.Category)); len(c) != 0 {
		fq = append(fq, fmt.Sprintf("{!tag=category}category:(%s)", strings.Join(c, " OR ")))
	}
	if r := p.Filter.Difficulty.ToRange(); r != "" {
		fq = append(fq, fmt.Sprintf("{!tag=difficulty}difficulty:%s", r))
	}
	if c := acs.SanitizeStrings(p.Filter.Color); len(c) != 0 {
		fq = append(fq, fmt.Sprintf("{!tag=color}color:%s", strings.Join(c, " OR ")))
	}

	return fq
}

type Response struct {
	ProblemID    string                `json:"problem_id"`
	ProblemTitle string                `json:"problem_title"`
	ProblemURL   string                `json:"problem_url"`
	ContestID    string                `json:"contest_id"`
	ContestTitle string                `json:"contest_title"`
	ContestURL   string                `json:"contest_url"`
	Difficulty   *int                  `json:"difficulty"`
	Color        *string               `json:"color"`
	StartAt      solr.FromSolrDateTime `json:"start_at"`
	Duration     int                   `json:"duration"`
	RateChange   string                `json:"rate_change"`
	Category     string                `json:"category"`
}

type FacetCounts struct {
	Category   *solr.TermFacetCount       `json:"category,omitempty"`
	Color      *solr.TermFacetCount       `json:"color,omitempty"`
	Difficulty *solr.RangeFacetCount[int] `json:"difficulty,omitempty"`
}

type FacetResponse struct {
	Category   []acs.FacetPart `json:"category,omitempty"`
	Color      []acs.FacetPart `json:"color,omitempty"`
	Difficulty []acs.FacetPart `json:"difficulty,omitempty"`
}

func (f *FacetCounts) Into(p FacetParams) FacetResponse {
	var category []acs.FacetPart
	if f.Category != nil {
		category = acs.ConvertBucket[string](f.Category.Buckets)
	}

	var color []acs.FacetPart
	if f.Color != nil {
		color = acs.ConvertBucket[string](f.Color.Buckets)
	}

	var difficulty []acs.FacetPart
	if f.Difficulty != nil {
		difficulty = acs.ConvertRangeBucket(f.Difficulty, p.Difficulty)
	}

	return FacetResponse{
		Category:   category,
		Color:      color,
		Difficulty: difficulty,
	}
}

type Searcher struct {
	core      *solr.Core
	validator *validator.Validate
	decoder   *schema.Decoder
}

func NewSearcher(baseURL string, coreName string) (Searcher, error) {
	core, err := solr.NewSolrCore(coreName, baseURL)
	if err != nil {
		return Searcher{}, failure.Translate(err, acs.SearcherInitializeError, failure.Context{"baseURL": baseURL, "coreName": coreName}, failure.Message("failed to create problem searcher"))
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

func (s *Searcher) search(r *http.Request, params SearchParams) (int, acs.SearchResultResponse[Response]) {
	startTime := time.Now()

	query := params.ToQuery()
	res, err := solr.Select[Response, FacetCounts](s.core, query)
	if err != nil {
		slog.Error("failed to request to solr", slog.String("url", r.URL.String()), slog.Any("query", query), slog.Any("params", params), slog.String("error", fmt.Sprintf("%+v", err)))
		return 500, NewErrorResponse("internal error", params)
	}

	rows, _ := strconv.Atoi(query.Get("rows"))

	result := acs.SearchResultResponse[Response]{
		Stats: acs.SearchResultStats{
			Time:   int(time.Since(startTime).Milliseconds()),
			Total:  res.Response.NumFound,
			Index:  (res.Response.Start / int(rows)) + 1,
			Count:  int(len(res.Response.Docs)),
			Pages:  (res.Response.NumFound + int(rows) - 1) / int(rows),
			Params: params,
			Facet:  res.FacetCounts.Into(params.Facet),
		},
		Items: res.Response.Docs,
	}
	slog.Info("querylog", slog.String("domain", "problem"), slog.Int("elapsed_time", result.Stats.Time), slog.Int("hits", res.Response.NumFound), slog.Any("params", params))

	return http.StatusOK, result
}

func (s *Searcher) Liveness() bool {
	ping, err := solr.Ping(s.core)
	if err != nil {
		return false
	}

	return ping.Status == "OK"
}

func (s *Searcher) Readiness() bool {
	status, err := solr.Status(s.core)
	if err != nil {
		return false
	}

	return status.Index.NumDocs != 0
}
