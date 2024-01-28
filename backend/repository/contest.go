//go:generate mockgen -source=$GOFILE -destination=./mock/mock_$GOFILE -package=$GOPACKAGE

package repository

import (
	"context"

	"github.com/goark/errs"
	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
)

type ContestRepository interface {
	Save(ctx context.Context, contests []Contest) error
	FetchContestIDs(ctx context.Context, categories []string) ([]string, error)
	FetchCategories(ctx context.Context) ([]string, error)
}

type Contest struct {
	bun.BaseModel    `bun:"table:contests,alias:c"`
	ContestID        string `bun:"contest_id,type:text"`
	StartEpochSecond int64  `bun:"start_epoch_second"`
	DurationSecond   int64  `bun:"duration_second"`
	Title            string `bun:"title,type:text"`
	RateChange       string `bun:"rate_change,type:text"`
	Category         string `bun:"category,type:text"`
}

type contestRepository struct {
	db *bun.DB
}

func NewContestRepository(db *bun.DB) ContestRepository {
	return &contestRepository{db: db}
}

func (r *contestRepository) Save(ctx context.Context, contests []Contest) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errs.New(
			"failed to start transaction to save contest information",
			errs.WithCause(err),
		)
	}
	defer tx.Rollback()

	for _, chunk := range Chunks(contests, 1000) {
		_, err = tx.NewMerge().
			Model(new(Contest)).
			With("contest", tx.NewValues(&chunk)).
			Using("contest").
			On("?TableAlias.? = ?.?", bun.Ident("contest_id"), bun.Ident("contest"), bun.Ident("contest_id")).
			WhenUpdate("MATCHED", func(q *bun.UpdateQuery) *bun.UpdateQuery {
				return q.
					SetColumn("contest_id", "?.?", bun.Ident("contest"), bun.Ident("contest_id")).
					SetColumn("start_epoch_second", "?.?", bun.Ident("contest"), bun.Ident("start_epoch_second")).
					SetColumn("duration_second", "?.?", bun.Ident("contest"), bun.Ident("duration_second")).
					SetColumn("title", "?.?", bun.Ident("contest"), bun.Ident("title")).
					SetColumn("rate_change", "?.?", bun.Ident("contest"), bun.Ident("rate_change")).
					SetColumn("category", "?.?", bun.Ident("contest"), bun.Ident("category")).
					SetColumn("updated_at", "NOW()")
			}).
			WhenInsert("NOT MATCHED", func(q *bun.InsertQuery) *bun.InsertQuery {
				return q.
					Value("contest_id", "?.?", bun.Ident("contest"), bun.Ident("contest_id")).
					Value("start_epoch_second", "?.?", bun.Ident("contest"), bun.Ident("start_epoch_second")).
					Value("duration_second", "?.?", bun.Ident("contest"), bun.Ident("duration_second")).
					Value("title", "?.?", bun.Ident("contest"), bun.Ident("title")).
					Value("rate_change", "?.?", bun.Ident("contest"), bun.Ident("rate_change")).
					Value("category", "?.?", bun.Ident("contest"), bun.Ident("category")).
					Value("updated_at", "NOW()")
			}).
			Exec(ctx)
		if err != nil {
			return errs.New(
				"failed to execute sql",
				errs.WithCause(err),
			)
		}
	}

	if err := tx.Commit(); err != nil {
		return errs.New(
			"failed to commit transaction",
			errs.WithCause(err),
		)
	}

	return nil
}

func (r *contestRepository) FetchContestIDs(ctx context.Context, targets []string) ([]string, error) {
	query := r.db.NewSelect().
		Model(new(Contest)).
		Column("contest_id").
		Order("start_epoch_second DESC")
	if len(targets) > 0 {
		query = query.Where("category IN (?)", bun.In(targets))
	}

	ids := make([]string, 0)
	if _, err := query.Exec(ctx, &ids); err != nil {
		return nil, errs.New(
			"failed to fetch contest ids",
			errs.WithCause(err),
		)
	}

	return ids, nil
}

func (r *contestRepository) FetchCategories(ctx context.Context) ([]string, error) {
	var categories []string

	err := r.db.NewSelect().
		Model(new(Contest)).
		Column("category").
		Distinct().
		Order("category").
		Scan(ctx, &categories)

	if err != nil {
		return nil, errs.New(
			"failed to fetch categories",
			errs.WithCause(err),
		)
	}

	return categories, nil
}
