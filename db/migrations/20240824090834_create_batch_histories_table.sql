-- migrate:up
CREATE TABLE IF NOT EXISTS "public"."batch_histories" (
    "id" bigserial NOT NULL,
    "name" TEXT NOT NULL,
    "started_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "finished_at" TIMESTAMPTZ NULL,
    "status" TEXT NOT NULL DEFAULT 'working',
    "options" json NULL,
    PRIMARY KEY ("id")
);

-- migrate:down
DROP TABLE IF EXISTS "public"."batch_histories";
