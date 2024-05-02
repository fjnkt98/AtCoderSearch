package solr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/goark/errs"
)

type MoreLikeThisQuery struct {
	client *http.Client
	uri    *url.URL
	params url.Values
	q      string
	qf     string
	minTF  *int
	minDF  *int
	maxDF  *int
	minWL  *int
	maxWL  *int
	maxQT  *int
	maxNTP *int
	boost  bool
}

func newMoreLikeThisQuery(client *http.Client, uri *url.URL) *MoreLikeThisQuery {
	return &MoreLikeThisQuery{
		client: client,
		uri:    uri,
		params: url.Values{},
	}
}

func (q *MoreLikeThisQuery) construct() {
	params := make([]string, 0)
	if q.qf != "" {
		params = append(params, fmt.Sprintf("qf=%s", Sanitize(q.qf)))
	}
	if q.minTF != nil {
		params = append(params, fmt.Sprintf("mintf=%d", *q.minTF))
	}
	if q.minDF != nil {
		params = append(params, fmt.Sprintf("mindf=%d", *q.minDF))
	}
	if q.maxDF != nil {
		params = append(params, fmt.Sprintf("maxdf=%d", *q.maxDF))
	}
	if q.minWL != nil {
		params = append(params, fmt.Sprintf("minwl=%d", *q.minWL))
	}
	if q.maxWL != nil {
		params = append(params, fmt.Sprintf("maxwl=%d", *q.maxWL))
	}
	if q.maxQT != nil {
		params = append(params, fmt.Sprintf("maxqt=%d", *q.maxQT))
	}
	if q.maxNTP != nil {
		params = append(params, fmt.Sprintf("maxntp=%d", *q.maxNTP))
	}
	if q.boost {
		params = append(params, "boost=true")
	}

	if len(params) > 0 {
		q.params.Set("q", fmt.Sprintf(
			"{!mlt %s}%s",
			strings.Join(params, " "),
			q.q,
		),
		)
	}
}

func (q *MoreLikeThisQuery) Exec(ctx context.Context) (*MoreLikeThisResult, error) {
	q.construct()
	q.uri.RawQuery = q.params.Encode()
	u := q.uri.String()
	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, errs.New(
			"failed to prepare more like this request",
			errs.WithCause(err),
			errs.WithContext("uri", u),
		)
	}
	res, err := q.client.Do(req)
	if err != nil {
		return nil, errs.New(
			"failed to execute more like this request",
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

	result := &MoreLikeThisResult{
		Raw: body,
	}
	if res.StatusCode != http.StatusOK {
		return result, errs.Wrap(ErrNonOKResponse, errs.WithContext("uri", u))
	}

	return result, nil
}

func (q *MoreLikeThisQuery) Raw() url.Values {
	q.construct()
	return q.params
}

func (q *MoreLikeThisQuery) Start(start int) *MoreLikeThisQuery {
	q.params.Set("start", strconv.Itoa(start))
	return q
}

func (q *MoreLikeThisQuery) Rows(rows int) *MoreLikeThisQuery {
	q.params.Set("rows", strconv.Itoa(rows))
	return q
}

func (q *MoreLikeThisQuery) Fq(fq ...string) *MoreLikeThisQuery {
	for _, fq := range fq {
		if fq != "" {
			q.params.Add("fq", fq)
		}
	}
	return q
}

func (q *MoreLikeThisQuery) Fl(fl string) *MoreLikeThisQuery {
	if fl != "" {
		q.params.Set("fl", fl)
	}
	return q
}

func (q *MoreLikeThisQuery) Wt(wt string) *MoreLikeThisQuery {
	if wt != "" {
		q.params.Set("wt", wt)
	}
	return q
}

func (q *MoreLikeThisQuery) Q(uniqueKey string) *MoreLikeThisQuery {
	if uniqueKey != "" {
		q.q = uniqueKey
	}
	return q
}

func (q *MoreLikeThisQuery) Qf(qf string) *MoreLikeThisQuery {
	if qf != "" {
		q.qf = qf
	}
	return q
}

func (q *MoreLikeThisQuery) MinTF(v int) *MoreLikeThisQuery {
	q.minTF = &v
	return q
}

func (q *MoreLikeThisQuery) MinDF(v int) *MoreLikeThisQuery {
	q.minDF = &v
	return q
}

func (q *MoreLikeThisQuery) MaxDF(v int) *MoreLikeThisQuery {
	q.maxDF = &v
	return q
}

func (q *MoreLikeThisQuery) MinWL(v int) *MoreLikeThisQuery {
	q.minWL = &v
	return q
}

func (q *MoreLikeThisQuery) MaxWL(v int) *MoreLikeThisQuery {
	q.maxWL = &v
	return q
}

func (q *MoreLikeThisQuery) MaxQT(v int) *MoreLikeThisQuery {
	q.maxQT = &v
	return q
}

func (q *MoreLikeThisQuery) MaxNTP(v int) *MoreLikeThisQuery {
	q.maxNTP = &v
	return q
}

func (q *MoreLikeThisQuery) Boost(v bool) *MoreLikeThisQuery {
	q.boost = v
	return q
}

type MoreLikeThisResult struct {
	Raw SelectResponse
}

func (r *MoreLikeThisResult) Scan(v any) error {
	if err := json.Unmarshal(r.Raw.Response.Docs, v); err != nil {
		return errs.New("failed to scan docs", errs.WithCause(err))
	}
	return nil
}
