package repository

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/goark/errs"
)

func TestOverwrite(t *testing.T) {
	startedAt := time.Date(2024, 6, 1, 0, 0, 0, 0, time.Local)
	finishedAt := time.Date(2024, 6, 1, 3, 0, 0, 0, time.Local)

	h := &BatchHistory{
		ID:         1,
		Name:       "Test",
		StartedAt:  startedAt,
		FinishedAt: nil,
		Status:     "working",
		Options:    nil,
	}

	r := BatchHistory{
		ID:         1,
		Name:       "Test",
		StartedAt:  startedAt,
		FinishedAt: &finishedAt,
		Status:     "finish",
		Options:    nil,
	}

	expected := &BatchHistory{
		ID:         1,
		Name:       "Test",
		StartedAt:  startedAt,
		FinishedAt: &finishedAt,
		Status:     "finish",
		Options:    nil,
	}

	h.overwrite(r)
	if !reflect.DeepEqual(h, expected) {
		t.Errorf("expected \n%+v\n, but got \n%+v\n", expected, h)
	}
}

func TestFailShouldBeSkip(t *testing.T) {
	startedAt := time.Date(2024, 6, 1, 0, 0, 0, 0, time.Local)
	finishedAt := time.Date(2024, 6, 1, 3, 0, 0, 0, time.Local)

	h := &BatchHistory{
		ID:         1,
		Name:       "Test",
		StartedAt:  startedAt,
		FinishedAt: &finishedAt,
		Status:     "finish",
		Options:    nil,
	}
	if err := h.Fail(context.Background(), nil); !errs.Is(err, ErrHistoryConfirmed) {
		t.Fatalf("test failed: %s", err.Error())
	}

	h = &BatchHistory{
		ID:         1,
		Name:       "Test",
		StartedAt:  startedAt,
		FinishedAt: &finishedAt,
		Status:     "canceled",
		Options:    nil,
	}
	if err := h.Fail(context.Background(), nil); !errs.Is(err, ErrHistoryConfirmed) {
		t.Fatalf("test failed: %s", err.Error())
	}
}
