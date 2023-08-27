package user

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/morikuni/failure"
	"golang.org/x/exp/slog"
)

type UserCrawler struct {
	targetURL string
	db        *sqlx.DB
	client    *http.Client
}

func NewUserCrawler(db *sqlx.DB) UserCrawler {
	return UserCrawler{
		targetURL: "https://atcoder.jp/ranking/all",
		db:        db,
		client:    &http.Client{},
	}
}

func (c *UserCrawler) Crawl(index int) ([]User, error) {
	u, _ := url.Parse(c.targetURL)
	v := url.Values{}
	v.Set("contestType", "algo")
	v.Set("page", strconv.Itoa(index))
	u.RawQuery = v.Encode()

	slog.Info(fmt.Sprintf("Crawling active user ranking page %s", u.String()))
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, failure.Translate(err, RequestCreationError, failure.Context{"url": u.String()}, failure.Message("failed to create request"))
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, failure.Translate(err, RequestExecutionError, failure.Context{"url": u.String()}, failure.Message("request failed"))
	}

	defer res.Body.Close()
	users, err := Scrape(res.Body)
	if err != nil {
		return nil, failure.Translate(err, ScrapeError, failure.Context{"url": u.String()}, failure.Message("failed to scrape active user information"))
	}
	return users, nil
}

func (c *UserCrawler) Save(users []User) error {
	tx, err := c.db.Beginx()
	if err != nil {
		return failure.Translate(err, DBError, failure.Message("failed to start transaction to save user information"))
	}
	defer tx.Rollback()

	for _, user := range users {
		_, err := tx.Exec(`
			MERGE INTO "users"
			USING
				(
					VALUES (
						$1::text,
						$2::integer,
						$3::integer,
						$4::text,
						$5::integer,
						$6::text,
						$7::text,
						$8::integer,
						$9::integer,
						$10::integer,
						$11::integer
					)
				) AS "user" (
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
					"wins"
				)
			ON
				"users"."user_name" = "user"."user_name"
			WHEN MATCHED THEN
				UPDATE SET (
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
					"updated_at"
				) = (
					"user"."rating",
					"user"."highest_rating",
					"user"."affiliation",
					"user"."birth_year",
					"user"."country",
					"user"."crown",
					"user"."join_count",
					"user"."rank",
					"user"."active_rank",
					"user"."wins",
					NOW()
				)
			WHEN NOT MATCHED THEN
				INSERT (
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
					"created_at",
					"updated_at"
				)
				VALUES (
					"user"."user_name",
					"user"."rating",
					"user"."highest_rating",
					"user"."affiliation",
					"user"."birth_year",
					"user"."country",
					"user"."crown",
					"user"."join_count",
					"user"."rank",
					"user"."active_rank",
					"user"."wins",
					NOW(),
					NOW()
				);	
		`,
			user.UserName,
			user.Rating,
			user.HighestRating,
			user.Affiliation,
			user.BirthYear,
			user.Country,
			user.Crown,
			user.JoinCount,
			user.Rank,
			user.ActiveRank,
			user.Wins,
		)
		if err != nil {
			return failure.Translate(err, DBError, failure.Context{"user": user.UserName}, failure.Message("failed to save user information"))
		}
	}
	if err := tx.Commit(); err != nil {
		return failure.Translate(err, DBError, failure.Message("failed to commit transaction to save user information"))
	}

	return nil
}

func (c *UserCrawler) Run(duration int) error {
	slog.Info("Start to crawl active user information.")

loop:
	for i := 0; ; i++ {
		users, err := c.Crawl(i)
		if err != nil {
			return failure.Wrap(err)
		}

		if len(users) == 0 {
			break loop
		}

		if err := c.Save(users); err != nil {
			return failure.Wrap(err)
		}

		time.Sleep(time.Duration(duration) * time.Millisecond)
	}

	return nil
}

type User struct {
	UserName      string  `json:"user_name" db:"user_name"`
	Rating        int     `json:"rating" db:"rating"`
	HighestRating int     `json:"highest_rating" db:"highest_rating"`
	Affiliation   *string `json:"affiliation" db:"affiliation"`
	BirthYear     *int    `json:"birth_year" db:"birth_year"`
	Country       *string `json:"country" db:"country"`
	Crown         *string `json:"crown" db:"crown"`
	JoinCount     int     `json:"join_count" db:"join_count"`
	Rank          int     `json:"rank" db:"rank"`
	ActiveRank    *int    `json:"active_rank" db:"active_rank"`
	Wins          int     `json:"wins" db:"wins"`
}
