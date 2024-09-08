package repository

import (
	"context"
	"fjnkt98/atcodersearch/internal/testutil"
	"testing"
)

func TestCreateAndUpdateBatchHistory(t *testing.T) {
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
		history, err := NewBatchHistory(ctx, pool, "test", nil)
		if err != nil {
			t.Fatal(err)
		}

		if history.Name != "test" {
			t.Errorf("history.Name = %s, want test", history.Name)
		}
		if history.Status != "working" {
			t.Errorf("history.Status = %s, want working", history.Status)
		}
		if history.FinishedAt != nil {
			t.Errorf("history.FinishedAt = %v, want nil", history.FinishedAt)
		}
	})

	t.Run("complete batch history", func(t *testing.T) {
		history, err := NewBatchHistory(ctx, pool, "test", nil)
		if err != nil {
			t.Fatal(err)
		}

		if err := history.Complete(ctx, pool); err != nil {
			t.Error(err)
		}
		if history.Status != "completed" {
			t.Errorf("history.Status = %s, want completed", history.Status)
		}
		if history.FinishedAt == nil {
			t.Errorf("history.FinishedAt must be registered")
		}
	})

	t.Run("abort batch history", func(t *testing.T) {
		history, err := NewBatchHistory(ctx, pool, "test", nil)
		if err != nil {
			t.Fatal(err)
		}

		if err := history.Abort(ctx, pool); err != nil {
			t.Error(err)
		}
		if history.Status != "aborted" {
			t.Errorf("history.Status = %s, want aborted", history.Status)
		}
		if history.FinishedAt == nil {
			t.Errorf("history.FinishedAt must be registered")
		}
	})

	t.Run("abort completed batch history", func(t *testing.T) {
		history, err := NewBatchHistory(ctx, pool, "test", nil)
		if err != nil {
			t.Fatal(err)
		}

		if err := history.Complete(ctx, pool); err != nil {
			t.Error(err)
		}

		if err := history.Abort(ctx, pool); err != ErrHistoryConfirmed {
			t.Errorf("err = %v, want ErrHistoryConfirmed", err)
		}
	})
}
