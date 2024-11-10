-- migrate:up
CREATE TABLE IF NOT EXISTS "public"."contests" (
    "contest_id" TEXT NOT NULL,
    "start_epoch_second" BIGINT NOT NULL,
    "duration_second" BIGINT NOT NULL,
    "title" TEXT NOT NULL,
    "rate_change" TEXT NOT NULL,
    "category" TEXT NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("contest_id")
);

-- migrate:down
DROP TABLE IF EXISTS "public"."contests";
