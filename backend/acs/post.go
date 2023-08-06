package acs

import (
	"context"
	"fjnkt98/atcodersearch/solr"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"
)

type DocumentUploader interface {
	PostDocument(core solr.SolrCore[any, any], saveDir string, optimize bool) error
}

type DefaultDocumentUploader struct {
	core    solr.SolrCore[any, any]
	saveDir string
}

func NewDefaultDocumentUploader(core solr.SolrCore[any, any], saveDir string) DefaultDocumentUploader {
	return DefaultDocumentUploader{
		core,
		saveDir,
	}
}

func (u *DefaultDocumentUploader) PostDocument(optimize bool, concurrent int) error {
	log.Printf("Start to post documents in `%s`", u.saveDir)
	paths, err := filepath.Glob(filepath.Join(u.saveDir, "doc-*.json"))
	if err != nil {
		return fmt.Errorf("failed to get document files at `%s`: %s", u.saveDir, err.Error())
	}

	ch := make(chan string, len(paths))

	eg, ctx := errgroup.WithContext(context.Background())
	for i := 0; i < concurrent; i++ {
		eg.Go(func() error {
		loop:
			for {
				select {
				case <-ctx.Done():
					return nil
				case path, ok := <-ch:
					if !ok {
						break loop
					}
					select {
					case <-ctx.Done():
						return nil
					default:
					}

					log.Printf("Post document `%s`", path)
					file, err := os.Open(path)
					if err != nil {
						return fmt.Errorf("failed to open file `%s`: %s", path, err.Error())
					}

					if _, err := u.core.Post(file, "application/json"); err != nil {
						return fmt.Errorf("failed to post file `%s`: %s", path, err.Error())
					}
				}
			}

			return nil
		})
	}

	for _, path := range paths {
		ch <- path
	}
	close(ch)

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("failed to post documents: %s", err.Error())
	}

	if optimize {
		if _, err := u.core.Optimize(); err != nil {
			return fmt.Errorf("failed to optimize index: %s", err.Error())
		}
	} else {
		if _, err := u.core.Commit(); err != nil {
			return fmt.Errorf("failed to commit index: %s", err.Error())
		}
	}

	return nil
}
