package acs

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
	"golang.org/x/sync/errgroup"
)

type ToDocument[D any] interface {
	ToDocument() (D, error)
}

type DocumentGenerator[R ToDocument[D], D any] interface {
	ReadRows(tx chan<- R) error
	Clean() error
	Generate(chunkSize int) error
}

type DefaultDocumentGenerator[R ToDocument[D], D any] struct {
	saveDir string
}

func (d *DefaultDocumentGenerator[R, D]) ReadRows(tx chan<- R) error {
	return nil
}

func (d *DefaultDocumentGenerator[R, D]) Clean() error {
	log.Printf("Start to delete existing document files in `%s`", d.saveDir)
	return filepath.WalkDir(d.saveDir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !entry.IsDir() && filepath.Ext(path) == ".json" {
			log.Printf("Delete existing file `%s`", path)
			if err := os.Remove(path); err != nil {
				return err
			}
		}
		return nil
	})
}

func generateDocument[R ToDocument[D], D any](rx <-chan R, tx chan<- D) error {
	for row := range rx {
		d, err := row.ToDocument()
		if err != nil {
			return fmt.Errorf("failed to convert from row into document: %s", err.Error())
		}

		tx <- d
	}
	return nil
}

func saveDocuments[D any](rx <-chan D, saveDir string, chunkSize int) error {
	suffix := 0
	buffer := make([]D, 0, chunkSize)

	writeDocument := func(documents []D, filePath string) error {
		file, err := os.Create(filePath)
		if err != nil {
			log.Printf("failed to open file `%s`: %s", filePath, err.Error())
			return err
		}
		defer file.Close()

		if err := json.NewEncoder(file).Encode(buffer); err != nil {
			log.Printf("failed to write into document file `%s`: %s", filePath, err.Error())
			return err
		}
		return nil
	}

	for doc := range rx {
		suffix++
		buffer = append(buffer, doc)

		if len(buffer) >= chunkSize {
			filePath := filepath.Join(saveDir, fmt.Sprintf("doc-%d.json", suffix))

			log.Printf("Generate document file: %s", filePath)
			writeDocument(buffer, filePath)

			buffer = buffer[:0]
		}
	}

	if len(buffer) != 0 {
		filePath := filepath.Join(saveDir, fmt.Sprintf("doc-%d.json", suffix))

		log.Printf("Generate document file: %s", filePath)
		writeDocument(buffer, filePath)

		buffer = buffer[:0]
	}

	return nil
}

func (d *DefaultDocumentGenerator[R, D]) Generate(chunkSize int) error {
	rowChannel := make(chan R, chunkSize)
	docChannel := make(chan D, chunkSize)

	eg := errgroup.Group{}

	eg.Go(func() error { return d.ReadRows(rowChannel) })
	eg.Go(func() error { return saveDocuments(docChannel, d.saveDir, chunkSize) })
	for i := 0; i < 4; i++ {
		eg.Go(func() error { return generateDocument(rowChannel, docChannel) })
	}

	if err := eg.Wait(); err != nil {
		log.Printf("failed to generate document: %s", err.Error())
		return err
	}

	return nil
}
