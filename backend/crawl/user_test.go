package crawl

import (
	"context"
	"errors"
	"fjnkt98/atcodersearch/atcoder"
	"fjnkt98/atcodersearch/internal/testutil"
	"fjnkt98/atcodersearch/repository"
	"testing"
	"time"

	"k8s.io/utils/ptr"
)

func TestUser(t *testing.T) {
	_, dsn, stop, err := testutil.CreateDBContainer()
	t.Cleanup(func() { stop() })

	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	pool, err := repository.NewPool(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("save empty users", func(t *testing.T) {

		users := make([]atcoder.User, 0)
		count, err := SaveUsers(ctx, pool, users)
		if err != nil {
			t.Fatal(err)
		}

		if count != 0 {
			t.Errorf("count = %d, want 0", count)
		}
	})

	t.Run("save single user", func(t *testing.T) {

		users := []atcoder.User{
			{
				UserID:        "tourist",
				Rating:        3863,
				HighestRating: 4229,
				Affiliation:   ptr.To("ITMO University"),
				BirthYear:     ptr.To(int32(1994)),
				Country:       ptr.To("BY"),
				Crown:         ptr.To("crown_champion"),
				JoinCount:     59,
				Rank:          1,
				ActiveRank:    ptr.To(int32(1)),
				Wins:          22,
			},
		}

		count, err := SaveUsers(ctx, pool, users)
		if err != nil {
			t.Fatal(err)
		}
		if count != 1 {
			t.Errorf("count = %d, want 1", count)
		}
	})

	t.Run("save multiple user", func(t *testing.T) {

		users := []atcoder.User{
			{
				UserID:        "tourist",
				Rating:        3863,
				HighestRating: 4229,
				Affiliation:   ptr.To("ITMO University"),
				BirthYear:     ptr.To(int32(1994)),
				Country:       ptr.To("BY"),
				Crown:         ptr.To("crown_champion"),
				JoinCount:     59,
				Rank:          1,
				ActiveRank:    ptr.To(int32(1)),
				Wins:          22,
			},
			{
				UserID:        "w4yneb0t",
				Rating:        3710,
				HighestRating: 3802,
				Affiliation:   ptr.To("ETH Zurich"),
				BirthYear:     nil,
				Country:       ptr.To("CH"),
				Crown:         nil,
				JoinCount:     21,
				Rank:          2,
				ActiveRank:    nil,
				Wins:          2,
			},
		}

		count, err := SaveUsers(ctx, pool, users)
		if err != nil {
			t.Fatal(err)
		}
		if count != 2 {
			t.Errorf("count = %d, want 2", count)
		}
	})

	t.Run("crawl users(success)", func(t *testing.T) {

		crawler := NewUserCrawler(
			&DummyAtCoderClientS{},
			pool,
			time.Second,
		)
		if err := crawler.Crawl(ctx); err != nil {
			t.Errorf("expected no error, but got %v", err)
		}
	})

	t.Run("crawl users(fail)", func(t *testing.T) {

		crawler := NewUserCrawler(
			&DummyAtCoderClientF{},
			pool,
			time.Second,
		)
		if err := crawler.Crawl(ctx); !errors.Is(err, ErrDummy) {
			t.Errorf("expected ErrDummy, but got %#v", err)
		}
	})
}
