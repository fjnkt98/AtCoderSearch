package submission

import (
	"database/sql"
	"errors"
	"fjnkt98/atcodersearch/atcoder"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/morikuni/failure"
	"golang.org/x/exp/slog"
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

func (c *Crawler) crawl(contestID string, period int64, duration int) error {
	submissions := make([]atcoder.Submission, 0)
	maxPage := 1
loop:
	for i := 1; i <= maxPage; i++ {
		slog.Info(fmt.Sprintf("fetch submissions at page %d / %d of the contest `%s`", i, maxPage, contestID))
		list, err := c.client.FetchSubmissionList(contestID, int(i))
		if err != nil {
			return failure.Translate(err, CrawlError, failure.Context{"contestID": contestID}, failure.Message("failed to crawl submissions"))
		}

		if list.Submissions[0].EpochSecond < period {
			slog.Info(fmt.Sprintf("All submissions after page `%d` have been crawled. Break crawling the contest `%s`", i, contestID))
			break loop
		}

		submissions = append(submissions, list.Submissions...)
		maxPage = int(list.MaxPage)
		time.Sleep(time.Duration(duration) * time.Millisecond)
	}

	if len(submissions) == 0 {
		slog.Info(fmt.Sprintf("No submissions to save for contest `%s`.", contestID))
		return nil
	}

	tx, err := c.db.Beginx()
	if err != nil {
		return failure.Translate(err, DBError, failure.Message("failed to start transaction to save submission"))
	}
	defer tx.Rollback()

	sql := `
	INSERT INTO "submissions" (
		"id",
		"epoch_second",
		"problem_id",
		"contest_id",
		"user_id",
		"language",
		"point",
		"length",
		"result",
		"execution_time",
		"crawled_at"
	)
	VALUES(
		:id,
		:epoch_second,
		:problem_id,
		:contest_id,
		:user_id,
		:language,
		:point,
		:length,
		:result,
		:execution_time,
		NOW()	
	)
	ON CONFLICT ("id") DO UPDATE SET
		"id" = EXCLUDED."id",
		"epoch_second" = EXCLUDED."epoch_second",
		"problem_id" = EXCLUDED."problem_id",
		"contest_id" = EXCLUDED."contest_id",
		"user_id" = EXCLUDED."user_id",
		"language" = EXCLUDED."language",
		"point" = EXCLUDED."point",
		"length" = EXCLUDED."length",
		"result" = EXCLUDED."result",
		"execution_time" = EXCLUDED."execution_time",
		"crawled_at" = NOW()
	;`
	affected := 0
	if result, err := tx.NamedExec(sql, submissions); err != nil {
		return failure.Translate(err, DBError, failure.Context{"contestID": contestID}, failure.Message("failed to exec sql to save submission"))
	} else {
		a, _ := result.RowsAffected()
		affected += int(a)
	}

	if err := tx.Commit(); err != nil {
		return failure.Translate(err, DBError, failure.Message("failed to commit transaction to save submissions"))
	} else {
		slog.Info(fmt.Sprintf("commit transaction. save %d rows.", affected))
	}

	return nil
}

func (c *Crawler) Run(duration int) error {
	ids, err := c.getContestIDs()
	if err != nil {
		return err
	}

	for _, id := range ids {
		history := NewCrawlHistory(c.db, id)
		period, err := history.GetLatestHistory()
		if err != nil {
			return failure.Wrap(err)
		}
		slog.Info(fmt.Sprintf("Start to crawl contest `%s` since period `%s`", id, time.Unix(int64(period), 0)))
		if err := c.crawl(id, int64(period), duration); err != nil {
			return failure.Wrap(err)
		}
		history.Finish()
	}

	return nil
}

type CrawlHistory struct {
	db        *sqlx.DB
	StartedAt int
	ContestID string
}

func NewCrawlHistory(db *sqlx.DB, contestID string) CrawlHistory {
	return CrawlHistory{
		db:        db,
		StartedAt: int(time.Now().Unix()),
		ContestID: contestID,
	}
}

func (h *CrawlHistory) GetLatestHistory() (int, error) {
	rows, err := h.db.Query(
		`SELECT "started_at" FROM "submission_crawl_history" WHERE "contest_id" = $1::text ORDER BY "started_at" DESC LIMIT 1;`,
		h.ContestID,
	)
	if err != nil {
		return 0, failure.Translate(err, DBError, failure.Message("failed to get latest crawl history"))
	}

	defer rows.Close()
	var startedAt int
	for rows.Next() {
		if err := rows.Scan(&startedAt); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				slog.Info(fmt.Sprintf("`submission_crawl_history` table is empty in term of contest id `%s`", h.ContestID))
				return 0, nil
			} else {
				return 0, failure.Translate(err, DBError, failure.Message("failed to get latest crawl history"))
			}
		}
	}

	return startedAt, nil
}

func (h *CrawlHistory) Finish() error {
	tx, err := h.db.Beginx()
	if err != nil {
		return failure.Translate(err, DBError, failure.Message("failed to start transaction to save submission crawl history"))
	}
	defer tx.Rollback()

	if _, err := tx.Exec(
		`INSERT INTO "submission_crawl_history" ("contest_id", "started_at") VALUES ($1::text, $2::bigint);`,
		h.ContestID,
		h.StartedAt,
	); err != nil {
		return failure.Translate(err, DBError, failure.Message("failed to exec sql to save submission crawl history"))
	}

	if err := tx.Commit(); err != nil {
		return failure.Translate(err, DBError, failure.Message("failed to commit transaction to save submission crawl history"))
	}

	return nil
}
