package problem

import (
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type FullTextExtractor struct {
}

func NewFullTextExtractor() FullTextExtractor {
	return FullTextExtractor{}
}

func (f *FullTextExtractor) Extract(html io.Reader) ([]string, []string, error) {
	doc, _ := goquery.NewDocumentFromReader(html)

	textJa := make([]string, 0)
	textEn := make([]string, 0)

	doc.Find("section").Each(func(_ int, section *goquery.Selection) {
		// For modern contest problem format
		if strings.Contains(section.Find("h3").Text(), "問題") {
			textJa = append(textJa, strings.Split(section.Text(), "\n")...)
		}

		// For legacy contest problem format
		if prev := section.Prev(); goquery.NodeName(prev) == "h3" {
			if text := prev.Text(); strings.Contains(text, "問題") {
				textJa = append(textJa, strings.Split(section.Text(), "\n")...)
			}
		}
	})

	doc.Find("span.lang-en").Find("section").Each(func(_ int, section *goquery.Selection) {
		if strings.Contains(section.NextAll().Find("h3").Text(), "Statement") || strings.Contains(section.Find("h3").Text(), "Statement") {
			textEn = append(textEn, strings.Split(section.Text(), "\n")...)
		}
	})
	return textJa, textEn, nil
}
