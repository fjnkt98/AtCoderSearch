package atcoder

import (
	"fmt"
	"net/http"
	"time"
)

const (
	UNRATED     = "unrated"
	ALL         = "all"
	RANGE       = "range"
	UPPER_BOUND = "upper_bound"
	LOWER_BOUND = "lower_bound"
)

type RatedTarget struct {
	Kind string
	From *int
	To   *int
}

type AtCoderProblemsClient struct {
	client *http.Client
}

func NewAtCoderProblemsClient() (*AtCoderProblemsClient, error) {
	client := &http.Client{
		Timeout: time.Duration(30) * time.Second,
	}

	return &AtCoderProblemsClient{client}, nil
}

func (c *AtCoderProblemsClient) FetchContests() ([]Contest, error) {
	panic("not implemented")
}

func (c *AtCoderProblemsClient) FetchProblems() ([]Problem, error) {
	panic("not implemented")
}

func (c *AtCoderProblemsClient) FetchDifficulties() (map[string]Difficulty, error) {
	panic("not implemented")
}

type Contest struct {
	ID               string `json:"id"`
	StartEpochSecond int64  `json:"start_epoch_second"`
	DurationSecond   int64  `json:"duration_second"`
	Title            string `json:"title"`
	RateChange       string `json:"rate_change"`
}

func (c *Contest) RatedTarget() RatedTarget {
	panic("not implemented")
}

func (c *Contest) Categorize() string {
	panic("not implemented")
}

type Problem struct {
	ID           string `json:"id"`
	ContestID    string `json:"contest_id"`
	ProblemIndex string `json:"problem_index"`
	Name         string `json:"name"`
	Title        string `json:"title"`
}

func (p *Problem) URL() string {
	return fmt.Sprintf("https://atcoder.jp/contests/%s/tasks/%s", p.ContestID, p.ID)
}

type Difficulty struct {
	Slope            *float64 `json:"slope"`
	Intercept        *float64 `json:"intercept"`
	Variance         *float64 `json:"variance"`
	Difficulty       *int64   `json:"difficulty"`
	Discrimination   *float64 `json:"discrimination"`
	IrtLogLikelihood *float64 `json:"irt_loglikelihood"`
	IrtUsers         *float64 `json:"irt_users"`
	IsExperimental   *bool    `json:"is_experimental"`
}
