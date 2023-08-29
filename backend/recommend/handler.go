package recommend

import (
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
	ProblemID string        `json:"problem_id,omitempty" schema:"problem_id" validate:"required"`
	Limit     int           `json:"limit,omitempty" schema:"limit" validate:"lte=200"`
	Page      int           `json:"page,omitempty" schema:"page"`
	Filter    *FilterParams `json:"filter,omitempty" schema:"filter"`
	Sort      string        `json:"sort,omitempty" schema:"sort" validate:"omitempty,oneof=-score"`
}

func (p *SearchParams) ToQuery() url.Values {
	return solr.NewEDisMaxQueryBuilder().
		Bq(p.bq()).
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
		return "score desc"
	}
	if strings.HasPrefix(p.Sort, "-") {
		return fmt.Sprintf("%s desc", p.Sort[1:])
	} else {
		return fmt.Sprintf("%s asc", p.Sort)
	}
}

func (p *SearchParams) fq() []string {
	fq := make([]string, 0)
	if p.Filter == nil {
		return fq
	}

	if !p.Filter.IncludeExperimental {
		fq = append(fq, "is_experimental:false")
	}
	if expr := strings.Join(acs.SanitizeStrings(p.Filter.Category), " OR "); expr != "" {
		fq = append(fq, fmt.Sprintf("category:(%s)", expr))
	}

	return fq
}

func (p *SearchParams) bq() []string {
	bq := make([]string, 0)

	bq = append(bq, fmt.Sprintf("{!boost b=10}{!join fromIndex=recommend from=problem_id to=problem_id score=max}{!payload_score f=difficulty_correlation func=sum operator=or includeSpanScore=false v=%s}", p.ProblemID))
	bq = append(bq, fmt.Sprintf("{!boost b=2}{!join fromIndex=recommend from=problem_id to=problem_id score=max}{!payload_score f=category_correlation func=sum operator=or includeSpanScore=false v=%s}", p.ProblemID))
	bq = append(bq, "{!join fromIndex=recommend from=problem_id to=problem_id score=max}{!func v=log(add(solved_count,1))}")
	bq = append(bq, "{!func}recip(ms(NOW,start_at),3.16e-11,2,1)")

	return bq
}

type FilterParams struct {
	IncludeExperimental bool     `json:"include_experimental,omitempty" schema:"include_experimental"`
	Category            []string `json:"category,omitempty" schema:"category"`
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
	SolvedCount  int                   `json:"solved_count"`
	Score        float64               `json:"score"`
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
	res, err := solr.Select[Response, any](s.core, query)
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
			Params: &params,
			Facet:  res.FacetCounts,
		},
		Items: res.Response.Docs,
	}
	slog.Info("querylog", slog.String("domain", "recommend"), slog.Int("elapsed_time", result.Stats.Time), slog.Int("hits", res.Response.NumFound), slog.Any("params", params))

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
