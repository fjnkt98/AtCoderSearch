package generate

import (
	"context"
	"encoding/json"
	"fjnkt98/atcodersearch/batch"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"log/slog"

	"github.com/goark/errs"
	_ "github.com/lib/pq"
	"golang.org/x/sync/errgroup"
)

type defaultGenerator struct {
	reader RowReader
	config defaultGeneratorConfig
}

type defaultGeneratorConfig struct {
	SaveDir    string `json:"save_dir"`
	ChunkSize  int    `json:"chunk_size"`
	Concurrent int    `json:"concurrent"`
}

func (g *defaultGenerator) Generate(ctx context.Context) error {
	slog.Info("Start to generate documents.")
	if err := GenerateDocument(
		ctx,
		g.reader,
		g.config.SaveDir,
		g.config.ChunkSize,
		g.config.Concurrent,
	); err != nil {
		return errs.Wrap(err)
	}
	slog.Info("Finished generating documents successfully.")
	return nil
}

func (g *defaultGenerator) Config() any {
	return g.config
}

type Documenter interface {
	Document(ctx context.Context) (map[string]any, error)
}

type RowReader interface {
	ReadRows(ctx context.Context, tx chan<- Documenter) error
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
		slog.Info(fmt.Sprintf("Delete existing file `%s`", file))
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

func convert(ctx context.Context, rx <-chan Documenter, tx chan<- map[string]any) error {
loop:
	for {
		select {
		case <-ctx.Done():
			slog.Info("convert canceled.")
			return batch.ErrInterrupt
		case row, ok := <-rx:
			if !ok {
				break loop
			}

			select {
			case <-ctx.Done():
				slog.Info("convert canceled.")
				return batch.ErrInterrupt
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

func save(ctx context.Context, rx <-chan map[string]any, saveDir string, chunkSize int) error {
	writeDocument := func(documents []any, filePath string) error {
		file, err := os.Create(filePath)
		if err != nil {
			return errs.New(
				"failed to open the file",
				errs.WithCause(err),
				errs.WithContext("file", filePath),
			)
		}
		defer file.Close()

		if err := json.NewEncoder(file).Encode(documents); err != nil {
			return errs.New(
				"failed to write the file",
				errs.WithCause(err),
				errs.WithContext("file", filePath),
			)
		}
		return nil
	}

	suffix := 0
	buffer := make([]any, 0, chunkSize)

loop:
	for {
		select {
		case <-ctx.Done():
			slog.Info("save canceled.")
			return batch.ErrInterrupt
		case doc, ok := <-rx:
			if !ok {
				break loop
			}
			select {
			case <-ctx.Done():
				slog.Info("SaveDocument canceled.")
				return batch.ErrInterrupt
			default:
			}

			suffix++
			buffer = append(buffer, doc)

			if len(buffer) >= chunkSize {
				filePath := filepath.Join(saveDir, fmt.Sprintf("doc-%d.json", suffix))

				slog.Info(fmt.Sprintf("Generate document file: %s", filePath))
				if err := writeDocument(buffer, filePath); err != nil {
					return errs.Wrap(err)
				}

				buffer = buffer[:0]
			}
		}
	}

	if len(buffer) != 0 {
		filePath := filepath.Join(saveDir, fmt.Sprintf("doc-%d.json", suffix))

		slog.Info(fmt.Sprintf("Generate document file: %s", filePath))
		if err := writeDocument(buffer, filePath); err != nil {
			return errs.Wrap(err)
		}
	}

	return nil
}

type msg struct{}

func GenerateDocument(
	ctx context.Context,
	reader RowReader,
	saveDir string,
	chunkSize int,
	concurrent int,
) error {
	if err := clean(saveDir); err != nil {
		return errs.Wrap(err)
	}

	rowChannel := make(chan Documenter, chunkSize)
	docChannel := make(chan map[string]any, chunkSize)

	eg, ctx := errgroup.WithContext(ctx)
	var wg sync.WaitGroup

	done := make(chan msg, 1)

	eg.Go(func() error {
		select {
		case <-ctx.Done():
			slog.Info("generate canceled")
			return batch.ErrInterrupt
		case <-done:
			return nil
		}
	})
	eg.Go(func() error { return reader.ReadRows(ctx, rowChannel) })
	eg.Go(func() error { return save(ctx, docChannel, saveDir, chunkSize) })
	for i := 0; i < concurrent; i++ {
		wg.Add(1)
		eg.Go(func() error {
			defer wg.Done()
			return convert(ctx, rowChannel, docChannel)
		})
	}
	eg.Go(func() error {
		wg.Wait()
		done <- msg{}
		defer close(docChannel)
		return nil
	})

	if err := eg.Wait(); err != nil {
		return errs.Wrap(err)
	}

	return nil
}

func StructToMap(doc any) map[string]any {
	result := make(map[string]any)
	ty := reflect.TypeOf(doc)
	if ty.Kind() != reflect.Struct {
		return nil
	}
	val := reflect.ValueOf(doc)

	for i := 0; i < ty.NumField(); i++ {
		fieldValue := val.Field(i)
		fieldType := ty.Field(i)

		if tag, ok := fieldType.Tag.Lookup("solr"); ok {
			if tag == "-" {
				continue
			}
			key, suffixes, ok := strings.Cut(tag, ",")
			if ok {
				for _, suffix := range strings.Split(suffixes, ",") {
					result[fmt.Sprintf("%s__%s", key, suffix)] = fieldValue.Interface()
				}
			}
			result[key] = fieldValue.Interface()
		} else {
			result[fieldType.Name] = fieldValue.Interface()
		}
	}

	return result
}
