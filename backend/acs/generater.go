package acs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	_ "github.com/lib/pq"
	"golang.org/x/sync/errgroup"
)

type ToDocument[D any] interface {
	ToDocument() (D, error)
}

type RowReader[R ToDocument[D], D any] interface {
	ReadRows(ctx context.Context, tx chan<- R) error
}

type DocumentGenerator[D any] interface {
	Clean() error
	Generate(chunkSize int) error
}

func CleanDocument(saveDir string) error {
	files, err := filepath.Glob(filepath.Join(saveDir, "doc-*.json"))
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return nil
	}

	log.Printf("Start to delete existing document files in `%s`", saveDir)
	for _, file := range files {
		log.Printf("Delete existing file `%s`", file)
		if err := os.Remove(file); err != nil {
			return fmt.Errorf("failed to remove file `%s`: %w", file, err)
		}
	}
	return nil
}

func ConvertDocument[R ToDocument[D], D any](ctx context.Context, rx <-chan R, tx chan<- D) error {
loop:
	for {
		select {
		case <-ctx.Done():
			log.Println("ConvertDocument canceled.")
			return nil
		case row, ok := <-rx:
			if !ok {
				break loop
			}

			select {
			default:
			case <-ctx.Done():
				log.Println("ConvertDocument canceled.")
				return nil
			}

			d, err := row.ToDocument()
			if err != nil {
				return fmt.Errorf("failed to convert from row into document: %w", err)
			}

			tx <- d
		}
	}
	return nil
}

func SaveDocument[D any](ctx context.Context, rx <-chan D, saveDir string, chunkSize int) error {
	suffix := 0
	buffer := make([]any, 0, chunkSize)

	writeDocument := func(documents []any, filePath string) error {
		file, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("failed to open file `%s`: %w", filePath, err)
		}
		defer file.Close()

		if err := json.NewEncoder(file).Encode(buffer); err != nil {
			return fmt.Errorf("failed to write into document file `%s`: %w", filePath, err)
		}
		return nil
	}

loop:
	for {
		select {
		case <-ctx.Done():
			log.Println("SaveDocument canceled.")
			return nil
		case doc, ok := <-rx:
			if !ok {
				break loop
			}
			select {
			case <-ctx.Done():
				log.Println("SaveDocument canceled.")
				return nil
			default:
			}

			suffix++
			buffer = append(buffer, doc)

			if len(buffer) >= chunkSize {
				filePath := filepath.Join(saveDir, fmt.Sprintf("doc-%d.json", suffix))

				log.Printf("Generate document file: %s", filePath)
				if err := writeDocument(buffer, filePath); err != nil {
					return err
				}

				buffer = buffer[:0]
			}
		}
	}

	if len(buffer) != 0 {
		filePath := filepath.Join(saveDir, fmt.Sprintf("doc-%d.json", suffix))

		log.Printf("Generate document file: %s", filePath)
		if err := writeDocument(buffer, filePath); err != nil {
			return err
		}

		buffer = buffer[:0]
	}

	return nil
}

func GenerateDocument[R ToDocument[D], D any](reader RowReader[R, D], saveDir string, chunkSize int, concurrent int) error {
	rowChannel := make(chan R, chunkSize)
	docChannel := make(chan D, chunkSize)

	eg, ctx := errgroup.WithContext(context.Background())
	var wg sync.WaitGroup

	eg.Go(func() error {
		wg.Wait()
		close(docChannel)
		return nil
	})
	eg.Go(func() error { return reader.ReadRows(ctx, rowChannel) })
	eg.Go(func() error { return SaveDocument(ctx, docChannel, saveDir, chunkSize) })
	for i := 0; i < concurrent; i++ {
		wg.Add(1)
		eg.Go(func() error {
			defer wg.Done()
			return ConvertDocument[R, D](ctx, rowChannel, docChannel)
		})
	}

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("failed to generate document: %w", err)
	}

	return nil
}
