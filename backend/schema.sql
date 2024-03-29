CREATE TABLE "contests" (
    "contest_id" text NOT NULL PRIMARY KEY,
    "start_epoch_second" bigint NOT NULL,
    "duration_second" bigint NOT NULL,
    "title" text NOT NULL,
    "rate_change" text NOT NULL,
    "category" text NOT NULL,
    "created_at" timestamp WITH time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updated_at" timestamp WITH time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE "problems" (
    "problem_id" text NOT NULL PRIMARY KEY,
    "contest_id" text NOT NULL,
    "problem_index" text NOT NULL,
    "name" text NOT NULL,
    "title" text NOT NULL,
    "url" text NOT NULL,
    "html" text NOT NULL,
    "created_at" timestamp WITH time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updated_at" timestamp WITH time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE "difficulties" (
    "problem_id" text NOT NULL PRIMARY KEY,
    "slope" double precision,
    "intercept" double precision,
    "variance" double precision,
    "difficulty" integer,
    "discrimination" double precision,
    "irt_loglikelihood" double precision,
    "irt_users" double precision,
    "is_experimental" boolean,
    "created_at" timestamp WITH time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updated_at" timestamp WITH time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE "users" (
    "user_name" text NOT NULL PRIMARY KEY,
    "rating" integer NOT NULL,
    "highest_rating" integer NOT NULL,
    "affiliation" text,
    "birth_year" integer,
    "country" text,
    "crown" text,
    "join_count" integer NOT NULL,
    "rank" integer NOT NULL,
    "active_rank" integer,
    "wins" integer NOT NULL,
    "created_at" timestamp WITH time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "updated_at" timestamp WITH time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE "submissions" (
    "id" bigint NOT NULL,
    "epoch_second" bigint NOT NULL,
    "problem_id" text NOT NULL,
    "contest_id" text,
    "user_id" text,
    "language" text,
    "point" double precision,
    "length" integer,
    "result" text,
    "execution_time" integer,
    "crawled_at" timestamp WITH time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    PRIMARY KEY ("id")
);

CREATE INDEX "submissions_contest_id_index" ON "submissions" ("contest_id");

CREATE INDEX "submissions_epoch_second_index" ON "submissions" ("epoch_second");

CREATE INDEX "submissions_problem_id_index" ON "submissions" ("problem_id");

CREATE INDEX "submissions_user_id_index" ON "submissions" ("user_id");

CREATE INDEX "submissions_language_index" ON "submissions" ("language");

CREATE INDEX "submissions_result_index" ON "submissions" ("result");

CREATE INDEX "submissions_crawled_at_index" ON "submissions" ("crawled_at");

CREATE TABLE "category_relationships" (
    "from" TEXT NOT NULL,
    "to" TEXT NOT NULL,
    "weight" DOUBLE PRECISION NOT NULL,
    PRIMARY KEY ("from", "to")
);

CREATE INDEX "category_relationships_from_index" ON "category_relationships" ("from");

CREATE INDEX "category_relationships_to_index" ON "category_relationships" ("to");

CREATE TABLE "update_history" (
    "id" bigserial NOT NULL,
    "domain" text NOT NULL,
    "started_at" timestamp WITH time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "finished_at" timestamp WITH time zone NOT NULL,
    "status" text NOT NULL,
    "options" json NOT NULL,
    PRIMARY KEY ("id")
);

CREATE INDEX "update_history_started_at_index" ON "update_history" ("started_at");

CREATE INDEX "update_history_domain_index" ON "update_history" ("domain");

CREATE TABLE "submission_crawl_history" (
    "id" bigserial NOT NULL,
    "contest_id" text NOT NULL,
    "started_at" bigint NOT NULL,
    PRIMARY KEY ("id")
);

CREATE INDEX "submission_crawl_history_contest_id_start_at_index" ON "submission_crawl_history" ("contest_id", "started_at");

CREATE TABLE "languages" (
    "language" text NOT NULL,
    "group" text,
    PRIMARY KEY ("language")
);

CREATE INDEX "languages_group_index" ON "languages" ("group");
