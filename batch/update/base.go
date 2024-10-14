package update

import (
	"context"
	"fmt"
	"log/slog"
	"time"

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
	Manager() meilisearch.ServiceManager
	Settings() *meilisearch.Settings
	PrimaryKey() string
	IndexName() string
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
		client := indexer.Manager()

		activeIndexName := indexer.IndexName()
		standbyIndexName := fmt.Sprintf("%s-standby", indexer.IndexName())

		standbyIndex := client.Index(standbyIndexName)

		stats, err := client.GetStatsWithContext(ctx)
		if err != nil {
			return fmt.Errorf("get stats: %w", err)
		}

		if _, ok := stats.Indexes[activeIndexName]; !ok {
			// アクティブインデックスが無かったら作る
			info, err := client.CreateIndexWithContext(ctx, &meilisearch.IndexConfig{
				Uid:        activeIndexName,
				PrimaryKey: indexer.PrimaryKey(),
			})
			if err != nil {
				return fmt.Errorf("create active index: %w", err)
			}

			task, err := client.WaitForTaskWithContext(ctx, info.TaskUID, 1*time.Second)
			if err != nil {
				return fmt.Errorf("task creating active index failed: %w", err)
			}
			slog.LogAttrs(ctx, slog.LevelInfo, "active index created", slog.Int64("uid", task.UID), slog.String("indexUID", task.IndexUID))
		}

		// スタンバイインデックス初期化
		if _, ok := stats.Indexes[standbyIndexName]; ok {
			// ドキュメントを消す
			info, err := standbyIndex.DeleteAllDocumentsWithContext(ctx)
			if err != nil {
				return fmt.Errorf("delete documents of the standby index: %w", err)
			}

			task, err := client.WaitForTaskWithContext(ctx, info.TaskUID, 1*time.Second)
			if err != nil {
				return fmt.Errorf("task deleting standby index failed: %w", err)
			}
			slog.LogAttrs(ctx, slog.LevelInfo, "standby index truncated", slog.Int64("uid", task.UID), slog.String("indexUID", task.IndexUID))
		} else {
			// 無いので作る
			info, err := client.CreateIndexWithContext(ctx, &meilisearch.IndexConfig{
				Uid:        standbyIndexName,
				PrimaryKey: indexer.PrimaryKey(),
			})
			if err != nil {
				return fmt.Errorf("create standby index: %w", err)
			}

			task, err := client.WaitForTaskWithContext(ctx, info.TaskUID, 1*time.Second)
			if err != nil {
				return fmt.Errorf("task creating standby index failed: %w", err)
			}
			slog.LogAttrs(ctx, slog.LevelInfo, "standby index created", slog.Int64("uid", task.UID), slog.String("indexUID", task.IndexUID))
		}

		// スタンバイインデックスsettings更新
		info, err := standbyIndex.UpdateIndexWithContext(ctx, indexer.PrimaryKey())
		if err != nil {
			return fmt.Errorf("update standby index primary key: %w", err)
		}
		if _, err := client.WaitForTaskWithContext(ctx, info.TaskUID, 1*time.Second); err != nil {
			return fmt.Errorf("task updating standby index failed: %w", err)
		}

		info, err = standbyIndex.UpdateSettingsWithContext(ctx, indexer.Settings())
		if err != nil {
			return fmt.Errorf("update standby index settings: %w", err)
		}

		task, err := client.WaitForTaskWithContext(ctx, info.TaskUID, 1*time.Second)
		if err != nil {
			return fmt.Errorf("task updating standby index failed: %w", err)
		}
		slog.LogAttrs(ctx, slog.LevelInfo, "settings of standby index updated", slog.Int64("uid", task.UID), slog.String("indexUID", task.IndexUID))

		// スタンバイインデックスにドキュメント追加
		buf := make([]D, 0, 2*chunkSize)
		tasks := make([]int64, 16)
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

				buf = append(buf, doc)

				if len(buf) >= chunkSize {
					info, err := standbyIndex.AddDocumentsWithContext(ctx, buf, indexer.PrimaryKey())
					if err != nil {
						return fmt.Errorf("add documents: %w", err)
					}

					tasks = append(tasks, info.TaskUID)

					buf = buf[:0]
				}
			}
		}
		if len(buf) >= 0 {
			info, err := standbyIndex.AddDocumentsWithContext(ctx, buf, indexer.PrimaryKey())
			if err != nil {
				return fmt.Errorf("add documents: %w", err)
			}

			tasks = append(tasks, info.TaskUID)
		}

		for _, task := range tasks {
			_, err := client.WaitForTaskWithContext(ctx, task, 1*time.Second)
			if err != nil {
				return fmt.Errorf("task adding documents failed: %w", err)
			}
		}

		// スワップ
		info, err = client.SwapIndexesWithContext(ctx, []*meilisearch.SwapIndexesParams{
			{
				Indexes: []string{activeIndexName, standbyIndexName},
			},
		})
		if err != nil {
			return fmt.Errorf("swap indexes: %w", err)
		}

		task, err = client.WaitForTaskWithContext(ctx, info.TaskUID, 1*time.Second)
		if err != nil {
			return fmt.Errorf("task swapping indexes failed: %w", err)
		}
		slog.LogAttrs(ctx, slog.LevelInfo, "index successfully updated", slog.Int64("uid", task.UID), slog.String("indexUID", indexer.IndexName()))

		return nil
	})

	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}
