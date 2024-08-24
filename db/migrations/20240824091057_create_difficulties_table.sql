-- migrate:up
CREATE TABLE IF NOT EXISTS "public"."difficulties" (
    "problem_id" TEXT NOT NULL,
    "slope" DOUBLE PRECISION NULL,
    "intercept" DOUBLE PRECISION NULL,
    "variance" DOUBLE PRECISION NULL,
    "difficulty" BIGINT NULL,
    "discrimination" DOUBLE PRECISION NULL,
    "irt_loglikelihood" DOUBLE PRECISION NULL,
    "irt_users" DOUBLE PRECISION NULL,
    "is_experimental" BOOLEAN NULL,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("problem_id")
);

-- migrate:down
DROP TABLE IF EXISTS "public"."difficulties";
