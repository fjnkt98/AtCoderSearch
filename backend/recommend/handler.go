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
	Model int `json:"model" schema:"model" validate:"required,model"`
	Rate  int `json:"rate" schema:"rate" validate:"required"`
	Limit int `json:"limit" schema:"limit" validate:"lte=200"`
	Page  int `json:"page" schema:"page"`
}

func (p *SearchParams) ToQuery() url.Values {
	return solr.NewEDisMaxQueryBuilder().
		Bq(p.bq()).
		Fl(acs.FieldList(Response{})).
		QAlt("*:*^=0").
		Rows(p.rows()).
		Start(p.start()).
		Sort("score desc,problem_id asc").
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

	return (p.Page - 1) * p.rows()
}

type Weights struct {
	Time       int
	Difficulty int
	ABC        int
	ARC        int
	AGC        int
	Other      int
}

func (p *SearchParams) bq() []string {
	bq := make([]string, 0)

	var w Weights
	var rate int

	switch p.Model {
	case 1000:
		w = Weights{Time: 10}
	case 3000:
		// TODO
		w = Weights{Time: 10}
	default:
		model := []rune(strconv.Itoa(p.Model))
		w = Weights{Time: 3, Difficulty: 10, ABC: 5, ARC: 5, AGC: 5, Other: 1}

		switch model[1] {
		case '0':
			rate = p.Rate - 200
		case '1':
			rate = p.Rate
		case '2':
			rate = p.Rate + 200
		}

		switch model[2] {
		case '0':
			w.ABC = 8
			w.ARC = 4
			w.AGC = 2
		case '1':
			w.ABC = 2
			w.ARC = 8
			w.AGC = 4
		case '2':
			w.ABC = 2
			w.ARC = 4
			w.AGC = 8
		}

		switch model[3] {
		case '0':
			w.Time = 3
		case '1':
			w.Time = 5
		case '2':
			w.Time = 7
		}
	}

	bq = append(bq, fmt.Sprintf("{!boost b=%d}{!func}pow(2,mul(-1,div(ms(NOW,start_at),2592000000)))", w.Time))
	bq = append(bq, fmt.Sprintf("{!boost b=%d}{!func}pow(2.71828182846,mul(-1,div(pow(sub(%d,difficulty),2),20000)))", w.Difficulty, rate))
	bq = append(bq, fmt.Sprintf(`{!boost b=%d}(category:"ABC" OR category:"ABC-Like"^0.5)`, w.ABC))
	bq = append(bq, fmt.Sprintf(`{!boost b=%d}(category:"ARC" OR category:"ARC-Like"^0.5)`, w.ARC))
	bq = append(bq, fmt.Sprintf(`{!boost b=%d}(category:"AGC" OR category:"AGC-Like"^0.5)`, w.AGC))
	bq = append(bq, fmt.Sprintf(`{!boost b=%d}category:("JOI" OR "Other Sponsored" OR "Other Contests" OR "PAST")`, w.Other))
	bq = append(bq, "{!boost b=1}{!join fromIndex=recommend from=problem_id to=problem_id score=max}{!func v=log(add(solved_count,1))}")

	return bq
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

func ValidateModel(fl validator.FieldLevel) bool {
	if !fl.Field().CanInt() {
		return false
	}

	if model := fl.Field().Int(); model == 1000 || (2000 <= model && model <= 2221) || model == 3000 {
		return true
	}
	return false
}

func NewSearcher(baseURL string, coreName string) (Searcher, error) {
	core, err := solr.NewSolrCore(coreName, baseURL)
	if err != nil {
		return Searcher{}, failure.Translate(err, SearcherInitializeError, failure.Context{"baseURL": baseURL, "coreName": coreName}, failure.Message("failed to create user searcher"))
	}

	validator := validator.New()
	validator.RegisterValidation("model", ValidateModel)
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
