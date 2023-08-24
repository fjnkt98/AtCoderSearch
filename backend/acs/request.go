package acs

import (
	"fjnkt98/atcodersearch/solr"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/constraints"
)

type IntegerRange[T constraints.Integer] struct {
	From *T `json:"from,omitempty" schema:"from"`
	To   *T `json:"to,omitempty" schema:"to"`
}

func (r *IntegerRange[T]) ToRange() string {
	if r.From == nil && r.To == nil {
		return ""
	}

	var from string
	if r.From == nil {
		from = "*"
	} else {
		from = strconv.Itoa(int(*r.From))
	}

	var to string
	if r.To == nil {
		to = "*"
	} else {
		to = strconv.Itoa(int(*r.To))
	}

	return fmt.Sprintf("[%s TO %s}", from, to)
}

type FloatRange[T constraints.Float] struct {
	From *T `json:"from,omitempty" schema:"from"`
	To   *T `json:"to,omitempty" schema:"to"`
}

func (r *FloatRange[T]) ToRange() string {
	if r.From == nil && r.To == nil {
		return ""
	}

	var from string
	if r.From == nil {
		from = "*"
	} else {
		from = strconv.FormatFloat(float64(*r.From), 'f', 6, 64)
	}

	var to string
	if r.To == nil {
		to = "*"
	} else {
		to = strconv.FormatFloat(float64(*r.To), 'f', 6, 64)
	}

	return fmt.Sprintf("[%s TO %s}", from, to)
}

type DateRange struct {
	From *time.Time `json:"from,omitempty" schema:"from"`
	To   *time.Time `json:"to,omitempty" schema:"to"`
}

func (r *DateRange) ToRange() string {
	if r.From == nil && r.To == nil {
		return ""
	}

	var from string
	if r.From == nil {
		from = "*"
	} else {
		from = solr.IntoSolrDateTime(*r.From).String()
	}

	var to string
	if r.To == nil {
		to = "*"
	} else {
		to = solr.IntoSolrDateTime(*r.To).String()
	}

	return fmt.Sprintf("[%s TO %s}", from, to)
}

func SanitizeStrings(s []string) []string {
	sanitized := make([]string, 0, len(s))
	for _, e := range s {
		if e := strings.TrimSpace(solr.Sanitize(e)); e != "" {
			sanitized = append(sanitized, e)
		}
	}
	return sanitized
}

func QuoteStrings(s []string) []string {
	ss := make([]string, len(s))
	for i, e := range s {
		ss[i] = fmt.Sprintf(`"%s"`, e)
	}
	return ss
}

type RangeFacetParam struct {
	From int `json:"from" schema:"from"`
	To   int `json:"to" schema:"to" validate:"gtecsfield=From"`
	Gap  int `json:"gap" schema:"gap"`
}
