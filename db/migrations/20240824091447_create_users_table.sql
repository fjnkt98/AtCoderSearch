-- migrate:up
CREATE TABLE IF NOT EXISTS "public"."users" (
    "user_id" TEXT NOT NULL,
    "rating" INTEGER NOT NULL,
    "highest_rating" INTEGER NOT NULL,
    "affiliation" TEXT NULL,
    "birth_year" INTEGER NULL,
    "country" TEXT NULL,
    "crown" TEXT NULL,
    "join_count" INTEGER NOT NULL,
    "rank" INTEGER NOT NULL,
    "active_rank" INTEGER NULL,
    "wins" INTEGER NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("user_id")
);

-- migrate:down
DROP TABLE IF EXISTS "public"."users";
