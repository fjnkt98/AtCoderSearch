// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package repository

import (
	"time"
)

type BatchHistory struct {
	ID         int64      `db:"id"`
	Name       string     `db:"name"`
	StartedAt  time.Time  `db:"started_at"`
	FinishedAt *time.Time `db:"finished_at"`
	Status     string     `db:"status"`
	Options    []byte     `db:"options"`
}

type Contest struct {
	ContestID        string    `bulk:"unique" db:"contest_id"`
	StartEpochSecond int64     `db:"start_epoch_second"`
	DurationSecond   int64     `db:"duration_second"`
	Title            string    `db:"title"`
	RateChange       string    `db:"rate_change"`
	Category         string    `db:"category"`
	UpdatedAt        time.Time `db:"updated_at"`
}

type Difficulty struct {
	ProblemID        string    `bulk:"unique" db:"problem_id"`
	Slope            *float64  `db:"slope"`
	Intercept        *float64  `db:"intercept"`
	Variance         *float64  `db:"variance"`
	Difficulty       *int64    `db:"difficulty"`
	Discrimination   *float64  `db:"discrimination"`
	IrtLoglikelihood *float64  `db:"irt_loglikelihood"`
	IrtUsers         *float64  `db:"irt_users"`
	IsExperimental   *bool     `db:"is_experimental"`
	UpdatedAt        time.Time `db:"updated_at"`
}

type Language struct {
	Language string  `db:"language"`
	Group    *string `db:"group"`
}

type Problem struct {
	ProblemID    string    `bulk:"unique" db:"problem_id"`
	ContestID    string    `db:"contest_id"`
	ProblemIndex string    `db:"problem_index"`
	Name         string    `db:"name"`
	Title        string    `db:"title"`
	Url          string    `db:"url"`
	Html         string    `db:"html"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type Submission struct {
	ID            int64     `bulk:"unique" db:"id"`
	EpochSecond   int64     `db:"epoch_second"`
	ProblemID     string    `db:"problem_id"`
	ContestID     *string   `db:"contest_id"`
	UserID        *string   `db:"user_id"`
	Language      *string   `db:"language"`
	Point         *float64  `db:"point"`
	Length        *int32    `db:"length"`
	Result        *string   `db:"result"`
	ExecutionTime *int32    `db:"execution_time"`
	UpdatedAt     time.Time `db:"updated_at"`
}

type SubmissionCrawlHistory struct {
	ID        int64  `db:"id"`
	ContestID string `db:"contest_id"`
	StartedAt int64  `db:"started_at"`
}

type User struct {
	UserID        string    `bulk:"unique" db:"user_id"`
	Rating        int32     `db:"rating"`
	HighestRating int32     `db:"highest_rating"`
	Affiliation   *string   `db:"affiliation"`
	BirthYear     *int32    `db:"birth_year"`
	Country       *string   `db:"country"`
	Crown         *string   `db:"crown"`
	JoinCount     int32     `db:"join_count"`
	Rank          int32     `db:"rank"`
	ActiveRank    *int32    `db:"active_rank"`
	Wins          int32     `db:"wins"`
	UpdatedAt     time.Time `db:"updated_at"`
}
