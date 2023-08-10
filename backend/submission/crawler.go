package submission

import (
	"fjnkt98/atcodersearch/atcoder"
	"log"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/morikuni/failure"
)

type Crawler struct {
	client atcoder.AtCoderClient
	db     *sqlx.DB
}

func NewCrawler(client atcoder.AtCoderClient, db *sqlx.DB) Crawler {
	return Crawler{
		client: client,
		db:     db,
	}
}

func (c *Crawler) getContestIDs() ([]string, error) {
	rows, err := c.db.Query(`
	SELECT
		"contest_id"
	FROM
		"contests"
	WHERE
		"category" IN ('ABC', 'ARC', 'AGC', 'ARC-Like', 'ABC-Like')
		;
	`)
	if err != nil {
		return nil, failure.Translate(err, DBError, failure.Message("failed to get contests id from `contests` table"))
	}

	defer rows.Close()
	ids := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, failure.Translate(err, DBError, failure.Message("failed to scan row"))
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (c *Crawler) crawl(contestID string, period int64, duration int64) error {
	log.Printf("fetch submissions at page %d of the contest `%s`", 1, contestID)
	list, err := c.client.FetchSubmissionList(contestID, 1)
	if err != nil {
		return failure.Translate(err, CrawlError, failure.Context{"contestID": contestID}, failure.Message("failed to crawl submissions"))
	}

	if err := c.save(list.Submissions); err != nil {
		return failure.Translate(err, DBError, failure.Context{"contestID": contestID}, failure.Message("failed to save submissions"))
	}

	time.Sleep(time.Duration(duration) * time.Millisecond)

	for i := 2; i <= int(list.MaxPage); i++ {
		log.Printf("fetch submissions at page %d of the contest `%s`", i, contestID)
		list, err = c.client.FetchSubmissionList(contestID, uint(i))
		if err != nil {
			return failure.Translate(err, CrawlError, failure.Context{"contestID": contestID}, failure.Message("failed to crawl submissions"))
		}

		if err := c.save(list.Submissions); err != nil {
			return failure.Translate(err, DBError, failure.Context{"contestID": contestID}, failure.Message("failed to save submissions"))
		}
		time.Sleep(time.Duration(duration) * time.Millisecond)
	}

	return nil
}

func (c *Crawler) save(submissions []atcoder.Submission) error {
	tx, err := c.db.Beginx()
	if err != nil {
		return failure.Translate(err, DBError, failure.Message("failed to start transaction to save submission"))
	}
	defer tx.Rollback()

	for _, s := range submissions {
		log.Printf("save submission %v", s)
		if _, err := tx.Exec(`
			INSERT INTO "submissions" VALUES(
				$1::bigint,
				$2::bigint,
				$3::text,
				$4::text,
				$5::text,
				$6::text,
				$7::double precision,
				$8::bigint,
				$9::text,
				$10::bigint
			)
			ON CONFLICT DO NOTHING;`,
			s.ID,
			s.EpochSecond,
			s.ProblemID,
			s.ContestID,
			s.UserID,
			s.Language,
			s.Point,
			s.Length,
			s.Result,
			s.ExecutionTime,
		); err != nil {
			return failure.Translate(err, DBError, failure.Context{"contestID": s.ContestID, "id": strconv.Itoa(int(s.ID))}, failure.Message("failed to exec sql to save submission"))
		}
	}

	if err := tx.Commit(); err != nil {
		return failure.Translate(err, DBError, failure.Message("failed to commit transaction to save submissions"))
	}

	return nil
}

func (c *Crawler) Run(duration int64, period int64) error {
	ids, err := c.getContestIDs()
	if err != nil {
		return err
	}

	for _, id := range ids {
		if err != nil {
			return failure.Wrap(err)
		}
		if err := c.crawl(id, period, duration); err != nil {
			return failure.Wrap(err)
		}
	}

	return nil
}
