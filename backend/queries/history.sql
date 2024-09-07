-- name: CreateBatchHistory :one
INSERT INTO
    "batch_histories" ("name", "started_at", "options")
VALUES
    ($1, NOW(), $2)
RETURNING
    *;

-- name: CompleteBatchHistory :one
UPDATE "batch_histories"
SET
    "status" = 'completed',
    "finished_at" = NOW()
WHERE
    "id" = $1
RETURNING
    *;

-- name: AbortBatchHistory :one
UPDATE "batch_histories"
SET
    "status" = 'aborted',
    "finished_at" = NOW()
WHERE
    "id" = $1
RETURNING
    *;

-- name: FetchLatestBatchHistory :one
SELECT
    *
FROM
    "batch_histories"
WHERE
    "name" = $1
    AND "status" = 'finished'
ORDER BY
    "started_at" DESC
LIMIT
    1;

-- name: CreateCrawlHistory :one
INSERT INTO
    "submission_crawl_histories" ("contest_id")
VALUES
    ($1)
RETURNING
    *;

-- name: CompleteCrawlHistory :one
UPDATE "submission_crawl_histories"
SET
    "status" = 'completed',
    "finished_at" = NOW()
WHERE
    "id" = $1
RETURNING
    *;

-- name: FetchLatestCrawlHistory :one
SELECT
    *
FROM
    "submission_crawl_histories"
WHERE
    "contest_id" = $1
    AND "status" = 'finished'
ORDER BY
    "started_at" DESC
LIMIT
    1;

-- name: AbortCrawlHistory :one
UPDATE "submission_crawl_histories"
SET
    "status" = 'aborted',
    "finished_at" = NOW()
WHERE
    "id" = $1
RETURNING
    *;
