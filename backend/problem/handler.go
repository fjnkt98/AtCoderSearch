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
	"github.com/morikuni/failure"
)

type SearchParams struct {
	Keyword string        `json:"keyword,omitempty" schema:"keyword" validate:"lte=200"`
	Limit   uint          `json:"limit,omitempty" schema:"limit" validate:"lte=200"`
	Page    uint          `json:"page,omitempty" schema:"page"`
	Filter  *FilterParams `json:"filter,omitempty" schema:"filter"`
	Sort    string        `json:"sort,omitempty" schema:"sort" validate:"omitempty,oneof=-score start_at -start_at difficulty -difficulty"`
	Facet   []string      `json:"facet,omitempty" schema:"facet" validate:"dive,oneof=category color"`
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

func (p *SearchParams) rows() uint {
	if p.Limit == 0 {
		return 20
	}
	return p.Limit
}

func (p *SearchParams) start() uint {
	if p.Page == 0 {
		return 0
	}

	return uint(int(p.Page)-1) * p.rows()
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

	for _, f := range p.Facet {
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

	facet, err := json.Marshal(facets)
	if err != nil {
		slog.Warn("failed to marshal json.facet parameter", slog.Any("facet", p.Facet))
		return ""
	}

	return string(facet)
}

func (p *SearchParams) fq() []string {
	if p.Filter == nil {
		return make([]string, 0)
	}

	fq := make([]string, 0)

	if c := acs.QuoteStrings(acs.SanitizeStrings(p.Filter.Category)); len(c) != 0 {
		fq = append(fq, fmt.Sprintf("{!tag=category}category:(%s)", strings.Join(c, " OR ")))
	}
	if p.Filter.Difficulty != nil {
		if r := p.Filter.Difficulty.ToRange(); r != "" {
			fq = append(fq, fmt.Sprintf("{!tag=difficulty}difficulty:%s", r))
		}
	}
	if c := acs.SanitizeStrings(p.Filter.Color); len(c) != 0 {
		fq = append(fq, fmt.Sprintf("{!tag=color}color:%s", strings.Join(c, " OR ")))
	}

	return fq
}

type FilterParams struct {
	Category   []string               `json:"category,omitempty" schema:"category"`
	Difficulty *acs.IntegerRange[int] `json:"difficulty,omitempty" schema:"difficulty"`
	Color      []string               `json:"color,omitempty" schema:"color"`
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
	// Count    uint                    `json:"count"`
	Category solr.TermFacetCount `json:"category"`
	Color    solr.TermFacetCount `json:"color"`
}

type FacetResponse struct {
	Category []acs.FacetPart `json:"category,omitempty"`
	Color    []acs.FacetPart `json:"color,omitempty"`
}

func (f *FacetCounts) Into() FacetResponse {
	return FacetResponse{
		Category: acs.ConvertBucket[string](f.Category.Buckets),
		Color:    acs.ConvertBucket[string](f.Color.Buckets),
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
		return Searcher{}, failure.Translate(err, SearcherInitializeError, failure.Context{"baseURL": baseURL, "coreName": coreName}, failure.Message("failed to create problem searcher"))
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

func (s *Searcher) HandleSearch(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json; charset=utf8")
		encoder := json.NewEncoder(w)

		query, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			slog.Error("failed to parse query string", slog.String("url", r.URL.String()), slog.String("error", fmt.Sprintf("%+v", err)))
			encoder.Encode(NewErrorResponse(fmt.Sprintf("failed to parse query string `%s`", r.URL.RawQuery), nil))
			return
		}

		var params SearchParams
		if err := s.decoder.Decode(&params, query); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			slog.Error("failed to decode request parameter", slog.String("url", r.URL.String()), slog.String("error", fmt.Sprintf("%+v", err)))
			encoder.Encode(NewErrorResponse(fmt.Sprintf("failed to decode request parameter `%s`", r.URL.RawQuery), nil))
			return
		}

		if err := s.validator.Struct(params); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			slog.Error("validation error", slog.String("url", r.URL.String()), slog.Any("params", params), slog.String("error", fmt.Sprintf("%+v", err)))
			encoder.Encode(NewErrorResponse(fmt.Sprintf("validation error `%s`: %s", r.URL.RawQuery, err.Error()), params))
			return
		}

		code, res := s.search(r, params)
		w.WriteHeader(code)
		encoder.Encode(res)
	default:

	}
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
			Time:   uint(time.Since(startTime).Milliseconds()),
			Total:  res.Response.NumFound,
			Index:  (res.Response.Start / uint(rows)) + 1,
			Count:  uint(len(res.Response.Docs)),
			Pages:  (res.Response.NumFound + uint(rows) - 1) / uint(rows),
			Params: params,
			Facet:  res.FacetCounts.Into(),
		},
		Items: res.Response.Docs,
	}
	slog.Info("querylog", slog.String("domain", "problem"), slog.Uint64("elapsed_time", uint64(result.Stats.Time)), slog.Uint64("hits", uint64(res.Response.NumFound)), slog.Any("params", params))

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
