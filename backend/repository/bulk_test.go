package repository

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ptr[T any](v T) *T {
	return &v
}

func newTestPool() *pgxpool.Pool {
	pool, err := NewPool(context.Background(), "postgres://atcodersearch:atcodersearch@localhost:5433/atcodersearch?sslmode=disable")
	if err != nil {
		panic(err)
	}
	return pool
}

func TestColumns(t *testing.T) {
	expected := []string{"id", "epoch_second", "problem_id", "contest_id", "user_id", "language", "point", "length", "result", "execution_time"}
	actual := Columns(new(Submission))

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected \n%+v\n , but got \n%+v\n", expected, actual)
	}
}

func TestUniqueKey(t *testing.T) {
	expected := []string{"id"}
	actual := UniqueKey(new(Submission))

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected \n%+v\n , but got \n%+v\n", expected, actual)
	}
}

func TestRows(t *testing.T) {
	data := []Submission{
		{0, 100, "abc300_a", ptr("abc300"), ptr("fjnkt98"), ptr("Go"), ptr(100.0), ptr(int32(200)), ptr("AC"), ptr(int32(5)), time.Date(2024, 4, 26, 22, 43, 0, 0, time.Local)},
		{1, 101, "abc300_b", nil, nil, nil, nil, nil, nil, nil, time.Date(2024, 4, 27, 22, 43, 0, 0, time.Local)},
	}

	expected := [][]any{
		{int64(0), int64(100), "abc300_a", ptr("abc300"), ptr("fjnkt98"), ptr("Go"), ptr(100.0), ptr(int32(200)), ptr("AC"), ptr(int32(5))},
		{int64(1), int64(101), "abc300_b", nil, nil, nil, nil, nil, nil, nil},
	}

	actual := Rows(data)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("expected \n%+v\n , but got \n%+v\n", expected, actual)
	}
}

func TestBulkUpdate(t *testing.T) {
	updatedAt := time.Now()
	data := []Submission{
		{0, 1714057200, "abc300_a", ptr("abc300"), ptr("fjnkt98"), ptr("Go"), ptr(100.0), ptr(int32(200)), ptr("AC"), ptr(int32(5)), updatedAt},
		{1, 1714057250, "abc300_b", ptr("abc300"), ptr("fjnkt98"), ptr("Go"), ptr(0.0), ptr(int32(200)), ptr("WA"), nil, updatedAt},
	}
	pool := newTestPool()
	if _, err := BulkUpdate(context.Background(), pool, "submissions", data); err != nil {
		t.Fatalf("failed to insert data: %s", err.Error())
	}
}
