//go:generate mockgen -source=$GOFILE -destination=./mock/mock_$GOFILE -package=$GOPACKAGE

package repository

import (
	"context"

	"github.com/goark/errs"
	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
)

type DifficultyRepository interface {
	Save(ctx context.Context, difficulties []Difficulty) error
}

type Difficulty struct {
	bun.BaseModel    `bun:"table:difficulties,alias:d"`
	ProblemID        string   `bun:"problem_id,type:text"`
	Slope            *float64 `bun:"slope"`
	Intercept        *float64 `bun:"intercept"`
	Variance         *float64 `bun:"variance"`
	Difficulty       *int64   `bun:"difficulty"`
	Discrimination   *float64 `bun:"discrimination"`
	IrtLogLikelihood *float64 `bun:"irt_loglikelihood"`
	IrtUsers         *float64 `bun:"irt_users"`
	IsExperimental   *bool    `bun:"is_experimental"`
}

type difficultyRepository struct {
	db *bun.DB
}

func NewDifficultyRepository(db *bun.DB) DifficultyRepository {
	return &difficultyRepository{db: db}
}

func (r *difficultyRepository) Save(ctx context.Context, difficulties []Difficulty) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errs.New(
			"failed to start transaction",
			errs.WithCause(err),
		)
	}
	defer tx.Rollback()

	for _, chunk := range Chunks(difficulties, 1000) {
		_, err = r.db.NewMerge().
			Model(new(Difficulty)).
			With("difficulty", r.db.NewValues(&chunk)).
			Using("difficulty").
			On("?TableAlias.? = ?.?", bun.Ident("problem_id"), bun.Ident("difficulty"), bun.Ident("problem_id")).
			WhenUpdate("MATCHED", func(q *bun.UpdateQuery) *bun.UpdateQuery {
				return q.
					SetColumn("problem_id", "?.?", bun.Ident("difficulty"), bun.Ident("problem_id")).
					SetColumn("slope", "?.?", bun.Ident("difficulty"), bun.Ident("slope")).
					SetColumn("intercept", "?.?", bun.Ident("difficulty"), bun.Ident("intercept")).
					SetColumn("variance", "?.?", bun.Ident("difficulty"), bun.Ident("variance")).
					SetColumn("difficulty", "?.?", bun.Ident("difficulty"), bun.Ident("difficulty")).
					SetColumn("discrimination", "?.?", bun.Ident("difficulty"), bun.Ident("discrimination")).
					SetColumn("irt_loglikelihood", "?.?", bun.Ident("difficulty"), bun.Ident("irt_loglikelihood")).
					SetColumn("irt_users", "?.?", bun.Ident("difficulty"), bun.Ident("irt_users")).
					SetColumn("is_experimental", "?.?", bun.Ident("difficulty"), bun.Ident("is_experimental")).
					SetColumn("updated_at", "NOW()")
			}).
			WhenInsert("NOT MATCHED", func(q *bun.InsertQuery) *bun.InsertQuery {
				return q.
					Value("problem_id", "?.?", bun.Ident("difficulty"), bun.Ident("problem_id")).
					Value("slope", "?.?", bun.Ident("difficulty"), bun.Ident("slope")).
					Value("intercept", "?.?", bun.Ident("difficulty"), bun.Ident("intercept")).
					Value("variance", "?.?", bun.Ident("difficulty"), bun.Ident("variance")).
					Value("difficulty", "?.?", bun.Ident("difficulty"), bun.Ident("difficulty")).
					Value("discrimination", "?.?", bun.Ident("difficulty"), bun.Ident("discrimination")).
					Value("irt_loglikelihood", "?.?", bun.Ident("difficulty"), bun.Ident("irt_loglikelihood")).
					Value("irt_users", "?.?", bun.Ident("difficulty"), bun.Ident("irt_users")).
					Value("is_experimental", "?.?", bun.Ident("difficulty"), bun.Ident("is_experimental")).
					Value("created_at", "NOW()").
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
