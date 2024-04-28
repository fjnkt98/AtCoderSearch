package generate

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"log/slog"

	"github.com/goark/errs"
	_ "github.com/lib/pq"
	"golang.org/x/sync/errgroup"
)

type msg struct{}

type Documenter[D any] interface {
	Document(ctx context.Context) (D, error)
}

type RowReader[R any] interface {
	ReadRows(ctx context.Context, tx chan<- R) error
}

type GenerateDocumentOptions struct {
	ChunkSize  int
	Concurrent int
}

type option func(*GenerateDocumentOptions)

func WithChunkSize(chunkSize int) option {
	return func(opt *GenerateDocumentOptions) {
		opt.ChunkSize = chunkSize
	}
}

func WithConcurrent(concurrent int) option {
	return func(opt *GenerateDocumentOptions) {
		opt.Concurrent = concurrent
	}
}

func GenerateDocument[D any, R Documenter[D]](ctx context.Context, reader RowReader[R], saveDir string, options ...option) error {
	option := &GenerateDocumentOptions{
		ChunkSize:  1000,
		Concurrent: 4,
	}
	for _, o := range options {
		o(option)
	}

	if err := clean(saveDir); err != nil {
		return errs.Wrap(err)
	}

	rowChannel := make(chan R, option.ChunkSize)
	docChannel := make(chan D, option.ChunkSize)

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if err := reader.ReadRows(ctx, rowChannel); err != nil {
			return errs.Wrap(err)
		}
		return nil
	})
	eg.Go(func() error {
		if err := save[D, R](ctx, docChannel, saveDir, option.ChunkSize); err != nil {
			return errs.Wrap(err)
		}
		return nil
	})

	var wg sync.WaitGroup
	done := make(chan msg, 1)
	for i := 0; i < option.Concurrent; i++ {
		wg.Add(1)
		eg.Go(func() error {
			defer wg.Done()
			if err := convert(ctx, rowChannel, docChannel); err != nil {
				return errs.Wrap(err)
			}
			return nil
		})
	}
	eg.Go(func() error {
		wg.Wait()
		done <- msg{}
		close(docChannel)
		return nil
	})

	if err := eg.Wait(); err != nil {
		return errs.Wrap(err)
	}
	return nil
}

func clean(saveDir string) error {
	files, err := filepath.Glob(filepath.Join(saveDir, "doc-*.json"))
	if err != nil {
		return errs.New(
			"failed to glob the directory",
			errs.WithCause(err),
			errs.WithContext("directory", saveDir),
		)
	}

	if len(files) == 0 {
		return nil
	}

	slog.Info(fmt.Sprintf("Start to delete existing document files in `%s`", saveDir))
	for _, file := range files {
		slog.Info(fmt.Sprintf("Delete file `%s`", file))
		if err := os.Remove(file); err != nil {
			return errs.New(
				"failed to delete the file",
				errs.WithCause(err),
				errs.WithContext("file", file),
			)
		}
	}
	return nil
}

func save[D any, R Documenter[D]](ctx context.Context, rx <-chan D, saveDir string, chunkSize int) (err error) {
	write := func(docs []D, path string) error {
		file, err := os.Create(path)
		if err != nil {
			return errs.New(
				"failed to create file",
				errs.WithCause(err),
				errs.WithContext("file", path),
			)
		}
		defer file.Close()
		defer func() {
			if closeErr := file.Close(); closeErr != nil {
				err = errs.Join(
					errs.New("failed to close doc file", errs.WithCause(closeErr)),
					err,
				)
			}
		}()

		if err := json.NewEncoder(file).Encode(docs); err != nil {
			return errs.New(
				"failed to write docs into the file",
				errs.WithCause(err),
				errs.WithContext("file", path),
			)
		}
		return nil
	}

	suffix := 0
	buffer := make([]D, 0, chunkSize)

loop:
	for {
		select {
		case <-ctx.Done():
			return nil
		case doc, ok := <-rx:
			if !ok {
				break loop
			}
			select {
			case <-ctx.Done():
				return nil
			default:
			}

			suffix++
			buffer = append(buffer, doc)

			if len(buffer) >= chunkSize {
				file := filepath.Join(saveDir, fmt.Sprintf("doc-%d.json", suffix))

				if err := write(buffer, file); err != nil {
					return errs.Wrap(err)
				}
				slog.Info("Generate file", slog.String("file", file))
				buffer = buffer[:0]
			}
		}
	}
	if len(buffer) > 0 {
		file := filepath.Join(saveDir, fmt.Sprintf("doc-%d.json", suffix))

		if err := write(buffer, file); err != nil {
			return errs.Wrap(err)
		}
		slog.Info("Generate file", slog.String("file", file))
	}
	return nil
}

func convert[D any, R Documenter[D]](ctx context.Context, rx <-chan R, tx chan<- D) error {
loop:
	for {
		select {
		case <-ctx.Done():
			return nil
		case row, ok := <-rx:
			if !ok {
				break loop
			}

			select {
			case <-ctx.Done():
				return nil
			default:
			}

			d, err := row.Document(ctx)
			if err != nil {
				return errs.New(
					"failed to convert document",
					errs.WithCause(err),
				)
			}

			tx <- d
		}
	}
	return nil
}
