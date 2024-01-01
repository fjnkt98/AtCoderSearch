package solr

import (
	"net/url"
	"regexp"
	"strconv"
)

var SOLR_SPECIAL_CHARACTERS, _ = regexp.Compile(`(\+|\-|&&|\|\||!|\(|\)|\{|\}|\[|\]|\^|"|\~|\*|\?|:|/|AND|OR)`)

func Sanitize(s string) string {
	return SOLR_SPECIAL_CHARACTERS.ReplaceAllString(s, `\${0}`)
}

type EDisMaxQueryBuilder struct {
	inner url.Values
}

func NewEDisMaxQueryBuilder() *EDisMaxQueryBuilder {
	params := url.Values{}
	params.Set("defType", "edismax")
	return &EDisMaxQueryBuilder{
		inner: params,
	}
}

func (b *EDisMaxQueryBuilder) Build() url.Values {
	return b.inner
}

func (b *EDisMaxQueryBuilder) Sort(sort string) *EDisMaxQueryBuilder {
	if sort == "" {
		return b
	}
	b.inner.Set("sort", sort)
	return b
}

func (b *EDisMaxQueryBuilder) Start(start int) *EDisMaxQueryBuilder {
	b.inner.Set("start", strconv.Itoa(int(start)))
	return b
}

func (b *EDisMaxQueryBuilder) Rows(rows int) *EDisMaxQueryBuilder {
	b.inner.Set("rows", strconv.Itoa(int(rows)))
	return b
}

func (b *EDisMaxQueryBuilder) Fq(fq []string) *EDisMaxQueryBuilder {
	for _, fq := range fq {
		b.inner.Add("fq", fq)
	}
	return b
}
func (b *EDisMaxQueryBuilder) Fl(fl string) *EDisMaxQueryBuilder {
	if fl == "" {
		return b
	}
	b.inner.Set("fl", fl)
	return b
}

func (b *EDisMaxQueryBuilder) Debug() *EDisMaxQueryBuilder {
	b.inner.Set("debug", "all")
	b.inner.Set("debug.explain.structured", "true")
	return b
}

func (b *EDisMaxQueryBuilder) Wt(wt string) *EDisMaxQueryBuilder {
	if wt == "" {
		return b
	}
	b.inner.Set("wt", wt)
	return b
}

func (b *EDisMaxQueryBuilder) Facet(facet string) *EDisMaxQueryBuilder {
	if facet == "" {
		return b
	}
	b.inner.Set("json.facet", facet)
	return b
}

func (b *EDisMaxQueryBuilder) Op(op string) *EDisMaxQueryBuilder {
	if op == "" {
		return b
	}
	b.inner.Set("q.op", op)
	return b
}

func (b *EDisMaxQueryBuilder) Df(df string) *EDisMaxQueryBuilder {
	if df == "" {
		return b
	}
	b.inner.Set("df", df)
	return b
}

func (b *EDisMaxQueryBuilder) Q(q string) *EDisMaxQueryBuilder {
	if q == "" {
		return b
	}
	b.inner.Set("q", q)
	return b
}

func (b *EDisMaxQueryBuilder) Qf(qf string) *EDisMaxQueryBuilder {
	if qf == "" {
		return b
	}
	b.inner.Set("qf", qf)
	return b
}

func (b *EDisMaxQueryBuilder) Qs(qs int) *EDisMaxQueryBuilder {
	b.inner.Set("qs", strconv.Itoa(int(qs)))
	return b
}

func (b *EDisMaxQueryBuilder) Pf(pf string) *EDisMaxQueryBuilder {
	if pf == "" {
		return b
	}
	b.inner.Set("pf", pf)
	return b
}

func (b *EDisMaxQueryBuilder) Mm(mm string) *EDisMaxQueryBuilder {
	if mm == "" {
		return b
	}
	b.inner.Set("mm", mm)
	return b
}

func (b *EDisMaxQueryBuilder) QAlt(q string) *EDisMaxQueryBuilder {
	if q == "" {
		return b
	}
	b.inner.Set("q.alt", q)
	return b
}

func (b *EDisMaxQueryBuilder) Tie(tie float64) *EDisMaxQueryBuilder {
	b.inner.Set("tie", strconv.FormatFloat(tie, 'f', 6, 64))
	return b
}

func (b *EDisMaxQueryBuilder) Bq(bq []string) *EDisMaxQueryBuilder {
	for _, bq := range bq {
		b.inner.Add("bq", bq)
	}
	return b
}

func (b *EDisMaxQueryBuilder) Bf(bf []string) *EDisMaxQueryBuilder {
	for _, bf := range bf {
		b.inner.Add("bf", bf)
	}
	return b
}

func (b *EDisMaxQueryBuilder) Sow(sow bool) *EDisMaxQueryBuilder {
	if sow {
		b.inner.Set("sow", "true")
	} else {
		b.inner.Set("sow", "false")
	}
	return b
}

func (b *EDisMaxQueryBuilder) Boost(boost []string) *EDisMaxQueryBuilder {
	for _, boost := range boost {
		b.inner.Add("boost", boost)
	}
	return b
}

func (b *EDisMaxQueryBuilder) LowerCaseOperators(flag bool) *EDisMaxQueryBuilder {
	if flag {
		b.inner.Set("lowercaseOperators", "true")
	} else {
		b.inner.Set("lowercaseOperators", "false")
	}
	return b
}

func (b *EDisMaxQueryBuilder) Pf2(pf2 string) *EDisMaxQueryBuilder {
	if pf2 == "" {
		return b
	}
	b.inner.Set("pf2", pf2)
	return b
}

func (b *EDisMaxQueryBuilder) Ps2(ps2 int) *EDisMaxQueryBuilder {
	b.inner.Set("ps2", strconv.Itoa(int(ps2)))
	return b
}

func (b *EDisMaxQueryBuilder) Pf3(pf3 string) *EDisMaxQueryBuilder {
	if pf3 == "" {
		return b
	}
	b.inner.Set("pf3", pf3)
	return b
}

func (b *EDisMaxQueryBuilder) Ps3(ps3 int) *EDisMaxQueryBuilder {
	b.inner.Set("ps3", strconv.Itoa(int(ps3)))
	return b
}

func (b *EDisMaxQueryBuilder) StopWords(flag bool) *EDisMaxQueryBuilder {
	if flag {
		b.inner.Set("stopwords", "true")
	} else {
		b.inner.Set("stopwords", "false")
	}
	return b
}

func (b *EDisMaxQueryBuilder) Uf(uf string) *EDisMaxQueryBuilder {
	if uf == "" {
		return b
	}
	b.inner.Set("uf", uf)
	return b
}
