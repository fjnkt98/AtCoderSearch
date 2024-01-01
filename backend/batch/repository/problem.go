package repository

import (
	"context"

	"github.com/goark/errs"
	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
)

type ProblemRepository interface {
	Save(ctx context.Context, problems []Problem) error
	FetchIDs(ctx context.Context) ([]string, error)
}

type Problem struct {
	bun.BaseModel `bun:"table:problems,alias:p"`
	ProblemID     string `bun:"problem_id,type:text"`
	ContestID     string `bun:"contest_id,type:text"`
	ProblemIndex  string `bun:"problem_index,type:text"`
	Name          string `bun:"name,type:text"`
	Title         string `bun:"title,type:text"`
	URL           string `bun:"url,type:text"`
	HTML          string `bun:"html,type:text"`
}

type problemRepository struct {
	db *bun.DB
}

func NewProblemRepository(db *bun.DB) ProblemRepository {
	return &problemRepository{db: db}
}

func (r *problemRepository) Save(ctx context.Context, problems []Problem) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errs.New(
			"failed to start transaction",
			errs.WithCause(err),
		)
	}
	defer tx.Rollback()

	for _, chunk := range Chunks(problems, 1000) {
		_, err := r.db.NewMerge().
			Model(new(Problem)).
			With("problem", r.db.NewValues(&chunk)).
			Using("problem").
			On("?TableAlias.? = ?.?", bun.Ident("problem_id"), bun.Ident("problem"), bun.Ident("problem_id")).
			WhenUpdate("MATCHED", func(q *bun.UpdateQuery) *bun.UpdateQuery {
				return q.
					SetColumn("problem_id", "?.?", bun.Ident("problem"), bun.Ident("problem_id")).
					SetColumn("contest_id", "?.?", bun.Ident("problem"), bun.Ident("contest_id")).
					SetColumn("problem_index", "?.?", bun.Ident("problem"), bun.Ident("problem_index")).
					SetColumn("name", "?.?", bun.Ident("problem"), bun.Ident("name")).
					SetColumn("title", "?.?", bun.Ident("problem"), bun.Ident("title")).
					SetColumn("url", "?.?", bun.Ident("problem"), bun.Ident("url")).
					SetColumn("html", "?.?", bun.Ident("problem"), bun.Ident("html")).
					SetColumn("updated_at", "NOW()")
			}).
			WhenInsert("NOT MATCHED", func(q *bun.InsertQuery) *bun.InsertQuery {
				return q.
					Value("problem_id", "?.?", bun.Ident("problem"), bun.Ident("problem_id")).
					Value("contest_id", "?.?", bun.Ident("problem"), bun.Ident("contest_id")).
					Value("problem_index", "?.?", bun.Ident("problem"), bun.Ident("problem_index")).
					Value("name", "?.?", bun.Ident("problem"), bun.Ident("name")).
					Value("title", "?.?", bun.Ident("problem"), bun.Ident("title")).
					Value("url", "?.?", bun.Ident("problem"), bun.Ident("url")).
					Value("html", "?.?", bun.Ident("problem"), bun.Ident("html")).
					Value("updated_at", "NOW()").
					Value("created_at", "NOW()")
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

func (r *problemRepository) FetchIDs(ctx context.Context) ([]string, error) {
	ids := make([]string, 0, 8192)
	if err := r.db.NewSelect().Model(new(Problem)).Column("problem_id").Scan(ctx, &ids); err != nil {
		return nil, errs.New(
			"failed to fetch problems ids",
			errs.WithCause(err),
		)
	}
	return ids, nil
}
