package acs

import (
	"fjnkt98/atcodersearch/solr"
	"strconv"
	"time"
)

type SearchResultResponse[D any] struct {
	Stats   SearchResultStats `json:"stats,omitempty"`
	Items   []D               `json:"items"`
	Message string            `json:"message,omitempty"`
}

type SearchResultStats struct {
	Time   uint `json:"time"`
	Total  uint `json:"total"`
	Index  uint `json:"index"`
	Pages  uint `json:"pages"`
	Count  uint `json:"count"`
	Params any  `json:"params,omitempty"`
	Facet  any  `json:"facet,omitempty"`
}

func NewErrorResponse[D any](msg string, params any) SearchResultResponse[D] {
	return SearchResultResponse[D]{
		Message: msg,
	}
}

type QueryLog struct {
	RequestAt time.Time `json:"request_at"`
	Domain    string    `json:"domain"`
	Time      uint      `json:"time"`
	Hits      uint      `json:"hits"`
	Params    any       `json:"params"`
}

type FacetPart struct {
	Label string `json:"label"`
	Count uint   `json:"count"`
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
		p[i] = FacetPart {
			Label: label,
			Count: b.Count,
		}
	}

	return nil
}
