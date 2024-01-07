//go:build test_generate

package generate

import (
	"context"
	"errors"
	"fjnkt98/atcodersearch/batch"
	"testing"
	"time"
)

func TestReadUsers(t *testing.T) {
	db, err := getTestDB()
	if err != nil {
		t.Fatalf("%+v", err)
	}
	reader := NewUserRowReader(db)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	ch := make(chan Documenter, 1)

	err = reader.ReadRows(ctx, ch)
	if err != nil && !errors.Is(err, batch.ErrInterrupt) {
		t.Errorf("an error occurred in read user rows: %s", err.Error())
	}
}
