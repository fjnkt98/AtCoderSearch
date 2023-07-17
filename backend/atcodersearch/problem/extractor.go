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

func (f *FullTextExtractor) DFS(element *goquery.Selection) string {
	buffer := make([]string, 0)

	element.Children().Each(func(_ int, e *goquery.Selection) {
		switch goquery.NodeName(e) {
		case "pre", "h3", "var":
		default:
			buffer = append(buffer, e.Text())
			buffer = append(buffer, f.DFS(e))
		}
	})

	return strings.Join(buffer, "")
}

func (f *FullTextExtractor) Extract(html io.Reader) ([]string, []string, error) {
	doc, _ := goquery.NewDocumentFromReader(html)

	textJa := make([]string, 0)
	textEn := make([]string, 0)

	if ja := doc.Find("span.lang-ja"); ja != nil {
		ja.Find("section").Each(func(_ int, section *goquery.Selection) {
			if h3 := section.Find("h3").Text(); strings.Contains(h3, "問題") {
				textJa = append(textJa, f.DFS(section))
			}
		})
	} else {
		doc.Find("section").Each(func(_ int, section *goquery.Selection) {
			if h3 := section.Find("h3").Text(); strings.Contains(h3, "問題") {
				textJa = append(textJa, f.DFS(section))
			}
		})
	}

	if en := doc.Find("span.lang-en"); en != nil {
		en.Find("section").Each(func(_ int, section *goquery.Selection) {
			if h3 := section.Find("h3").Text(); strings.Contains(h3, "Statement") {
				textEn = append(textEn, f.DFS(section))
			}
		})
	}
	return textJa, textEn, nil
}
