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

func ParseSort(sort []string) string {
	orders := make([]string, len(sort))
	for i, s := range sort {
		if strings.HasPrefix(s, "-") {
			orders[i] = fmt.Sprintf("%s desc", solr.Sanitize(s[1:]))
		} else {
			orders[i] = fmt.Sprintf("%s asc", solr.Sanitize(s))
		}
	}

	return strings.Join(orders, ",")
}

func ParseQ(q string) string {
	s := make([]string, 0)
	for _, w := range solr.Parse(q) {
		s = append(s, w.String())
	}
	return strings.Join(s, " ")
}

type ResultResponse struct {
	Stats   ResultStats `json:"stats"`
	Items   any         `json:"items"`
	Message string      `json:"message"`
}

func NewErrorResponse(message string, params any) ResultResponse {
	return ResultResponse{
		Stats: ResultStats{
			Params: params,
		},
		Items:   []any{},
		Message: message,
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

func NewFacetCountsFromStringBucket(buckets []solr.StringBucket) []FacetCount {
	res := make([]FacetCount, len(buckets))
	for i, b := range buckets {
		res[i] = FacetCount{
			Label: b.Val,
			Count: b.Count,
		}
	}
	return res
}

func NewFacetCountsFromRangeBucket(buckets []solr.RangeBucket) []FacetCount {
	res := make([]FacetCount, len(buckets))
	for i, b := range buckets {
		if b.Begin == nil {
			res[i] = FacetCount{
				Label: fmt.Sprintf("~ %d", b.End),
				Count: b.Count,
			}
		} else if b.End == nil {
			res[i] = FacetCount{
				Label: fmt.Sprintf("%d ~", b.Begin),
				Count: b.Count,
			}
		} else {
			res[i] = FacetCount{
				Label: fmt.Sprintf("%d ~ %d", b.Begin, b.End),
				Count: b.Count,
			}
		}
	}
	return res
}
