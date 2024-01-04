//go:generate mockgen -source=$GOFILE -destination=./mock/mock_$GOFILE -package=$GOPACKAGE

package repository

import (
	"context"

	"github.com/goark/errs"
	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
)

type SubmissionRepository interface {
	Save(ctx context.Context, submissions []Submission) error
	FetchLanguages(ctx context.Context) ([]string, error)
}

type Submission struct {
	bun.BaseModel `bun:"table:submissions,alias:s"`
	ID            int     `bun:"id"`
	EpochSecond   int64   `bun:"epoch_second"`
	ProblemID     string  `bun:"problem_id,type:text"`
	ContestID     string  `bun:"contest_id,type:text"`
	UserID        string  `bun:"user_id,type:text"`
	Language      string  `bun:"language,type:text"`
	Point         float64 `bun:"point"`
	Length        int     `bun:"length"`
	Result        string  `bun:"result,type:text"`
	ExecutionTime *int    `bun:"execution_time"`
}

type submissionRepository struct {
	db *bun.DB
}

func NewSubmissionRepository(db *bun.DB) SubmissionRepository {
	return &submissionRepository{db: db}
}

func (r *submissionRepository) Save(ctx context.Context, submissions []Submission) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errs.New(
			"failed to start transaction to save submissions",
			errs.WithCause(err),
		)
	}
	defer tx.Rollback()

	for _, chunk := range Chunks(submissions, 1000) {
		_, err := r.db.NewMerge().
			Model(new(Submission)).
			With("submission", r.db.NewValues(&chunk)).
			Using("submission").
			On("?TableAlias.? = ?.?", bun.Ident("id"), bun.Ident("submission"), bun.Ident("id")).
			WhenUpdate("MATCHED", func(q *bun.UpdateQuery) *bun.UpdateQuery {
				return q.
					SetColumn("id", "?.?", bun.Ident("submission"), bun.Ident("id")).
					SetColumn("epoch_second", "?.?", bun.Ident("submission"), bun.Ident("epoch_second")).
					SetColumn("problem_id", "?.?", bun.Ident("submission"), bun.Ident("problem_id")).
					SetColumn("contest_id", "?.?", bun.Ident("submission"), bun.Ident("contest_id")).
					SetColumn("user_id", "?.?", bun.Ident("submission"), bun.Ident("user_id")).
					SetColumn("language", "?.?", bun.Ident("submission"), bun.Ident("language")).
					SetColumn("point", "?.?", bun.Ident("submission"), bun.Ident("point")).
					SetColumn("length", "?.?", bun.Ident("submission"), bun.Ident("length")).
					SetColumn("result", "?.?", bun.Ident("submission"), bun.Ident("result")).
					SetColumn("execution_time", "?.?", bun.Ident("submission"), bun.Ident("execution_time"))
			}).
			WhenInsert("NOT MATCHED", func(q *bun.InsertQuery) *bun.InsertQuery {
				return q.
					Value("id", "?.?", bun.Ident("submission"), bun.Ident("id")).
					Value("epoch_second", "?.?", bun.Ident("submission"), bun.Ident("epoch_second")).
					Value("problem_id", "?.?", bun.Ident("submission"), bun.Ident("problem_id")).
					Value("contest_id", "?.?", bun.Ident("submission"), bun.Ident("contest_id")).
					Value("user_id", "?.?", bun.Ident("submission"), bun.Ident("user_id")).
					Value("language", "?.?", bun.Ident("submission"), bun.Ident("language")).
					Value("point", "?.?", bun.Ident("submission"), bun.Ident("point")).
					Value("length", "?.?", bun.Ident("submission"), bun.Ident("length")).
					Value("result", "?.?", bun.Ident("submission"), bun.Ident("result")).
					Value("execution_time", "?.?", bun.Ident("submission"), bun.Ident("execution_time"))
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

func (r *submissionRepository) FetchLanguages(ctx context.Context) ([]string, error) {
	var languages []string
	err := r.db.NewSelect().
		Model(new(Submission)).
		Column("language").
		Distinct().
		Scan(ctx, &languages)
	if err != nil {
		return nil, errs.New(
			"failed to fetch languages from submissions table",
			errs.WithCause(err),
		)
	}

	return languages, nil
}
