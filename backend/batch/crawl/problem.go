package crawl

import (
	"bytes"
	"context"
	"fjnkt98/atcodersearch/pkg/atcoder"
	"fjnkt98/atcodersearch/repository"
	"fmt"
	"strings"
	"time"

	"log/slog"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/goark/errs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
)

type ProblemCrawler struct {
	problemsClient *atcoder.AtCoderProblemsClient
	atcoderClient  *atcoder.AtCoderClient
	pool           *pgxpool.Pool
	minifier       *minify.M
	duration       time.Duration
	all            bool
}

func NewProblemCrawler(
	problemsClient *atcoder.AtCoderProblemsClient,
	atcoderClient *atcoder.AtCoderClient,
	pool *pgxpool.Pool,
	duration time.Duration,
	all bool,
) *ProblemCrawler {
	m := minify.New()
	m.AddFunc("text/html", html.Minify)

	return &ProblemCrawler{
		problemsClient: problemsClient,
		atcoderClient:  atcoderClient,
		pool:           pool,
		minifier:       m,
		duration:       duration,
		all:            all,
	}
}

func (c *ProblemCrawler) DetectDiff(ctx context.Context) ([]atcoder.Problem, error) {
	q := repository.New(c.pool)
	ids, err := q.FetchProblemIDs(ctx)
	if err != nil {
		return nil, errs.New(
			"failed to fetch existing problem ids",
			errs.WithCause(err),
		)
	}
	exists := mapset.NewSet(ids...)

	problems, err := c.problemsClient.FetchProblems(ctx)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	targets := make([]atcoder.Problem, 0, len(problems))
	for _, problem := range problems {
		if !exists.Contains(problem.ID) {
			targets = append(targets, problem)
		}
	}

	return targets, nil
}

func (c *ProblemCrawler) Crawl(ctx context.Context) error {
	var targets []atcoder.Problem
	var err error
	if c.all {
		slog.Info("Start to fetch all problems.")
		targets, err = c.problemsClient.FetchProblems(ctx)
		if err != nil {
			return errs.Wrap(err)
		}
		slog.Info("Finish fetching all problems.")
	} else {
		slog.Info("Start to fetch new problems.")
		targets, err = c.DetectDiff(ctx)
		if err != nil {
			return errs.Wrap(err)
		}
		slog.Info("Finish fetching new problems.")
	}

	for _, target := range targets {
		html, err := c.atcoderClient.FetchProblem(ctx, target.ContestID, target.ID)
		if err != nil {
			return errs.Wrap(err)
		}
		var buf bytes.Buffer
		if err := c.minifier.Minify("text/html", &buf, strings.NewReader(html)); err != nil {
			return errs.New(
				"failed to minify html",
				errs.WithCause(err),
				errs.WithContext("problem id", target.ID),
			)
		}

		tx, err := c.pool.Begin(ctx)
		if err != nil {
			return errs.New("failed to start transaction", errs.WithCause(err))
		}
		q := repository.New(tx)
		_, err = q.InsertProblem(
			ctx,
			repository.InsertProblemParams{
				ProblemID:    target.ID,
				ContestID:    target.ContestID,
				ProblemIndex: target.ProblemIndex,
				Name:         target.Name,
				Title:        target.Title,
				Url:          fmt.Sprintf("https://atcoder.jp/contests/%s/tasks/%s", target.ContestID, target.ID),
				Html:         buf.String(),
			})
		if err != nil {
			return errs.New("failed to insert problem", errs.WithCause(err), errs.WithContext("problem", target))
		}
		if err := tx.Commit(ctx); err != nil {
			return errs.New("failed to commit transaction", errs.WithCause(err))
		}
		slog.Info("Finish crawling the problem successfully", slog.String("target", target.ID))
		time.Sleep(c.duration)
	}

	return nil
}
