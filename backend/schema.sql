CREATE TABLE "contests" (
    "contest_id" TEXT NOT NULL PRIMARY KEY,
    "start_epoch_second" BIGINT NOT NULL,
    "duration_second" BIGINT NOT NULL,
    "title" TEXT NOT NULL,
    "rate_change" TEXT NOT NULL,
    "category" TEXT NOT NULL,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE "problems" (
    "problem_id" TEXT NOT NULL PRIMARY KEY,
    "contest_id" TEXT NOT NULL,
    "problem_index" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "title" TEXT NOT NULL,
    "url" TEXT NOT NULL,
    "html" TEXT NOT NULL,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE "difficulties" (
    "problem_id" TEXT NOT NULL PRIMARY KEY,
    "slope" DOUBLE PRECISION,
    "intercept" DOUBLE PRECISION,
    "variance" DOUBLE PRECISION,
    "difficulty" BIGINT,
    "discrimination" DOUBLE PRECISION,
    "irt_loglikelihood" DOUBLE PRECISION,
    "irt_users" DOUBLE PRECISION,
    "is_experimental" BOOLEAN,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE "users" (
    "user_name" TEXT NOT NULL PRIMARY KEY,
    "rating" INTEGER NOT NULL,
    "highest_rating" INTEGER NOT NULL,
    "affiliation" TEXT,
    "birth_year" INTEGER,
    "country" TEXT,
    "crown" TEXT,
    "join_count" INTEGER NOT NULL,
    "rank" INTEGER NOT NULL,
    "active_rank" INTEGER,
    "wins" INTEGER NOT NULL,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE "submissions" (
    "id" BIGINT NOT NULL PRIMARY KEY,
    "epoch_second" BIGINT NOT NULL,
    "problem_id" TEXT NOT NULL,
    "contest_id" TEXT,
    "user_id" TEXT,
    "language" TEXT,
    "point" DOUBLE PRECISION,
    "length" INTEGER,
    "result" TEXT,
    "execution_time" INTEGER,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX "submissions_contest_id_index" ON "submissions" ("contest_id");

CREATE INDEX "submissions_epoch_second_index" ON "submissions" ("epoch_second");

CREATE INDEX "submissions_problem_id_index" ON "submissions" ("problem_id");

CREATE INDEX "submissions_user_id_index" ON "submissions" ("user_id");

CREATE INDEX "submissions_language_index" ON "submissions" ("language");

CREATE INDEX "submissions_result_index" ON "submissions" ("result");

CREATE INDEX "submissions_updated_at_index" ON "submissions" ("updated_at");

CREATE TABLE "batch_history" (
    "id" BIGSERIAL NOT NULL PRIMARY KEY,
    "name" TEXT NOT NULL,
    "started_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "finished_at" TIMESTAMP WITH TIME ZONE,
    "status" TEXT,
    "options" JSON NOT NULL
);

CREATE TABLE "submission_crawl_history" (
    "id" BIGSERIAL NOT NULL PRIMARY KEY,
    "contest_id" TEXT NOT NULL,
    "started_at" BIGINT NOT NULL
);

CREATE INDEX "submission_crawl_history_contest_id_start_at_index" ON "submission_crawl_history" ("contest_id", "started_at");

CREATE TABLE "languages" (
    "language" TEXT NOT NULL PRIMARY KEY,
    "group" TEXT
);

CREATE INDEX "languages_group_index" ON "languages" ("group");
