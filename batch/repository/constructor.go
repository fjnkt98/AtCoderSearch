package repository

import (
	"slices"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/fjnkt98/atcodersearch-batch/pkg/atcoder"
)

func NewContest(c atcoder.Contest) Contest {
	return Contest{
		ContestID:        c.ID,
		StartEpochSecond: c.StartEpochSecond,
		DurationSecond:   c.DurationSecond,
		Title:            c.Title,
		RateChange:       c.RateChange,
		Category:         c.Categorize(),
	}
}

func NewContests(contests []atcoder.Contest) []Contest {
	result := make([]Contest, len(contests))
	for i, c := range contests {
		result[i] = NewContest(c)
	}
	return result
}

func NewDifficulty(problemID string, d atcoder.Difficulty) Difficulty {
	return Difficulty{
		ProblemID:        problemID,
		Slope:            d.Slope,
		Intercept:        d.Intercept,
		Variance:         d.Variance,
		Difficulty:       d.Difficulty,
		Discrimination:   d.Discrimination,
		IrtLoglikelihood: d.IrtLogLikelihood,
		IrtUsers:         d.IrtUsers,
		IsExperimental:   d.IsExperimental,
	}
}

func NewDifficulties(difficulties map[string]atcoder.Difficulty) []Difficulty {
	result := make([]Difficulty, 0, len(difficulties))
	for problemID, d := range difficulties {
		result = append(result, NewDifficulty(problemID, d))
	}
	return result
}

func NewSubmission(s atcoder.Submission) Submission {
	return Submission{
		ID:            s.ID,
		EpochSecond:   s.EpochSecond,
		ProblemID:     s.ProblemID,
		ContestID:     &s.ContestID,
		UserID:        &s.UserID,
		Language:      &s.Language,
		Point:         &s.Point,
		Length:        &s.Length,
		Result:        &s.Result,
		ExecutionTime: s.ExecutionTime,
	}
}

func NewSubmissions(submissions []atcoder.Submission) []Submission {
	result := make([]Submission, 0, len(submissions))
	ids := mapset.NewSet[int64]()
	for i := len(submissions) - 1; i >= 0; i-- {
		if ids.Contains(submissions[i].ID) {
			continue
		}
		ids.Add(submissions[i].ID)
		result = append(result, NewSubmission(submissions[i]))
	}
	slices.Reverse(result)
	return result
}

func NewUser(u atcoder.User) User {
	return User{
		UserID:        u.UserID,
		Rating:        u.Rating,
		HighestRating: u.HighestRating,
		Affiliation:   u.Affiliation,
		BirthYear:     u.BirthYear,
		Country:       u.Country,
		Crown:         u.Crown,
		JoinCount:     u.JoinCount,
		Rank:          u.Rank,
		ActiveRank:    u.ActiveRank,
		Wins:          u.Wins,
	}
}

func NewUsers(users []atcoder.User) []User {
	result := make([]User, len(users))
	for i, u := range users {
		result[i] = NewUser(u)
	}
	return result
}
