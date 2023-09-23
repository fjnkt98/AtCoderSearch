package submission

import (
	"context"
	"database/sql"
	"errors"
	"fjnkt98/atcodersearch/acs"
	"fjnkt98/atcodersearch/atcoder"
	"fmt"
	"strconv"
	"strings"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
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

func (c *Crawler) getContestIDs(ctx context.Context, targets []string) ([]string, error) {
	var rows *sqlx.Rows
	var err error
	if len(targets) == 0 {
		rows, err = c.db.QueryxContext(
			ctx,
			`SELECT
				"contest_id"
			FROM
				"contests"
			ORDER BY
				"start_epoch_second" DESC`,
		)
		if err != nil {
			return nil, failure.Translate(err, acs.DBError, failure.Message("failed to get contests id from `contests` table"))
		}
	} else {
		sql, args, err := sqlx.In(`
			SELECT
				"contest_id"
			FROM
				"contests"
			WHERE
				"category" IN (?)
			ORDER BY
				"start_epoch_second" DESC
		`, targets)
		if err != nil {
			return nil, failure.Translate(err, acs.DBError, failure.Context{"targets": strings.Join(targets, ",")}, failure.Message("failed to build sql query"))
		}
		sql = c.db.Rebind(sql)
		rows, err = c.db.QueryxContext(ctx, sql, args...)
		if err != nil {
			return nil, failure.Translate(err, acs.DBError, failure.Context{"sql": sql}, failure.Message("failed to get contests id from `contests` table"))
		}
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
		if len(list.Submissions) == 0 {
			slog.Info(fmt.Sprintf("There is no submissions in contest `%s`.", contestID))
			break loop
		}
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

	noDupSubmissions := make([]atcoder.Submission, 0, len(submissions))
	ids := mapset.NewSet[int64]()
	for _, s := range submissions {
		if ids.Contains(s.ID) {
			continue
		}
		ids.Add(s.ID)
		noDupSubmissions = append(noDupSubmissions, s)
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
	for _, submission := range noDupSubmissions {
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

func (c *Crawler) Run(ctx context.Context, targets []string, duration int, retry int) error {
	ids, err := c.getContestIDs(ctx, targets)
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

type Target struct {
	ID        int64  `db:"id"`
	ContestID string `db:"contest_id"`
	Result    string `db:"result"`
}

func (c *Crawler) ReCrawl(ctx context.Context, duration int) error {
	rows, err := c.db.QueryxContext(
		ctx,
		`
		SELECT
			"id",
			"contest_id",
			"result"
		FROM
			"submissions"
		WHERE
			"result" NOT IN ('AC', 'CE', 'IE', 'MLE', 'NG', 'OLE', 'QLE', 'RE', 'TLE', 'WA', 'WJ', 'WR')
		`,
	)
	if err != nil {
		return failure.Translate(err, acs.DBError, failure.Message("failed to get submissions `submissions` table"))
	}

	update := func(ctx context.Context, db *sqlx.DB, target Target) error {
		tx, err := db.Beginx()
		if err != nil {
			return failure.Translate(err, acs.DBError, failure.Message("failed to start transaction to save submission"))
		}
		defer tx.Rollback()

		sql := `
		UPDATE
			"submissions"
		SET
			"result" = :result
		WHERE
			"id" = :id
			AND "contest_id" = :contest_id
		`

		if _, err := tx.NamedExecContext(ctx, sql, target); err != nil {
			return failure.Translate(err, acs.DBError, failure.Context{"submissionID": strconv.Itoa(int(target.ID)), "contestID": target.ContestID}, failure.Message("failed to exec sql to update submission result"))
		}

		if err := tx.Commit(); err != nil {
			return failure.Translate(err, acs.DBError, failure.Context{"submissionID": strconv.Itoa(int(target.ID)), "contestID": target.ContestID}, failure.Message("failed to commit transaction to update submission result"))
		}
		slog.Info(fmt.Sprintf("commit transaction updating submission `%d` result.", target.ID))

		return nil
	}

	defer rows.Close()
	for rows.Next() {
		var t Target
		if err := rows.StructScan(&t); err != nil {
			return failure.Translate(err, acs.DBError, failure.Message("failed to scan row"))
		}

		var result string
		if result, err = c.client.FetchSubmissionResult(ctx, t.ContestID, t.ID); err != nil {
			return failure.Translate(err, acs.CrawlError, failure.Context{"contestID": t.ContestID, "submissionID": strconv.Itoa(int(t.ID))}, failure.Message("failed to crawl submission"))
		}
		slog.Info(fmt.Sprintf("update submission result of `%d` in contest `%s` from `%s` into `%s`", t.ID, t.ContestID, t.Result, result))
		t.Result = result

		if err := update(ctx, c.db, t); err != nil {
			return failure.Wrap(err)
		}

		time.Sleep(time.Duration(duration) * time.Millisecond)
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
