package submission

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

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"github.com/morikuni/failure"
)

type SearchParams struct {
	Limit  uint          `json:"limit,omitempty" schema:"limit" validate:"lte=200"`
	Page   uint          `json:"page,omitempty" schema:"page"`
	Filter *FilterParams `json:"filter,omitempty" schema:"filter"`
	Sort   string        `json:"sort,omitempty" schema:"sort" validate:"omitempty,oneof=-score execution_time -execution_time submitted_at -submitted_at point -point length -length"`
	Facet  []string      `json:"facet,omitempty" schema:"facet" validate:"dive,oneof=problem_id user_id language length result execution_time"`
}

func (p *SearchParams) ToQuery() url.Values {
	return solr.NewEDisMaxQueryBuilder().
		Facet(p.facet()).
		Fl(acs.FieldList(Response{})).
		Fq(p.fq()).
		Op("AND").
		Qf("text_unigram").
		Q("*:*").
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

	for _, f := range p.Facet {
		if f == "length" {
			facets[f] = map[string]any{
				"type":     "range",
				"field":    f,
				"mincount": 0,
				"start":    0,
				"end":      60000,
				"gap":      1000,
				"domain": map[string]any{
					"excludeTags": []string{f},
				},
			}
		} else if f == "execution_time" {
			facets[f] = map[string]any{
				"type":     "range",
				"field":    f,
				"mincount": 0,
				"start":    0,
				"end":      10000,
				"gap":      100,
				"other":    "after",
				"domain": map[string]any{
					"excludeTags": []string{f},
				},
			}
		} else {
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
	}

	facet, err := json.Marshal(facets)
	if err != nil {
		log.Printf("WARN: failed to marshal json.facet parameter from %+v", p.Facet)
		return ""
	}

	return string(facet)
}

func (p *SearchParams) fq() []string {
	fq := make([]string, 0)
	if p.Filter == nil {
		return fq
	}

	ids := make([]string, 0, len(p.Filter.SubmissionID))
	for _, i := range p.Filter.SubmissionID {
		ids = append(ids, strconv.Itoa(i))
	}
	if expr := strings.Join(ids, " OR "); expr != "" {
		fq = append(fq, fmt.Sprintf("{!tag=submission_id}submission_id:(%s)", expr))
	}

	if p.Filter.EpochSecond != nil {
		if r := p.Filter.EpochSecond.ToRange(); r != "" {
			fq = append(fq, fmt.Sprintf("{!tag=epoch_second}epoch_second:%s", r))
		}
	}
	if p.Filter.Point != nil {
		if r := p.Filter.Point.ToRange(); r != "" {
			fq = append(fq, fmt.Sprintf("{!tag=point}point:%s", r))
		}
	}
	if p.Filter.Length != nil {
		if r := p.Filter.Length.ToRange(); r != "" {
			fq = append(fq, fmt.Sprintf("{!tag=length}length:%s", r))
		}
	}
	if p.Filter.ExecutionTime != nil {
		if r := p.Filter.ExecutionTime.ToRange(); r != "" {
			fq = append(fq, fmt.Sprintf("{!tag=execution_time}execution_time:%s", r))
		}
	}

	if expr := strings.Join(acs.SanitizeStrings(p.Filter.ProblemID), " OR "); expr != "" {
		fq = append(fq, fmt.Sprintf("{!tag=problem_id}problem_id:(%s)", expr))
	}
	if expr := strings.Join(acs.SanitizeStrings(p.Filter.ContestID), " OR "); expr != "" {
		fq = append(fq, fmt.Sprintf("{!tag=contest_id}contest_id:(%s)", expr))
	}
	if expr := strings.Join(acs.SanitizeStrings(p.Filter.Category), " OR "); expr != "" {
		fq = append(fq, fmt.Sprintf("{!tag=category}category:(%s)", expr))
	}
	if expr := strings.Join(acs.SanitizeStrings(p.Filter.UserID), " OR "); expr != "" {
		fq = append(fq, fmt.Sprintf("{!tag=user_id}user_id:(%s)", expr))
	}
	if expr := strings.Join(acs.SanitizeStrings(p.Filter.Language), " OR "); expr != "" {
		fq = append(fq, fmt.Sprintf("{!tag=language}language:(%s)", expr))
	}
	if expr := strings.Join(acs.SanitizeStrings(p.Filter.Result), " OR "); expr != "" {
		fq = append(fq, fmt.Sprintf("{!tag=result}result:(%s)", expr))
	}

	return fq
}

type FilterParams struct {
	SubmissionID  []int                    `json:"submission_id,omitempty" schema:"submission_id"`
	EpochSecond   *acs.IntegerRange[int]   `json:"epoch_second,omitempty" schema:"epoch_second"`
	ProblemID     []string                 `json:"problem_id,omitempty" schema:"problem_id"`
	ContestID     []string                 `json:"contest_id,omitempty" schema:"contest_id"`
	Category      []string                 `json:"category,omitempty" schema:"category"`
	UserID        []string                 `json:"user_id,omitempty" schema:"user_id"`
	Language      []string                 `json:"language,omitempty" schema:"language"`
	Point         *acs.FloatRange[float64] `json:"point,omitempty" schema:"point"`
	Length        *acs.IntegerRange[int]   `json:"length,omitempty" schema:"length"`
	Result        []string                 `json:"result,omitempty" schema:"result"`
	ExecutionTime *acs.IntegerRange[int]   `json:"execution_time,omitempty" schema:"execution_time"`
}

type Response struct {
	SubmissionID  int64                 `json:"submission_id" solr:"submission_id"`
	SubmittedAt   solr.FromSolrDateTime `json:"submitted_at" solr:"submitted_at"`
	ProblemID     string                `json:"problem_id" solr:"problem_id"`
	ContestID     string                `json:"contest_id" solr:"contest_id"`
	UserID        string                `json:"user_id" solr:"user_id"`
	Language      string                `json:"language" solr:"language"`
	Point         float64               `json:"point" solr:"point"`
	Length        uint64                `json:"length" solr:"length"`
	Result        string                `json:"result" solr:"result"`
	ExecutionTime *uint64               `json:"execution_time" solr:"execution_time"`
}

type FacetCounts struct {
	ProblemID     solr.SolrTermFacetCount       `json:"problem_id,omitempty"`
	UserID        solr.SolrTermFacetCount       `json:"user_id,omitempty"`
	Language      solr.SolrTermFacetCount       `json:"language,omitempty"`
	Result        solr.SolrTermFacetCount       `json:"result,omitempty"`
	Length        solr.SolrRangeFacetCount[int] `json:"length,omitempty"`
	ExecutionTime solr.SolrRangeFacetCount[int] `json:"execution_time,omitempty"`
}

type Searcher struct {
	core      *solr.SolrCore[Response, FacetCounts]
	validator *validator.Validate
	decoder   *schema.Decoder
}

func NewSearcher(baseURL string, coreName string) (Searcher, error) {
	core, err := solr.NewSolrCore[Response, FacetCounts](coreName, baseURL)
	if err != nil {
		return Searcher{}, failure.Translate(err, SearcherInitializeError, failure.Context{"baseURL": baseURL, "coreName": coreName}, failure.Message("failed to create user searcher"))
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
			encoder.Encode(NewErrorResponse(fmt.Sprintf("validation error, `%s`: %s", r.URL.RawQuery, err.Error()), params))
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
			Params: &params,
			Facet:  res.FacetCounts,
		},
		Items: res.Response.Docs,
	}
	querylog := acs.QueryLog{
		RequestAt: startTime,
		Domain:    "submissions",
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