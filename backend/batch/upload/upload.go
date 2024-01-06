package upload

import (
	"context"
	"fjnkt98/atcodersearch/batch"
	"fjnkt98/atcodersearch/pkg/solr"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/goark/errs"
	"golang.org/x/sync/errgroup"
)

type DocumentUploader interface {
	batch.Batch
	Upload(ctx context.Context) error
}

type documentUploader struct {
	core   solr.SolrCore
	config documentUploaderConfig
}

type documentUploaderConfig struct {
	SaveDir    string `json:"save_dir"`
	Concurrent int    `json:"concurrent"`
	Optimize   bool   `json:"optimize"`
	Truncate   bool   `json:"truncate"`
}

func NewDocumentUploader(core solr.SolrCore, saveDir string, concurrent int, optimize bool, truncate bool) DocumentUploader {
	return &documentUploader{
		core: core,
		config: documentUploaderConfig{
			SaveDir:    saveDir,
			Concurrent: concurrent,
			Optimize:   optimize,
			Truncate:   truncate,
		},
	}
}

func (u *documentUploader) Name() string {
	return "DocumentUploader"
}

func (u *documentUploader) Config() any {
	return u.config
}

func (u *documentUploader) Run(ctx context.Context) error {
	return u.Upload(ctx)
}

func (u *documentUploader) Upload(ctx context.Context) error {
	slog.Info(fmt.Sprintf("Start to post documents in `%s`", u.config.SaveDir))

	paths, err := filepath.Glob(filepath.Join(u.config.SaveDir, "doc-*.json"))
	if err != nil {
		return errs.New(
			"failed to get files from the directory",
			errs.WithCause(err),
			errs.WithContext("directory", u.config.SaveDir),
		)
	}

	ch := make(chan string, len(paths))

	eg, ctx := errgroup.WithContext(ctx)
	var wg sync.WaitGroup

	f := func(ctx context.Context, p string) error {
		file, err := os.Open(p)
		if err != nil {
			return errs.New(
				"failed to open file",
				errs.WithCause(err),
				errs.WithContext("filepath", p),
			)
		}
		defer file.Close()
		if _, err := solr.PostWithContext(ctx, u.core, file, "application/json"); err != nil {
			return errs.New(
				"failed to post file",
				errs.WithCause(err),
				errs.WithContext("filepath", p),
			)
		}

		return nil
	}

	for i := 0; i < u.config.Concurrent; i++ {
		wg.Add(1)
		workerNum := i
		eg.Go(func() error {
			defer wg.Done()

		loop:
			for {
				select {
				case <-ctx.Done():
					slog.Info(fmt.Sprintf("post worker `%d` canceled", workerNum))
					return batch.ErrInterrupt
				case path, ok := <-ch:
					if !ok {
						break loop
					}
					select {
					case <-ctx.Done():
						slog.Info(fmt.Sprintf("post worker `%d` canceled", workerNum))
						return batch.ErrInterrupt
					default:
					}

					slog.Info(fmt.Sprintf("Post document `%s` by worker `%d`", path, workerNum))
					if err := f(ctx, path); err != nil {
						return errs.Wrap(err)
					}
				}
			}
			return nil
		})
	}

	eg.Go(func() error {
		if u.config.Truncate {
			slog.Info("Start to truncate index...")
			if _, err := solr.TruncateWithContext(ctx, u.core); err != nil {
				return errs.New(
					"failed to truncate index",
					errs.WithCause(err),
				)
			}
			slog.Info("Finished truncating index successfully.")
		}

		for _, path := range paths {
			ch <- path
		}
		close(ch)

		wg.Wait()
		select {
		case <-ctx.Done():
			slog.Info("post canceled. start rollback...")
			if _, err := solr.Rollback(u.core); err != nil {
				return errs.New(
					"failed to rollback index",
					errs.WithCause(err),
				)
			}
			slog.Info("rollback finished successfully.")
			return batch.ErrInterrupt
		default:
			if u.config.Optimize {
				slog.Info("Start to optimize index...")
				if _, err := solr.OptimizeWithContext(ctx, u.core); err != nil {
					return errs.New(
						"failed to optimize index",
						errs.WithCause(err),
					)
				}
				slog.Info("Finished optimize index successfully.")
			} else {
				slog.Info("Start to commit index...")
				if _, err := solr.CommitWithContext(ctx, u.core); err != nil {
					return errs.New(
						"failed to commit index",
						errs.WithCause(err),
					)
				}
				slog.Info("Finished commit index successfully.")
			}
		}

		return nil
	})

	if err := eg.Wait(); err != nil {
		return errs.Wrap(err)
	}

	return nil
}
