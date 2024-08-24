-- migrate:up
CREATE TABLE IF NOT EXISTS "public"."problems" (
    "problem_id" TEXT NOT NULL,
    "contest_id" TEXT NOT NULL,
    "problem_index" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "title" TEXT NOT NULL,
    "url" TEXT NOT NULL,
    "html" TEXT NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("problem_id")
);

-- migrate:down
DROP TABLE IF EXISTS "public"."problems";
