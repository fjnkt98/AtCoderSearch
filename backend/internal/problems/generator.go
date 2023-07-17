package atcoder_search

import (
	"fjnkt98/atcodersearch/pkg/indexing/rate"
	"fmt"
	"log"
	"strings"
	"time"
)

var extractor = NewFullTextExtractor()

type Row struct {
	ProblemID      string
	ProblemTitle   string
	ProblemURL     string
	ContestID      string
	ContestTitle   string
	StartAt        int64
	Duration       int64
	RateChange     string
	Category       string
	HTML           string
	Difficulty     *int
	IsExperimental bool
}

func (r *Row) ToDocument() (ProblemIndex, error) {
	statementJa, statementEn, err := extractor.Extract(strings.NewReader(r.HTML))
	if err != nil {
		log.Printf("failed to extract statement at problem `%s`: %s", r.ProblemID, err.Error())
		return ProblemIndex{}, err
	}

	contestURL := fmt.Sprintf("https://atcoder.jp/contests/%s", r.ContestID)
	startAt := time.Unix(r.StartAt, 0)
	rate.RateToColor

	return ProblemIndex{
		ProblemID:      r.ProblemID,
		ProblemURL:     r.ProblemURL,
		ProblemTitle:   r.ProblemTitle,
		ContestID:      r.ContestID,
		ContestURL:     contestURL,
		ContestTitle:   r.ContestTitle,
		Color:          color,
		StartAt:        startAt,
		Duration:       r.Duration,
		RateChange:     r.RateChange,
		Category:       r.Category,
		Difficulty:     r.Difficulty,
		IsExperimental: r.IsExperimental,
		StatementJa:    statementJa,
		StatementEn:    statementEn,
	}, nil
}

type ProblemIndex struct {
	ProblemID      string
	ProblemTitle   string
	ProblemURL     string
	ContestID      string
	ContestURL     string
	ContestTitle   string
	Color          string
	StartAt        time.Time
	Duration       int64
	RateChange     string
	Category       string
	Difficulty     *int
	IsExperimental bool
	StatementJa    []string
	StatementEn    []string
}
