package api

import (
	"fjnkt98/atcodersearch/pkg/solr"
	"fmt"
	"strings"
)

type ParameterBase struct {
	Limit *int `json:"limit" query:"limit"`
	Page  int  `json:"page" query:"page"`
}

func (r *ParameterBase) Rows() int {
	if r.Limit == nil {
		return 20
	}
	return *r.Limit
}

func (p *ParameterBase) Start() int {
	if p.Page == 0 || p.Rows() == 0 {
		return 0
	}

	return (p.Page - 1) * p.Rows()
}

func ParseSort(sort []string, defaults ...string) []string {
	orders := make([]string, 0, len(sort)+len(defaults))
	for _, s := range sort {
		if strings.HasPrefix(s, "-") {
			orders = append(orders, fmt.Sprintf("%s desc", solr.Sanitize(s[1:])))
		} else {
			orders = append(orders, fmt.Sprintf("%s asc", solr.Sanitize(s)))
		}
	}
	orders = append(orders, defaults...)

	return orders
}

func ParseQ(q string) string {
	s := make([]string, 0)
	for _, w := range solr.Parse(q) {
		s = append(s, w.String())
	}
	return strings.Join(s, " ")
}

type ResultResponse[T any] struct {
	Stats   ResultStats `json:"stats"`
	Items   []T         `json:"items"`
	Message string      `json:"message,omitempty"`
}

func NewErrorResponse(message string, params any) ResultResponse[any] {
	return ResultResponse[any]{
		Stats: ResultStats{
			Params: params,
		},
		Items:   []any{},
		Message: message,
	}
}

func NewEmptyResponse() ResultResponse[any] {
	return ResultResponse[any]{
		Stats: ResultStats{},
		Items: []any{},
	}
}

type ResultStats struct {
	Time   int                     `json:"time"`
	Total  int                     `json:"total"`
	Index  int                     `json:"index"`
	Pages  int                     `json:"pages"`
	Count  int                     `json:"count"`
	Params any                     `json:"params,omitempty"`
	Facet  map[string][]FacetCount `json:"facet,omitempty"`
}

type FacetCount struct {
	Label string `json:"label"`
	Count int    `json:"count"`
}

func NewFacetCount(f *solr.JSONFacetResponse) map[string][]FacetCount {
	if f == nil {
		return nil
	}
	counts := make(map[string][]FacetCount)
	for k, v := range f.Terms {
		counts[k] = FacetCountsFromStringBucket(v.Buckets)
	}
	for k, v := range f.Range {
		counts[k] = FacetCountsFromRangeBucket(v.Buckets)
	}
	return counts
}

func FacetCountsFromStringBucket(buckets []solr.StringBucket) []FacetCount {
	res := make([]FacetCount, len(buckets))
	for i, b := range buckets {
		res[i] = FacetCount{
			Label: b.Val,
			Count: b.Count,
		}
	}
	return res
}

func FacetCountsFromRangeBucket(buckets []solr.RangeBucket) []FacetCount {
	res := make([]FacetCount, len(buckets))
	for i, b := range buckets {
		if b.Begin == nil {
			res[i] = FacetCount{
				Label: fmt.Sprintf(" ~ %d", *b.End),
				Count: b.Count,
			}
		} else if b.End == nil {
			res[i] = FacetCount{
				Label: fmt.Sprintf("%d ~ ", *b.Begin),
				Count: b.Count,
			}
		} else {
			res[i] = FacetCount{
				Label: fmt.Sprintf("%d ~ %d", *b.Begin, *b.End),
				Count: b.Count,
			}
		}
	}
	return res
}
