package update

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/meilisearch/meilisearch-go"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

type UserIndexer struct {
	client meilisearch.ServiceManager
}

func NewUserIndexer(client meilisearch.ServiceManager) *UserIndexer {
	return &UserIndexer{
		client: client,
	}
}

func (ix *UserIndexer) Manager() meilisearch.ServiceManager {
	return ix.client
}

func (ix *UserIndexer) Settings() *meilisearch.Settings {
	return &meilisearch.Settings{
		Dictionary: []string{},
		DisplayedAttributes: []string{
			"*",
		},
		FilterableAttributes: []string{
			"userID",
			"rating",
			"highestRating",
			"affiliation",
			"birthYear",
			"country",
			"crown",
			"joinCount",
			"rank",
			"activeRank",
			"wins",
		},
		SearchableAttributes: []string{
			"userId",
			"affiliation",
		},
		SortableAttributes: []string{
			"rating",
			"birthYear",
		},
		Synonyms: map[string][]string{},
	}
}

func (ix *UserIndexer) PrimaryKey() string {
	return "userID"
}

func (ix *UserIndexer) IndexName() string {
	return "users"
}

type UserRowReader struct {
	pool *pgxpool.Pool
}

func NewUserRowReader(pool *pgxpool.Pool) *UserRowReader {
	return &UserRowReader{
		pool: pool,
	}
}

func (r *UserRowReader) ReadRows(ctx context.Context, tx chan<- UserRow) error {
	db := bun.NewDB(stdlib.OpenDBFromPool(r.pool), pgdialect.New())
	rows, err := db.NewSelect().
		ColumnExpr("u.user_id").
		ColumnExpr("u.rating").
		ColumnExpr("u.highest_rating").
		ColumnExpr("u.affiliation").
		ColumnExpr("u.birth_year").
		ColumnExpr("u.country").
		ColumnExpr("u.crown").
		ColumnExpr("u.join_count").
		ColumnExpr("u.rank").
		ColumnExpr("u.active_rank").
		ColumnExpr("u.wins").
		TableExpr("users AS u").
		Rows(ctx)

	if err != nil {
		return fmt.Errorf("execute sql: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			var row UserRow
			if err := db.ScanRow(ctx, rows, &row); err != nil {
				return fmt.Errorf("scan row: %w", err)
			}
			tx <- row
		}
	}

	return nil
}

type UserRow struct {
	UserID        string  `bun:"user_id"`
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

func (r UserRow) Document(ctx context.Context) (UserDocument, error) {
	return UserDocument{
		UserID:        r.UserID,
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
		Color:         RateToRangeLabel(&r.Rating),
		HighestColor:  RateToRangeLabel(&r.HighestRating),
		UserURL:       fmt.Sprintf("https://atcoder.jp/users/%s", r.UserID),
	}, nil
}

type UserDocument struct {
	UserID        string  `json:"userID"`
	Rating        int     `json:"rating"`
	HighestRating int     `json:"highestRating"`
	Affiliation   *string `json:"affiliation"`
	BirthYear     *int    `json:"birthYear"`
	Country       *string `json:"country"`
	Crown         *string `json:"crown"`
	JoinCount     int     `json:"joinCount"`
	Rank          int     `json:"rank"`
	ActiveRank    *int    `json:"activeRank"`
	Wins          int     `json:"wins" `
	Color         string  `json:"color"`
	HighestColor  string  `json:"highestColor"`
	UserURL       string  `json:"userURL"`
}
