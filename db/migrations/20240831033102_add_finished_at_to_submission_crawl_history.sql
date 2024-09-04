-- migrate:up
ALTER TABLE "submission_crawl_histories"
ADD COLUMN "finished_at" TIMESTAMPTZ;

-- migrate:down
ALTER TABLE "submission_crawl_histories"
DROP COLUMN "finished_at";
