package acs

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	_ "github.com/lib/pq"
	"github.com/morikuni/failure"
	"golang.org/x/exp/slog"
	"golang.org/x/sync/errgroup"
)

type ToDocFunc[R, D any] func(ctx context.Context, row R) (D, error)

type ReadRowsFunc[R any] func(ctx context.Context, tx chan<- R) error

func CleanDocument(saveDir string) error {
	files, err := filepath.Glob(filepath.Join(saveDir, "doc-*.json"))
	if err != nil {
		return failure.Translate(err, FileOperationError, failure.Context{"directory": saveDir}, failure.Message("failed to glob the directory"))
	}

	if len(files) == 0 {
		return nil
	}

	slog.Info(fmt.Sprintf("Start to delete existing document files in `%s`", saveDir))
	for _, file := range files {
		slog.Info(fmt.Sprintf("Delete existing file `%s`", file))
		if err := os.Remove(file); err != nil {
			return failure.Translate(err, FileOperationError, failure.Context{"file": file}, failure.Message("failed to delete the file"))
		}
	}
	return nil
}

func ConvertDocument[R, D any](ctx context.Context, rx <-chan R, tx chan<- D, convert ToDocFunc[R, D]) error {
loop:
	for {
		select {
		case <-ctx.Done():
			slog.Info("ConvertDocument canceled.")
			return failure.New(Interrupt, failure.Message("ConvertDocument canceled"))
		case row, ok := <-rx:
			if !ok {
				break loop
			}

			select {
			default:
			case <-ctx.Done():
				slog.Info("ConvertDocument canceled.")
				return failure.New(Interrupt, failure.Message("ConvertDocument canceled"))
			}

			d, err := convert(ctx, row)
			if err != nil {
				return failure.Translate(err, ConvertError, failure.Message("failed to convert document"))
			}

			tx <- d
		}
	}
	return nil
}

func SaveDocument[D any](ctx context.Context, rx <-chan D, saveDir string, chunkSize int) error {
	writeDocument := func(documents []any, filePath string) error {
		file, err := os.Create(filePath)
		if err != nil {
			return failure.Translate(err, FileOperationError, failure.Context{"file": filePath}, failure.Message("failed to open the file"))
		}
		defer file.Close()

		if err := json.NewEncoder(file).Encode(documents); err != nil {
			return failure.Translate(err, FileOperationError, failure.Context{"file": filePath}, failure.Message("failed to write the file"))
		}
		return nil
	}

	suffix := 0
	buffer := make([]any, 0, chunkSize)

loop:
	for {
		select {
		case <-ctx.Done():
			slog.Info("SaveDocument canceled.")
			return failure.New(Interrupt, failure.Message("SaveDocument canceled"))
		case doc, ok := <-rx:
			if !ok {
				break loop
			}
			select {
			case <-ctx.Done():
				slog.Info("SaveDocument canceled.")
				return failure.New(Interrupt, failure.Message("SaveDocument canceled"))
			default:
			}

			suffix++
			buffer = append(buffer, doc)

			if len(buffer) >= chunkSize {
				filePath := filepath.Join(saveDir, fmt.Sprintf("doc-%d.json", suffix))

				slog.Info(fmt.Sprintf("Generate document file: %s", filePath))
				if err := writeDocument(buffer, filePath); err != nil {
					return failure.Wrap(err)
				}

				buffer = buffer[:0]
			}
		}
	}

	if len(buffer) != 0 {
		filePath := filepath.Join(saveDir, fmt.Sprintf("doc-%d.json", suffix))

		slog.Info(fmt.Sprintf("Generate document file: %s", filePath))
		if err := writeDocument(buffer, filePath); err != nil {
			return failure.Wrap(err)
		}
	}

	return nil
}

type msg struct{}

func GenerateDocument[R, D any](ctx context.Context, saveDir string, chunkSize int, concurrent int, readFunc ReadRowsFunc[R], convertFunc ToDocFunc[R, D]) error {
	rowChannel := make(chan R, chunkSize)
	docChannel := make(chan D, chunkSize)

	eg, ctx := errgroup.WithContext(ctx)
	var wg sync.WaitGroup

	done := make(chan msg, 1)

	eg.Go(func() error {
		select {
		case <-ctx.Done():
			slog.Info("generate interrupted")
			return failure.New(Interrupt, failure.Message("generate interrupted"))
		case <-done:
			return nil
		}
	})
	eg.Go(func() error {
		wg.Wait()
		done <- msg{}
		defer close(docChannel)
		return nil
	})
	eg.Go(func() error { return readFunc(ctx, rowChannel) })
	eg.Go(func() error { return SaveDocument(ctx, docChannel, saveDir, chunkSize) })
	for i := 0; i < concurrent; i++ {
		wg.Add(1)
		eg.Go(func() error {
			defer wg.Done()
			return ConvertDocument[R, D](ctx, rowChannel, docChannel, convertFunc)
		})
	}

	if err := eg.Wait(); err != nil {
		return failure.Wrap(err)
	}

	return nil
}
