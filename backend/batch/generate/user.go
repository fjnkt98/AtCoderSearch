package generate

import (
	"context"
	"fmt"

	"github.com/goark/errs"
	"github.com/uptrace/bun"
)

type UserRow struct {
	UserName      string  `bun:"user_name"`
	Rating        int     `bun:"rating"`
	HighestRating int     `bun:"highest_rating"`
	Affiliation   *string `bun:"affiliation"`
	BirthYear     *int    `bun:"birth_year"`
	Country       *string `bun:"country"`
	Crown         *string `bun:"crown"`
	JoinCount     int     `bun:"join_count"`
	Rank          int     `bun:"rank"`
	ActiveRank    *int    `bun:"active_rank"`
	Wins          int     `bun:"wins"`
}

type UserDoc struct {
	UserName      string  `json:"user_name"`
	Rating        int     `json:"rating"`
	HighestRating int     `json:"highest_rating"`
	Affiliation   *string `json:"affiliation"`
	BirthYear     *int    `json:"birth_year"`
	Country       *string `json:"country"`
	Crown         *string `json:"crown"`
	JoinCount     int     `json:"join_count"`
	Rank          int     `json:"rank"`
	ActiveRank    *int    `json:"active_rank"`
	Wins          int     `json:"wins" `
	Color         string  `json:"color"`
	HighestColor  string  `json:"highest_color"`
	UserURL       string  `json:"user_url"`
}

func (r *UserRow) Document(ctx context.Context) (*UserDoc, error) {
	return &UserDoc{
		UserName:      r.UserName,
		Rating:        r.Rating,
		HighestRating: r.HighestRating,
		Affiliation:   r.Affiliation,
		BirthYear:     r.BirthYear,
		Country:       r.Country,
		Crown:         r.Crown,
		JoinCount:     r.JoinCount,
		Rank:          r.Rank,
		ActiveRank:    r.ActiveRank,
		Wins:          r.Wins,
		Color:         RateToColor(r.Rating),
		HighestColor:  RateToColor(r.HighestRating),
		UserURL:       fmt.Sprintf("https://atcoder.jp/users/%s", r.UserName),
	}, nil
}

type userRowReader struct {
	db *bun.DB
}

func NewUserRowReader(db *bun.DB) RowReader[*UserRow] {
	return &userRowReader{
		db: db,
	}
}

func (r *userRowReader) ReadRows(ctx context.Context, tx chan<- *UserRow) error {
	rows, err := r.db.NewSelect().
		Column(
			"user_name",
			"rating",
			"highest_rating",
			"affiliation",
			"birth_year",
			"country",
			"crown",
			"join_count",
			"rank",
			"active_rank",
			"wins",
		).
		Table("users").
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
			var row UserRow
			err := r.db.ScanRow(ctx, rows, &row)
			if err != nil {
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

func NewUserGenerator(reader RowReader[*UserRow], saveDir string, chunkSize, concurrent int) DocumentGenerator {
	return NewDocumentGenerator(reader, saveDir, chunkSize, concurrent)
}
