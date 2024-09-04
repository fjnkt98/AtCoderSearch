-- migrate:up
CREATE TABLE IF NOT EXISTS "public"."submission_crawl_histories" (
    "id" bigserial NOT NULL,
    "contest_id" TEXT NOT NULL,
    "started_at" BIGINT NOT NULL,
    PRIMARY KEY ("id")
);

CREATE INDEX IF NOT EXISTS "submission_crawl_histories_contest_id_start_at_index" ON "public"."submission_crawl_histories" ("contest_id", "started_at");

-- migrate:down
DROP TABLE IF EXISTS "public"."submission_crawl_history";
