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
    "problems"
ORDER BY
    "problem_id" ASC;
