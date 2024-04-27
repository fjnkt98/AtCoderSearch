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
        "updated_at"
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7, NOW())
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

-- name: CreateBatchHistory :one
INSERT INTO
    "batch_history" ("name", "started_at", "options")
VALUES
    ($1, NOW(), $2)
RETURNING
    "id";

-- name: UpdateBatchHistory :exec
UPDATE "batch_history"
SET
    "finished_at" = NOW(),
    "status" = $1
WHERE
    "id" = $2;

-- name: FetchLatestBatchHistory :one
SELECT
    "id",
    "started_at",
    "finished_at"
FROM
    "batch_history"
WHERE
    "name" = $1
    AND "status" = 'finished'
ORDER BY
    "started_at" DESC
LIMIT
    1;
