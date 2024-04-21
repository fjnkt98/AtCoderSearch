// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package repository

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Contest struct {
	ContestID        string
	StartEpochSecond int64
	DurationSecond   int64
	Title            string
	RateChange       string
	Category         string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type Difficulty struct {
	ProblemID        string
	Slope            sql.NullFloat64
	Intercept        sql.NullFloat64
	Variance         sql.NullFloat64
	Difficulty       sql.NullInt32
	Discrimination   sql.NullFloat64
	IrtLoglikelihood sql.NullFloat64
	IrtUsers         sql.NullFloat64
	IsExperimental   sql.NullBool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type Language struct {
	Language string
	Group    sql.NullString
}

type Problem struct {
	ProblemID    string
	ContestID    string
	ProblemIndex string
	Name         string
	Title        string
	Url          string
	Html         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Submission struct {
	ID            int64
	EpochSecond   int64
	ProblemID     string
	ContestID     sql.NullString
	UserID        sql.NullString
	Language      sql.NullString
	Point         sql.NullFloat64
	Length        sql.NullInt32
	Result        sql.NullString
	ExecutionTime sql.NullInt32
	CrawledAt     time.Time
}

type SubmissionCrawlHistory struct {
	ID        int64
	ContestID string
	StartedAt int64
}

type UpdateHistory struct {
	ID         int64
	Domain     string
	StartedAt  time.Time
	FinishedAt sql.NullTime
	Status     sql.NullString
	Options    json.RawMessage
}

type User struct {
	UserName      string
	Rating        int32
	HighestRating int32
	Affiliation   sql.NullString
	BirthYear     sql.NullInt32
	Country       sql.NullString
	Crown         sql.NullString
	JoinCount     int32
	Rank          int32
	ActiveRank    sql.NullInt32
	Wins          int32
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
