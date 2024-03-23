package post

import (
	"context"
	"fjnkt98/atcodersearch/pkg/solr"
	"os"
	"path/filepath"
	"sync"

	"github.com/goark/errs"
	"golang.org/x/sync/errgroup"
)

type DocumentPoster interface {
	Post(ctx context.Context) error
}

type documentPoster struct {
	core       solr.SolrCore
	saveDir    string
	concurrent int
	optimize   bool
	truncate   bool
}

func NewDocumentPoster(core solr.SolrCore, saveDir string, concurrent int, optimize, truncate bool) DocumentPoster {
	return &documentPoster{
		core:       core,
		saveDir:    saveDir,
		concurrent: concurrent,
		optimize:   optimize,
		truncate:   truncate,
	}
}

func (p *documentPoster) Post(ctx context.Context) error {
	files, err := filepath.Glob(filepath.Join(p.saveDir, "doc-*.json"))
	if err != nil {
		return errs.New(
			"failed to get files from the directory",
			errs.WithCause(err),
			errs.WithContext("directory", p.saveDir),
		)
	}

	ch := make(chan string, len(files))
	eg, ctx := errgroup.WithContext(ctx)
	var wg sync.WaitGroup

	for i := 0; i < p.concurrent; i++ {
		wg.Add(1)
		eg.Go(func() error {
			defer wg.Done()

		loop:
			for {
				select {
				case <-ctx.Done():
					return nil
				case p, ok := <-ch:
					if !ok {
						break loop
					}
					select {
					case <-ctx.Done():
						return nil
					default:
					}

					file, err := os.Open(p)
					if err != nil {
						return errs.New(
							"failed to open the file",
							errs.WithCause(err),
							errs.WithContext("file", p),
						)
					}
					defer file.Close()
					panic("TODO")
				}
			}
			return nil
		})
	}
	eg.Go(func() error {
		if p.truncate {
			panic("TODO")
		}

		for _, p := range files {
			ch <- p
		}
		close(ch)

		wg.Wait()
		select {
		case <-ctx.Done():
			panic("TODO")
		default:
			if p.optimize {
				panic("TODO")
			} else {
				panic("TODO")
			}
		}
		return nil
	})

	return nil
}
