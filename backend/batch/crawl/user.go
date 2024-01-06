package crawl

import (
	"context"
	"fjnkt98/atcodersearch/batch"
	"fjnkt98/atcodersearch/pkg/atcoder"
	"fjnkt98/atcodersearch/repository"
	"time"

	"github.com/goark/errs"
	"golang.org/x/exp/slog"
)

type UserCrawler interface {
	batch.Batch
	CrawlUser(ctx context.Context) error
}

type userCrawler struct {
	client atcoder.AtCoderClient
	repo   repository.UserRepository
	config userCrawlerConfig
}

type userCrawlerConfig struct {
	Duration int `json:"duration"`
}

func NewUserCrawler(
	client atcoder.AtCoderClient,
	repo repository.UserRepository,
	duration int,
) *userCrawler {
	return &userCrawler{
		client: client,
		repo:   repo,
		config: userCrawlerConfig{
			Duration: duration,
		},
	}
}

func (c *userCrawler) Name() string {
	return "UserCrawler"
}

func (c *userCrawler) Config() any {
	return c.config
}

func (c *userCrawler) CrawlUser(ctx context.Context) error {
	slog.Info("Start to crawl users.")

	allUsers := make([]repository.User, 0)

loop:
	for i := 0; ; i++ {
		users, err := c.client.FetchUsers(ctx, i)
		if err != nil {
			return errs.Wrap(err)
		}

		if len(users) == 0 {
			break loop
		}

		allUsers = append(allUsers, convertUsers(users)...)

		time.Sleep(time.Duration(c.config.Duration) * time.Millisecond)
	}

	if err := c.repo.Save(ctx, allUsers); err != nil {
		return errs.Wrap(err)
	}

	return nil
}

func (c *userCrawler) Run(ctx context.Context) error {
	return c.CrawlUser(ctx)
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
