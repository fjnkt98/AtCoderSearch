package crawl

import (
	"bytes"
	"context"
	"database/sql"
	"fjnkt98/atcodersearch/atcoder"
	"fjnkt98/atcodersearch/repository"
	"fmt"
	"log/slog"
	"maps"
	"slices"
	"strings"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

type ProblemCrawler struct {
	atcoderClient  atcoder.AtCoderClient
	problemsClient atcoder.AtCoderProblemsClient
	pool           *pgxpool.Pool
	duration       time.Duration
	all            bool
	minifier       *minify.M
}

func NewProblemCrawler(
	atcoderClient atcoder.AtCoderClient,
	problemsClient atcoder.AtCoderProblemsClient,
	pool *pgxpool.Pool,
	duration time.Duration,
	all bool,
) *ProblemCrawler {
	m := minify.New()
	m.AddFunc("text/html", html.Minify)

	return &ProblemCrawler{
		atcoderClient:  atcoderClient,
		problemsClient: problemsClient,
		pool:           pool,
		minifier:       m,
		duration:       duration,
		all:            all,
	}
}

func (c *ProblemCrawler) CrawlContests(ctx context.Context) error {
	slog.LogAttrs(ctx, slog.LevelInfo, "start to crawl contests.")

	contests, err := c.problemsClient.FetchContests(ctx)
	if err != nil {
		return fmt.Errorf("fetch contests: %w", err)
	}

	count, err := SaveContests(ctx, c.pool, contests)
	if err != nil {
		return fmt.Errorf("save contests: %w", err)
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "finish to crawl contests.", slog.Int64("count", count))
	return nil
}

func (c *ProblemCrawler) CrawlDifficulties(ctx context.Context) error {
	slog.LogAttrs(ctx, slog.LevelInfo, "start to crawl difficulties.")

	difficulties, err := c.problemsClient.FetchDifficulties(ctx)
	if err != nil {
		return fmt.Errorf("fetch difficulties: %w", err)
	}

	count, err := SaveDifficulties(ctx, c.pool, difficulties)
	if err != nil {
		return fmt.Errorf("save difficulties: %w", err)
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "finish to crawl difficulties.", slog.Int64("count", count))
	return nil
}

func (c *ProblemCrawler) CrawlProblems(ctx context.Context) error {
	targets, err := c.problemsClient.FetchProblems(ctx)
	if err != nil {
		return fmt.Errorf("fetch problems: %w", err)
	}

	if !c.all {
		targets, err = DetectDiff(ctx, c.pool, targets)
		if err != nil {
			return fmt.Errorf("detect diff: %w", err)
		}
	}
	slog.LogAttrs(ctx, slog.LevelInfo, "start to crawl problems.", slog.Int("targets", len(targets)))

	save := func(ctx context.Context, p atcoder.Problem) (int64, error) {
		html, err := c.atcoderClient.FetchProblemHTML(ctx, p.ContestID, p.ID)
		if err != nil {
			return 0, fmt.Errorf("fetch problem html: %w", err)
		}

		tx, err := c.pool.Begin(ctx)
		if err != nil {
			return 0, fmt.Errorf("begin transaction: %w", err)
		}
		defer tx.Rollback(ctx)

		q := repository.New(c.pool).WithTx(tx)

		var buf bytes.Buffer
		if err := c.minifier.Minify("text/html", &buf, strings.NewReader(html)); err != nil {
			return 0, fmt.Errorf("minify html: %w", err)
		}

		res, err := q.InsertProblem(ctx, repository.InsertProblemParams{
			ProblemID:    p.ID,
			ContestID:    p.ContestID,
			ProblemIndex: p.ProblemIndex,
			Name:         p.Name,
			Title:        p.Title,
			Url:          p.URL(),
			Html:         buf.String(),
		})
		if err != nil {
			return 0, fmt.Errorf("insert problem: %w", err)
		}

		if err := tx.Commit(ctx); err != nil {
			return 0, fmt.Errorf("commit transaction: %w", err)
		}

		return res.RowsAffected(), nil
	}

	var count int64 = 0
	for _, p := range targets {
		if c, err := save(ctx, p); err != nil {
			return err
		} else {
			count += c
		}
		slog.LogAttrs(ctx, slog.LevelInfo, "save problem", slog.String("problemID", p.ID))
		time.Sleep(c.duration)
	}

	slog.LogAttrs(ctx, slog.LevelInfo, "finish to crawl problems.", slog.Int64("count", count))

	return nil
}

func DetectDiff(ctx context.Context, pool *pgxpool.Pool, problems []atcoder.Problem) ([]atcoder.Problem, error) {
	q := repository.New(pool)
	ids, err := q.FetchProblemIDs(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch problem ids: %w", err)
	}

	exists := mapset.NewSet(ids...)

	diff := make([]atcoder.Problem, 0, len(problems))
	for _, p := range problems {
		if !exists.Contains(p.ID) {
			diff = append(diff, p)
		}
	}

	return diff, nil
}

type Contest struct {
	bun.BaseModel    `bun:"table:contests,alias:c"`
	ContestID        string `bun:"contest_id"`
	StartEpochSecond int64  `bun:"start_epoch_second"`
	DurationSecond   int64  `bun:"duration_second"`
	Title            string `bun:"title"`
	RateChange       string `bun:"rate_change"`
	Category         string `bun:"category"`
}

func NewContests(contests []atcoder.Contest) []Contest {
	result := make([]Contest, len(contests))

	for i, c := range contests {
		result[i] = Contest{
			ContestID:        c.ID,
			StartEpochSecond: c.StartEpochSecond,
			DurationSecond:   c.DurationSecond,
			Title:            c.Title,
			RateChange:       c.RateChange,
			Category:         c.Categorize(),
		}
	}

	return result
}

func SaveContests(ctx context.Context, pool *pgxpool.Pool, contests []atcoder.Contest) (int64, error) {
	if len(contests) == 0 {
		return 0, nil
	}

	tx, err := bun.NewDB(stdlib.OpenDBFromPool(pool), pgdialect.New()).BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return 0, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	var count int64 = 0

	for chunk := range slices.Chunk(NewContests(contests), 1000) {
		res, err := tx.NewInsert().
			Model(&chunk).
			On("CONFLICT (contest_id) DO UPDATE").
			Set("contest_id = EXCLUDED.contest_id").
			Set("start_epoch_second = EXCLUDED.start_epoch_second").
			Set("duration_second = EXCLUDED.duration_second").
			Set("title = EXCLUDED.title").
			Set("rate_change = EXCLUDED.rate_change").
			Set("category = EXCLUDED.category").
			Set("updated_at = NOW()").
			Exec(ctx)
		if err != nil {
			return 0, fmt.Errorf("insert: %w", err)
		}

		if c, err := res.RowsAffected(); err != nil {
			return 0, fmt.Errorf("rows affected: %w", err)
		} else {
			count += c
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit transaction: %w", err)
	}

	return count, nil
}

type Difficulty struct {
	bun.BaseModel    `bun:"table:difficulties,alias:d"`
	ProblemID        string   `bun:"problem_id"`
	Slope            *float64 `bun:"slope"`
	Intercept        *float64 `bun:"intercept"`
	Variance         *float64 `bun:"variance"`
	Difficulty       *int64   `bun:"difficulty"`
	Discrimination   *float64 `bun:"discrimination"`
	IrtLoglikelihood *float64 `bun:"irt_loglikelihood"`
	IrtUsers         *float64 `bun:"irt_users"`
	IsExperimental   *bool    `bun:"is_experimental"`
}

func NewDifficulties(difficulties map[string]atcoder.Difficulty) []Difficulty {
	result := make([]Difficulty, len(difficulties))

	for i, problemID := range slices.Sorted(maps.Keys(difficulties)) {
		d := difficulties[problemID]

		result[i] = Difficulty{
			ProblemID:        problemID,
			Slope:            d.Slope,
			Intercept:        d.Intercept,
			Variance:         d.Variance,
			Difficulty:       d.Difficulty,
			Discrimination:   d.Discrimination,
			IrtLoglikelihood: d.IrtLoglikelihood,
			IrtUsers:         d.IrtUsers,
			IsExperimental:   d.IsExperimental,
		}
	}

	return result
}

func SaveDifficulties(ctx context.Context, pool *pgxpool.Pool, difficulties map[string]atcoder.Difficulty) (int64, error) {
	if len(difficulties) == 0 {
		return 0, nil
	}

	tx, err := bun.NewDB(stdlib.OpenDBFromPool(pool), pgdialect.New()).BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return 0, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	var count int64 = 0

	for chunk := range slices.Chunk(NewDifficulties(difficulties), 1000) {
		res, err := tx.NewInsert().
			Model(&chunk).
			On("CONFLICT (problem_id) DO UPDATE").
			Set("problem_id = EXCLUDED.problem_id").
			Set("slope = EXCLUDED.slope").
			Set("intercept = EXCLUDED.intercept").
			Set("variance = EXCLUDED.variance").
			Set("difficulty = EXCLUDED.difficulty").
			Set("discrimination = EXCLUDED.discrimination").
			Set("irt_loglikelihood = EXCLUDED.irt_loglikelihood").
			Set("irt_users = EXCLUDED.irt_users").
			Set("is_experimental = EXCLUDED.is_experimental").
			Set("updated_at = NOW()").
			Exec(ctx)
		if err != nil {
			return 0, fmt.Errorf("insert: %w", err)
		}

		if c, err := res.RowsAffected(); err != nil {
			return 0, fmt.Errorf("rows affected: %w", err)
		} else {
			count += c
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit transaction: %w", err)
	}

	return count, nil
}
