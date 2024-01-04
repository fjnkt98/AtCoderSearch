package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/goark/errs"
	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
)

type UserRepository interface {
	Save(ctx context.Context, users []User) error
	FetchRatingByUserName(ctx context.Context, username string) (int, error)
}

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	UserName      string  `bun:"user_name,type:text"`
	Rating        int     `bun:"rating"`
	HighestRating int     `bun:"highest_rating"`
	Affiliation   *string `bun:"affiliation,type:text"`
	BirthYear     *int    `bun:"birth_year"`
	Country       *string `bun:"country,type:text"`
	Crown         *string `bun:"crown,type:text"`
	JoinCount     int     `bun:"join_count"`
	Rank          int     `bun:"rank"`
	ActiveRank    *int    `bun:"active_rank"`
	Wins          int     `bun:"wins"`
}

type userRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Save(ctx context.Context, users []User) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errs.New(
			"failed to start transaction to save user information",
			errs.WithCause(err),
		)
	}
	defer tx.Rollback()

	for _, chunk := range Chunks(users, 1000) {
		_, err = r.db.NewMerge().
			Model(new(User)).
			With("user", r.db.NewValues(&chunk)).
			Using("?", bun.Ident("user")).
			On("?TableAlias.? = ?.?", bun.Ident("user_name"), bun.Ident("user"), bun.Ident("user_name")).
			WhenUpdate("MATCHED", func(q *bun.UpdateQuery) *bun.UpdateQuery {
				return q.
					Set("? = ?.?", bun.Ident("user_name"), bun.Ident("user"), bun.Ident("user_name")).
					Set("? = ?.?", bun.Ident("rating"), bun.Ident("user"), bun.Ident("rating")).
					Set("? = ?.?", bun.Ident("highest_rating"), bun.Ident("user"), bun.Ident("highest_rating")).
					Set("? = ?.?", bun.Ident("affiliation"), bun.Ident("user"), bun.Ident("affiliation")).
					Set("? = ?.?", bun.Ident("birth_year"), bun.Ident("user"), bun.Ident("birth_year")).
					Set("? = ?.?", bun.Ident("country"), bun.Ident("user"), bun.Ident("country")).
					Set("? = ?.?", bun.Ident("crown"), bun.Ident("user"), bun.Ident("crown")).
					Set("? = ?.?", bun.Ident("join_count"), bun.Ident("user"), bun.Ident("join_count")).
					Set("? = ?.?", bun.Ident("rank"), bun.Ident("user"), bun.Ident("rank")).
					Set("? = ?.?", bun.Ident("active_rank"), bun.Ident("user"), bun.Ident("active_rank")).
					Set("? = ?.?", bun.Ident("wins"), bun.Ident("user"), bun.Ident("wins")).
					Set("? = NOW()", bun.Ident("updated_at"))
			}).
			WhenInsert("NOT MATCHED", func(q *bun.InsertQuery) *bun.InsertQuery {
				return q.
					Value("user_name", "?.?", bun.Ident("user"), bun.Ident("user_name")).
					Value("rating", "?.?", bun.Ident("user"), bun.Ident("rating")).
					Value("highest_rating", "?.?", bun.Ident("user"), bun.Ident("highest_rating")).
					Value("affiliation", "?.?", bun.Ident("user"), bun.Ident("affiliation")).
					Value("birth_year", "?.?", bun.Ident("user"), bun.Ident("birth_year")).
					Value("country", "?.?", bun.Ident("user"), bun.Ident("country")).
					Value("crown", "?.?", bun.Ident("user"), bun.Ident("crown")).
					Value("join_count", "?.?", bun.Ident("user"), bun.Ident("join_count")).
					Value("rank", "?.?", bun.Ident("user"), bun.Ident("rank")).
					Value("active_rank", "?.?", bun.Ident("user"), bun.Ident("active_rank")).
					Value("wins", "?.?", bun.Ident("user"), bun.Ident("wins")).
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

func (r *userRepository) FetchRatingByUserName(ctx context.Context, username string) (int, error) {
	var rating int
	err := r.db.NewSelect().
		Model(new(User)).
		Column("rating").
		Where("? = ?", bun.Ident("user_name"), username).
		Limit(1).
		Scan(ctx, &rating)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errs.Wrap(
				err,
				errs.WithContext("user_name", username),
			)
		} else {
			return 0, errs.New(
				"failed to get the rating of the specified user",
				errs.WithCause(err),
				errs.WithContext("user_name", username),
			)
		}
	}

	return rating, nil
}
