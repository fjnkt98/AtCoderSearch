package problem

import (
	"encoding/json"
	"fjnkt98/atcodersearch/acs"
	"fjnkt98/atcodersearch/solr"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

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

	return (p.Page - 1) / p.rows()
}

func (p *SearchParams) sort() string {
	if p.Sort == "" {
		return "score desc"
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
		log.Printf("WARN: failed to marshal json.facet parameter from %v", p.Facet)
		return ""
	}

	return string(facet)
}

func (p *SearchParams) fq() []string {
	if p.Filter == nil {
		return make([]string, 0)
	}

	fq := make([]string, 0)

	if c := acs.SanitizeStrings(p.Filter.Category); len(c) != 0 {
		fq = append(fq, fmt.Sprintf("{!tag=category}category:(%s)", strings.Join(c, " OR ")))
	}
	if p.Filter.Difficulty != nil {
		if r := p.Filter.Difficulty.ToRange(); r != "" {
			fq = append(fq, fmt.Sprintf("{!tag=difficulty}difficulty:%s", r))
		}
	}

	return fq
}

type FilterParams struct {
	Category   []string               `json:"category,omitempty" schema:"category"`
	Difficulty *acs.IntegerRange[int] `json:"difficulty,omitempty" schema:"difficulty"`
}

type Response struct {
	ProblemID    string                `json:"problem_id" solr:"problem_id"`
	ProblemTitle string                `json:"problem_title" solr:"problem_title"`
	ProblemURL   string                `json:"problem_url" solr:"problem_url"`
	ContestID    string                `json:"contest_id" solr:"contest_id"`
	ContestTitle string                `json:"contest_title" solr:"contest_title"`
	ContestURL   string                `json:"contest_url" solr:"contest_url"`
	Difficulty   *int                  `json:"difficulty" solr:"difficulty"`
	Color        *string               `json:"color" solr:"color"`
	StartAt      solr.FromSolrDateTime `json:"start_at" solr:"start_at"`
	Duration     int                   `json:"duration" solr:"duration"`
	RateChange   string                `json:"rate_change" solr:"rate_change"`
	Category     string                `json:"category" solr:"category"`
}

type FacetCounts struct {
	// Count    uint                    `json:"count"`
	Category solr.SolrTermFacetCount `json:"category"`
	Color    solr.SolrTermFacetCount `json:"color"`
}

type FacetResponse struct {
	Category []FacetPart `json:"category,omitempty"`
	Color    []FacetPart `json:"color,omitempty"`
}

type FacetPart struct {
	Label string `json:"label"`
	Count uint   `json:"count"`
}

func (f *FacetCounts) Into() FacetResponse {
	category := make([]FacetPart, len(f.Category.Buckets))
	for i, b := range f.Category.Buckets {
		category[i] = FacetPart{
			Label: b.Val,
			Count: b.Count,
		}
	}

	color := make([]FacetPart, len(f.Color.Buckets))
	for i, b := range f.Color.Buckets {
		color[i] = FacetPart{
			Label: b.Val,
			Count: b.Count,
		}
	}

	return FacetResponse{
		Category: category,
		Color:    color,
	}
}

type Searcher struct {
	core      *solr.SolrCore[Response, FacetCounts]
	validator *validator.Validate
	decoder   *schema.Decoder
}

func NewSearcher(baseURL string, coreName string) (Searcher, error) {
	core, err := solr.NewSolrCore[Response, FacetCounts](coreName, baseURL)
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
		core:      &core,
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
			log.Printf("ERROR: failed to parse query string `%s`: %s", r.URL.RawQuery, err.Error())
			encoder.Encode(NewErrorResponse(fmt.Sprintf("failed to parse query string `%s`", r.URL.RawQuery), nil))
			return
		}

		var params SearchParams
		if err := s.decoder.Decode(&params, query); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("ERROR: failed to decode request parameter `%s`: %s", r.URL.RawQuery, err.Error())
			encoder.Encode(NewErrorResponse(fmt.Sprintf("failed to decode request parameter `%s`", r.URL.RawQuery), nil))
			return
		}

		if err := s.validator.Struct(params); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("ERROR: validation error: %+v, `%s`: %s", params, r.URL.RawQuery, err.Error())
			encoder.Encode(NewErrorResponse(fmt.Sprintf("validation error `%s`: %s", r.URL.RawQuery, err.Error()), params))
			return
		}

		code, res := search(s.core, params)
		w.WriteHeader(code)
		encoder.Encode(res)
	default:

	}
}

func search(core *solr.SolrCore[Response, FacetCounts], params SearchParams) (int, acs.SearchResultResponse[Response]) {
	startTime := time.Now()

	query := params.ToQuery()
	res, err := core.Select(query)
	if err != nil {
		log.Printf("ERROR: failed to request to solr with %+v, from %+v: %s", query, params, err.Error())
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
	querylog := acs.QueryLog{
		RequestAt: startTime,
		Domain:    "problem",
		Time:      result.Stats.Time,
		Hits:      res.Response.NumFound,
		Params:    params,
	}
	encoder := json.NewEncoder(log.Writer())
	encoder.Encode(querylog)

	return http.StatusOK, result
}

func (s *Searcher) Liveness() bool {
	ping, err := s.core.Ping()
	if err != nil {
		return false
	}

	return ping.Status == "OK"
}

func (s *Searcher) Readiness() bool {
	status, err := s.core.Status()
	if err != nil {
		return false
	}

	return status.Index.NumDocs != 0
}
