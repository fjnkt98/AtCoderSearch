package atcoder

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

const AGC001_STARTED_AT = 1468670400

var joiPattern = regexp.MustCompile(`^(jag|JAG)`)
var marathonPattern1 = regexp.MustCompile(`(^Chokudai self|ハーフマラソン|^HACK TO THE FUTURE|Asprova|Heuristics Contest)`)
var marathonPattern2 = regexp.MustCompile(`(^future-meets-you-self|^hokudai-hitachi)`)
var marathonPattern3 = []string{"genocon2021", "stage0-2021", "caddi2019", "pakencamp-2019-day2", "kuronekoyamato-self2019", "wn2017_1"}

var sponsoredPattern1 = regexp.MustCompile(`ドワンゴ|^Mujin|SoundHound|^codeFlyer|^COLOCON|みんなのプロコン|CODE THANKS FESTIVAL`)
var sponsoredPattern2 = regexp.MustCompile(`(CODE FESTIVAL|^DISCO|日本最強プログラマー学生選手権|全国統一プログラミング王|Indeed)`)
var sponsoredPattern3 = regexp.MustCompile(`(^Donuts|^dwango|^DigitalArts|^Code Formula|天下一プログラマーコンテスト)`)

type AtCoderProblemsClient struct {
	client *http.Client
}

func NewAtCoderProblemsClient() *AtCoderProblemsClient {
	client := &http.Client{}
	return NewAtCoderProblemsClientWithHTTPClient(client)
}

func NewAtCoderProblemsClientWithHTTPClient(client *http.Client) *AtCoderProblemsClient {
	return &AtCoderProblemsClient{
		client,
	}
}

func (c *AtCoderProblemsClient) FetchContests(ctx context.Context) ([]Contest, error) {
	uri := "https://kenkoooo.com/atcoder/resources/contests.json"
	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("ACCEPT_ENCODING", "gzip")

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to fetch contests: %w", err)
	}
	defer res.Body.Close()

	var contests []Contest
	if err := json.NewDecoder(res.Body).Decode(&contests); err != nil {
		return nil, fmt.Errorf("decode contests json: %w", err)
	}

	return contests, nil
}

func (c *AtCoderProblemsClient) FetchProblems(ctx context.Context) ([]Problem, error) {
	uri := "https://kenkoooo.com/atcoder/resources/problems.json"
	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("ACCEPT_ENCODING", "gzip")

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to fetch problems: %w", err)
	}
	defer res.Body.Close()

	var problems []Problem
	if err := json.NewDecoder(res.Body).Decode(&problems); err != nil {
		return nil, fmt.Errorf("decode problems json: %w", err)
	}

	return problems, nil
}

func (c *AtCoderProblemsClient) FetchDifficulties(ctx context.Context) (map[string]Difficulty, error) {
	uri := "https://kenkoooo.com/atcoder/resources/problem-models.json"
	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("ACCEPT_ENCODING", "gzip")

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to fetch difficulties: %w", err)
	}
	defer res.Body.Close()

	var difficulties map[string]Difficulty
	if err := json.NewDecoder(res.Body).Decode(&difficulties); err != nil {
		return nil, fmt.Errorf("decode difficulties json: %w", err)
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

const (
	UNRATED     = "unrated"
	ALL         = "all"
	RANGE       = "range"
	UPPER_BOUND = "upper_bound"
	LOWER_BOUND = "lower_bound"
)

type RatedTarget struct {
	Type string
	From *int
	To   *int
}

func (c *Contest) RatedTarget() RatedTarget {
	if c.StartEpochSecond < AGC001_STARTED_AT {
		return RatedTarget{UNRATED, nil, nil}
	}

	switch c.RateChange {
	case "-":
		return RatedTarget{UNRATED, nil, nil}
	case "All":
		return RatedTarget{ALL, nil, nil}
	default:
		rateRange := strings.Split(c.RateChange, " ~ ")
		if len(rateRange) != 2 {
			return RatedTarget{UNRATED, nil, nil}
		}

		target := RatedTarget{}
		if from, err := strconv.Atoi(rateRange[0]); err == nil {
			target.From = &from
		}
		if to, err := strconv.Atoi(rateRange[1]); err == nil {
			target.To = &to
		}

		if target.From == nil {
			if target.To == nil {
				target.Type = ALL
			} else {
				target.Type = UPPER_BOUND
			}
		} else {
			if target.To == nil {
				target.Type = LOWER_BOUND
			} else {
				target.Type = RANGE
			}
		}
		return target
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

	if target := c.RatedTarget(); target.Type == ALL {
		return "AGC-Like"
	} else if target.Type == UPPER_BOUND {
		return "ABC-Like"
	} else if target.Type == LOWER_BOUND {
		return "ARC-Like"
	} else if target.Type == UNRATED {
		if strings.HasPrefix(c.ID, "past") {
			return "PAST"
		}
		if strings.HasPrefix(c.ID, "joi") {
			return "JOI"
		}
		if joiPattern.MatchString(c.ID) {
			return "JOI"
		}
		isMarathon1 := marathonPattern1.MatchString(c.ID)
		isMarathon2 := marathonPattern2.MatchString(c.ID)
		isMarathon3 := slices.Contains(marathonPattern3, c.ID)
		if isMarathon1 || isMarathon2 || isMarathon3 {
			return "Marathon"
		}

		isSponsored1 := sponsoredPattern1.MatchString(c.Title)
		isSponsored2 := sponsoredPattern2.MatchString(c.Title)
		isSponsored3 := sponsoredPattern3.MatchString(c.Title)

		if isSponsored1 || isSponsored2 || isSponsored3 {
			return "Other Sponsored"
		}

		return "Other Contests"
	} else {
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
