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

type DocumentGenerator interface {
	Generate(ctx context.Context) error
}

type documentGenerator[D any, R Documenter[D]] struct {
	reader     RowReader[R]
	saveDir    string
	chunkSize  int
	concurrent int
}

func NewDocumentGenerator[D any, R Documenter[D]](reader RowReader[R], saveDir string, chunkSize, concurrent int) *documentGenerator[D, R] {
	return &documentGenerator[D, R]{
		reader:     reader,
		saveDir:    saveDir,
		chunkSize:  chunkSize,
		concurrent: concurrent,
	}
}

func (g *documentGenerator[D, R]) Generate(ctx context.Context) error {
	if err := g.clean(); err != nil {
		return errs.Wrap(err)
	}

	rowChannel := make(chan R, g.chunkSize)
	docChannel := make(chan D, g.chunkSize)

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if err := g.reader.ReadRows(ctx, rowChannel); err != nil {
			return errs.Wrap(err)
		}
		return nil
	})
	eg.Go(func() error {
		if err := g.save(ctx, docChannel); err != nil {
			return errs.Wrap(err)
		}
		return nil
	})

	var wg sync.WaitGroup
	done := make(chan msg, 1)
	for i := 0; i < g.concurrent; i++ {
		wg.Add(1)
		eg.Go(func() error {
			defer wg.Done()
			if err := g.convert(ctx, rowChannel, docChannel); err != nil {
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

func (g *documentGenerator[D, R]) clean() error {
	files, err := filepath.Glob(filepath.Join(g.saveDir, "doc-*.json"))
	if err != nil {
		return errs.New(
			"failed to glob the directory",
			errs.WithCause(err),
			errs.WithContext("directory", g.saveDir),
		)
	}

	if len(files) == 0 {
		return nil
	}

	slog.Info(fmt.Sprintf("Start to delete existing document files in `%s`", g.saveDir))
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

func (g *documentGenerator[D, R]) save(ctx context.Context, rx <-chan D) error {
	write := func(docs []D, p string) error {
		file, err := os.Create(p)
		if err != nil {
			return errs.New(
				"failed to create file",
				errs.WithCause(err),
				errs.WithContext("file", p),
			)
		}
		defer file.Close()

		if err := json.NewEncoder(file).Encode(docs); err != nil {
			return errs.New(
				"failed to write docs into the file",
				errs.WithCause(err),
				errs.WithContext("file", p),
			)
		}
		return nil
	}

	suffix := 0
	buffer := make([]D, 0, g.chunkSize)

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

			if len(buffer) >= g.chunkSize {
				file := filepath.Join(g.saveDir, fmt.Sprintf("doc-%d.json", suffix))

				if err := write(buffer, file); err != nil {
					return errs.Wrap(err)
				}
				slog.Info("Generate file", slog.String("file", file))
				buffer = buffer[:0]
			}
		}
	}
	if len(buffer) > 0 {
		file := filepath.Join(g.saveDir, fmt.Sprintf("doc-%d.json", suffix))

		if err := write(buffer, file); err != nil {
			return errs.Wrap(err)
		}
		slog.Info("Generate file", slog.String("file", file))
	}
	return nil
}

func (g *documentGenerator[D, R]) convert(ctx context.Context, rx <-chan R, tx chan<- D) error {
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
