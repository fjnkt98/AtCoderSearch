package repository

import (
	"context"
	"errors"
	"fjnkt98/atcodersearch/internal/testutil"
	"testing"
	"time"
)

func TestHistory(t *testing.T) {
	_, dsn, stop, err := testutil.CreateDBContainer()
	t.Cleanup(func() { stop() })

	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	pool, err := NewPool(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("create batch history", func(t *testing.T) {
		t.Parallel()

		history, err := NewBatchHistory(ctx, pool, "test", nil)
		if err != nil {
			t.Fatal(err)
		}

		if history.Name != "test" {
			t.Errorf("history.Name = %s, want test", history.Name)
		}
		if history.Status != HistoryStatusWorking {
			t.Errorf("history.Status = %s, want %s", history.Status, HistoryStatusWorking)
		}
		if history.FinishedAt != nil {
			t.Errorf("history.FinishedAt = %v, want nil", history.FinishedAt)
		}
	})

	t.Run("complete batch history", func(t *testing.T) {
		t.Parallel()

		history, err := NewBatchHistory(ctx, pool, "test", nil)
		if err != nil {
			t.Fatal(err)
		}

		if err := history.Complete(ctx, pool); err != nil {
			t.Error(err)
		}
		if history.Status != HistoryStatusCompleted {
			t.Errorf("history.Status = %s, want %s", history.Status, HistoryStatusCompleted)
		}
		if history.FinishedAt == nil {
			t.Errorf("history.FinishedAt must be registered")
		}
	})

	t.Run("abort batch history", func(t *testing.T) {
		t.Parallel()

		history, err := NewBatchHistory(ctx, pool, "test", nil)
		if err != nil {
			t.Fatal(err)
		}

		if err := history.Abort(ctx, pool); err != nil {
			t.Error(err)
		}
		if history.Status != HistoryStatusAborted {
			t.Errorf("history.Status = %s, want %s", history.Status, HistoryStatusAborted)
		}
		if history.FinishedAt == nil {
			t.Errorf("history.FinishedAt must be registered")
		}
	})

	t.Run("abort completed batch history", func(t *testing.T) {
		t.Parallel()

		history, err := NewBatchHistory(ctx, pool, "test", nil)
		if err != nil {
			t.Fatal(err)
		}

		if err := history.Complete(ctx, pool); err != nil {
			t.Error(err)
		}

		if err := history.Abort(ctx, pool); !errors.Is(err, ErrHistoryConfirmed) {
			t.Errorf("err = %v, want ErrHistoryConfirmed", err)
		}
	})

	t.Run("create crawl history", func(t *testing.T) {
		t.Parallel()

		history, err := NewCrawlHistory(ctx, pool, "test001")
		if err != nil {
			t.Fatal(err)
		}

		if history.ContestID != "test001" {
			t.Errorf("history.ContestID = %s, want test001", history.ContestID)
		}
		if history.Status != HistoryStatusWorking {
			t.Errorf("history.Status = %s, want %s", history.Status, HistoryStatusWorking)
		}
		if history.FinishedAt != nil {
			t.Errorf("history.FinishedAt = %v, want nil", history.FinishedAt)
		}
	})

	t.Run("complete crawl history", func(t *testing.T) {
		t.Parallel()

		history, err := NewCrawlHistory(ctx, pool, "test002")
		if err != nil {
			t.Fatal(err)
		}

		if err := history.Complete(ctx, pool); err != nil {
			t.Error(err)
		}
		if history.Status != HistoryStatusCompleted {
			t.Errorf("history.Status = %s, want %s", history.Status, HistoryStatusCompleted)
		}
		if history.FinishedAt == nil {
			t.Errorf("history.FinishedAt must be registered")
		}
	})

	t.Run("abort crawl history", func(t *testing.T) {
		t.Parallel()

		history, err := NewCrawlHistory(ctx, pool, "test003")
		if err != nil {
			t.Fatal(err)
		}

		if err := history.Abort(ctx, pool); err != nil {
			t.Error(err)
		}
		if history.Status != HistoryStatusAborted {
			t.Errorf("history.Status = %s, want %s", history.Status, HistoryStatusAborted)
		}
		if history.FinishedAt == nil {
			t.Errorf("history.FinishedAt must be registered")
		}
	})

	t.Run("abort completed crawl history", func(t *testing.T) {
		t.Parallel()

		history, err := NewCrawlHistory(ctx, pool, "test004")
		if err != nil {
			t.Fatal(err)
		}

		if err := history.Complete(ctx, pool); err != nil {
			t.Error(err)
		}

		if err := history.Abort(ctx, pool); !errors.Is(err, ErrHistoryConfirmed) {
			t.Errorf("err = %v, want ErrHistoryConfirmed", err)
		}
	})

	t.Run("fetch latest crawl history", func(t *testing.T) {
		t.Parallel()

		history1, err := NewCrawlHistory(ctx, pool, "test005")
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(1 * time.Second)

		_, err = NewCrawlHistory(ctx, pool, "test005")
		if err != nil {
			t.Fatal(err)
		}

		if err := history1.Complete(ctx, pool); err != nil {
			t.Error(err)
		}

		latest, err := FetchLatestCrawlHistory(ctx, pool, "test005")
		if err != nil {
			t.Error(err)
		}

		if latest.ID != history1.ID {
			t.Errorf("latest.ID = %d, want %d", latest.ID, history1.ID)
		}
	})
}
