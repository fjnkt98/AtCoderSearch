package utility

import (
	"fjnkt98/atcodersearch/pkg/solr"
	"fmt"
	"strconv"
	"time"
)

type SearchResultResponse[I any] struct {
	Stats   SearchResultStats `json:"stats"`
	Items   []I               `json:"items"`
	Message string            `json:"message,omitempty"`
}

type SearchResultStats struct {
	Time   int                    `json:"time"`
	Total  int                    `json:"total"`
	Index  int                    `json:"index"`
	Pages  int                    `json:"pages"`
	Count  int                    `json:"count"`
	Params any                    `json:"params,omitempty"`
	Facet  map[string][]FacetPart `json:"facet,omitempty"`
}

func NewErrorResponse[I any](msg string) SearchResultResponse[I] {
	return SearchResultResponse[I]{
		Message: msg,
	}
}

type FacetPart struct {
	Label string `json:"label"`
	Count int    `json:"count"`
}

func ConvertBucket[T solr.BucketElement](b []solr.Bucket[T]) []FacetPart {
	p := make([]FacetPart, len(b))
	for i, b := range b {
		var label string
		switch v := any(b.Val).(type) {
		case int:
			label = strconv.Itoa(v)
		case uint:
			label = strconv.Itoa(int(v))
		case float64:
			label = strconv.FormatFloat(v, 'f', 6, 64)
		case time.Time:
			label = v.String()
		case string:
			label = v
		}
		p[i] = FacetPart{
			Label: label,
			Count: b.Count,
		}
	}

	return p
}

func ConvertRangeBucket(r solr.RangeFacetCount[int], p RangeFacetParam) []FacetPart {
	if p.From == nil || p.To == nil || p.Gap == nil {
		return nil
	}

	parts := make([]FacetPart, 0, len(r.Buckets)+2)

	parts = append(parts, FacetPart{Label: fmt.Sprintf("~ %d", *p.From), Count: r.Before.Count})
	end := *p.To
	for _, b := range r.Buckets {
		parts = append(parts, FacetPart{
			Label: fmt.Sprintf("%d ~ %d", b.Val, b.Val+*p.Gap),
			Count: b.Count,
		})
		end = max(end, b.Val+*p.Gap)
	}
	parts = append(parts, FacetPart{Label: fmt.Sprintf("%d ~", end), Count: r.After.Count})

	return parts
}
