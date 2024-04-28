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

type PostDocumentOptions struct {
	Concurrent int
	Truncate   bool
	Optimize   bool
}

type option func(*PostDocumentOptions)

func WithConcurrent(concurrent int) option {
	return func(opt *PostDocumentOptions) {
		opt.Concurrent = concurrent
	}
}

func WithTruncate(truncate bool) option {
	return func(opt *PostDocumentOptions) {
		opt.Truncate = truncate
	}
}

func WithOptimize(optimize bool) option {
	return func(opt *PostDocumentOptions) {
		opt.Optimize = optimize
	}
}

func PostDocument(ctx context.Context, core *solr.SolrCore, saveDir string, options ...option) error {
	option := &PostDocumentOptions{
		Concurrent: 4,
		Truncate:   false,
		Optimize:   false,
	}
	for _, opt := range options {
		opt(option)
	}

	files, err := filepath.Glob(filepath.Join(saveDir, "doc-*.json"))
	if err != nil {
		return errs.New(
			"failed to get files from the directory",
			errs.WithCause(err),
			errs.WithContext("directory", saveDir),
		)
	}

	f := func(ctx context.Context, path string) error {
		file, err := os.Open(path)
		if err != nil {
			return errs.New("failed to open the file", errs.WithCause(err), errs.WithContext("file", path))
		}
		defer file.Close()

		if _, err := core.Post(ctx, file, "application/json"); err != nil {
			return errs.New("failed to post the file", errs.WithCause(err), errs.WithContext("file", path), errs.WithContext("core", core.Name()))
		}
		return nil
	}

	ch := make(chan string, len(files))
	eg, ctx := errgroup.WithContext(ctx)
	var wg sync.WaitGroup

	for i := 0; i < option.Concurrent; i++ {
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
					slog.Info("Post document", slog.String("file", path), slog.String("core", core.Name()), slog.Int("worker", i))
				}
			}
			return nil
		})
	}
	eg.Go(func() error {
		defer core.Rollback()

		if option.Truncate {
			if _, err := core.Delete(ctx); err != nil {
				return errs.New("failed to truncate documents", errs.WithCause(err), errs.WithContext("core", core.Name()))
			}
			slog.Info("Finished truncating core successfully.", slog.String("core", core.Name()))
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
			if option.Optimize {
				if _, err := core.Optimize(ctx); err != nil {
					return errs.New("failed to optimize core", errs.WithCause(err), errs.WithContext("core", core.Name()))
				}
				slog.Info("Finished optimizing core successfully.", slog.String("core", core.Name()))
			} else {
				if _, err := core.Commit(ctx); err != nil {
					return errs.New("failed to commit core", errs.WithCause(err), errs.WithContext("core", core.Name()))
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
