-- migrate:up
CREATE TABLE IF NOT EXISTS "public"."submissions" (
    "id" BIGINT NOT NULL,
    "epoch_second" BIGINT NOT NULL,
    "problem_id" TEXT NOT NULL,
    "contest_id" TEXT NULL,
    "user_id" TEXT NULL,
    "language" TEXT NULL,
    "point" DOUBLE PRECISION NULL,
    "length" INTEGER NULL,
    "result" TEXT NULL,
    "execution_time" INTEGER NULL,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
)
PARTITION BY
    RANGE ("epoch_second");

CREATE UNIQUE INDEX IF NOT EXISTS "submissions_unique" ON "public"."submissions" ("epoch_second", "id");

CREATE INDEX IF NOT EXISTS "submissions_contest_id_epoch_second_index" ON "public"."submissions" ("contest_id", "epoch_second");

CREATE INDEX IF NOT EXISTS "submissions_contest_id_execution_time_index" ON "public"."submissions" ("contest_id", "execution_time");

CREATE INDEX IF NOT EXISTS "submissions_contest_id_length_index" ON "public"."submissions" ("contest_id", "length");

CREATE INDEX IF NOT EXISTS "submissions_contest_id_point_index" ON "public"."submissions" ("contest_id", "point");

CREATE INDEX IF NOT EXISTS "submissions_language_epoch_second_index" ON "public"."submissions" ("language", "epoch_second");

CREATE INDEX IF NOT EXISTS "submissions_language_execution_time_index" ON "public"."submissions" ("language", "execution_time");

CREATE INDEX IF NOT EXISTS "submissions_language_length_index" ON "public"."submissions" ("language", "length");

CREATE INDEX IF NOT EXISTS "submissions_language_point_index" ON "public"."submissions" ("language", "point");

CREATE INDEX IF NOT EXISTS "submissions_problem_id_epoch_second_index" ON "public"."submissions" ("problem_id", "epoch_second");

CREATE INDEX IF NOT EXISTS "submissions_problem_id_execution_time_index" ON "public"."submissions" ("problem_id", "execution_time");

CREATE INDEX IF NOT EXISTS "submissions_problem_id_length_index" ON "public"."submissions" ("problem_id", "length");

CREATE INDEX IF NOT EXISTS "submissions_problem_id_point_index" ON "public"."submissions" ("problem_id", "point");

CREATE INDEX IF NOT EXISTS "submissions_result_epoch_second_index" ON "public"."submissions" ("result", "epoch_second");

CREATE INDEX IF NOT EXISTS "submissions_result_execution_time_index" ON "public"."submissions" ("result", "execution_time");

CREATE INDEX IF NOT EXISTS "submissions_result_length_index" ON "public"."submissions" ("result", "length");

CREATE INDEX IF NOT EXISTS "submissions_result_point_index" ON "public"."submissions" ("result", "point");

CREATE INDEX IF NOT EXISTS "submissions_user_id_epoch_second_index" ON "public"."submissions" ("user_id", "epoch_second");

CREATE INDEX IF NOT EXISTS "submissions_user_id_execution_time_index" ON "public"."submissions" ("user_id", "execution_time");

CREATE INDEX IF NOT EXISTS "submissions_user_id_length_index" ON "public"."submissions" ("user_id", "length");

CREATE INDEX IF NOT EXISTS "submissions_user_id_point_index" ON "public"."submissions" ("user_id", "point");

CREATE INDEX IF NOT EXISTS "submissions_user_id_result_index" ON "public"."submissions" ("user_id", "result");

CREATE TABLE IF NOT EXISTS "submissions_2010" PARTITION OF "public"."submissions" FOR
VALUES
FROM
    (0) TO (1293807600);

CREATE TABLE IF NOT EXISTS "submissions_2011" PARTITION OF "public"."submissions" FOR
VALUES
FROM
    (1293807600) TO (1325343600);

CREATE TABLE IF NOT EXISTS "submissions_2012" PARTITION OF "public"."submissions" FOR
VALUES
FROM
    (1325343600) TO (1356966000);

CREATE TABLE IF NOT EXISTS "submissions_2013" PARTITION OF "public"."submissions" FOR
VALUES
FROM
    (1356966000) TO (1388502000);

CREATE TABLE IF NOT EXISTS "submissions_2014" PARTITION OF "public"."submissions" FOR
VALUES
FROM
    (1388502000) TO (1420038000);

CREATE TABLE IF NOT EXISTS "submissions_2015" PARTITION OF "public"."submissions" FOR
VALUES
FROM
    (1420038000) TO (1451574000);

CREATE TABLE IF NOT EXISTS "submissions_2016" PARTITION OF "public"."submissions" FOR
VALUES
FROM
    (1451574000) TO (1483196400);

CREATE TABLE IF NOT EXISTS "submissions_2017" PARTITION OF "public"."submissions" FOR
VALUES
FROM
    (1483196400) TO (1514732400);

CREATE TABLE IF NOT EXISTS "submissions_2018" PARTITION OF "public"."submissions" FOR
VALUES
FROM
    (1514732400) TO (1546268400);

CREATE TABLE IF NOT EXISTS "submissions_2019" PARTITION OF "public"."submissions" FOR
VALUES
FROM
    (1546268400) TO (1577804400);

CREATE TABLE IF NOT EXISTS "submissions_2020" PARTITION OF "public"."submissions" FOR
VALUES
FROM
    (1577804400) TO (1609426800);

CREATE TABLE IF NOT EXISTS "submissions_2021" PARTITION OF "public"."submissions" FOR
VALUES
FROM
    (1609426800) TO (1640962800);

CREATE TABLE IF NOT EXISTS "submissions_2022" PARTITION OF "public"."submissions" FOR
VALUES
FROM
    (1640962800) TO (1672498800);

CREATE TABLE IF NOT EXISTS "submissions_2023" PARTITION OF "public"."submissions" FOR
VALUES
FROM
    (1672498800) TO (1704034800);

CREATE TABLE IF NOT EXISTS "submissions_2024" PARTITION OF "public"."submissions" FOR
VALUES
FROM
    (1704034800) TO (1735657200);

CREATE TABLE IF NOT EXISTS "submissions_2025" PARTITION OF "public"."submissions" FOR
VALUES
FROM
    (1735657200) TO (1767193200);

CREATE TABLE IF NOT EXISTS "submissions_2026" PARTITION OF "public"."submissions" FOR
VALUES
FROM
    (1767193200) TO (1798729200);

CREATE TABLE IF NOT EXISTS "submissions_2027" PARTITION OF "public"."submissions" FOR
VALUES
FROM
    (1798729200) TO (1830265200);

CREATE TABLE IF NOT EXISTS "submissions_2028" PARTITION OF "public"."submissions" FOR
VALUES
FROM
    (1830265200) TO (1861887600);

CREATE TABLE IF NOT EXISTS "submissions_2029" PARTITION OF "public"."submissions" FOR
VALUES
FROM
    (1861887600) TO (1893423600);

-- migrate:down
DROP TABLE IF EXISTS "public"."submissions";

DROP INDEX IF EXISTS "submissions_unique";

DROP INDEX IF EXISTS "submissions_contest_id_epoch_second_index";

DROP INDEX IF EXISTS "submissions_contest_id_execution_time_index";

DROP INDEX IF EXISTS "submissions_contest_id_length_index";

DROP INDEX IF EXISTS "submissions_contest_id_point_index";

DROP INDEX IF EXISTS "submissions_language_epoch_second_index";

DROP INDEX IF EXISTS "submissions_language_execution_time_index";

DROP INDEX IF EXISTS "submissions_language_length_index";

DROP INDEX IF EXISTS "submissions_language_point_index";

DROP INDEX IF EXISTS "submissions_problem_id_epoch_second_index";

DROP INDEX IF EXISTS "submissions_problem_id_execution_time_index";

DROP INDEX IF EXISTS "submissions_problem_id_length_index";

DROP INDEX IF EXISTS "submissions_problem_id_point_index";

DROP INDEX IF EXISTS "submissions_result_epoch_second_index";

DROP INDEX IF EXISTS "submissions_result_execution_time_index";

DROP INDEX IF EXISTS "submissions_result_length_index";

DROP INDEX IF EXISTS "submissions_result_point_index";

DROP INDEX IF EXISTS "submissions_user_id_epoch_second_index";

DROP INDEX IF EXISTS "submissions_user_id_execution_time_index";

DROP INDEX IF EXISTS "submissions_user_id_length_index";

DROP INDEX IF EXISTS "submissions_user_id_point_index";

DROP INDEX IF EXISTS "submissions_user_id_result_index";

DROP TABLE IF EXISTS "submissions_2010";

DROP TABLE IF EXISTS "submissions_2011";

DROP TABLE IF EXISTS "submissions_2012";

DROP TABLE IF EXISTS "submissions_2013";

DROP TABLE IF EXISTS "submissions_2014";

DROP TABLE IF EXISTS "submissions_2015";

DROP TABLE IF EXISTS "submissions_2016";

DROP TABLE IF EXISTS "submissions_2017";

DROP TABLE IF EXISTS "submissions_2018";

DROP TABLE IF EXISTS "submissions_2019";

DROP TABLE IF EXISTS "submissions_2020";

DROP TABLE IF EXISTS "submissions_2021";

DROP TABLE IF EXISTS "submissions_2022";

DROP TABLE IF EXISTS "submissions_2023";

DROP TABLE IF EXISTS "submissions_2024";

DROP TABLE IF EXISTS "submissions_2025";

DROP TABLE IF EXISTS "submissions_2026";

DROP TABLE IF EXISTS "submissions_2027";

DROP TABLE IF EXISTS "submissions_2028";

DROP TABLE IF EXISTS "submissions_2029";
