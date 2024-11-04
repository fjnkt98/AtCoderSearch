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
			"userId",
			"rating",
			"ratingFacet",
			"highestRating",
			"affiliation",
			"birthYear",
			"birthYearFacet",
			"country",
			"crown",
			"joinCount",
			"joinCountFacet",
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
			"userId",
			"rank",
			"submissionCount",
			"accepted",
		},
		Synonyms: map[string][]string{},
	}
}

func (ix *UserIndexer) PrimaryKey() string {
	return "userId"
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
	submissionCounts := db.NewSelect().
		ColumnExpr("user_id").
		ColumnExpr("COUNT(*) AS submission_count").
		TableExpr("submissions").
		GroupExpr("user_id")

	accepts := db.NewSelect().
		ColumnExpr("user_id").
		ColumnExpr("COUNT(DISTINCT problem_id) AS accepted").
		TableExpr("submissions").
		Where("result = ?", "AC").
		GroupExpr("user_id")

	rows, err := db.NewSelect().
		With("sc", submissionCounts).
		With("ac", accepts).
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
		ColumnExpr("sc.submission_count").
		ColumnExpr("ac.accepted").
		TableExpr("users AS u").
		Join("LEFT JOIN sc ON u.user_id = sc.user_id").
		Join("LEFT JOIN ac ON u.user_id = ac.user_id").
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
	UserID          string  `bun:"user_id"`
	Rating          int     `bun:"rating"`
	HighestRating   int     `bun:"highest_rating"`
	Affiliation     *string `bun:"affiliation"`
	BirthYear       *int    `bun:"birth_year"`
	Country         *string `bun:"country"`
	Crown           *string `bun:"crown"`
	JoinCount       int     `bun:"join_count"`
	Rank            int     `bun:"rank"`
	ActiveRank      *int    `bun:"active_rank"`
	Wins            int     `bun:"wins"`
	SubmissionCount int     `bun:"submission_count"`
	Accepted        int     `bun:"accepted"`
}

func (r UserRow) Document(ctx context.Context) (UserDocument, error) {
	var ratingFacet string
	var birthYearFacet string
	var joinCountFacet string

	if rating := (r.Rating / 400) * 400; rating < 4000 {
		ratingFacet = fmt.Sprintf("%4d ~ %4d", rating, rating+400)
	} else {
		ratingFacet = "4000 ~     "
	}

	if r.BirthYear != nil {
		if birth := (*r.BirthYear / 10) * 10; birth < 1970 {
			birthYearFacet = "     ~ 1970"
		} else if birth >= 2020 {
			birthYearFacet = "2020 ~     "
		} else {
			birthYearFacet = fmt.Sprintf("%4d ~ %4d", birth, birth+10)
		}
	}

	if join := (r.JoinCount / 20) * 20; join < 100 {
		joinCountFacet = fmt.Sprintf("%4d ~ %4d", join, join+20)
	} else {
		joinCountFacet = " 100 ~     "
	}

	return UserDocument{
		UserID:          r.UserID,
		Rating:          r.Rating,
		RatingFacet:     ratingFacet,
		HighestRating:   r.HighestRating,
		Affiliation:     r.Affiliation,
		BirthYear:       r.BirthYear,
		BirthYearFacet:  birthYearFacet,
		Country:         r.Country,
		Crown:           r.Crown,
		JoinCount:       r.JoinCount,
		JoinCountFacet:  joinCountFacet,
		Rank:            r.Rank,
		ActiveRank:      r.ActiveRank,
		Wins:            r.Wins,
		UserURL:         fmt.Sprintf("https://atcoder.jp/users/%s", r.UserID),
		SubmissionCount: r.SubmissionCount,
		Accepted:        r.Accepted,
	}, nil
}

type UserDocument struct {
	UserID          string  `json:"userId"`
	Rating          int     `json:"rating"`
	RatingFacet     string  `json:"ratingFacet,omitempty"`
	HighestRating   int     `json:"highestRating"`
	Affiliation     *string `json:"affiliation"`
	BirthYear       *int    `json:"birthYear"`
	BirthYearFacet  string  `json:"birthYearFacet,omitempty"`
	Country         *string `json:"country"`
	Crown           *string `json:"crown"`
	JoinCount       int     `json:"joinCount"`
	JoinCountFacet  string  `json:"joinCountFacet,omitempty"`
	Rank            int     `json:"rank"`
	ActiveRank      *int    `json:"activeRank"`
	Wins            int     `json:"wins" `
	UserURL         string  `json:"userUrl"`
	Accepted        int     `json:"accepted"`
	SubmissionCount int     `json:"submissionCount"`
}
