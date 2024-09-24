package repository

import (
	"fjnkt98/atcodersearch/pkg/atcoder"
	"iter"
	"maps"
	"slices"
	"time"
)

func Map[S, T any, C func(S, time.Time) T](c C, src iter.Seq[S], updateAt time.Time) iter.Seq[T] {
	return func(yield func(T) bool) {
		for s := range src {
			if !yield(c(s, updateAt)) {
				break
			}
		}
	}
}

func NewContest(c atcoder.Contest, updatedAt time.Time) Contest {
	return Contest{
		ContestID:        c.ID,
		StartEpochSecond: c.StartEpochSecond,
		DurationSecond:   c.DurationSecond,
		Title:            c.Title,
		RateChange:       c.RateChange,
		Category:         c.Categorize(),
		UpdatedAt:        updatedAt,
	}
}

func NewDifficulties(d map[string]atcoder.Difficulty, updatedAt time.Time) []Difficulty {
	difficulties := make([]Difficulty, 0, len(d))
	for _, problemID := range slices.Sorted(maps.Keys(d)) {
		difficulty := d[problemID]

		difficulties = append(difficulties, Difficulty{
			ProblemID:        problemID,
			Slope:            difficulty.Slope,
			Intercept:        difficulty.Intercept,
			Variance:         difficulty.Variance,
			Difficulty:       difficulty.Difficulty,
			Discrimination:   difficulty.Discrimination,
			IrtLoglikelihood: difficulty.IrtLoglikelihood,
			IrtUsers:         difficulty.IrtUsers,
			IsExperimental:   difficulty.IsExperimental,
			UpdatedAt:        updatedAt,
		})
	}

	return difficulties
}

func NewSubmission(s atcoder.Submission, updatedAt time.Time) Submission {
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
		UpdatedAt:     updatedAt,
	}
}

func NewUser(u atcoder.User, updatedAt time.Time) User {
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
		UpdatedAt:     updatedAt,
	}
}
