//go:build test_generate

package generate

import (
	"context"
	"errors"
	"fjnkt98/atcodersearch/batch"
	"fjnkt98/atcodersearch/batch/repository"
	"fjnkt98/atcodersearch/config"
	"testing"
	"time"
)

func TestReadSubmissions(t *testing.T) {
	db := getTestDB()
	repo := repository.NewUpdateHistoryRepository(db)
	cfg := config.ReadSubmissionConfig{
		Interval: 30,
		All:      false,
	}
	reader := NewSubmissionRowReader(db, repo, cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	ch := make(chan Documenter, 1)

	err := reader.ReadRows(ctx, ch)
	if err != nil && !errors.Is(err, batch.ErrInterrupt) {
		t.Errorf("an error occurred in read submission rows: %s", err.Error())
	}
}
