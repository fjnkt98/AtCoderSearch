package generate

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/goark/errs"
)

func TestGenerateDocumentOptions(t *testing.T) {
	cases := []struct {
		name     string
		options  []option
		expected GenerateDocumentOptions
	}{
		{name: "default", options: nil, expected: GenerateDocumentOptions{ChunkSize: 0, Concurrent: 0}},
		{name: "with chunk size", options: []option{WithChunkSize(6)}, expected: GenerateDocumentOptions{ChunkSize: 6, Concurrent: 0}},
		{name: "with concurrent", options: []option{WithConcurrent(10)}, expected: GenerateDocumentOptions{ChunkSize: 0, Concurrent: 10}},
		{name: "with chunk size and concurrent", options: []option{WithChunkSize(4), WithConcurrent(4)}, expected: GenerateDocumentOptions{ChunkSize: 4, Concurrent: 4}},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := &GenerateDocumentOptions{}
			for _, opt := range tt.options {
				opt(actual)
			}

			if !reflect.DeepEqual(tt.expected, *actual) {
				t.Errorf("expected \n%+v\n, but got \n%+v\n", tt.expected, *actual)
			}
		})
	}
}

type SuccessRow struct{}
type FailRow struct{}
type Doc struct{ Data string }

func (r *SuccessRow) Document(ctx context.Context) (*Doc, error) {
	return &Doc{Data: "doc"}, nil
}

func (r *FailRow) Document(ctx context.Context) (*Doc, error) {
	return nil, fmt.Errorf("fail")
}

func TestClean(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		dir := t.TempDir()

		err := clean(dir)
		if err != nil {
			t.Errorf("expected nil, but got %s", err.Error())
		}
	})

	t.Run("delete", func(t *testing.T) {
		dir := t.TempDir()

		_, err := os.Create(filepath.Join(dir, "doc-1.json"))
		if err != nil {
			t.Fatalf("failed to create test file: %s", err.Error())
		}

		before, err := filepath.Glob(filepath.Join(dir, "doc-1.json"))
		if err != nil {
			t.Fatalf("failed to glob test dir: %s", err.Error())
		}

		if err := clean(dir); err != nil {
			t.Fatalf("failed to clean test dir: %s", err.Error())
		}

		after, err := filepath.Glob(filepath.Join(dir, "doc-1.json"))
		if err != nil {
			t.Fatalf("failed to glob test dir: %s", err.Error())
		}

		if (len(before) == len(after)) || (len(after) != 0) {
			t.Errorf("expected len(before) == 1 and len(after) == 0, but got len(before) == %d and len(after) == %d", len(before), len(after))
		}
	})

}

func TestSave(t *testing.T) {

}

func TestConvert(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		rx := make(chan *SuccessRow, 1)
		tx := make(chan *Doc, 1)

		rx <- &SuccessRow{}
		close(rx)

		if err := convert(ctx, rx, tx); err != nil {
			t.Errorf("failed to convert: %s", err.Error())
		}

		d := <-tx
		close(tx)
		expected := Doc{Data: "doc"}
		if !reflect.DeepEqual(expected, *d) {
			t.Errorf("expected \n%+v\n, but got \n+%v\n", expected, *d)
		}
	})

	t.Run("fail", func(t *testing.T) {
		ctx := context.Background()
		rx := make(chan *FailRow, 1)
		tx := make(chan *Doc, 1)

		rx <- &FailRow{}
		close(rx)

		if err := convert(ctx, rx, tx); err == nil {
			t.Errorf("expected err, but got nil")
		}
	})

	t.Run("canceled", func(t *testing.T) {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		cancel()
		rx := make(chan *SuccessRow, 1)
		tx := make(chan *Doc, 1)

		err := convert(ctx, rx, tx)
		if !errs.Is(err, context.Canceled) {
			t.Errorf("expected context.Canceled, but got %+v", err)
		}
	})
}
