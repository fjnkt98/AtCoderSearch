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
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/morikuni/failure"
	"golang.org/x/exp/slog"
)

const (
	RECENT  = 1
	RATING  = 2
	HISTORY = 3
)

type SearchParams struct {
	Model    int    `json:"model" schema:"model" validate:"required,model"`
	Option   string `json:"option" schema:"option" validate:"omitempty,option"`
	UserID   string `json:"user_id" schema:"user_id"`
	Rating   int    `json:"rating" schema:"rating"`
	Limit    *int   `json:"limit" schema:"limit" validate:"omitempty,lte=200"`
	Page     int    `json:"page" schema:"page"`
	Unsolved bool   `json:"unsolved" schema:"unsolved"`
}

func (p *SearchParams) ToQuery() url.Values {
	return solr.NewEDisMaxQueryBuilder().
		Bq(p.bq()).
		Fq(p.fq()).
		Fl(acs.FieldList(Response{})).
		QAlt("*:*^=0").
		Rows(p.rows()).
		Start(p.start()).
		Sort("score desc,problem_id asc").
		Build()
}

func (p *SearchParams) rows() int {
	if p.Limit == nil {
		return 0
	}
	return *p.Limit
}

func (p *SearchParams) start() int {
	if p.Page == 0 || p.rows() == 0 {
		return 0
	}

	return (p.Page - 1) * p.rows()
}

func (p *SearchParams) fq() []string {
	fq := make([]string, 1)
	if p.Unsolved {
		fq = append(fq, fmt.Sprintf(`-{!join fromIndex=submission from=problem_id to=problem_id v="+user_id:%s +result:AC"}`, solr.Sanitize(p.UserID)))
	}
	return fq
}

type Weights struct {
	Trend           int
	Difficulty      int
	ABC             int
	ARC             int
	AGC             int
	Other           int
	NotExperimental int
}

func (p *SearchParams) bq() []string {
	bq := make([]string, 0)

	var w Weights
	var rate int

	switch p.Model {
	case RECENT:
		w = Weights{Trend: 10}
	case RATING:
		w = Weights{Trend: 3, Difficulty: 10, ABC: 5, ARC: 5, AGC: 5, Other: 1, NotExperimental: 0}

		if p.Option != "" {
			opt := []rune(p.Option)

			switch opt[0] {
			case '0':
				rate = p.Rating - 200
			case '1':
				rate = p.Rating
			case '2':
				rate = p.Rating + 200
			}

			switch opt[1] {
			case '1':
				w.ABC = 16
				w.ARC = 4
				w.AGC = 2
			case '2':
				w.ABC = 2
				w.ARC = 16
				w.AGC = 4
			case '3':
				w.ABC = 2
				w.ARC = 4
				w.AGC = 16
			default:
			}

			switch opt[2] {
			case '0':
				w.Trend = 3
			case '1':
				w.Trend = 7
			}

			switch opt[3] {
			case '0':
				w.NotExperimental = 0
			case '1':
				w.NotExperimental = 10
			}
		}
	case HISTORY:
		// TODO
		w = Weights{Trend: 10}
	}

	bq = append(bq, fmt.Sprintf("{!boost b=%d}{!func}pow(2,mul(-1,div(ms(NOW,start_at),2592000000)))", w.Trend))
	bq = append(bq, fmt.Sprintf("{!boost b=%d}{!func}pow(2.71828182846,mul(-1,div(pow(sub(%d,difficulty),2),20000)))", w.Difficulty, rate))
	bq = append(bq, fmt.Sprintf(`{!boost b=%d}(category:"ABC" OR category:"ABC-Like"^0.5)`, w.ABC))
	bq = append(bq, fmt.Sprintf(`{!boost b=%d}(category:"ARC" OR category:"ARC-Like"^0.5)`, w.ARC))
	bq = append(bq, fmt.Sprintf(`{!boost b=%d}(category:"AGC" OR category:"AGC-Like"^0.5)`, w.AGC))
	bq = append(bq, fmt.Sprintf(`{!boost b=%d}category:("JOI" OR "Other Sponsored" OR "Other Contests" OR "PAST")`, w.Other))
	bq = append(bq, fmt.Sprintf(`{!boost b=%d}is_experimental:false`, w.NotExperimental))
	bq = append(bq, "{!boost b=0.2}{!join fromIndex=recommend from=problem_id to=problem_id score=max}{!func v=log(add(solved_count,1))}")

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
	Score        float64               `json:"score"`
}

type Searcher struct {
	core      *solr.Core
	db        *sqlx.DB
	validator *validator.Validate
	decoder   *schema.Decoder
}

func ValidateModel(fl validator.FieldLevel) bool {
	if !fl.Field().CanInt() {
		return false
	}

	if model := fl.Field().Int(); model == RECENT || model == RATING || model == HISTORY {
		return true
	}
	return false
}

func ValidateOption(fl validator.FieldLevel) bool {
	s := fl.Field().String()

	if _, err := strconv.Atoi(s); err != nil {
		return false
	}

	opt := []rune(s)
	if len(opt) != 4 {
		return false
	}

	if !('0' <= opt[0] && opt[0] <= '2') {
		return false
	}
	if !('0' <= opt[1] && opt[1] <= '3') {
		return false
	}
	if !('0' <= opt[2] && opt[2] <= '1') {
		return false
	}
	if !('0' <= opt[3] && opt[3] <= '1') {
		return false
	}

	return true
}

func NewSearcher(baseURL string, coreName string, db *sqlx.DB) (Searcher, error) {
	core, err := solr.NewSolrCore(coreName, baseURL)
	if err != nil {
		return Searcher{}, failure.Translate(err, acs.SearcherInitializeError, failure.Context{"baseURL": baseURL, "coreName": coreName}, failure.Message("failed to create user searcher"))
	}

	validator := validator.New()
	validator.RegisterValidation("model", ValidateModel)
	validator.RegisterValidation("option", ValidateOption)
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	decoder.RegisterConverter([]string{}, func(input string) reflect.Value {
		return reflect.ValueOf(strings.Split(input, ","))
	})

	searcher := Searcher{
		core:      core,
		db:        db,
		validator: validator,
		decoder:   decoder,
	}
	return searcher, nil
}

func NewErrorResponse(msg string, params any) acs.SearchResultResponse[Response] {
	return acs.NewErrorResponse[Response](msg, params)
}

func (s *Searcher) getRating(userID string) (int, error) {
	row := s.db.QueryRow(`
	SELECT
		"rating"
	FROM
		"users"
	WHERE
		"user_name" = $1::text
	`,
		userID,
	)
	var rating int
	if err := row.Scan(&rating); err != nil {
		return 0, failure.Translate(err, acs.DBError, failure.Context{"user_id": userID}, failure.Messagef("failed to get rating of the user `%s`", userID))
	}

	return rating, nil
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

func (s *Searcher) HandlePOST(c echo.Context) error {
	var params SearchParams
	if err := c.Bind(&params); err != nil {
		slog.Error("failed to decode request parameter", slog.String("uri", c.Request().RequestURI), slog.String("error", fmt.Sprintf("%+v", err)))
		return c.JSON(http.StatusBadRequest, NewErrorResponse("failed to decode request parameter", nil))
	}

	if err := s.validator.Struct(params); err != nil {
		slog.Error("validation error", slog.String("uri", c.Request().RequestURI), slog.Any("params", params), slog.String("error", fmt.Sprintf("%+v", err)))
		return c.JSON(http.StatusBadRequest, NewErrorResponse(fmt.Sprintf("validation error: %s", err.Error()), params))
	}

	code, res := s.search(c.Request(), params)
	return c.JSON(code, res)
}

func (s *Searcher) search(r *http.Request, params SearchParams) (int, acs.SearchResultResponse[Response]) {
	startTime := time.Now()

	var msg string
	if params.UserID != "" {
		if rating, err := s.getRating(params.UserID); err != nil {
			msg = "invalid user id"
		} else {
			params.Rating = rating
		}
	}

	query := params.ToQuery()
	res, err := solr.Select[Response, any](s.core, query)
	if err != nil {
		slog.Error("failed to request to solr", slog.String("url", r.URL.String()), slog.Any("query", query), slog.Any("params", params), slog.String("error", fmt.Sprintf("%+v", err)))
		return 500, NewErrorResponse("internal error", params)
	}

	rows, _ := strconv.Atoi(query.Get("rows"))
	var pages = 0
	var index = 0
	if rows != 0 {
		pages = (res.Response.NumFound + rows) / rows
		index = (res.Response.Start / rows) + 1
	}

	result := acs.SearchResultResponse[Response]{
		Stats: acs.SearchResultStats{
			Time:   int(time.Since(startTime).Milliseconds()),
			Total:  res.Response.NumFound,
			Index:  index,
			Count:  len(res.Response.Docs),
			Pages:  pages,
			Params: &params,
		},
		Items:   res.Response.Docs,
		Message: msg,
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
