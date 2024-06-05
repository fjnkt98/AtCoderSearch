-- Add new schema named "public"
CREATE SCHEMA IF NOT EXISTS "public";
-- Set comment to schema: "public"
COMMENT ON SCHEMA "public" IS 'standard public schema';
-- Create "batch_history" table
CREATE TABLE "public"."batch_history" ("id" bigserial NOT NULL, "name" text NOT NULL, "started_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, "finished_at" timestamptz NULL, "status" text NOT NULL DEFAULT 'working', "options" json NOT NULL, PRIMARY KEY ("id"));
-- Create "contests" table
CREATE TABLE "public"."contests" ("contest_id" text NOT NULL, "start_epoch_second" bigint NOT NULL, "duration_second" bigint NOT NULL, "title" text NOT NULL, "rate_change" text NOT NULL, "category" text NOT NULL, "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY ("contest_id"));
-- Create "difficulties" table
CREATE TABLE "public"."difficulties" ("problem_id" text NOT NULL, "slope" double precision NULL, "intercept" double precision NULL, "variance" double precision NULL, "difficulty" bigint NULL, "discrimination" double precision NULL, "irt_loglikelihood" double precision NULL, "irt_users" double precision NULL, "is_experimental" boolean NULL, "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY ("problem_id"));
-- Create "languages" table
CREATE TABLE "public"."languages" ("language" text NOT NULL, "group" text NULL, PRIMARY KEY ("language"));
-- Create index "languages_group_index" to table: "languages"
CREATE INDEX "languages_group_index" ON "public"."languages" ("group");
-- Create "problems" table
CREATE TABLE "public"."problems" ("problem_id" text NOT NULL, "contest_id" text NOT NULL, "problem_index" text NOT NULL, "name" text NOT NULL, "title" text NOT NULL, "url" text NOT NULL, "html" text NOT NULL, "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY ("problem_id"));
-- Create "submission_crawl_history" table
CREATE TABLE "public"."submission_crawl_history" ("id" bigserial NOT NULL, "contest_id" text NOT NULL, "started_at" bigint NOT NULL, PRIMARY KEY ("id"));
-- Create index "submission_crawl_history_contest_id_start_at_index" to table: "submission_crawl_history"
CREATE INDEX "submission_crawl_history_contest_id_start_at_index" ON "public"."submission_crawl_history" ("contest_id", "started_at");
-- Create "submissions" table
CREATE TABLE "public"."submissions" ("id" bigint NOT NULL, "epoch_second" bigint NOT NULL, "problem_id" text NOT NULL, "contest_id" text NULL, "user_id" text NULL, "language" text NULL, "point" double precision NULL, "length" integer NULL, "result" text NULL, "execution_time" integer NULL, "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP) PARTITION BY RANGE ("epoch_second");
-- Create index "submissions_contest_id_index" to table: "submissions"
CREATE INDEX "submissions_contest_id_index" ON "public"."submissions" ("contest_id");
-- Create index "submissions_epoch_second_index" to table: "submissions"
CREATE INDEX "submissions_epoch_second_index" ON "public"."submissions" ("epoch_second");
-- Create index "submissions_execution_time_index" to table: "submissions"
CREATE INDEX "submissions_execution_time_index" ON "public"."submissions" ("execution_time");
-- Create index "submissions_language_index" to table: "submissions"
CREATE INDEX "submissions_language_index" ON "public"."submissions" ("language");
-- Create index "submissions_length_index" to table: "submissions"
CREATE INDEX "submissions_length_index" ON "public"."submissions" ("length");
-- Create index "submissions_point_index" to table: "submissions"
CREATE INDEX "submissions_point_index" ON "public"."submissions" ("point");
-- Create index "submissions_problem_id_index" to table: "submissions"
CREATE INDEX "submissions_problem_id_index" ON "public"."submissions" ("problem_id");
-- Create index "submissions_result_index" to table: "submissions"
CREATE INDEX "submissions_result_index" ON "public"."submissions" ("result");
-- Create index "submissions_updated_at_index" to table: "submissions"
CREATE INDEX "submissions_updated_at_index" ON "public"."submissions" ("epoch_second", "updated_at");
-- Create index "submissions_user_id_index" to table: "submissions"
CREATE INDEX "submissions_user_id_index" ON "public"."submissions" ("user_id");
-- Create "users" table
CREATE TABLE "public"."users" ("user_id" text NOT NULL, "rating" integer NOT NULL, "highest_rating" integer NOT NULL, "affiliation" text NULL, "birth_year" integer NULL, "country" text NULL, "crown" text NULL, "join_count" integer NOT NULL, "rank" integer NOT NULL, "active_rank" integer NULL, "wins" integer NOT NULL, "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY ("user_id"));
