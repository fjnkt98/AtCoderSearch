package update

import (
	"context"
	"fmt"

	"github.com/meilisearch/meilisearch-go"
	"golang.org/x/sync/errgroup"
)

type Documenter[D any] interface {
	Document(ctx context.Context) (D, error)
}

type RowReader[R any] interface {
	ReadRows(ctx context.Context, tx chan<- R) error
}

type Indexer interface {
	Settings() meilisearch.Settings
	Manager() meilisearch.IndexManager
}

func UpdateIndex[D any, R Documenter[D]](
	ctx context.Context,
	reader RowReader[R],
	indexer Indexer,
	chunkSize int,
	concurrent int,
) error {
	eg, ctx := errgroup.WithContext(ctx)

	rows := make(chan R, chunkSize)
	docs := make(chan D, chunkSize)

	eg.Go(func() error {
		defer close(rows)

		if err := reader.ReadRows(ctx, rows); err != nil {
			return fmt.Errorf("read rows: %w", err)
		}
		return nil
	})

	eg.Go(func() error {
		defer close(docs)

		eg, ctx := errgroup.WithContext(ctx)
		eg.SetLimit(concurrent)

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
						return fmt.Errorf("document: %w", err)
					}

					docs <- doc
					return nil
				})
			}
		}

		if err := eg.Wait(); err != nil {
			return err
		}

		return nil
	})

	eg.Go(func() error {
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

				fmt.Println(doc)
			}
		}

		return nil
	})

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}
