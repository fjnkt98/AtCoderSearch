-- name: InsertContest :execresult
INSERT INTO
    "contests" (
        "contest_id",
        "start_epoch_second",
        "duration_second",
        "title",
        "rate_change",
        "category",
        "updated_at"
    )
VALUES
    ($1, $2, $3, $4, $5, $6, NOW())
ON CONFLICT ("contest_id") DO
UPDATE
SET
    "start_epoch_second" = EXCLUDED."start_epoch_second",
    "duration_second" = EXCLUDED."duration_second",
    "title" = EXCLUDED."title",
    "rate_change" = EXCLUDED."rate_change",
    "category" = EXCLUDED."category",
    "updated_at" = NOW();

-- name: FetchContestIDs :many
SELECT
    "contest_id"
FROM
    "contests"
ORDER BY
    "start_epoch_second" DESC;

-- name: FetchContestIDsByCategory :many
SELECT
    "contest_id"
FROM
    "contests"
WHERE
    "category" = ANY (@category::TEXT[])
ORDER BY
    "start_epoch_second" DESC;

-- name: FetchCategories :many
SELECT DISTINCT
    "category"
FROM
    "contests"
ORDER BY
    "category" ASC;

-- name: InsertDifficulty :execresult
INSERT INTO
    "difficulties" (
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
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
ON CONFLICT ("problem_id") DO
UPDATE
SET
    "slope" = EXCLUDED."slope",
    "intercept" = EXCLUDED."intercept",
    "variance" = EXCLUDED."variance",
    "difficulty" = EXCLUDED."difficulty",
    "discrimination" = EXCLUDED."discrimination",
    "irt_loglikelihood" = EXCLUDED."irt_loglikelihood",
    "irt_users" = EXCLUDED."irt_users",
    "is_experimental" = EXCLUDED."is_experimental",
    "updated_at" = NOW();

-- name: FetchLanguages :many
SELECT
    "language"
FROM
    "languages"
ORDER BY
    "language";

-- name: FetchLanguagesByGroup :many
SELECT
    "language"
FROM
    "languages"
WHERE
    "group" = ANY (@groups::TEXT[])
ORDER BY
    "language";

-- name: FetchLanguageGroups :many
SELECT DISTINCT
    "group"
FROM
    "languages"
ORDER BY
    "group" DESC;

-- name: InsertProblem :execresult
INSERT INTO
    "problems" (
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
VALUES
    ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
ON CONFLICT ("problem_id") DO
UPDATE
SET
    "contest_id" = EXCLUDED."contest_id",
    "problem_index" = EXCLUDED."problem_index",
    "name" = EXCLUDED."name",
    "title" = EXCLUDED."title",
    "url" = EXCLUDED."url",
    "html" = EXCLUDED."html",
    "updated_at" = NOW();

-- name: FetchProblemIDs :many
SELECT
    "problem_id"
FROM
    "problems";

-- name: FetchProblemIDsByContestID :many
SELECT
    "problem_id"
FROM
    "problems"
WHERE
    "contest_id" = ANY (@contest_id::TEXT[])
ORDER BY
    "problem_id" ASC;

-- name: InsertSubmission :execresult
INSERT INTO
    "submissions" (
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
VALUES
    ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW())
ON CONFLICT ("id") DO
UPDATE
SET
    "epoch_second" = EXCLUDED."epoch_second",
    "problem_id" = EXCLUDED."problem_id",
    "contest_id" = EXCLUDED."contest_id",
    "user_id" = EXCLUDED."user_id",
    "language" = EXCLUDED."language",
    "point" = EXCLUDED."point",
    "length" = EXCLUDED."length",
    "result" = EXCLUDED."result",
    "execution_time" = EXCLUDED."execution_time",
    "crawled_at" = NOW();

-- name: InsertUser :execresult
INSERT INTO
    "users" (
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
VALUES
    (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8,
        $9,
        $10,
        $11,
        NOW(),
        NOW()
    )
ON CONFLICT ("user_name") DO
UPDATE
SET
    "rating" = EXCLUDED."rating",
    "highest_rating" = EXCLUDED."highest_rating",
    "affiliation" = EXCLUDED."affiliation",
    "birth_year" = EXCLUDED."birth_year",
    "country" = EXCLUDED."country",
    "crown" = EXCLUDED."crown",
    "join_count" = EXCLUDED."join_count",
    "rank" = EXCLUDED."rank",
    "active_rank" = EXCLUDED."active_rank",
    "wins" = EXCLUDED."wins",
    "created_at" = NOW(),
    "updated_at" = NOW();

-- name: FetchRatingByUserName :one
SELECT
    "rating"
FROM
    "users"
WHERE
    "user_name" = $1
LIMIT
    1;

-- name: CreateCrawlHistory :one
INSERT INTO
    "submission_crawl_history" ("started_at", "contest_id")
VALUES
    (NOW(), $1)
RETURNING
    "started_at",
    "contest_id";

-- name: FetchLatestCrawlHistory :one
SELECT
    "started_at"
FROM
    "submission_crawl_history"
WHERE
    "contest_id" = $1
LIMIT
    1;

-- name: CreateUpdateHistory :one
INSERT INTO
    "update_history" ("domain", "started_at", "options")
VALUES
    ($1, NOW(), $2)
RETURNING
    "id",
    "domain",
    "started_at";

-- name: UpdateUpdateHistory :exec
UPDATE "update_history"
SET
    "finished_at" = NOW(),
    "status" = $1
WHERE
    "id" = $2;

-- name: FetchLatestUpdateHistory :one
SELECT
    "id",
    "started_at",
    "finished_at"
FROM
    "update_history"
WHERE
    "domain" = $1
    AND "status" = 'finished'
ORDER BY
    "started_at" DESC
LIMIT
    1;
