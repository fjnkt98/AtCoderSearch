package generate

import (
	"context"
	"fmt"

	"github.com/goark/errs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

type SolutionRow struct {
	SubmissionID string `bun:"submission_id"`
	UserID       string `bun:"user_id"`
}

type SolutionDoc struct {
	UniqueKey    string `json:"uniqueKey"`
	SubmissionID string `json:"submissionId"`
	UserID       string `json:"userId"`
}

func (r *SolutionRow) Document(ctx context.Context) (*SolutionDoc, error) {
	return &SolutionDoc{
		UniqueKey:    fmt.Sprintf("%s-%s", r.UserID, r.SubmissionID),
		SubmissionID: r.SubmissionID,
		UserID:       r.UserID,
	}, nil
}

type SolutionRowReader struct {
	pool *pgxpool.Pool
}

func NewSolutionRowReader(pool *pgxpool.Pool) *SolutionRowReader {
	return &SolutionRowReader{
		pool: pool,
	}
}

func (r *SolutionRowReader) ReadRows(ctx context.Context, tx chan<- *SolutionRow) error {
	db := bun.NewDB(stdlib.OpenDBFromPool(r.pool), pgdialect.New())
	rows, err := db.NewSelect().
		ColumnExpr("s.id AS submission_id").
		ColumnExpr("s.user_id").
		TableExpr("submissions AS s").
		Where("s.result = ?", "AC").
		Rows(ctx)

	if err != nil {
		return errs.New(
			"failed to read rows",
			errs.WithCause(err),
		)
	}
	defer rows.Close()
	defer close(tx)

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
		}
	}

	return nil
}

func GenerateSolutionDocument(ctx context.Context, reader RowReader[*SolutionRow], saveDir string, options ...option) error {
	return GenerateDocument(ctx, reader, saveDir, options...)
}
