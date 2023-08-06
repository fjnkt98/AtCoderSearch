package problem

import (
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

const AGC001_STARTED_AT = 1468670400

type ContestJSON struct {
	ID               string `json:"id"`
	StartEpochSecond int64  `json:"start_epoch_second"`
	DurationSecond   int64  `json:"duration_second"`
	Title            string `json:"title"`
	RateChange       string `json:"rate_change"`
}

func (c *ContestJSON) RatedTarget() string {
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

func (c *ContestJSON) Categorize() string {
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

type Contest struct {
	ContestID        string `db:"contest_id"`
	StartEpochSecond int64  `db:"start_epoch_second"`
	DurationSecond   int64  `db:"duration_second"`
	Title            string `db:"title"`
	RateChange       string `db:"rate_change"`
	Category         string `db:"category"`
}

// JSON schema of problem information which can be retrieved from AtCoder Problems.
//
// - id: problem id which be used in problem URL.
// - contest_id: contest id in which the problem belong.
// - problem_index: problem index, such as `A`, `B`, or `Ex`.
// - name: problem name
// - title: string consisting of the problem index and problem name
type ProblemJSON struct {
	ID           string `json:"id"`
	ContestID    string `json:"contest_id"`
	ProblemIndex string `json:"problem_index"`
	Name         string `json:"name"`
	Title        string `json:"title"`
}

type DifficultyJSON struct {
	Slope            *float64 `json:"slope"`
	Intercept        *float64 `json:"intercept"`
	Variance         *float64 `json:"variance"`
	Difficulty       *int64   `json:"difficulty"`
	Discrimination   *float64 `json:"discrimination"`
	IrtLogLikelihood *float64 `json:"irt_loglikelihood"`
	IrtUsers         *float64 `json:"irt_users"`
	IsExperimental   *bool    `json:"is_experimental"`
}

type Problem struct {
	ProblemID    string `db:"problem_id"`
	ContestID    string `db:"contest_id"`
	ProblemIndex string `db:"problem_index"`
	Name         string `db:"name"`
	Title        string `db:"title"`
	URL          string `db:"url"`
	HTML         string `db:"html"`
}

type Difficulty struct {
	ProblemID        string   `db:"problem_id"`
	Slope            *float64 `db:"slope"`
	Intercept        *float64 `db:"intercept"`
	Variance         *float64 `db:"variance"`
	Difficulty       *int64   `db:"difficulty"`
	Discrimination   *float64 `db:"discrimination"`
	IrtLogLikelihood *float64 `db:"irt_loglikelihood"`
	IrtUsers         *float64 `db:"irt_users"`
	IsExperimental   *bool    `db:"is_experimental"`
}

type User struct {
	UserName      string  `db:"user_name"`
	Rating        int     `db:"rating"`
	HighestRating int     `db:"highest_rating"`
	Affiliation   *string `db:"affiliation"`
	BirthYear     *int    `db:"birth_year"`
	Country       *string `db:"country"`
	Crown         *string `db:"crown"`
	JoinCount     int     `db:"join_count"`
	Rank          int     `db:"rank"`
	Wins          int     `db:"wins"`
}

type Submission struct {
	ID            int      `db:"id"`
	EpochSecond   int      `db:"epoch_second"`
	ProblemID     string   `db:"problem_id"`
	ContestID     *string  `db:"contest_id"`
	UserID        *string  `db:"user_id"`
	Language      *string  `db:"language"`
	Point         *float64 `db:"point"`
	Length        *int     `db:"length"`
	Result        *string  `db:"result"`
	ExecutionTime *int     `db:"execution_time"`
}
