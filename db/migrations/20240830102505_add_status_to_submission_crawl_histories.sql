-- migrate:up
ALTER TABLE submission_crawl_histories
ADD COLUMN status TEXT NOT NULL DEFAULT 'working';

-- migrate:down
ALTER TABLE submission_crawl_histories
DROP COLUMN status;
