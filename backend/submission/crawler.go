package submission

import (
	"context"
	"database/sql"
	"errors"
	"fjnkt98/atcodersearch/acs"
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

func (c *Crawler) getContestIDs(ctx context.Context) ([]string, error) {
	rows, err := c.db.QueryContext(ctx, `
	SELECT
		"contest_id"
	FROM
		"contests"
	ORDER BY
		"start_epoch_second" DESC
	`)
	if err != nil {
		return nil, failure.Translate(err, acs.DBError, failure.Message("failed to get contests id from `contests` table"))
	}

	defer rows.Close()
	ids := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, failure.Translate(err, acs.DBError, failure.Message("failed to scan row"))
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (c *Crawler) crawl(ctx context.Context, contestID string, period int64, duration int, retry int) error {
	submissions := make([]atcoder.Submission, 0)
	maxPage := 1
loop:
	for i := 1; i <= maxPage; i++ {
		slog.Info(fmt.Sprintf("fetch submissions at page %d / %d of the contest `%s`", i, maxPage, contestID))
		list, err := c.client.FetchSubmissionList(ctx, contestID, i)
		if err != nil {
		retryLoop:
			for j := 0; j < retry; j++ {
				select {
				case <-ctx.Done():
					return failure.New(acs.Interrupt, failure.Message("retry to crawl submission has been canceled"))
				default:
					slog.Error("failed to crawl submission", slog.String("contestID", contestID), slog.String("error", fmt.Sprintf("%+v", err)))
					slog.Info("retry to crawl submission after 1 minutes...")
					time.Sleep(time.Duration(60) * time.Second)
					list, err = c.client.FetchSubmissionList(ctx, contestID, i)
					if err == nil {
						break retryLoop
					}
				}
			}

			if err != nil {
				return failure.Translate(err, acs.CrawlError, failure.Context{"contestID": contestID}, failure.Message("failed to crawl submissions"))
			}
		}

		submissions = append(submissions, list.Submissions...)
		if list.Submissions[0].EpochSecond < period {
			slog.Info(fmt.Sprintf("All submissions after page `%d` have been crawled. Break crawling the contest `%s`", i, contestID))
			time.Sleep(time.Duration(duration) * time.Millisecond)
			break loop
		}
		maxPage = int(list.MaxPage)
		time.Sleep(time.Duration(duration) * time.Millisecond)
	}

	if len(submissions) == 0 {
		slog.Info(fmt.Sprintf("No submissions to save for contest `%s`.", contestID))
		return nil
	}

	tx, err := c.db.Beginx()
	if err != nil {
		return failure.Translate(err, acs.DBError, failure.Message("failed to start transaction to save submission"))
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
	ON CONFLICT DO NOTHING;
	`
	affected := 0
	for _, submission := range submissions {
		if result, err := tx.NamedExecContext(ctx, sql, submission); err != nil {
			return failure.Translate(err, acs.DBError, failure.Context{"contestID": contestID}, failure.Message("failed to exec sql to save submission"))
		} else {
			a, _ := result.RowsAffected()
			affected += int(a)
		}
	}

	if err := tx.Commit(); err != nil {
		return failure.Translate(err, acs.DBError, failure.Message("failed to commit transaction to save submissions"))
	} else {
		slog.Info(fmt.Sprintf("commit transaction. save %d rows.", affected))
	}

	return nil
}

func (c *Crawler) Run(ctx context.Context, duration int, retry int) error {
	ids, err := c.getContestIDs(ctx)
	if err != nil {
		return err
	}

	for _, id := range ids {
		history := NewCrawlHistory(c.db, id)
		period, err := history.GetLatestHistory(ctx)
		if err != nil {
			return failure.Wrap(err)
		}
		slog.Info(fmt.Sprintf("Start to crawl contest `%s` since period `%s`", id, time.Unix(int64(period), 0)))
		if err := c.crawl(ctx, id, int64(period), duration, retry); err != nil {
			return failure.Wrap(err)
		}
		history.Finish(ctx)
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

func (h *CrawlHistory) GetLatestHistory(ctx context.Context) (int, error) {
	rows, err := h.db.QueryContext(
		ctx,
		`SELECT "started_at" FROM "submission_crawl_history" WHERE "contest_id" = $1::text ORDER BY "started_at" DESC LIMIT 1;`,
		h.ContestID,
	)
	if err != nil {
		return 0, failure.Translate(err, acs.DBError, failure.Message("failed to get latest crawl history"))
	}

	defer rows.Close()
	var startedAt int
	for rows.Next() {
		if err := rows.Scan(&startedAt); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				slog.Info(fmt.Sprintf("`submission_crawl_history` table is empty in term of contest id `%s`", h.ContestID))
				return 0, nil
			} else {
				return 0, failure.Translate(err, acs.DBError, failure.Message("failed to get latest crawl history"))
			}
		}
	}

	return startedAt, nil
}

func (h *CrawlHistory) Finish(ctx context.Context) error {
	tx, err := h.db.Beginx()
	if err != nil {
		return failure.Translate(err, acs.DBError, failure.Message("failed to start transaction to save submission crawl history"))
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO "submission_crawl_history" ("contest_id", "started_at") VALUES ($1::text, $2::bigint);`,
		h.ContestID,
		h.StartedAt,
	); err != nil {
		return failure.Translate(err, acs.DBError, failure.Message("failed to exec sql to save submission crawl history"))
	}

	if err := tx.Commit(); err != nil {
		return failure.Translate(err, acs.DBError, failure.Message("failed to commit transaction to save submission crawl history"))
	}

	return nil
}
