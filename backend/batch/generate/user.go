package generate

import (
	"context"
	"fjnkt98/atcodersearch/batch"
	"fjnkt98/atcodersearch/config"
	"fmt"

	"github.com/goark/errs"
	"github.com/uptrace/bun"
	"golang.org/x/exp/slog"
)

type UserGenerator interface {
	batch.Batch
	GenerateUser(ctx context.Context) error
}

type userGenerator struct {
	defaultGenerator
}

func NewUserGenerator(cfg config.GenerateUserConfig, reader RowReader) UserGenerator {
	return &userGenerator{
		defaultGenerator{
			cfg: config.GenerateCommonConfig{
				SaveDir:    cfg.SaveDir,
				ChunkSize:  cfg.ChunkSize,
				Concurrent: cfg.Concurrent,
			},
			reader: reader,
		},
	}
}

func (g *userGenerator) Name() string {
	return "UserGenerator"
}

func (g *userGenerator) Run(ctx context.Context) error {
	return g.GenerateUser(ctx)
}

func (g *userGenerator) GenerateUser(ctx context.Context) error {
	if err := g.Generate(ctx); err != nil {
		return errs.Wrap(err)
	}
	return nil
}

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

type UserDocument struct {
	UserName      string  `solr:"user_name"`
	Rating        int     `solr:"rating"`
	HighestRating int     `solr:"highest_rating"`
	Affiliation   *string `solr:"affiliation"`
	BirthYear     *int    `solr:"birth_year"`
	Country       *string `solr:"country"`
	Crown         *string `solr:"crown"`
	JoinCount     int     `solr:"join_count"`
	Rank          int     `solr:"rank"`
	ActiveRank    *int    `solr:"active_rank"`
	Wins          int     `solr:"wins" `
	Color         string  `solr:"color"`
	HighestColor  string  `solr:"highest_color"`
	UserURL       string  `solr:"user_url"`
}

func (r *UserRow) Document(ctx context.Context) (map[string]any, error) {
	return StructToMap(UserDocument{
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
	}), nil
}

type userRowReader struct {
	db *bun.DB
}

func NewUserRowReader(db *bun.DB) RowReader {
	return &userRowReader{
		db: db,
	}
}

func (r *userRowReader) ReadRows(ctx context.Context, tx chan<- Documenter) error {
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
			slog.Info("read rows canceled.")
			return batch.ErrInterrupt
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
