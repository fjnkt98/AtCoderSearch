package acs

import (
	"context"
	"fjnkt98/atcodersearch/solr"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"

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
	var wg sync.WaitGroup

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	finish := make(chan msg, 1)

	eg.Go(func() error {
		select {
		case <-quit:
			log.Print("post interrupted.")
			return failure.New(Interrupt, failure.Message("post interrupted"))
		case <-ctx.Done():
			log.Print("interrupt observer canceled.")
			return nil
		case <-finish:
			return nil
		}
	})
	f := func(p string) error {
		file, err := os.Open(p)
		if err != nil {
			return failure.Translate(err, FileOperationError, failure.Messagef("failed to open file `%s`", p))
		}
		defer file.Close()
		if _, err := u.core.Post(file, "application/json"); err != nil {
			return failure.Translate(err, PostError, failure.Messagef("failed to open file `%s`", p))
		}

		return nil
	}

	for i := 0; i < concurrent; i++ {
		wg.Add(1)
		eg.Go(func() error {
			defer wg.Done()

		loop:
			for {
				select {
				case <-ctx.Done():
					log.Print("post canceled")
					return nil
				case path, ok := <-ch:
					if !ok {
						break loop
					}
					select {
					case <-ctx.Done():
						log.Print("post canceled")
						return nil
					default:
					}

					log.Printf("Post document `%s`", path)
					if err := f(path); err != nil {
						return failure.Wrap(err)
					}
				}
			}
			return nil
		})
	}

	eg.Go(func() error {
		for _, path := range paths {
			ch <- path
		}
		close(ch)
		return nil
	})

	eg.Go(func() error {
		wg.Wait()

		select {
		case <-ctx.Done():
			log.Print("post canceled. start rollback")
			if _, err := u.core.Rollback(); err != nil {
				return failure.Translate(err, PostError, failure.Message("failed to rollback index"))
			}
		default:
			if optimize {
				if _, err := u.core.Optimize(); err != nil {
					return failure.Translate(err, PostError, failure.Message("failed to optimize index"))
				}
			} else {
				if _, err := u.core.Commit(); err != nil {
					return failure.Translate(err, PostError, failure.Message("failed to commit index"))
				}
			}
		}

		finish <- msg{}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return failure.Wrap(err)
	}

	return nil
}
