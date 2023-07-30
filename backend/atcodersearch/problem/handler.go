package problem

import (
	"encoding/json"
	"fjnkt98/atcodersearch/atcodersearch/common"
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
)

var validate = validator.New()
var decoder = schema.NewDecoder()
var lister = common.NewFieldLister()

type ProblemSearchParams struct {
	Keyword string        `validate:"lte=200" json:"keyword,omitempty"`
	Limit   uint          `validate:"lte=200" json:"limit,omitempty"`
	Page    uint          `json:"page,omitempty"`
	Filter  *FilterParams `json:"filter,omitempty"`
	Sort    string        `validate:"omitempty,oneof=-score start_at -start_at difficulty -difficulty" json:"sort,omitempty"`
	Facet   []string      `validate:"dive,oneof=category color" json:"facet,omitempty"`
}

func (p *ProblemSearchParams) ToQuery() url.Values {
	return solr.NewEDisMaxQueryBuilder().
		Facet(p.facet()).
		Fl(lister.FieldList(ProblemResponse{})).
		Fq(p.fq()).
		Op("AND").
		Q(solr.Sanitize(norm.NFKC.String(p.Keyword))).
		Qf("text_ja text_en text_reading").
		Rows(p.rows()).
		Sort(p.sort()).
		Sow(true).
		Start(p.start()).
		Build()
}

func (p *ProblemSearchParams) rows() uint {
	if p.Limit == 0 {
		return 20
	}
	return p.Limit
}

func (p *ProblemSearchParams) start() uint {
	if p.Page == 0 {
		return 0
	}

	return (p.Page - 1) / p.rows()
}

func (p *ProblemSearchParams) sort() string {
	if p.Sort == "" {
		return "score desc"
	}
	if strings.HasPrefix(p.Sort, "-") {
		return fmt.Sprintf("%s desc", p.Sort[1:])
	} else {
		return fmt.Sprintf("%s asc", p.Sort)
	}
}

func (p *ProblemSearchParams) facet() string {
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
		log.Printf("WARN: failed to marshal json.facet parameter from %+v", p.Facet)
		return ""
	}

	return string(facet)
}

func (p *ProblemSearchParams) fq() []string {
	if p.Filter == nil {
		return make([]string, 0)
	}

	fq := make([]string, 0)

	categories := make([]string, 0, len(p.Filter.Category))
	for _, c := range p.Filter.Category {
		category := solr.Sanitize(c)
		if c == "" {
			continue
		}
		categories = append(categories, category)
	}

	fq = append(fq, fmt.Sprintf("{!tag=category}category:(%s)", strings.Join(categories, " OR ")))
	fq = append(fq, fmt.Sprintf("{!tag=difficulty}difficulty:%s", p.Filter.Difficulty.ToRange()))

	return fq
}

type FilterParams struct {
	Category   []string                `json:"category,omitempty"`
	Difficulty common.RangeFilterParam `json:"difficulty,omitempty"`
}

type ProblemResponse struct {
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
	Category []FacetPart `json:"category"`
	Color    []FacetPart `json:"color"`
}

type FacetPart struct {
	Label string `json:"label"`
	Count uint   `json:"count"`
}

func NewFacetResponse(facet FacetCounts) FacetResponse {
	category := make([]FacetPart, len(facet.Category.Buckets))
	for i, b := range facet.Category.Buckets {
		category[i] = FacetPart{
			Label: b.Val,
			Count: b.Count,
		}
	}

	color := make([]FacetPart, len(facet.Color.Buckets))
	for i, b := range facet.Color.Buckets {
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

type ProblemSearcher struct {
	core *solr.SolrCore[ProblemResponse, FacetCounts]
}

func NewProblemSearcher(baseURL string, coreName string) (ProblemSearcher, error) {
	core, err := solr.NewSolrCore[ProblemResponse, FacetCounts](coreName, baseURL)
	if err != nil {
		return ProblemSearcher{}, fmt.Errorf("failed to create problem searcher: %w", err)
	}

	decoder.IgnoreUnknownKeys(true)
	decoder.RegisterConverter([]string{}, func(input string) reflect.Value {
		return reflect.ValueOf(strings.Split(input, ","))
	})
	searcher := ProblemSearcher{
		core: &core,
	}
	return searcher, nil
}

func NewErrorResponse(msg string, params any) common.SearchResultResponse[ProblemResponse] {
	return common.NewErrorResponse[ProblemResponse](msg, params)
}

func (s *ProblemSearcher) HandleSearchProblem(w http.ResponseWriter, r *http.Request) {
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

		var params ProblemSearchParams
		if err := decoder.Decode(&params, query); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("ERROR: failed to decode request parameter `%s`: %s", r.URL.RawQuery, err.Error())
			encoder.Encode(NewErrorResponse(fmt.Sprintf("failed to decode request parameter `%s`", r.URL.RawQuery), nil))
			return
		}

		if err := validate.Struct(params); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Printf("ERROR: validation error: %+v, `%s`: %s", params, r.URL.RawQuery, err.Error())
			encoder.Encode(NewErrorResponse(fmt.Sprintf("validation error: %+v, `%s`: %s", params, r.URL.RawQuery, err.Error()), nil))
			return
		}

		code, res := searchProblem(s.core, params)
		w.WriteHeader(code)
		encoder.Encode(res)
	default:

	}
}

func searchProblem(core *solr.SolrCore[ProblemResponse, FacetCounts], params ProblemSearchParams) (int, common.SearchResultResponse[ProblemResponse]) {
	startTime := time.Now()

	query := params.ToQuery()
	res, err := core.Select(query)
	if err != nil {
		log.Printf("ERROR: failed to request to solr with %+v, from %+v: %s", query, params, err.Error())
		return 500, NewErrorResponse("internal error", params)
	}

	rows, _ := strconv.Atoi(query.Get("rows"))

	result := common.SearchResultResponse[ProblemResponse]{
		Stats: common.SearchResultStats{
			Time:   uint(time.Since(startTime).Milliseconds()),
			Total:  res.Response.NumFound,
			Index:  (res.Response.Start / uint(rows)) + 1,
			Count:  uint(len(res.Response.Docs)),
			Pages:  (res.Response.NumFound + uint(rows) - 1) / uint(rows),
			Params: &params,
			Facet:  NewFacetResponse(*res.FacetCounts),
		},
		Items: res.Response.Docs,
	}
	querylog := common.QueryLog{
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
