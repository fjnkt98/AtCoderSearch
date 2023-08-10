package user

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
)

type SearchParams struct {
	Keyword string        `validate:"lte=200" json:"keyword,omitempty"`
	Limit   uint          `validate:"lte=200" json:"limit,omitempty"`
	Page    uint          `json:"page,omitempty"`
	Filter  *FilterParams `json:"filter,omitempty"`
	Sort    string        `validate:"omitempty,oneof=-score rating -rating birth_year -birth_year" json:"sort,omitempty"`
	Facet   []string      `validate:"dive,oneof=color birth_year join_count country" json:"facet,omitempty"`
}

func (p *SearchParams) ToQuery() url.Values {
	return solr.NewEDisMaxQueryBuilder().
		Facet(p.facet()).
		Fl(acs.FieldList(Response{})).
		Fq(p.fq()).
		Op("AND").
		Q(solr.Sanitize(norm.NFKC.String(p.Keyword))).
		QAlt("*:*").
		Qf("text_unigram").
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
		return "rating desc"
	}
	if strings.HasPrefix(p.Sort, "-") {
		return fmt.Sprintf("%s desc", p.Sort[1:])
	} else {
		return fmt.Sprintf("%s asc", p.Sort)
	}
}

var FACET_MAP = map[string]string{
	"birth_year": "period",
	"join_count": "join_count_grade",
}

func (p *SearchParams) facet() string {
	facets := make(map[string]any)

	for _, f := range p.Facet {
		field, ok := FACET_MAP[f]
		if !ok {
			field = f
		}
		facets[f] = map[string]any{
			"type":     "terms",
			"field":    field,
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

func (p *SearchParams) fq() []string {
	if p.Filter == nil {
		return make([]string, 0)
	}

	fq := make([]string, 0)

	countries := make([]string, 0, len(p.Filter.Country))
	for _, c := range p.Filter.Country {
		category := solr.Sanitize(c)
		if c == "" {
			continue
		}
		countries = append(countries, category)
	}

	fq = append(fq, fmt.Sprintf("{!tag=country}country:(%s)", strings.Join(countries, " OR ")))
	fq = append(fq, fmt.Sprintf("{!tag=rating}rating:%s", p.Filter.Rating.ToRange()))
	fq = append(fq, fmt.Sprintf("{!tag=birth_year}birth_year:%s", p.Filter.BirthYear.ToRange()))
	fq = append(fq, fmt.Sprintf("{!tag=join_count}join_count:%s", p.Filter.JoinCount.ToRange()))

	return fq
}

type FilterParams struct {
	Rating    acs.RangeFilterParam `json:"rating,omitempty"`
	BirthYear acs.RangeFilterParam `json:"birth_year,omitempty"`
	JoinCount acs.RangeFilterParam `json:"join_count,omitempty"`
	Country   []string             `json:"country,omitempty"`
}

type Response struct {
	UserName      string  `json:"user_name" solr:"user_name"`
	Rating        int     `json:"rating" solr:"rating"`
	HighestRating int     `json:"highest_rating" solr:"highest_rating"`
	Affiliation   *string `json:"affiliation" solr:"affiliation"`
	BirthYear     *uint   `json:"birth_year" solr:"birth_year"`
	Country       *string `json:"country" solr:"country"`
	Crown         *string `json:"crown" solr:"crown"`
	JoinCount     uint    `json:"join_count" solr:"join_count"`
	Rank          uint    `json:"rank" solr:"rank"`
	ActiveRank    *uint   `json:"active_rank" solr:"active_rank"`
	Wins          uint    `json:"wins"  solr:"wins"`
	Color         string  `json:"color" solr:"color"`
}

type FacetCounts struct {
	Color     solr.SolrTermFacetCount `json:"color"`
	BirthYear solr.SolrTermFacetCount `json:"birth_year"`
	JoinCount solr.SolrTermFacetCount `json:"join_count"`
	Country   solr.SolrTermFacetCount `json:"country"`
}

type FacetResponse struct {
	Color     []FacetPart `json:"color,omitempty"`
	BirthYear []FacetPart `json:"birth_year,omitempty"`
	JoinCount []FacetPart `json:"join_count,omitempty"`
	Country   []FacetPart `json:"country,omitempty"`
}

type FacetPart struct {
	Label string `json:"label"`
	Count uint   `json:"count"`
}

func NewFacetResponse(facet FacetCounts) FacetResponse {
	color := make([]FacetPart, len(facet.Color.Buckets))
	for i, b := range facet.Color.Buckets {
		color[i] = FacetPart{
			Label: b.Val,
			Count: b.Count,
		}
	}

	birthYear := make([]FacetPart, len(facet.BirthYear.Buckets))
	for i, b := range facet.BirthYear.Buckets {
		birthYear[i] = FacetPart{
			Label: b.Val,
			Count: b.Count,
		}
	}

	joinCount := make([]FacetPart, len(facet.JoinCount.Buckets))
	for i, b := range facet.JoinCount.Buckets {
		joinCount[i] = FacetPart{
			Label: b.Val,
			Count: b.Count,
		}
	}

	country := make([]FacetPart, len(facet.Country.Buckets))
	for i, b := range facet.Country.Buckets {
		country[i] = FacetPart{
			Label: b.Val,
			Count: b.Count,
		}
	}

	return FacetResponse{
		Color:     color,
		BirthYear: birthYear,
		JoinCount: joinCount,
		Country:   country,
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
		return Searcher{}, fmt.Errorf("failed to create user searcher: %w", err)
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
			encoder.Encode(NewErrorResponse(fmt.Sprintf("validation error: %+v, `%s`: %s", params, r.URL.RawQuery, err.Error()), nil))
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
			Facet:  NewFacetResponse(*res.FacetCounts),
		},
		Items: res.Response.Docs,
	}
	querylog := acs.QueryLog{
		RequestAt: startTime,
		Domain:    "user",
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