//go:build test_generate

package generate

import (
	"context"
	"errors"
	"fjnkt98/atcodersearch/batch"
	"testing"
	"time"
)

func TestReadSubmissions(t *testing.T) {
	db, err := getTestDB()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	reader := NewSubmissionRowReader(db, 30, false)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	ch := make(chan Documenter, 1)

	err = reader.ReadRows(ctx, ch)
	if err != nil && !errors.Is(err, batch.ErrInterrupt) {
		t.Errorf("an error occurred in read submission rows: %s", err.Error())
	}
}

func TestReadAllSubmissions(t *testing.T) {
	db, err := getTestDB()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	reader := NewSubmissionRowReader(db, 30, true)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	ch := make(chan Documenter, 1)

	err = reader.ReadRows(ctx, ch)
	if err != nil && !errors.Is(err, batch.ErrInterrupt) {
		t.Errorf("an error occurred in read submission rows: %s", err.Error())
	}
}
