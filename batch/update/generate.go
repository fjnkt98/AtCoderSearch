package update

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type Documenter[D any] interface {
	Document(ctx context.Context) (D, error)
}

type RowReader[R any] interface {
	ReadRows(ctx context.Context, tx chan<- R) error
}

type generateDocumentOptions struct {
	ChunkSize  int
	Concurrent int
}

type generateOptionFunc func(*generateDocumentOptions)

func WithChunkSize(chunkSize int) generateOptionFunc {
	return func(opt *generateDocumentOptions) {
		opt.ChunkSize = chunkSize
	}
}

func WithConcurrent(concurrent int) generateOptionFunc {
	return func(opt *generateDocumentOptions) {
		opt.Concurrent = concurrent
	}
}

func GenerateDocument[D any, R Documenter[D]](ctx context.Context, reader RowReader[R], options ...generateOptionFunc) error {
	option := &generateDocumentOptions{
		ChunkSize:  1000,
		Concurrent: 4,
	}
	for _, opt := range options {
		opt(option)
	}

	rows := make(chan R, option.ChunkSize)
	docs := make(chan D, option.ChunkSize)

	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(option.Concurrent)

	// read rows and send it
	eg.Go(func() error {
		defer close(rows)

		if err := reader.ReadRows(ctx, rows); err != nil {
			return err
		}
		return nil
	})

	// collect and index docs
	eg.Go(func() error {
		buffer := make([]D, 0, option.ChunkSize)

	loop:
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case doc, ok := <-docs:
				if !ok {
					break loop
				}

				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
				}

				buffer = append(buffer, doc)

				if len(buffer) >= option.ChunkSize {
					// TODO

					buffer = buffer[:0]
				}
			}
		}

		if len(buffer) > 0 {
			// TODO
		}
		return nil
	})

	// receive row and send doc
	eg.Go(func() error {
		defer close(docs)
	loop:
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case row, ok := <-rows:
				if !ok {
					break loop
				}
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
				}

				eg.Go(func() error {
					doc, err := row.Document(ctx)
					if err != nil {
						return err
					}

					docs <- doc
					return nil
				})
			}
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}
