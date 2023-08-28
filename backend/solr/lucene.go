package solr

import (
	"net/url"
	"strconv"
)

type LuceneQueryBuilder struct {
	inner url.Values
}

func NewLuceneQueryBuilder() *LuceneQueryBuilder {
	params := url.Values{}
	params.Set("defType", "lucene")
	return &LuceneQueryBuilder{
		inner: params,
	}
}

func (b *LuceneQueryBuilder) Build() url.Values {
	return b.inner
}

func (b *LuceneQueryBuilder) Sort(sort string) *LuceneQueryBuilder {
	if sort == "" {
		return b
	}
	b.inner.Set("sort", sort)
	return b
}

func (b *LuceneQueryBuilder) Start(start int) *LuceneQueryBuilder {
	b.inner.Set("start", strconv.Itoa(int(start)))
	return b
}

func (b *LuceneQueryBuilder) Rows(rows int) *LuceneQueryBuilder {
	b.inner.Set("rows", strconv.Itoa(int(rows)))
	return b
}

func (b *LuceneQueryBuilder) Fq(fq []string) *LuceneQueryBuilder {
	for _, fq := range fq {
		b.inner.Add("fq", fq)
	}
	return b
}
func (b *LuceneQueryBuilder) Fl(fl string) *LuceneQueryBuilder {
	if fl == "" {
		return b
	}
	b.inner.Set("fl", fl)
	return b
}

func (b *LuceneQueryBuilder) Debug() *LuceneQueryBuilder {
	b.inner.Set("debug", "all")
	b.inner.Set("debug.explain.structured", "true")
	return b
}

func (b *LuceneQueryBuilder) Wt(wt string) *LuceneQueryBuilder {
	if wt == "" {
		return b
	}
	b.inner.Set("wt", wt)
	return b
}

func (b *LuceneQueryBuilder) Facet(facet string) *LuceneQueryBuilder {
	if facet == "" {
		return b
	}
	b.inner.Set("json.facet", facet)
	return b
}

func (b *LuceneQueryBuilder) Op(op string) *LuceneQueryBuilder {
	if op == "" {
		return b
	}
	b.inner.Set("q.op", op)
	return b
}

func (b *LuceneQueryBuilder) Df(df string) *LuceneQueryBuilder {
	if df == "" {
		return b
	}
	b.inner.Set("df", df)
	return b
}

func (b *LuceneQueryBuilder) Q(q string) *LuceneQueryBuilder {
	if q == "" {
		return b
	}
	b.inner.Set("q", q)
	return b
}
