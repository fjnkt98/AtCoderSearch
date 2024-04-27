package post

import (
	"context"
	"fjnkt98/atcodersearch/pkg/solr"
	"log/slog"
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
	core       *solr.SolrCore
	saveDir    string
	concurrent int
	optimize   bool
	truncate   bool
}

func NewDocumentPoster(core *solr.SolrCore, saveDir string, concurrent int, optimize, truncate bool) DocumentPoster {
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

	f := func(ctx context.Context, path string) error {
		file, err := os.Open(path)
		if err != nil {
			return errs.New("failed to open the file", errs.WithCause(err), errs.WithContext("file", path))
		}
		defer file.Close()

		if _, err := p.core.Post(ctx, file, "application/json"); err != nil {
			return errs.New("failed to post the file", errs.WithCause(err), errs.WithContext("file", path), errs.WithContext("core", p.core.Name()))
		}
		return nil
	}

	ch := make(chan string, len(files))
	eg, ctx := errgroup.WithContext(ctx)
	var wg sync.WaitGroup

	for i := 0; i < p.concurrent; i++ {
		wg.Add(1)
		i := i
		eg.Go(func() error {
			defer wg.Done()

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

					if err := f(ctx, path); err != nil {
						return errs.Wrap(err)
					}
					slog.Info("Post document", slog.String("file", path), slog.String("core", p.core.Name()), slog.Int("worker", i))
				}
			}
			return nil
		})
	}
	eg.Go(func() error {
		defer p.core.Rollback()

		if p.truncate {
			slog.Info("Start to truncate core.")
			if _, err := p.core.Delete(ctx); err != nil {
				return errs.New("failed to truncate documents", errs.WithCause(err), errs.WithContext("core", p.core.Name()))
			}
			slog.Info("Finished truncating core successfully.")
		}

		for _, path := range files {
			ch <- path
		}
		close(ch)

		wg.Wait()
		select {
		case <-ctx.Done():
			return nil
		default:
			if p.optimize {
				if _, err := p.core.Optimize(ctx); err != nil {
					return errs.New("failed to optimize core", errs.WithCause(err), errs.WithContext("core", p.core.Name()))
				}
			} else {
				if _, err := p.core.Commit(ctx); err != nil {
					return errs.New("failed to commit core", errs.WithCause(err), errs.WithContext("core", p.core.Name()))
				}
			}
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return errs.Wrap(err)
	}

	return nil
}
