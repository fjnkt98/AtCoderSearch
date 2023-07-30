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
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
)

var validate = validator.New()
var decoder = schema.NewDecoder()

type ProblemSearchParams struct {
	Keyword string        `validate:"lte=200" json:"keyword,omitempty"`
	Limit   uint          `validate:"lte=200" json:"limit,omitempty"`
	Page    uint          `json:"page,omitempty"`
	Filter  *FilterParams `json:"filter,omitempty"`
	Sort    string        `validate:"omitempty,oneof=-score start_at -start_at difficulty -difficulty" json:"sort,omitempty"`
	Facet   []string      `validate:"dive,oneof=category difficulty" json:"facet,omitempty"`
}

type FilterParams struct {
	Category   []string                `json:"category,omitempty"`
	Difficulty common.RangeFilterParam `json:"difficulty,omitempty"`
}

type ProblemResponse struct {
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
	Count      uint
	Category   solr.SolrTermFacetCount
	Difficulty solr.SolrTermFacetCount
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

		// code, res := searchProblem(s.core, params)
		code := 200
		w.WriteHeader(code)
		encoder.Encode(params)
	default:

	}
}

// func searchProblem(core *solr.SolrCore[ProblemResponse, FacetCounts], params ProblemSearchParams) (int, common.SearchResultResponse[ProblemResponse]) {
// 	startTime := time.Now()

// 	query := params.ToQuery()
// 	res, err := core.Select(query)
// 	if err != nil {
// 		msg := fmt.Sprintf("ERROR: failed to request to solr with %+v, from %+v: %s", query, params, err.Error())
// 		log.Printf(msg)
// 		return 500, NewErrorResponse(msg, params)
// 	}

// 	rows, _ := strconv.Atoi(query.Get("rows"))

// 	result := common.SearchResultResponse[ProblemResponse]{
// 		Stats: common.SearchResultStats{
// 			Time:   uint(time.Since(startTime).Milliseconds()),
// 			Total:  res.Response.NumFound,
// 			Index:  (res.Response.Start / uint(rows)) + 1,
// 			Count:  uint(len(res.Response.Docs)),
// 			Pages:  (res.Response.NumFound + uint(rows) - 1) / uint(rows),
// 			Params: &params,
// 			Facet:  res.FacetCounts,
// 		},
// 		Items: res.Response.Docs,
// 	}
// 	querylog := common.QueryLog{
// 		Domain: "problem",
// 		Time:   result.Stats.Time,
// 		Hits:   res.Response.NumFound,
// 		Params: params,
// 	}
// 	encoder := json.NewEncoder(log.Writer())
// 	encoder.Encode(querylog)

// 	return http.StatusOK, result
// }
