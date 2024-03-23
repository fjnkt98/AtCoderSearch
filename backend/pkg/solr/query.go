package solr

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/goark/errs"
)

type SelectQuery struct {
	client *http.Client
	uri    *url.URL
	params url.Values
	facet  *JSONFacetQuery
}

func NewSelectQuery(client *http.Client, uri *url.URL) *SelectQuery {
	params := url.Values{}
	params.Set("defType", "edismax")
	return &SelectQuery{
		client: client,
		uri:    uri,
		params: params,
	}
}

func (q *SelectQuery) Exec(ctx context.Context) (*SelectResult, error) {
	q.uri.RawQuery = q.params.Encode()
	u := q.uri.String()
	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, errs.New(
			"failed to prepare select request",
			errs.WithCause(err),
			errs.WithContext("uri", u),
		)
	}
	res, err := q.client.Do(req)
	if err != nil {
		return nil, errs.New(
			"failed to execute select request",
			errs.WithCause(err),
			errs.WithContext("uri", u),
		)
	}
	defer res.Body.Close()
	var body SelectResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, errs.New(
			"failed to decode select response",
			errs.WithCause(err),
			errs.WithContext("uri", u),
		)
	}

	result := &SelectResult{
		query: q,
		Raw:   body,
	}

	if res.StatusCode != http.StatusOK {
		return result, errs.Wrap(ErrNonOKResponse, errs.WithContext("uri", u))
	}
	return result, nil
}

func (q *SelectQuery) Sort(sort string) *SelectQuery {
	if sort != "" {
		q.params.Set("sort", sort)
	}
	return q
}

func (q *SelectQuery) Start(start int) *SelectQuery {
	q.params.Set("start", strconv.Itoa(start))
	return q
}

func (q *SelectQuery) Rows(rows int) *SelectQuery {
	q.params.Set("rows", strconv.Itoa(rows))
	return q
}

func (q *SelectQuery) Fq(fq []string) *SelectQuery {
	for _, fq := range fq {
		q.params.Add("fq", fq)
	}
	return q
}

func (q *SelectQuery) Fl(fl string) *SelectQuery {
	if fl != "" {
		q.params.Set("fl", fl)
	}
	return q
}

func (q *SelectQuery) Wt(wt string) *SelectQuery {
	if wt != "" {
		q.params.Set("wt", wt)
	}
	return q
}

func (q *SelectQuery) JsonFacet(facet *JSONFacetQuery) *SelectQuery {
	f, err := json.Marshal(facet)
	if err == nil {
		q.params.Set("json.facet", string(f))
	}
	return q
}

func (q *SelectQuery) Op(op string) *SelectQuery {
	if op != "" {
		q.params.Set("q.op", op)
	}
	return q
}

func (q *SelectQuery) Df(df string) *SelectQuery {
	if df != "" {
		q.params.Set("df", df)
	}
	return q
}

func (sq *SelectQuery) Q(q string) *SelectQuery {
	if q != "" {
		sq.params.Set("q", q)
	}
	return sq
}

func (q *SelectQuery) Qf(qf string) *SelectQuery {
	if qf != "" {
		q.params.Set("qf", qf)
	}
	return q
}

func (q *SelectQuery) Qs(qs int) *SelectQuery {
	q.params.Set("qs", strconv.Itoa(qs))
	return q
}

func (q *SelectQuery) Pf(pf string) *SelectQuery {
	if pf != "" {
		q.params.Set("pf", pf)
	}
	return q
}

func (q *SelectQuery) Mm(mm string) *SelectQuery {
	if mm != "" {
		q.params.Set("mm", mm)
	}
	return q
}

func (sq *SelectQuery) QAlt(q string) *SelectQuery {
	if q != "" {
		sq.params.Set("q.alt", q)
	}
	return sq
}

func (q *SelectQuery) Tie(tie float64) *SelectQuery {
	q.params.Set("tie", strconv.FormatFloat(tie, 'f', 6, 64))
	return q
}

func (q *SelectQuery) Bq(bq []string) *SelectQuery {
	for _, bq := range bq {
		q.params.Add("bq", bq)
	}
	return q
}

func (q *SelectQuery) Bf(bf []string) *SelectQuery {
	for _, bf := range bf {
		q.params.Add("bf", bf)
	}
	return q
}

func (q *SelectQuery) Sow(sow bool) *SelectQuery {
	if sow {
		q.params.Set("sow", "true")
	} else {
		q.params.Set("sow", "false")
	}
	return q
}

func (q *SelectQuery) Boost(boost []string) *SelectQuery {
	for _, boost := range boost {
		q.params.Add("boost", boost)
	}
	return q
}

func (q *SelectQuery) LowerCaseOperators(flag bool) *SelectQuery {
	if flag {
		q.params.Set("lowercaseOperators", "true")
	} else {
		q.params.Set("lowercaseOperators", "false")
	}
	return q
}

func (q *SelectQuery) Pf2(pf2 string) *SelectQuery {
	if pf2 != "" {
		q.params.Set("pf2", pf2)
	}
	return q
}

func (q *SelectQuery) Ps2(ps2 int) *SelectQuery {
	q.params.Set("ps2", strconv.Itoa(ps2))
	return q
}

func (q *SelectQuery) Pf3(pf3 string) *SelectQuery {
	if pf3 != "" {
		q.params.Set("pf3", pf3)
	}
	return q
}

func (q *SelectQuery) Ps3(ps3 int) *SelectQuery {
	q.params.Set("ps3", strconv.Itoa(ps3))
	return q
}

func (q *SelectQuery) StopWords(flag bool) *SelectQuery {
	if flag {
		q.params.Set("stopwords", "true")
	} else {
		q.params.Set("stopwords", "false")
	}
	return q
}

func (q *SelectQuery) Uf(uf string) *SelectQuery {
	if uf != "" {
		q.params.Set("uf", uf)
	}
	return q
}

func (q *SelectQuery) Some(key, value string) *SelectQuery {
	if key != "" && value != "" {
		q.params.Add(key, value)
	}
	return q
}

type SelectResult struct {
	query *SelectQuery
	Raw   SelectResponse
}

func (r *SelectResult) Scan(v any) error {
	if err := json.Unmarshal(r.Raw.Response.Docs, v); err != nil {
		return errs.New("failed to scan docs", errs.WithCause(err))
	}
	return nil
}

func (r *SelectResult) Facet() (*JSONFacetResponse, error) {
	res, err := r.Raw.Facets.Parse(r.query.facet)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return res, nil
}
