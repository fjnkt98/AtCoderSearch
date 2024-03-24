package api

import (
	"fjnkt98/atcodersearch/pkg/solr"
	"fmt"
	"strconv"
	"strings"
)

type localParam struct {
	Key   string
	Value string
}

func (p localParam) String() string {
	if p.Key == "" || p.Value == "" {
		return ""
	}
	return fmt.Sprintf("%s=%s", p.Key, p.Value)
}

func LocalParam(key, value string) localParam {
	return localParam{
		Key:   key,
		Value: value,
	}
}

func Quotes(s []string) []string {
	res := make([]string, 0, len(s))
	for _, s := range s {
		if s != "" {
			res = append(res, fmt.Sprintf(`"%s"`, s))
		}
	}
	return res
}

func TermsFilter(values []string, field string, params ...localParam) string {
	lp := make([]string, 0, len(params))
	for _, p := range params {
		if s := p.String(); s != "" {
			lp = append(lp, p.String())
		}
	}
	var p string
	if len(lp) > 0 {
		p = fmt.Sprintf(`{!%s}`, strings.Join(lp, " "))
	}

	v := Quotes(solr.Sanitizes(values))
	if len(v) == 0 {
		return ""
	}
	return fmt.Sprintf(`%s%s:(%s)`, p, field, strings.Join(v, " OR "))
}

func RangeFilter(from *int, to *int, field string, params ...localParam) string {
	if from == nil && to == nil {
		return ""
	}

	lp := make([]string, 0, len(params))
	for _, p := range params {
		if s := p.String(); s != "" {
			lp = append(lp, p.String())
		}
	}
	var p string
	if len(lp) > 0 {
		p = fmt.Sprintf(`{!%s}`, strings.Join(lp, " "))
	}

	var f string
	if from == nil {
		f = "*"
	} else {
		f = strconv.Itoa(*from)
	}

	var t string
	if to == nil {
		t = "*"
	} else {
		t = strconv.Itoa(*to)
	}

	return fmt.Sprintf(`%s%s:[%s TO %s]`, p, field, f, t)
}
