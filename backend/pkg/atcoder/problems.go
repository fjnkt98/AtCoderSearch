package atcoder

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/goark/errs"
	"golang.org/x/exp/slices"
)

const AGC001_STARTED_AT = 1468670400

type AtCoderProblemsClient interface {
	FetchContests(ctx context.Context) ([]Contest, error)
	FetchProblems(ctx context.Context) ([]Problem, error)
	FetchDifficulties(ctx context.Context) (map[string]Difficulty, error)
}

type atcoderProblemsClient struct {
	client *http.Client
}

func NewAtCoderProblemsClient() AtCoderProblemsClient {
	return &atcoderProblemsClient{
		client: &http.Client{},
	}
}

func (c *atcoderProblemsClient) FetchContests(ctx context.Context) ([]Contest, error) {
	uri := "https://kenkoooo.com/atcoder/resources/contests.json"
	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		return nil, errs.New(
			"failed to create request",
			errs.WithCause(err),
			errs.WithContext("uri", uri),
		)
	}
	req.Header.Set("ACCEPT_ENCODING", "gzip")

	res, err := c.client.Do(req)
	if err != nil {
		return nil, errs.New(
			"request failed",
			errs.WithCause(err),
			errs.WithContext("uri", uri),
		)
	}
	defer res.Body.Close()

	var contests []Contest
	if err := json.NewDecoder(res.Body).Decode(&contests); err != nil {
		return nil, errs.New(
			"failed to decode JSON into Contest",
			errs.WithCause(err),
			errs.WithContext("uri", uri),
		)
	}

	return contests, nil
}

func (c *atcoderProblemsClient) FetchProblems(ctx context.Context) ([]Problem, error) {
	uri := "https://kenkoooo.com/atcoder/resources/problems.json"
	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		return nil, errs.New(
			"failed to create request",
			errs.WithCause(err),
			errs.WithContext("uri", uri),
		)
	}
	req.Header.Set("ACCEPT_ENCODING", "gzip")

	res, err := c.client.Do(req)
	if err != nil {
		return nil, errs.New(
			"request failed",
			errs.WithCause(err),
			errs.WithContext("uri", uri),
		)
	}
	defer res.Body.Close()

	var problems []Problem
	if err := json.NewDecoder(res.Body).Decode(&problems); err != nil {
		return nil, errs.New(
			"failed to decode JSON into Contest",
			errs.WithCause(err),
			errs.WithContext("uri", uri),
		)
	}

	return problems, nil
}

func (c *atcoderProblemsClient) FetchDifficulties(ctx context.Context) (map[string]Difficulty, error) {
	uri := "https://kenkoooo.com/atcoder/resources/problem-models.json"
	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		return nil, errs.New(
			"failed to create request",
			errs.WithCause(err),
			errs.WithContext("uri", uri),
		)
	}
	req.Header.Set("ACCEPT_ENCODING", "gzip")

	res, err := c.client.Do(req)
	if err != nil {
		return nil, errs.New(
			"request failed",
			errs.WithCause(err),
			errs.WithContext("uri", uri),
		)
	}
	defer res.Body.Close()

	var difficulties map[string]Difficulty
	if err := json.NewDecoder(res.Body).Decode(&difficulties); err != nil {
		return nil, errs.New(
			"failed to decode JSON into Contest",
			errs.WithCause(err),
			errs.WithContext("uri", uri),
		)
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

func (c *Contest) RatedTarget() string {
	if c.StartEpochSecond < AGC001_STARTED_AT {
		return "unrated"
	}

	switch c.RateChange {
	case "-":
		return "unrated"
	case "All":
		return "all"
	default:
		rateRange := strings.Split(c.RateChange, "~")
		for i, word := range rateRange {
			rateRange[i] = strings.TrimSpace(word)
		}
		if len(rateRange) != 2 {
			return "unrated"
		}

		if _, err := strconv.Atoi(rateRange[0]); err == nil {
			return "lowerbound"
		}
		if _, err := strconv.Atoi(rateRange[1]); err == nil {
			return "upperbound"
		}
	}

	return "unrated"
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

	switch c.RatedTarget() {
	case "all":
		return "AGC-Like"
	case "upperbound":
		return "ABC-Like"
	case "lowerbound":
		return "ARC-Like"
	case "unrated":
		if strings.HasPrefix(c.ID, "past") {
			return "PAST"
		}
		if strings.HasPrefix(c.ID, "joi") {
			return "JOI"
		}
		if matched, _ := regexp.Match(`^(jag|JAG)`, []byte(c.ID)); matched {
			return "JOI"
		}
		isMarathon1, _ := regexp.Match(`(^Chokudai self|ハーフマラソン|^HACK TO THE FUTURE|Asprova|Heuristics Contest)`, []byte(c.ID))
		isMarathon2, _ := regexp.Match(`(^future-meets-you-self|^hokudai-hitachi)`, []byte(c.ID))
		isMarathon3 := slices.Contains([]string{"genocon2021", "stage0-2021", "caddi2019", "pakencamp-2019-day2", "kuronekoyamato-self2019", "wn2017_1"}, c.ID)
		if isMarathon1 || isMarathon2 || isMarathon3 {
			return "Marathon"
		}

		isSponsored1, _ := regexp.Match(`ドワンゴ|^Mujin|SoundHound|^codeFlyer|^COLOCON|みんなのプロコン|CODE THANKS FESTIVAL`, []byte(c.Title))
		isSponsored2, _ := regexp.Match(`(CODE FESTIVAL|^DISCO|日本最強プログラマー学生選手権|全国統一プログラミング王|Indeed)`, []byte(c.Title))
		isSponsored3, _ := regexp.Match(`(^Donuts|^dwango|^DigitalArts|^Code Formula|天下一プログラマーコンテスト)`, []byte(c.Title))
		if isSponsored1 || isSponsored2 || isSponsored3 {
			return "Other Sponsored"
		}

		return "Other Contests"
	default:
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
