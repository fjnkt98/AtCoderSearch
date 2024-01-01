package utility

import (
	"encoding/json"
	"fjnkt98/atcodersearch/pkg/solr"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type IntegerRange struct {
	From *int `json:"from" schema:"from"`
	To   *int `json:"to" schema:"to"`
}

func (r *IntegerRange) ToRange() string {
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

type FloatRange struct {
	From *float64 `json:"from" schema:"from"`
	To   *float64 `json:"to" schema:"to"`
}

func (r *FloatRange) ToRange() string {
	if r.From == nil && r.To == nil {
		return ""
	}

	var from string
	if r.From == nil {
		from = "*"
	} else {
		from = strconv.FormatFloat(*r.From, 'f', 6, 64)
	}

	var to string
	if r.To == nil {
		to = "*"
	} else {
		to = strconv.FormatFloat(*r.To, 'f', 6, 64)
	}

	return fmt.Sprintf("[%s TO %s}", from, to)
}

type DateRange struct {
	From *time.Time `json:"from" schema:"from"`
	To   *time.Time `json:"to" schema:"to"`
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
	From *int `json:"from" schema:"from"`
	To   *int `json:"to" schema:"to"`
	Gap  *int `json:"gap" schema:"gap"`
}

func (p *RangeFacetParam) ToFacet(name, field string) map[string]map[string]any {
	if p.From == nil || p.To == nil || p.Gap == nil {
		return nil
	}

	return map[string]map[string]any{
		name: {
			"type":  "range",
			"field": field,
			"start": *p.From,
			"end":   *p.To,
			"gap":   *p.Gap,
			"other": "all",
			"domain": map[string]any{
				"excludeTags": []string{field},
			},
		},
	}
}

type TermFacetParam []string

func (p *TermFacetParam) ToFacet(m map[string]string) map[string]map[string]any {
	facet := make(map[string]map[string]any)
	for _, name := range *p {
		field, ok := m[name]
		if !ok {
			field = name
		}
		facet[name] = map[string]any{
			"type":     "terms",
			"field":    field,
			"limit":    -1,
			"mincount": 0,
			"sort":     "index",
			"domain": map[string]any{
				"excludeTags": []string{field},
			},
		}
	}

	return facet
}

func Facet(p any) string {
	facets := make(map[string]any)

	ty := reflect.TypeOf(p)
	if ty.Kind() != reflect.Struct {
		return ""
	}

	val := reflect.ValueOf(p)
	for i := 0; i < ty.NumField(); i++ {
		fieldValue := val.Field(i)
		fieldType := ty.Field(i)

		mapping := make(map[string]string)
		if tag, ok := fieldType.Tag.Lookup("facet"); ok {
			for _, t := range strings.Split(tag, ",") {
				if name, field, ok := strings.Cut(t, ":"); ok {
					mapping[name] = field
				}
			}
		}

		switch v := fieldValue.Interface().(type) {
		case TermFacetParam:
			for name, facet := range v.ToFacet(mapping) {
				facets[name] = facet
			}
		case RangeFacetParam:
			for name, field := range mapping {
				for name, facet := range v.ToFacet(name, field) {
					facets[name] = facet
				}
				break
			}
		}
	}

	if facet, err := json.Marshal(facets); err != nil {
		return ""
	} else {
		return string(facet)
	}
}

func Filter(p any) []string {
	ty := reflect.TypeOf(p)
	if ty.Kind() != reflect.Struct {
		return nil
	}
	val := reflect.ValueOf(p)

	fq := make([]string, 0, ty.NumField())
loop:
	for i := 0; i < ty.NumField(); i++ {
		fieldType := ty.Field(i)
		fieldValue := val.Field(i)

		quote := false
		field := ""

		if tag, ok := fieldType.Tag.Lookup("filter"); ok {
			if tag == "-" {
				continue loop
			}
			if f, opt, ok := strings.Cut(tag, ","); ok {
				field = f
				if opt == "quote" {
					quote = true
				}
			} else {
				field = tag
			}
		}

		switch v := fieldValue.Interface().(type) {
		case IntegerRange:
			if r := v.ToRange(); r != "" {
				fq = append(fq, fmt.Sprintf("%s:%s", field, r))
			}
		case FloatRange:
			if r := v.ToRange(); r != "" {
				fq = append(fq, fmt.Sprintf("%s:%s", field, r))
			}
		case DateRange:
			if r := v.ToRange(); r != "" {
				fq = append(fq, fmt.Sprintf("%s:%s", field, r))
			}
		case []string:
			values := SanitizeStrings(v)
			if quote {
				values = QuoteStrings(v)
			}
			if len(values) == 0 {
				continue loop
			}
			fq = append(fq, fmt.Sprintf("%s:(%s)", field, strings.Join(values, " OR ")))
		}
	}

	return fq
}

type Validator interface {
	Validate() bool
}

type SearchParams[F, C Validator] struct {
	Limit  *int     `json:"limit" schema:"limit"`
	Page   int      `json:"page" schema:"page"`
	Filter F        `json:"filter" schema:"filter"`
	Sort   []string `json:"sort" schema:"sort"`
	Facet  C        `json:"facet" schema:"facet"`
}

func (p *SearchParams[F, C]) GetRows() int {
	if p.Limit == nil {
		return 20
	}

	return *p.Limit
}

func (p *SearchParams[F, C]) GetStart() int {
	if p.Page == 0 || p.GetRows() == 0 {
		return 0
	}

	return (p.Page - 1) * p.GetRows()
}

func (p *SearchParams[F, C]) GetSort() string {
	orders := make([]string, 0, len(p.Sort))
	for _, s := range p.Sort {
		if strings.HasPrefix(s, "-") {
			orders = append(orders, fmt.Sprintf("%s desc", s[1:]))
		} else {
			orders = append(orders, fmt.Sprintf("%s asc", s))
		}
	}

	return strings.Join(orders, ",")
}

func (p *SearchParams[F, C]) GetFacet() string {
	return Facet(p.Facet)
}

func (p *SearchParams[F, C]) GetFilter() []string {
	return Filter(p.Filter)
}

func FieldList(doc any) []string {
	ty := reflect.TypeOf(doc)
	if ty.Kind() != reflect.Pointer {
		return nil
	}

	ty = ty.Elem()
	if ty.Kind() != reflect.Struct {
		return nil
	}

	fl := make([]string, 0, ty.NumField())
	for i := 0; i < ty.NumField(); i++ {
		f := ty.Field(i)

		var name string
		if tag, ok := f.Tag.Lookup("json"); ok {
			if tag == "-" {
				continue
			}
			n, _, _ := strings.Cut(tag, ",")
			name = n
		} else {
			name = f.Name
		}
		fl = append(fl, name)
	}
	return fl
}
