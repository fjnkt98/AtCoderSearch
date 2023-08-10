package acs

import (
	"context"
	"fjnkt98/atcodersearch/solr"
	"log"
	"os"
	"path/filepath"

	"github.com/morikuni/failure"
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
		return failure.Translate(err, PostError, failure.Messagef("failed to get document files at `%s`", u.saveDir))
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
						return failure.Translate(err, FileOperationError, failure.Messagef("failed to open file `%s`", path))
					}

					if _, err := u.core.Post(file, "application/json"); err != nil {
						return failure.Translate(err, PostError, failure.Messagef("failed to open file `%s`", path))
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
		return failure.Wrap(err)
	}

	if optimize {
		if _, err := u.core.Optimize(); err != nil {
			return failure.Translate(err, PostError, failure.Message("failed to optimize index"))
		}
	} else {
		if _, err := u.core.Commit(); err != nil {
			return failure.Translate(err, PostError, failure.Message("failed to commit index"))
		}
	}

	return nil
}
