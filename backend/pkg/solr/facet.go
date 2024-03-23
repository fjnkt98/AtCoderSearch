package solr

import (
	"encoding/json"

	"github.com/goark/errs"
)

const (
	TERMS_FACET_TYPE = "terms"
	RANGE_FACET_TYPE = "range"
	QUERY_FACET_TYPE = "query"
)

type FacetQueryer interface {
	FacetQuery() (string, map[string]any)
}

type JSONFacetQuery struct {
	termsFacets map[string]*termsFacetQuery
	rangeFacets map[string]*rangeFacetQuery
}

func NewJSONFacetQuery(facets ...FacetQueryer) *JSONFacetQuery {
	termsFacets := make(map[string]*termsFacetQuery)
	rangeFacets := make(map[string]*rangeFacetQuery)

	for _, f := range facets {
		switch f := f.(type) {
		case *termsFacetQuery:
			termsFacets[f.name] = f
		case *rangeFacetQuery:
			rangeFacets[f.name] = f
		}
	}

	return &JSONFacetQuery{
		termsFacets: termsFacets,
		rangeFacets: rangeFacets,
	}
}

func (q JSONFacetQuery) MarshalJSON() ([]byte, error) {
	body := make(map[string]any)

	for k, v := range q.termsFacets {
		body[k] = v.params
	}
	for k, v := range q.rangeFacets {
		body[k] = v.params
	}

	return json.Marshal(body)
}

type termsFacetQuery struct {
	name   string
	params map[string]any
}

func NewTermsFacetQuery(field string) *termsFacetQuery {
	return &termsFacetQuery{
		name: field,
		params: map[string]any{
			"type":  TERMS_FACET_TYPE,
			"field": field,
		},
	}
}

func (q *termsFacetQuery) FacetQuery() (string, map[string]any) {
	return q.name, q.params
}

func (q *termsFacetQuery) Type() string {
	return TERMS_FACET_TYPE
}

func (q *termsFacetQuery) Name(v string) *termsFacetQuery {
	q.name = v
	return q
}

func (q *termsFacetQuery) Limit(v int) *termsFacetQuery {
	q.params["limit"] = v
	return q
}

func (q *termsFacetQuery) MinCount(v int) *termsFacetQuery {
	q.params["mincount"] = v
	return q
}

func (q *termsFacetQuery) ExcludeTags(tags ...string) *termsFacetQuery {
	q.params["domain"] = map[string]any{
		"excludeTags": tags,
	}
	return q
}

type rangeFacetQuery struct {
	name   string
	start  int
	end    int
	gap    int
	params map[string]any
}

func NewRangeFacetQuery(field string, start, end, gap int) *rangeFacetQuery {
	return &rangeFacetQuery{
		name:  field,
		start: start,
		end:   end,
		gap:   gap,
		params: map[string]any{
			"type":  RANGE_FACET_TYPE,
			"field": field,
			"start": start,
			"end":   end,
			"gap":   gap,
		},
	}
}

func (q *rangeFacetQuery) FacetQuery() (string, map[string]any) {
	return q.name, q.params
}

func (q *rangeFacetQuery) Type() string {
	return RANGE_FACET_TYPE
}

func (q *rangeFacetQuery) Name(v string) *rangeFacetQuery {
	q.name = v
	return q
}

func (q *rangeFacetQuery) MinCount(v int) *rangeFacetQuery {
	q.params["mincount"] = v
	return q
}

func (q *rangeFacetQuery) HardEnd(v bool) *rangeFacetQuery {
	q.params["hardend"] = v
	return q
}

func (q *rangeFacetQuery) Other(v string) *rangeFacetQuery {
	q.params["other"] = v
	return q
}

func (q *rangeFacetQuery) Include(v string) *rangeFacetQuery {
	q.params["include"] = v
	return q
}

func (q *rangeFacetQuery) ExcludeTags(tags ...string) *rangeFacetQuery {
	q.params["domain"] = map[string]any{
		"excludeTags": tags,
	}
	return q
}

type RawJSONFacetResponse map[string]json.RawMessage

func (f *RawJSONFacetResponse) UnmarshalJSON(data []byte) error {
	raw := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	res := make(map[string]json.RawMessage)
	for k, v := range raw {
		if k == "count" {
			continue
		}
		res[k] = v
	}
	*f = res
	return nil
}

type JSONFacetResponse struct {
	Terms map[string]TermsFacetCount
	Range map[string]RangeFacetCount
}

var ErrNoFacetCounts = errs.New("no facet counts")

func (f *RawJSONFacetResponse) Parse(query *JSONFacetQuery) (*JSONFacetResponse, error) {
	if query == nil {
		return nil, ErrNoFacetCounts
	}

	termsFacets := make(map[string]TermsFacetCount)
	rangeFacets := make(map[string]RangeFacetCount)

	for k := range query.termsFacets {
		raw, ok := (*f)[k]
		if !ok {
			continue
		}

		var t TermsFacetCount
		if err := json.Unmarshal(raw, &t); err != nil {
			return nil, errs.New(
				"failed to unmarshal terms facet count",
				errs.WithCause(err),
				errs.WithContext("key", k),
			)
		}
		termsFacets[k] = t
	}

	for k, v := range query.rangeFacets {
		raw, ok := (*f)[k]
		if !ok {
			continue
		}

		var r RawRangeFacetCount
		if err := json.Unmarshal(raw, &r); err != nil {
			return nil, errs.New(
				"failed to unmarshal range facet count",
				errs.WithCause(err),
				errs.WithContext("key", k),
			)
		}
		buckets := make([]RangeBucket, 0, len(r.Buckets)+2)
		if before := r.Before; before != nil {
			start := v.start
			buckets = append(buckets, RangeBucket{
				End:   &start,
				Count: before.Count,
			})
		}
		for _, b := range r.Buckets {
			begin := b.Val
			end := b.Val + v.gap
			buckets = append(buckets, RangeBucket{
				Begin: &begin,
				End:   &end,
				Count: b.Count,
			})
		}
		if after := r.After; after != nil {
			end := v.end
			buckets = append(buckets, RangeBucket{
				Begin: &end,
				Count: after.Count,
			})
		}

		rangeFacets[k] = RangeFacetCount{Buckets: buckets}
	}

	return &JSONFacetResponse{
		Terms: termsFacets,
		Range: rangeFacets,
	}, nil
}

type TermsFacetCount struct {
	Buckets []StringBucket `json:"buckets"`
}

type OthersCount struct {
	Count int `json:"count"`
}

type RawRangeFacetCount struct {
	Buckets []IntBucket  `json:"buckets"`
	Before  *OthersCount `json:"before"`
	After   *OthersCount `json:"after"`
	All     *OthersCount `json:"all"`
}

type RangeFacetCount struct {
	Buckets []RangeBucket `json:"buckets"`
}

type StringBucket struct {
	Val   string `json:"val"`
	Count int    `json:"count"`
}

type IntBucket struct {
	Val   int `json:"val"`
	Count int `json:"count"`
}

type RangeBucket struct {
	Begin *int `json:"begin"`
	End   *int `json:"end"`
	Count int  `json:"count"`
}
