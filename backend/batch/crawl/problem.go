package crawl

import (
	"bytes"
	"context"
	"fjnkt98/atcodersearch/batch"
	"fjnkt98/atcodersearch/config"
	"fjnkt98/atcodersearch/pkg/atcoder"
	"fjnkt98/atcodersearch/repository"
	"fmt"
	"strings"
	"time"

	"golang.org/x/exp/slog"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/goark/errs"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
)

type ProblemCrawler interface {
	batch.Batch
	CrawlProblem(ctx context.Context) error
}

type problemCrawler struct {
	problemsClient atcoder.AtCoderProblemsClient
	atcoderClient  atcoder.AtCoderClient
	repo           repository.ProblemRepository
	minifier       *minify.M
	cfg            config.CrawlProblemConfig
}

func NewProblemCrawler(
	problemsClient atcoder.AtCoderProblemsClient,
	atcoderClient atcoder.AtCoderClient,
	repo repository.ProblemRepository,
	cfg config.CrawlProblemConfig,
) ProblemCrawler {
	m := minify.New()
	m.AddFunc("text/html", html.Minify)

	return &problemCrawler{
		problemsClient: problemsClient,
		atcoderClient:  atcoderClient,
		repo:           repo,
		minifier:       m,
		cfg:            cfg,
	}
}

func (c *problemCrawler) Name() string {
	return "ProblemCrawler"
}

func (c *problemCrawler) DetectDiff(ctx context.Context) ([]atcoder.Problem, error) {
	ids, err := c.repo.FetchIDs(ctx)
	if err != nil {
		return nil, errs.New(
			"failed to fetch existing problem ids",
			errs.WithCause(err),
		)
	}
	exists := mapset.NewSet[string](ids...)

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

func (c *problemCrawler) CrawlProblem(ctx context.Context) error {
	var targets []atcoder.Problem
	var err error
	if c.cfg.All {
		slog.Info("Start to fetch all problems.")
		targets, err = c.problemsClient.FetchProblems(ctx)
		slog.Info("Finish fetching all problems.")
	} else {
		slog.Info("Start to fetch new problems.")
		targets, err = c.DetectDiff(ctx)
		slog.Info("Finish fetching new problems.")
	}

	if err != nil {
		return errs.Wrap(err)
	}

	problems := make([]repository.Problem, 0, len(targets))
	for _, target := range targets {
		slog.Info("Start to crawl problem `%s`", target)
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
		problems = append(problems, repository.Problem{
			ProblemID:    target.ID,
			ContestID:    target.ContestID,
			ProblemIndex: target.ProblemIndex,
			Name:         target.Name,
			Title:        target.Title,
			URL:          fmt.Sprintf("https://atcoder.jp/contests/%s/tasks/%s", target.ContestID, target.ID),
			HTML:         buf.String(),
		})
		slog.Info("Finish crawling problem `%s` successfully", target)
		time.Sleep(time.Duration(c.cfg.Duration) * time.Millisecond)
	}

	if err := c.repo.Save(ctx, problems); err != nil {
		return errs.Wrap(err)
	}

	return nil
}

func (c *problemCrawler) Run(ctx context.Context) error {
	return c.CrawlProblem(ctx)
}
