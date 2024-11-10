-- migrate:up
ALTER TABLE "submission_crawl_histories"
ALTER COLUMN "started_at" TYPE TIMESTAMP WITH TIME ZONE USING TO_TIMESTAMP("started_at") AT TIME ZONE 'UTC';

ALTER TABLE "submission_crawl_histories"
ALTER COLUMN "started_at"
SET DEFAULT CURRENT_TIMESTAMP;

-- migrate:down
ALTER TABLE "submission_crawl_histories"
ALTER COLUMN "started_at"
DROP DEFAULT;

ALTER TABLE "submission_crawl_histories"
ALTER COLUMN "started_at" TYPE BIGINT USING EXTRACT(
    epoch
    FROM
        "started_at"
)::BIGINT;
