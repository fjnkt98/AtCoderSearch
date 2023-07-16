package atcoder_search

import (
	"fjnkt98/atcodersearch/pkg/solr"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"
)

type DocumentUploader[D any, F any] interface {
	PostDocument(core solr.SolrCore[D, F], saveDir string, optimize bool) error
}

type DefaultDocumentUploader struct {
	saveDir string
}

func NewDefaultDocumentUploader(saveDir string) DefaultDocumentUploader {
	return DefaultDocumentUploader{
		saveDir,
	}
}

func (d *DefaultDocumentUploader) PostDocument(core solr.SolrCore[any, any], saveDir string, optimize bool) error {
	log.Printf("Start to post documents in `%s`", d.saveDir)
	paths := make([]string, 0, 1)
	if err := filepath.WalkDir(d.saveDir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !entry.IsDir() && filepath.Ext(path) == ".json" {
			paths = append(paths, path)
		}
		return nil
	}); err != nil {
		log.Printf("failed to get file list: %s", err.Error())
		return err
	}

	eg := errgroup.Group{}
	for _, path := range paths {
		path := path
		eg.Go(func() error {
			file, err := os.Open(path)
			if err != nil {
				log.Printf("failed to open file `%s`: %s", path, err.Error())
				return err
			}

			if _, err := core.Post(file, "json"); err != nil {
				log.Printf("failed to post file `%s`: %s", path, err.Error())
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		log.Printf("failed to post documents: %s", err.Error())
		return err
	}

	if optimize {
		if _, err := core.Optimize(); err != nil {
			log.Printf("failed to optimize index: %s", err.Error())
		}
	} else {
		if _, err := core.Commit(); err != nil {
			log.Printf("failed to optimize index: %s", err.Error())
		}
	}

	return nil
}
