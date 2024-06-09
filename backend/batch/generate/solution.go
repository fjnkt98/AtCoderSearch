package generate

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/goark/errs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

type SolutionRow struct {
	ProblemID string `bun:"problem_id"`
	UserID    string `bun:"user_id"`
}

type SolutionDoc struct {
	UniqueKey string `json:"uniqueKey"`
	ProblemID string `json:"problemId"`
	UserID    string `json:"userId"`
}

func (r *SolutionRow) Document(ctx context.Context) (*SolutionDoc, error) {
	return &SolutionDoc{
		UniqueKey: fmt.Sprintf("%s-%s", r.UserID, r.ProblemID),
		ProblemID: r.ProblemID,
		UserID:    r.UserID,
	}, nil
}

type SolutionRowReader struct {
	pool        *pgxpool.Pool
	interval    int
	lastUpdated *time.Time
}

func NewSolutionRowReader(pool *pgxpool.Pool, interval int, lastUpdated *time.Time) *SolutionRowReader {
	return &SolutionRowReader{
		pool:        pool,
		interval:    interval,
		lastUpdated: lastUpdated,
	}
}

func (r *SolutionRowReader) ReadRows(ctx context.Context, tx chan<- *SolutionRow) error {
	db := bun.NewDB(stdlib.OpenDBFromPool(r.pool), pgdialect.New())

	q := db.NewSelect().
		Distinct().
		ColumnExpr("s.problem_id").
		ColumnExpr("s.user_id").
		TableExpr("submissions AS s").
		Where("s.result = ?", "AC").
		Where("s.epoch_second > EXTRACT(EPOCH FROM CURRENT_DATE - CAST(? || ' day' AS INTERVAL))", r.interval)

	if r.lastUpdated != nil {
		slog.Info("reading solutions differentially", slog.Any("since", r.lastUpdated))
		q = q.Where("s.updated_at > ?", r.lastUpdated)
	}

	rows, err := q.Rows(ctx)
	if err != nil {
		return errs.New(
			"failed to read rows",
			errs.WithCause(err),
		)
	}
	defer rows.Close()
	defer close(tx)

	var count = 0
	for rows.Next() {
		select {
		case <-ctx.Done():
			return nil
		default:
			var row SolutionRow
			if err := db.ScanRow(ctx, rows, &row); err != nil {
				return errs.New(
					"failed to scan row",
					errs.WithCause(err),
				)
			}
			tx <- &row
			count++
		}
	}
	slog.Info("Finished reading row", slog.Int("count", count))

	return nil
}

func GenerateSolutionDocument(ctx context.Context, reader RowReader[*SolutionRow], saveDir string, options ...option) error {
	return GenerateDocument(ctx, reader, saveDir, options...)
}
