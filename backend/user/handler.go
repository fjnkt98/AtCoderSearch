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
	"github.com/labstack/echo/v4"
	"github.com/morikuni/failure"
)

type SearchParams struct {
	Keyword string       `json:"keyword" schema:"keyword" validate:"lte=200"`
	Limit   int          `json:"limit" schema:"limit" validate:"lte=200"`
	Page    int          `json:"page" schema:"page"`
	Filter  FilterParams `json:"filter" schema:"filter"`
	Sort    string       `json:"sort" schema:"sort" validate:"omitempty,oneof=-score rating -rating birth_year -birth_year"`
	Facet   FacetParams  `json:"facet" schema:"facet"`
}

type FilterParams struct {
	Rating    acs.IntegerRange `json:"rating" schema:"rating"`
	BirthYear acs.IntegerRange `json:"birth_year" schema:"birth_year"`
	JoinCount acs.IntegerRange `json:"join_count" schema:"join_count"`
	Country   []string         `json:"country" schema:"country"`
	Color     []string         `json:"color" schema:"color"`
}

type FacetParams struct {
	Term      []string            `json:"term" schema:"term" validate:"dive,oneof=color country"`
	Rating    acs.RangeFacetParam `json:"rating" schema:"rating"`
	BirthYear acs.RangeFacetParam `json:"birth_year" schema:"birth_year"`
	JoinCount acs.RangeFacetParam `json:"join_count" schema:"join_count"`
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
		return "rating desc"
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

	if f := p.Facet.Rating.ToFacet("rating"); f != nil {
		facets["rating"] = f
	}
	if f := p.Facet.BirthYear.ToFacet("birth_year"); f != nil {
		facets["birth_year"] = f
	}
	if f := p.Facet.JoinCount.ToFacet("join_count"); f != nil {
		facets["join_count"] = f
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

	if c := acs.SanitizeStrings(p.Filter.Country); len(c) != 0 {
		fq = append(fq, fmt.Sprintf("{!tag=country}country:(%s)", strings.Join(c, " OR ")))
	}
	if c := acs.SanitizeStrings(p.Filter.Color); len(c) != 0 {
		fq = append(fq, fmt.Sprintf("{!tag=color}color:(%s)", strings.Join(c, " OR ")))
	}

	if r := p.Filter.Rating.ToRange(); r != "" {
		fq = append(fq, fmt.Sprintf("{!tag=rating}rating:%s", r))
	}
	if r := p.Filter.BirthYear.ToRange(); r != "" {
		fq = append(fq, fmt.Sprintf("{!tag=birth_year}birth_year:%s", p.Filter.BirthYear.ToRange()))
	}
	if r := p.Filter.JoinCount.ToRange(); r != "" {
		fq = append(fq, fmt.Sprintf("{!tag=join_count}join_count:%s", p.Filter.JoinCount.ToRange()))
	}

	return fq
}

type Response struct {
	UserName      string  `json:"user_name"`
	Rating        int     `json:"rating"`
	HighestRating int     `json:"highest_rating"`
	Affiliation   *string `json:"affiliation"`
	BirthYear     *int    `json:"birth_year"`
	Country       *string `json:"country"`
	Crown         *string `json:"crown"`
	JoinCount     int     `json:"join_count"`
	Rank          int     `json:"rank"`
	ActiveRank    *int    `json:"active_rank"`
	Wins          int     `json:"wins" `
	Color         string  `json:"color"`
	UserURL       string  `json:"user_url"`
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
	Count int    `json:"count"`
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
			Facet:  res.FacetCounts.Into(),
		},
		Items: res.Response.Docs,
	}
	slog.Info("querylog", slog.String("domain", "user"), slog.Int("elapsed_time", result.Stats.Time), slog.Int("hits", res.Response.NumFound), slog.Any("params", params))

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
