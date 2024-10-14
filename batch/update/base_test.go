package update

import (
	"context"
	"fjnkt98/atcodersearch/internal/testutil"
	"fjnkt98/atcodersearch/repository"
	"testing"

	"github.com/meilisearch/meilisearch-go"
)

func TestUpdateIndex(t *testing.T) {
	_, dsn, stopDB, err := testutil.CreateDBContainer()
	if err != nil {
		t.Fatal(err)
	}
	_, url, key, stopEngine, err := testutil.CreateEngineContainer()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		stopDB()
		stopEngine()
	})

	ctx := context.Background()
	pool, err := repository.NewPool(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	client := meilisearch.New(url, meilisearch.WithAPIKey(key))

	t.Run("test update problem", func(t *testing.T) {
		sql := `
INSERT INTO "problems" ("problem_id", "contest_id", "problem_index", "name", "title", "url", "html")
VALUES
    ('sample_contest_A', 'sample_contest', 'A', 'sample problem 1', 'A. sample problem 1', '', '<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"><title>sample problem 1</title></head><body></body></html>'),
    ('sample_contest_B', 'sample_contest', 'B', 'sample problem 2', 'B. sample problem 2', '', '<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"><title>sample problem 2</title></head><body></body></html>');

INSERT INTO "contests" ("contest_id", "start_epoch_second", "duration_second", "title", "rate_change", "category")
VALUES
	('sample_contest', 0, 0, 'sample contest', '', '');
`

		if _, err := pool.Exec(ctx, sql); err != nil {
			t.Fatal(err)
		}

		reader := NewProblemRowReader(pool)
		indexer := NewProblemIndexer(client)

		if err := UpdateIndex(
			ctx,
			reader,
			indexer,
			10,
			1,
		); err != nil {
			t.Error(err)
		}

		stats, err := client.GetStats()
		if err != nil {
			t.Fatal(err)
		}

		indexStats, ok := stats.Indexes["problems"]
		if !ok {
			t.Error("problems index not found")
		}
		if n := indexStats.NumberOfDocuments; n != 2 {
			t.Errorf("expected number of documents is 2, but got %d", n)
		}
	})

	t.Run("test update user", func(t *testing.T) {
		sql := `
INSERT INTO "users" ("user_id", "rating", "highest_rating", "affiliation", "birth_year", "country", "crown", "join_count", "rank", "active_rank", "wins")
VALUES
	('wzp', 2563, 2563, NULL, NULL, 'CN', 'user-orange-2', 23, 470, NULL, 0), 
	('tourist', 3774, 4229, 'ITMOUniversity', 1994, 'BY', 'crown_champion', 61, 1, 1, 22);
`
		if _, err := pool.Exec(ctx, sql); err != nil {
			t.Fatal(err)
		}

		reader := NewUserRowReader(pool)
		indexer := NewUserIndexer(client)

		if err := UpdateIndex(
			ctx,
			reader,
			indexer,
			10,
			1,
		); err != nil {
			t.Error(err)
		}

		stats, err := client.GetStats()
		if err != nil {
			t.Fatal(err)
		}

		indexStats, ok := stats.Indexes["users"]
		if !ok {
			t.Error("users index not found")
		}
		if n := indexStats.NumberOfDocuments; n != 2 {
			t.Errorf("expected number of documents is 2, but got %d", n)
		}
	})
}
