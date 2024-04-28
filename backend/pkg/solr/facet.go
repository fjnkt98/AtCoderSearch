package solr

import (
	"encoding/json"

	"github.com/goark/errs"
)

const (
	FACET_TYPE_TERMS = "terms"
	FACET_TYPE_RANGE = "range"
	FACET_TYPE_QUERY = "query"
)

type JSONFacetQuery struct {
	termsFacets map[string]*TermsFacetQuery
	rangeFacets map[string]*RangeFacetQuery
}

func NewJSONFacetQuery() *JSONFacetQuery {
	return &JSONFacetQuery{
		termsFacets: make(map[string]*TermsFacetQuery),
		rangeFacets: make(map[string]*RangeFacetQuery),
	}
}

func (q *JSONFacetQuery) Terms(f *TermsFacetQuery) *JSONFacetQuery {
	q.termsFacets[f.name] = f
	return q
}

func (q *JSONFacetQuery) Range(f *RangeFacetQuery) *JSONFacetQuery {
	q.rangeFacets[f.name] = f
	return q
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

type TermsFacetQuery struct {
	name   string
	params map[string]any
}

func NewTermsFacetQuery(field string) *TermsFacetQuery {
	return &TermsFacetQuery{
		name: field,
		params: map[string]any{
			"type":  FACET_TYPE_TERMS,
			"field": field,
		},
	}
}

func (q *TermsFacetQuery) Name(v string) *TermsFacetQuery {
	q.name = v
	return q
}

func (q *TermsFacetQuery) Limit(v int) *TermsFacetQuery {
	q.params["limit"] = v
	return q
}

func (q *TermsFacetQuery) MinCount(v int) *TermsFacetQuery {
	q.params["mincount"] = v
	return q
}
func (q *TermsFacetQuery) Sort(v string) *TermsFacetQuery {
	q.params["sort"] = v
	return q
}

func (q *TermsFacetQuery) ExcludeTags(tags ...string) *TermsFacetQuery {
	q.params["domain"] = map[string]any{
		"excludeTags": tags,
	}
	return q
}

type RangeFacetQuery struct {
	name   string
	start  int
	end    int
	gap    int
	params map[string]any
}

func NewRangeFacetQuery(field string, start, end, gap int) *RangeFacetQuery {
	return &RangeFacetQuery{
		name:  field,
		start: start,
		end:   end,
		gap:   gap,
		params: map[string]any{
			"type":  FACET_TYPE_RANGE,
			"field": field,
			"start": start,
			"end":   end,
			"gap":   gap,
		},
	}
}

func (q *RangeFacetQuery) Name(v string) *RangeFacetQuery {
	q.name = v
	return q
}

func (q *RangeFacetQuery) MinCount(v int) *RangeFacetQuery {
	q.params["mincount"] = v
	return q
}

func (q *RangeFacetQuery) HardEnd(v bool) *RangeFacetQuery {
	q.params["hardend"] = v
	return q
}

func (q *RangeFacetQuery) Other(v string) *RangeFacetQuery {
	q.params["other"] = v
	return q
}

func (q *RangeFacetQuery) Include(v string) *RangeFacetQuery {
	q.params["include"] = v
	return q
}

func (q *RangeFacetQuery) ExcludeTags(tags ...string) *RangeFacetQuery {
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
	termsFacets := make(map[string]TermsFacetCount)
	rangeFacets := make(map[string]RangeFacetCount)

	if query == nil {
		return &JSONFacetResponse{
			Terms: termsFacets,
			Range: rangeFacets,
		}, ErrNoFacetCounts
	}

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
