package common

import (
	"fmt"
	"strconv"
)

type RangeFilterParameter struct {
	From *int `json:"from"`
	To   *int `json:"to"`
}

func (r *RangeFilterParameter) ToRange() string {
	if r.From == nil && r.To == nil {
		return ""
	}

	var from string
	if r.From == nil {
		from = "*"
	} else {
		from = strconv.Itoa(*r.From)
	}

	var to string
	if r.To == nil {
		to = "*"
	} else {
		to = strconv.Itoa(*r.To)
	}

	return fmt.Sprintf("[%s TO %s]", from, to)
}

type SearchResultResponse[P any, D any, F any] struct {
	Stats   SearchResultStats[P, F] `json:"status,omitempty"`
	Items   []D                     `json:"items"`
	Message string                  `json:"message,omitempty"`
}

type SearchResultStats[P any, F any] struct {
	Time   uint32 `json:"time"`
	Total  uint32 `json:"total"`
	Index  uint32 `json:"index"`
	Pages  uint32 `json:"pages"`
	Count  uint32 `json:"count"`
	Params P      `json:"params,omitempty"`
	Facet  F      `json:"facet,omitempty"`
}
