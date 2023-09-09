package acs

import (
	"context"
	"fjnkt98/atcodersearch/solr"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/morikuni/failure"
	"golang.org/x/exp/slog"
	"golang.org/x/sync/errgroup"
)

func PostDocument(ctx context.Context, core *solr.Core, saveDir string, optimize bool, truncate bool, concurrent int) error {
	slog.Info(fmt.Sprintf("Start to post documents in `%s`", saveDir))
	paths, err := filepath.Glob(filepath.Join(saveDir, "doc-*.json"))
	if err != nil {
		return failure.Translate(err, PostError, failure.Messagef("failed to get document files at `%s`", saveDir))
	}

	ch := make(chan string, len(paths))

	eg, ctx := errgroup.WithContext(ctx)
	var wg sync.WaitGroup

	f := func(ctx context.Context, p string) error {
		file, err := os.Open(p)
		if err != nil {
			return failure.Translate(err, FileOperationError, failure.Messagef("failed to open file `%s`", p))
		}
		defer file.Close()
		if _, err := solr.PostWithContext(ctx, core, file, "application/json"); err != nil {
			return failure.Translate(err, PostError, failure.Messagef("failed to open file `%s`", p))
		}

		return nil
	}

	for i := 0; i < concurrent; i++ {
		wg.Add(1)
		workerNum := i
		eg.Go(func() error {
			defer wg.Done()

		loop:
			for {
				select {
				case <-ctx.Done():
					slog.Info(fmt.Sprintf("post worker `%d` canceled", workerNum))
					return nil
				case path, ok := <-ch:
					if !ok {
						break loop
					}
					select {
					case <-ctx.Done():
						slog.Info(fmt.Sprintf("post worker `%d` canceled", workerNum))
						return nil
					default:
					}

					slog.Info(fmt.Sprintf("Post document `%s` by worker `%d`", path, workerNum))
					if err := f(ctx, path); err != nil {
						return failure.Wrap(err)
					}
				}
			}
			return nil
		})
	}

	eg.Go(func() error {
		if truncate {
			if _, err := solr.TruncateWithContext(ctx, core); err != nil {
				return failure.Translate(err, PostError, failure.Message("failed to truncate index"))
			}
		}

		for _, path := range paths {
			ch <- path
		}
		close(ch)

		wg.Wait()
		select {
		case <-ctx.Done():
			slog.Info("post canceled. start rollback...")
			if _, err := solr.RollbackWithContext(ctx, core); err != nil {
				return failure.Translate(err, PostError, failure.Message("failed to rollback index"))
			}
		default:
			if optimize {
				if _, err := solr.OptimizeWithContext(ctx, core); err != nil {
					return failure.Translate(err, PostError, failure.Message("failed to optimize index"))
				}
			} else {
				if _, err := solr.CommitWithContext(ctx, core); err != nil {
					return failure.Translate(err, PostError, failure.Message("failed to commit index"))
				}
			}
		}

		return nil
	})

	if err := eg.Wait(); err != nil {
		return failure.Wrap(err)
	}

	return nil
}
