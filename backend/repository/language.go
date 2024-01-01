package repository

import (
	"context"

	"github.com/goark/errs"
	"github.com/uptrace/bun"
)

type Language struct {
	bun.BaseModel `bun:"table:languages,alias:l"`
	Language      string `bun:"language,pk"`
	Group         string `bun:"group"`
}

type LanguageRepository interface {
	Save(ctx context.Context, languages []Language) error
}

type languageRepository struct {
	db *bun.DB
}

func NewLanguageRepository(db *bun.DB) LanguageRepository {
	return &languageRepository{
		db: db,
	}
}

func (r *languageRepository) Save(ctx context.Context, languages []Language) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errs.New(
			"failed to start transaction to save submissions",
			errs.WithCause(err),
		)
	}
	defer tx.Rollback()

	_, err = r.db.NewInsert().
		Model(&languages).
		On("CONFLICT (?PKs) DO UPDATE").
		Set("? = EXCLUDED.?", bun.Ident("group"), bun.Ident("group")).
		Exec(ctx)
	if err != nil {
		return errs.New(
			"failed to insert languages",
			errs.WithCause(err),
		)
	}

	if err := tx.Commit(); err != nil {
		return errs.New(
			"failed to commit transaction to save languages",
			errs.WithCause(err),
		)
	}

	return nil
}
