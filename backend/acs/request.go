package acs

import (
	"fmt"
	"strconv"
)

type RangeFilterParam struct {
	From *int `json:"from"`
	To   *int `json:"to"`
}

func (r *RangeFilterParam) ToRange() string {
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
