package user

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
	Keyword string        `validate:"lte=200" json:"keyword,omitempty" schema:"keyword"`
	Limit   uint          `validate:"lte=200" json:"limit,omitempty" schema:"limit"`
	Page    uint          `json:"page,omitempty" schema:"page"`
	Filter  *FilterParams `json:"filter,omitempty" schema:"filter"`
	Sort    string        `validate:"omitempty,oneof=-score rating -rating birth_year -birth_year" json:"sort,omitempty" schema:"sort"`
	Facet   []string      `validate:"dive,oneof=color birth_year join_count country" json:"facet,omitempty" schema:"facet"`
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

	return uint(int(p.Page)-1) * p.rows()
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

	if c := acs.SanitizeStrings(p.Filter.Country); len(c) != 0 {
		fq = append(fq, fmt.Sprintf("{!tag=country}country:(%s)", strings.Join(c, " OR ")))
	}
	if p.Filter.Rating != nil {
		if r := p.Filter.Rating.ToRange(); r != "" {
			fq = append(fq, fmt.Sprintf("{!tag=rating}rating:%s", r))
		}
	}
	if p.Filter.BirthYear != nil {
		if r := p.Filter.BirthYear.ToRange(); r != "" {
			fq = append(fq, fmt.Sprintf("{!tag=birth_year}birth_year:%s", p.Filter.BirthYear.ToRange()))
		}
	}
	if p.Filter.JoinCount != nil {
		if r := p.Filter.JoinCount.ToRange(); r != "" {
			fq = append(fq, fmt.Sprintf("{!tag=join_count}join_count:%s", p.Filter.JoinCount.ToRange()))
		}
	}

	return fq
}

type FilterParams struct {
	Rating    *acs.IntegerRange[int] `json:"rating,omitempty" schema:"rating"`
	BirthYear *acs.IntegerRange[int] `json:"birth_year,omitempty" schema:"birth_year"`
	JoinCount *acs.IntegerRange[int] `json:"join_count,omitempty" schema:"join_count"`
	Country   []string               `json:"country,omitempty" schema:"country"`
}

type Response struct {
	UserName      string  `json:"user_name"`
	Rating        int     `json:"rating"`
	HighestRating int     `json:"highest_rating"`
	Affiliation   *string `json:"affiliation"`
	BirthYear     *uint   `json:"birth_year"`
	Country       *string `json:"country"`
	Crown         *string `json:"crown"`
	JoinCount     uint    `json:"join_count"`
	Rank          uint    `json:"rank"`
	ActiveRank    *uint   `json:"active_rank"`
	Wins          uint    `json:"wins" `
	Color         string  `json:"color"`
}

type FacetCounts struct {
	Color     solr.TermFacetCount `json:"color"`
	BirthYear solr.TermFacetCount `json:"birth_year"`
	JoinCount solr.TermFacetCount `json:"join_count"`
	Country   solr.TermFacetCount `json:"country"`
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

func (f *FacetCounts) Into() FacetResponse {
	color := make([]FacetPart, len(f.Color.Buckets))
	for i, b := range f.Color.Buckets {
		color[i] = FacetPart{
			Label: b.Val,
			Count: b.Count,
		}
	}

	birthYear := make([]FacetPart, len(f.BirthYear.Buckets))
	for i, b := range f.BirthYear.Buckets {
		birthYear[i] = FacetPart{
			Label: b.Val,
			Count: b.Count,
		}
	}

	joinCount := make([]FacetPart, len(f.JoinCount.Buckets))
	for i, b := range f.JoinCount.Buckets {
		joinCount[i] = FacetPart{
			Label: b.Val,
			Count: b.Count,
		}
	}

	country := make([]FacetPart, len(f.Country.Buckets))
	for i, b := range f.Country.Buckets {
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
	core      *solr.Core
	validator *validator.Validate
	decoder   *schema.Decoder
}

func NewSearcher(baseURL string, coreName string) (Searcher, error) {
	core, err := solr.NewSolrCore(coreName, baseURL)
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
	slog.Info("querylog", slog.String("domain", "user"), slog.Uint64("elapsed_time", uint64(result.Stats.Time)), slog.Uint64("hits", uint64(res.Response.NumFound)), slog.Any("params", params))

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
