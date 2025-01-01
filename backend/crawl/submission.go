package crawl

import (
	"context"
	"database/sql"
	"errors"
	"fjnkt98/atcodersearch/atcoder"
	"fjnkt98/atcodersearch/repository"
	"fmt"
	"log/slog"
	"slices"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

type SubmissionCrawler struct {
	client        atcoder.AtCoderClient
	pool          *pgxpool.Pool
	duration      time.Duration
	retry         int
	retryInterval time.Duration
	targets       []string
	endless       bool
}

func NewSubmissionCrawler(
	client atcoder.AtCoderClient,
	pool *pgxpool.Pool,
	duration time.Duration,
	retry int,
	retryInterval time.Duration,
	targets []string,
	endless bool,
) *SubmissionCrawler {
	return &SubmissionCrawler{
		client:        client,
		pool:          pool,
		duration:      duration,
		retry:         retry,
		retryInterval: retryInterval,
		targets:       targets,
		endless:       endless,
	}
}

func (c *SubmissionCrawler) Crawl(ctx context.Context) error {
	if c.endless {
		for {
			if err := c.crawl(ctx); err != nil {
				return err
			}
		}
	} else {
		if err := c.crawl(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (c *SubmissionCrawler) crawl(ctx context.Context) error {
	contests, err := FetchContestIDs(ctx, c.pool, c.targets)
	if err != nil {
		return fmt.Errorf("fetch contest ids: %w", err)
	}

	for _, contestID := range contests {
		if err := c.crawlContest(ctx, contestID); err != nil {
			return fmt.Errorf("crawl contest %s: %w", contestID, err)
		}

		time.Sleep(c.duration)
	}

	return nil
}

func (c *SubmissionCrawler) crawlContest(ctx context.Context, contestID string) (err error) {
	latest, err := repository.FetchLatestCrawlHistory(ctx, c.pool, contestID)
	if err != nil {
		return fmt.Errorf("fetch latest crawl history: %w", err)
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "start to crawl submissions", slog.String("contestID", contestID), slog.Time("lastCrawled", latest.StartedAt))

	history, err := repository.NewCrawlHistory(ctx, c.pool, contestID)
	if err != nil {
		return fmt.Errorf("create crawl history: %w", err)
	}
	defer func() {
		if historyErr := history.Abort(ctx, c.pool); historyErr != nil {
			if !errors.Is(historyErr, repository.ErrHistoryConfirmed) {
				err = errors.Join(
					err,
					historyErr,
				)
			}
		}
	}()

	saver, err := NewSubmissionSaver(ctx, c.pool)
	if err != nil {
		return fmt.Errorf("new submission saver: %w", err)
	}
	defer saver.Close()

	queue := []int{1}
	remain := c.retry

loop:
	for len(queue) > 0 {
		page := queue[0]
		queue = queue[1:]

		slog.LogAttrs(ctx, slog.LevelInfo, "fetch submissions", slog.String("contestID", contestID), slog.Int("page", page))
		submissions, err := c.client.FetchSubmissions(ctx, contestID, page)
		if err != nil {
			if errors.Is(err, atcoder.ErrNotFound) {
				slog.LogAttrs(ctx, slog.LevelWarn, "the submission page of the contest not found. break crawling.", slog.String("contestID", contestID))
				break loop
			}
			if remain <= 0 {
				return fmt.Errorf("fetch submissions: %w", err)
			}

			slog.LogAttrs(ctx, slog.LevelWarn, "failed to crawl submissions. retry to crawl...", slog.String("contestID", contestID), slog.Int("page", page), slog.Any("error", err))
			queue = append(queue, page)
			remain -= 1
			time.Sleep(c.retryInterval)
			continue loop
		}

		if len(submissions) == 0 {
			slog.LogAttrs(ctx, slog.LevelInfo, "there is no submissions", slog.String("contestID", contestID), slog.Int("page", page))
			break loop
		}

		if err := saver.Append(ctx, submissions); err != nil {
			return fmt.Errorf("append submissions: %w", err)
		}

		if submissions[0].EpochSecond < latest.StartedAt.Unix()-86400 {
			slog.LogAttrs(
				ctx,
				slog.LevelInfo,
				"all submissions after here have been crawled",
				slog.String("contestID", contestID),
				slog.Int("page", page),
				slog.Int64("id", submissions[0].ID),
				slog.Time("epochSecond", time.Unix(submissions[0].EpochSecond, 0)),
			)
			time.Sleep(c.duration)
			break loop
		}

		queue = append(queue, page+1)
		remain = c.retry

		time.Sleep(c.duration)
	}

	count, err := saver.Save(ctx)
	if err != nil {
		return fmt.Errorf("save submissions: %w", err)
	}
	if err := history.Complete(ctx, c.pool); err != nil {
		return fmt.Errorf("complete crawl history: %w", err)
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "finish to crawl submissions", slog.String("contestID", contestID), slog.Int64("count", count))

	return nil
}

func FetchContestIDs(ctx context.Context, pool *pgxpool.Pool, category []string) ([]string, error) {
	db := bun.NewDB(stdlib.OpenDBFromPool(pool), pgdialect.New())

	q := db.NewSelect().Table("contests").Column("contest_id").Order("start_epoch_second DESC")

	if len(category) > 0 {
		q = q.Where("category IN (?)", bun.In(category))
	}

	var result []string
	if err := q.Scan(ctx, &result); err != nil {
		return nil, fmt.Errorf("exec query: %w", err)
	}

	return result, nil
}

type Submission struct {
	bun.BaseModel `bun:"table:submissions,alias:s"`
	ID            int64    `bun:"id"`
	EpochSecond   int64    `bun:"epoch_second"`
	ProblemID     string   `bun:"problem_id"`
	ContestID     *string  `bun:"contest_id"`
	UserID        *string  `bun:"user_id"`
	Language      *string  `bun:"language"`
	Point         *float64 `bun:"point"`
	Length        *int32   `bun:"length"`
	Result        *string  `bun:"result"`
	ExecutionTime *int32   `bun:"execution_time"`
}

func NewSubmissions(submissions []atcoder.Submission) []Submission {
	result := make([]Submission, 0, len(submissions))

	set := mapset.NewSet[int64]()

	for _, s := range submissions {
		s := s
		if set.Contains(s.ID) {
			continue
		}

		result = append(result, Submission{
			ID:            s.ID,
			EpochSecond:   s.EpochSecond,
			ProblemID:     s.ProblemID,
			ContestID:     &s.ContestID,
			UserID:        &s.UserID,
			Language:      &s.Language,
			Point:         &s.Point,
			Length:        &s.Length,
			Result:        &s.Result,
			ExecutionTime: s.ExecutionTime,
		})

		set.Add(s.ID)
	}

	return result
}

func NewSubmission(s atcoder.Submission) Submission {
	return Submission{
		ID:            s.ID,
		EpochSecond:   s.EpochSecond,
		ProblemID:     s.ProblemID,
		ContestID:     &s.ContestID,
		UserID:        &s.UserID,
		Language:      &s.Language,
		Point:         &s.Point,
		Length:        &s.Length,
		Result:        &s.Result,
		ExecutionTime: s.ExecutionTime,
	}
}

type SubmissionSaver struct {
	tx    bun.Tx
	set   mapset.Set[int64]
	count int64
}

func NewSubmissionSaver(ctx context.Context, pool *pgxpool.Pool) (*SubmissionSaver, error) {
	tx, err := bun.NewDB(stdlib.OpenDBFromPool(pool), pgdialect.New()).BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}

	return &SubmissionSaver{
		tx:    tx,
		set:   mapset.NewSet[int64](),
		count: 0,
	}, nil
}

func (s *SubmissionSaver) Append(ctx context.Context, submissions []atcoder.Submission) error {
	data := make([]Submission, 0, len(submissions))
	for _, submission := range submissions {
		if s.set.Contains(submission.ID) {
			continue
		}

		data = append(data, NewSubmission(submission))
		s.set.Add(submission.ID)
	}

	for chunk := range slices.Chunk(data, 1000) {
		res, err := s.tx.NewInsert().
			Model(&chunk).
			On("CONFLICT (epoch_second, id) DO UPDATE").
			Set("id = EXCLUDED.id").
			Set("epoch_second = EXCLUDED.epoch_second").
			Set("problem_id = EXCLUDED.problem_id").
			Set("contest_id = EXCLUDED.contest_id").
			Set("user_id = EXCLUDED.user_id").
			Set("language = EXCLUDED.language").
			Set("point = EXCLUDED.point").
			Set("length = EXCLUDED.length").
			Set("result = EXCLUDED.result").
			Set("execution_time = EXCLUDED.execution_time").
			Set("updated_at = NOW()").
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("insert: %w", err)
		}

		if c, err := res.RowsAffected(); err != nil {
			return fmt.Errorf("get affected rows count: %w", err)
		} else {
			s.count += c
		}
	}

	return nil
}

func (s *SubmissionSaver) Save(ctx context.Context) (int64, error) {
	if err := s.tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit transaction: %w", err)
	}
	return s.count, nil
}

func (s *SubmissionSaver) Close() error {
	return s.tx.Rollback()
}
