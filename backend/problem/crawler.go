package problem

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/jmoiron/sqlx"
	"github.com/morikuni/failure"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
)

type ContestCrawler struct {
	targetURL string
	db        *sqlx.DB
	client    *http.Client
}

func NewContestCrawler(db *sqlx.DB) ContestCrawler {
	return ContestCrawler{
		targetURL: "https://kenkoooo.com/atcoder/resources/contests.json",
		db:        db,
		client:    &http.Client{},
	}
}

func (c *ContestCrawler) FetchContestList() ([]ContestJSON, error) {
	req, err := http.NewRequest("GET", c.targetURL, nil)
	if err != nil {
		return nil, failure.Translate(err, RequestCreationError, failure.Context{"url": c.targetURL}, failure.Message("failed to create request"))
	}
	req.Header.Set("ACCEPT_ENCODING", "gzip")

	res, err := c.client.Do(req)
	if err != nil {
		return nil, failure.Translate(err, RequestExecutionError, failure.Context{"url": c.targetURL}, failure.Message("request failed"))
	}

	defer res.Body.Close()
	var contests []ContestJSON
	if err := json.NewDecoder(res.Body).Decode(&contests); err != nil {
		return nil, failure.Translate(err, DecodeError, failure.Message("failed to decode JSON into ContestJSON"))
	}

	return contests, nil
}

func (c *ContestCrawler) Save(contests []ContestJSON) error {
	tx, err := c.db.Beginx()
	if err != nil {
		return failure.Translate(err, DBError, failure.Message("failed to start transaction to save contest information"))
	}
	defer tx.Rollback()

	for _, contestJSON := range contests {
		contest := Contest{
			ContestID:        contestJSON.ID,
			StartEpochSecond: contestJSON.StartEpochSecond,
			DurationSecond:   contestJSON.DurationSecond,
			Title:            contestJSON.Title,
			RateChange:       contestJSON.RateChange,
			Category:         contestJSON.Categorize(),
		}

		_, err := tx.Exec(`
			MERGE INTO "contests"
			USING
				(
					VALUES(
						$1::text,
						$2::bigint,
						$3::bigint,
						$4::text,
						$5::text,
						$6::text
					)
				) AS "contest"(
					"contest_id",
					"start_epoch_second",
					"duration_second",
					"title",
					"rate_change",
					"category"
				)
			ON
				"contests"."contest_id" = "contest"."contest_id"
			WHEN MATCHED THEN
				UPDATE SET (
					"contest_id",
					"start_epoch_second",
					"duration_second",
					"title",
					"rate_change",
					"category",
					"updated_at"
				) = (
					"contest"."contest_id",
					"contest"."start_epoch_second",
					"contest"."duration_second",
					"contest"."title",
					"contest"."rate_change",
					"contest"."category",
					NOW()
				)
			WHEN NOT MATCHED THEN
				INSERT (
					"contest_id",
					"start_epoch_second",
					"duration_second",
					"title",
					"rate_change",
					"category",
					"created_at",
					"updated_at"
				)
				VALUES (
					"contest"."contest_id",
					"contest"."start_epoch_second",
					"contest"."duration_second",
					"contest"."title",
					"contest"."rate_change",
					"contest"."category",
					NOW(),
					NOW()
				);
			`,
			contest.ContestID,
			contest.StartEpochSecond,
			contest.DurationSecond,
			contest.Title,
			contest.RateChange,
			contest.Category,
		)
		if err != nil {
			return failure.Translate(err, DBError, failure.Context{"contestID": contest.ContestID}, failure.Message("failed to save contest information"))
		}
	}

	if err := tx.Commit(); err != nil {
		return failure.Translate(err, DBError, failure.Message("failed to commit transaction"))
	}

	return nil
}

func (c *ContestCrawler) Run() error {
	log.Println("Start to fetch contest list.")
	contests, err := c.FetchContestList()
	if err != nil {
		return failure.Wrap(err)
	}
	log.Println("Finish fetching contest list.")

	log.Println("Start to save contest list.")
	if err := c.Save(contests); err != nil {
		return failure.Wrap(err)
	}
	log.Println("Finish saving contest list.")

	return nil
}

type ProblemCrawler struct {
	targetURL string
	db        *sqlx.DB
	client    *http.Client
	minifier  *minify.M
}

func NewProblemCrawler(db *sqlx.DB) ProblemCrawler {
	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	return ProblemCrawler{
		targetURL: "https://kenkoooo.com/atcoder/resources/problems.json",
		db:        db,
		client:    &http.Client{},
		minifier:  m,
	}
}

func (c *ProblemCrawler) FetchProblemList() ([]ProblemJSON, error) {
	req, err := http.NewRequest("GET", c.targetURL, nil)
	if err != nil {
		return nil, failure.Translate(err, RequestCreationError, failure.Context{"url": c.targetURL}, failure.Message("failed to create request"))
	}
	req.Header.Set("ACCEPT_ENCODING", "gzip")

	res, err := c.client.Do(req)
	if err != nil {
		return nil, failure.Translate(err, RequestExecutionError, failure.Context{"url": c.targetURL}, failure.Message("request failed"))
	}

	defer res.Body.Close()
	var problems []ProblemJSON
	if err := json.NewDecoder(res.Body).Decode(&problems); err != nil {
		return nil, failure.Translate(err, DecodeError, failure.Message("failed to decode JSON into ProblemJSON"))
	}

	return problems, nil
}

func (c *ProblemCrawler) Crawl(URL string) (string, error) {
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return "", failure.Translate(err, RequestCreationError, failure.Context{"url": c.targetURL}, failure.Message("failed to create request"))
	}
	res, err := c.client.Do(req)
	if err != nil {
		return "", failure.Translate(err, RequestExecutionError, failure.Context{"url": c.targetURL}, failure.Message("request failed"))
	}

	var buf bytes.Buffer
	defer res.Body.Close()
	if err := c.minifier.Minify("text/html", &buf, res.Body); err != nil {
		return "", failure.Translate(err, MinifyHTMLError, failure.Context{"url": c.targetURL}, failure.Message("error occurred in minifying html"))
	}

	return buf.String(), nil
}

func (c *ProblemCrawler) DetectDiff() ([]ProblemJSON, error) {
	exists := mapset.NewSet[string]()

	rows, err := c.db.Queryx(`SELECT "problem_id" FROM "problems"`)
	if err != nil {
		return nil, failure.Translate(err, DBError, failure.Message("failed to select problems id"))
	}

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, failure.Translate(err, DBError, failure.Message("failed to scan problem_id"))
		}
		exists.Add(id)
	}

	problems, err := c.FetchProblemList()
	if err != nil {
		return nil, failure.Wrap(err)
	}

	targetProblems := make([]ProblemJSON, 0, len(problems))
	for _, problem := range problems {
		if !exists.Contains(problem.ID) {
			targetProblems = append(targetProblems, problem)
		}
	}

	return targetProblems, nil
}

func (c *ProblemCrawler) Save(problemJSONs []ProblemJSON, duration int) error {
	for _, problemJSON := range problemJSONs {
		tx, err := c.db.Beginx()
		if err != nil {
			return failure.Translate(err, DBError, failure.Message("failed to start transaction"))
		}
		defer tx.Rollback()

		url := fmt.Sprintf("https://atcoder.jp/contests/%s/tasks/%s", problemJSON.ContestID, problemJSON.ID)
		HTML, err := c.Crawl(url)
		if err != nil {
			return failure.Translate(err, DBError, failure.Context{"url": url}, failure.Message("failed to crawl problem"))
		}

		problem := Problem{
			ProblemID:    problemJSON.ID,
			ContestID:    problemJSON.ContestID,
			ProblemIndex: problemJSON.ProblemIndex,
			Name:         problemJSON.Name,
			Title:        problemJSON.Title,
			URL:          url,
			HTML:         HTML,
		}

		log.Printf("save `%s`", problem.ProblemID)
		if _, err := tx.Exec(`
			MERGE INTO "problems"
			USING
				(
					VALUES(
						$1::text,
						$2::text,
						$3::text,
						$4::text,
						$5::text,
						$6::text,
						$7::text
					)
				) AS "problem"(
					"problem_id",
					"contest_id",
					"problem_index",
					"name",
					"title",
					"url",
					"html"
				)
			ON
				"problems"."problem_id" = "problem"."problem_id"
			WHEN MATCHED THEN
				UPDATE SET (
					"problem_id",
					"contest_id",
					"problem_index",
					"name",
					"title",
					"url",
					"html",
					"updated_at"
				) = (
					"problem"."problem_id",
					"problem"."contest_id",
					"problem"."problem_index",
					"problem"."name",
					"problem"."title",
					"problem"."url",
					"problem"."html",
					NOW()
				)
			WHEN NOT MATCHED THEN
				INSERT (
					"problem_id",
					"contest_id",
					"problem_index",
					"name",
					"title",
					"url",
					"html",
					"created_at",
					"updated_at"
				)
				VALUES (
					"problem"."problem_id",
					"problem"."contest_id",
					"problem"."problem_index",
					"problem"."name",
					"problem"."title",
					"problem"."url",
					"problem"."html",
					NOW(),
					NOW()
				);
			`,
			problem.ProblemID,
			problem.ContestID,
			problem.ProblemIndex,
			problem.Name,
			problem.Title,
			problem.URL,
			problem.HTML,
		); err != nil {
			return failure.Translate(err, DBError, failure.Context{"problemID": problem.ProblemID}, failure.Message("failed to save problem"))
		}
		if err := tx.Commit(); err != nil {
			return failure.Translate(err, DBError, failure.Context{"problemID": problem.ProblemID}, failure.Message("failed to commit to save problem `%s`"))
		}

		time.Sleep(time.Duration(duration) * time.Millisecond)
	}

	return nil
}

func (c *ProblemCrawler) Run(all bool, duration int) error {
	var targets []ProblemJSON
	var err error
	if all {
		log.Println("Start to fetch all problem list.")
		targets, err = c.FetchProblemList()
		log.Println("Finish fetching all problem list.")
	} else {
		log.Println("Start to fetch new problem list.")
		targets, err = c.DetectDiff()
		log.Println("Finish fetching new problem list.")
	}
	if err != nil {
		return failure.Wrap(err)
	}

	log.Println("Start to save problem list.")
	if err := c.Save(targets, duration); err != nil {
		return failure.Wrap(err)
	}
	log.Println("Finish saving problem list.")
	return nil
}

type DifficultyCrawler struct {
	targetURL string
	db        *sqlx.DB
	client    *http.Client
}

func NewDifficultyCrawler(db *sqlx.DB) DifficultyCrawler {
	return DifficultyCrawler{
		targetURL: "https://kenkoooo.com/atcoder/resources/problem-models.json",
		db:        db,
		client:    &http.Client{},
	}
}

func (c *DifficultyCrawler) FetchDifficulties() (map[string]DifficultyJSON, error) {
	req, err := http.NewRequest("GET", c.targetURL, nil)
	if err != nil {
		return nil, failure.Translate(err, RequestCreationError, failure.Context{"url": c.targetURL}, failure.Message("failed to create request"))
	}
	req.Header.Set("ACCEPT_ENCODING", "gzip")

	res, err := c.client.Do(req)
	if err != nil {
		return nil, failure.Translate(err, RequestExecutionError, failure.Context{"url": c.targetURL}, failure.Message("request failed"))
	}

	defer res.Body.Close()
	var difficulties map[string]DifficultyJSON
	if err := json.NewDecoder(res.Body).Decode(&difficulties); err != nil {
		return nil, failure.Translate(err, DecodeError, failure.Message("failed to decode JSON into DifficultyJSON"))
	}

	return difficulties, nil
}

func (c *DifficultyCrawler) Save(difficulties map[string]DifficultyJSON) error {
	tx, err := c.db.Beginx()
	if err != nil {
		return failure.Translate(err, DBError, failure.Message("failed to start transaction"))
	}
	defer tx.Rollback()

	for problemID, difficultyJSON := range difficulties {
		difficulty := Difficulty{
			ProblemID:        problemID,
			Slope:            difficultyJSON.Slope,
			Intercept:        difficultyJSON.Intercept,
			Variance:         difficultyJSON.Variance,
			Difficulty:       difficultyJSON.Difficulty,
			Discrimination:   difficultyJSON.Discrimination,
			IrtLogLikelihood: difficultyJSON.IrtLogLikelihood,
			IrtUsers:         difficultyJSON.IrtUsers,
			IsExperimental:   difficultyJSON.IsExperimental,
		}

		if _, err := tx.Exec(`
			MERGE INTO "difficulties"
			USING
				(
					VALUES (
						$1::text,
						$2::double precision,
						$3::double precision,
						$4::double precision,
						$5::integer,
						$6::double precision,
						$7::double precision,
						$8::double precision,
						$9::boolean
					)
				) AS "difficulty"(
					"problem_id",
					"slope",
					"intercept",
					"variance",
					"difficulty",
					"discrimination",
					"irt_loglikelihood",
					"irt_users",
					"is_experimental"
				)
			ON
				"difficulties"."problem_id" = "difficulty"."problem_id"
			WHEN MATCHED THEN
				UPDATE SET (
					"problem_id",
					"slope",
					"intercept",
					"variance",
					"difficulty",
					"discrimination",
					"irt_loglikelihood",
					"irt_users",
					"is_experimental",
					"updated_at"
				) = (
					"difficulty"."problem_id",
					"difficulty"."slope",
					"difficulty"."intercept",
					"difficulty"."variance",
					"difficulty"."difficulty",
					"difficulty"."discrimination",
					"difficulty"."irt_loglikelihood",
					"difficulty"."irt_users",
					"difficulty"."is_experimental",
					NOW()
				)
			WHEN NOT MATCHED THEN
				INSERT (
					"problem_id",
					"slope",
					"intercept",
					"variance",
					"difficulty",
					"discrimination",
					"irt_loglikelihood",
					"irt_users",
					"is_experimental",
					"created_at",
					"updated_at"
				)
				VALUES (
					"difficulty"."problem_id",
					"difficulty"."slope",
					"difficulty"."intercept",
					"difficulty"."variance",
					"difficulty"."difficulty",
					"difficulty"."discrimination",
					"difficulty"."irt_loglikelihood",
					"difficulty"."irt_users",
					"difficulty"."is_experimental",
					NOW(),
					NOW()
				);
			`,
			difficulty.ProblemID,
			difficulty.Slope,
			difficulty.Intercept,
			difficulty.Variance,
			difficulty.Difficulty,
			difficulty.Discrimination,
			difficulty.IrtLogLikelihood,
			difficulty.IrtUsers,
			difficulty.IsExperimental,
		); err != nil {
			return failure.Translate(err, DBError, failure.Context{"problemID": difficulty.ProblemID}, failure.Message("failed to save difficulty"))
		}
	}

	if err := tx.Commit(); err != nil {
		return failure.Translate(err, DBError, failure.Message("failed to commit transaction"))
	}

	return nil
}

func (c *DifficultyCrawler) Run() error {
	log.Println("Start to fetch difficulty list.")
	difficulties, err := c.FetchDifficulties()
	if err != nil {
		return failure.Wrap(err)
	}
	log.Println("Finish fetching difficulty list.")

	log.Println("Start to save difficulty list.")
	if err := c.Save(difficulties); err != nil {
		return failure.Wrap(err)
	}
	log.Println("Finish saving difficulty list.")

	return nil
}
