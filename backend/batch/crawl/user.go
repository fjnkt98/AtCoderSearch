package crawl

import (
	"context"
	"fjnkt98/atcodersearch/pkg/atcoder"
	"fjnkt98/atcodersearch/repository"
	"time"

	"log/slog"

	"github.com/goark/errs"
)

type UserCrawler interface {
	CrawlUser(ctx context.Context) error
}

type userCrawler struct {
	client   atcoder.AtCoderClient
	repo     repository.UserRepository
	duration time.Duration
}

func NewUserCrawler(
	client atcoder.AtCoderClient,
	repo repository.UserRepository,
	duration time.Duration,
) *userCrawler {
	return &userCrawler{
		client:   client,
		repo:     repo,
		duration: duration,
	}
}

func (c *userCrawler) CrawlUser(ctx context.Context) error {
	slog.Info("Start to crawl users.")

	allUsers := make([]repository.User, 0)

loop:
	for i := 1; ; i++ {
		slog.Info("Crawl users", slog.Int("page", i))
		users, err := c.client.FetchUsers(ctx, i)
		if err != nil {
			return errs.Wrap(err)
		}

		if len(users) == 0 {
			slog.Info("There is no more crawl target.")
			break loop
		}

		allUsers = append(allUsers, convertUsers(users)...)

		time.Sleep(c.duration)
	}

	if err := c.repo.Save(ctx, allUsers); err != nil {
		return errs.Wrap(err)
	}

	slog.Info("Finish crawling users successfully.")
	return nil
}

func convertUser(user atcoder.User) repository.User {
	return repository.User{
		UserName:      user.UserName,
		Rating:        user.Rating,
		HighestRating: user.HighestRating,
		Affiliation:   user.Affiliation,
		BirthYear:     user.BirthYear,
		Country:       user.Country,
		Crown:         user.Crown,
		JoinCount:     user.JoinCount,
		Rank:          user.Rank,
		ActiveRank:    user.ActiveRank,
		Wins:          user.Wins,
	}
}

func convertUsers(users []atcoder.User) []repository.User {
	result := make([]repository.User, len(users))
	for i, user := range users {
		result[i] = convertUser(user)
	}

	return result
}
