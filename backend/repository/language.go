//go:generate mockgen -source=$GOFILE -destination=./mock/mock_$GOFILE -package=$GOPACKAGE

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
	FetchLanguagesByGroup(ctx context.Context, groups []string) ([]string, error)
	FetchLanguageGroups(ctx context.Context) ([]string, error)
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

	_, err = tx.NewInsert().
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

func (r *languageRepository) FetchLanguagesByGroup(ctx context.Context, groups []string) ([]string, error) {
	var languages []string
	query := r.db.NewSelect().
		Model(new(Language)).
		Column("language").
		Order("language")
	if len(groups) > 0 {
		query = query.Where("? IN (?)", bun.Ident("group"), bun.In(groups))
	}

	if err := query.Scan(ctx, &languages); err != nil {
		return nil, errs.New(
			"failed to fetch languages",
			errs.WithCause(err),
		)
	}

	return languages, nil
}

func (r *languageRepository) FetchLanguageGroups(ctx context.Context) ([]string, error) {
	var languages []string
	err := r.db.NewSelect().
		Model(new(Language)).
		Column("group").
		Distinct().
		Order("group").
		Scan(ctx, &languages)

	if err != nil {
		return nil, errs.New(
			"failed to fetch language groups",
			errs.WithCause(err),
		)
	}

	return languages, nil
}
