package atcoder

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	UNRATED     = "unrated"
	ALL         = "all"
	RANGE       = "range"
	UPPER_BOUND = "upper_bound"
	LOWER_BOUND = "lower_bound"
)
const AGC001_STARTED_AT = 1468670400

var JOIPattern = regexp.MustCompile(`^(jag|JAG)`)
var MarathonPattern1 = regexp.MustCompile(`(^Chokudai self|ハーフマラソン|^HACK TO THE FUTURE|Asprova|Heuristics Contest)`)
var MarathonPattern2 = regexp.MustCompile(`(^future-meets-you-self|^hokudai-hitachi)`)
var MarathonPattern3 = regexp.MustCompile(`^(genocon2021|stage0-2021|caddi2019|pakencamp-2019-day2|kuronekoyamato-self2019|wn2017_1)$`)
var SponsoredPattern1 = regexp.MustCompile(`ドワンゴ|^Mujin|SoundHound|^codeFlyer|^COLOCON|みんなのプロコン|CODE THANKS FESTIVAL`)
var SponsoredPattern2 = regexp.MustCompile(`(CODE FESTIVAL|^DISCO|日本最強プログラマー学生選手権|全国統一プログラミング王|Indeed)`)
var SponsoredPattern3 = regexp.MustCompile(`(^Donuts|^dwango|^DigitalArts|^Code Formula|天下一プログラマーコンテスト)`)

type RatedTarget struct {
	Kind string
	From *int
	To   *int
}

type AtCoderProblemsClient interface {
	FetchContests(ctx context.Context) ([]Contest, error)
	FetchProblems(ctx context.Context) ([]Problem, error)
	FetchDifficulties(ctx context.Context) (map[string]Difficulty, error)
}

var _ AtCoderProblemsClient = (*atCoderProblemsClient)(nil)

type atCoderProblemsClient struct {
	client *http.Client
}

func NewAtCoderProblemsClient() *atCoderProblemsClient {
	client := &http.Client{
		Timeout: time.Duration(30) * time.Second,
	}

	return &atCoderProblemsClient{client}
}

func (c *atCoderProblemsClient) FetchContests(ctx context.Context) ([]Contest, error) {
	uri := "https://kenkoooo.com/atcoder/resources/contests.json"
	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("create new request: %w", err)
	}
	req.Header.Set("ACCEPT_ENCODING", "gzip")

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer res.Body.Close()

	var contests []Contest
	if err := json.NewDecoder(res.Body).Decode(&contests); err != nil {
		return nil, fmt.Errorf("decode response body: %w", err)
	}

	return contests, nil
}

func (c *atCoderProblemsClient) FetchProblems(ctx context.Context) ([]Problem, error) {
	uri := "https://kenkoooo.com/atcoder/resources/problems.json"
	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("create new request: %w", err)
	}
	req.Header.Set("ACCEPT_ENCODING", "gzip")

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer res.Body.Close()

	var problems []Problem
	if err := json.NewDecoder(res.Body).Decode(&problems); err != nil {
		return nil, fmt.Errorf("decode response body: %w", err)
	}

	return problems, nil
}

func (c *atCoderProblemsClient) FetchDifficulties(ctx context.Context) (map[string]Difficulty, error) {
	uri := "https://kenkoooo.com/atcoder/resources/problem-models.json"
	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("create new request: %w", err)
	}
	req.Header.Set("ACCEPT_ENCODING", "gzip")

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer res.Body.Close()

	var difficulties map[string]Difficulty
	if err := json.NewDecoder(res.Body).Decode(&difficulties); err != nil {
		return nil, fmt.Errorf("decode response body: %w", err)
	}

	return difficulties, nil
}

type Contest struct {
	ID               string `json:"id"`
	StartEpochSecond int64  `json:"start_epoch_second"`
	DurationSecond   int64  `json:"duration_second"`
	Title            string `json:"title"`
	RateChange       string `json:"rate_change"`
}

func (c *Contest) RatedTarget() RatedTarget {
	if c.StartEpochSecond < AGC001_STARTED_AT {
		return RatedTarget{Kind: UNRATED}
	}

	switch c.RateChange {
	case "-":
		return RatedTarget{Kind: UNRATED}
	case "All":
		return RatedTarget{Kind: ALL}
	default:
		splitted := strings.Split(c.RateChange, " ~ ")
		if len(splitted) != 2 {
			return RatedTarget{Kind: UNRATED}
		}

		var from *int
		if f, err := strconv.Atoi(splitted[0]); err == nil {
			from = &f
		}
		var to *int
		if t, err := strconv.Atoi(splitted[1]); err == nil {
			to = &t
		}

		if from == nil {
			if to == nil {
				return RatedTarget{Kind: UNRATED}
			} else {
				return RatedTarget{Kind: UPPER_BOUND, To: to}
			}
		} else {
			if to == nil {
				return RatedTarget{Kind: LOWER_BOUND, From: from}
			} else {
				return RatedTarget{Kind: RANGE, From: from, To: to}
			}
		}
	}
}

func (c *Contest) Categorize() string {
	if strings.HasPrefix(c.ID, "abc") {
		return "ABC"
	}
	if strings.HasPrefix(c.ID, "arc") {
		return "ARC"
	}
	if strings.HasPrefix(c.ID, "agc") {
		return "AGC"
	}
	if strings.HasPrefix(c.ID, "ahc") {
		return "AHC"
	}

	t := c.RatedTarget()
	switch t.Kind {
	case ALL:
		return "AGC-Like"
	case UPPER_BOUND:
		return "ABC-Like"
	case LOWER_BOUND:
		return "ARC-Like"
	default:
		if strings.HasPrefix(c.ID, "past") {
			return "PAST"
		}
		if strings.HasPrefix(c.ID, "joi") || JOIPattern.MatchString(c.ID) {
			return "JOI"
		}

		if MarathonPattern1.MatchString(c.ID) || MarathonPattern2.MatchString(c.ID) || MarathonPattern3.MatchString(c.ID) {
			return "Marathon"
		}

		if SponsoredPattern1.MatchString(c.Title) || SponsoredPattern2.MatchString(c.Title) || SponsoredPattern3.MatchString(c.Title) {
			return "Other Sponsored"
		}

		return "Other Contests"
	}
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
	IrtLoglikelihood *float64 `json:"irt_loglikelihood"`
	IrtUsers         *float64 `json:"irt_users"`
	IsExperimental   *bool    `json:"is_experimental"`
}
