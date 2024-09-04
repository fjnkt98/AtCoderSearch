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

CREATE INDEX IF NOT EXISTS "submissions_contest_id_epoch_second_index" ON "public"."submissions" ("contest_id", "epoch_second");

CREATE INDEX IF NOT EXISTS "submissions_contest_id_execution_time_index" ON "public"."submissions" ("contest_id", "execution_time");

CREATE INDEX IF NOT EXISTS "submissions_contest_id_length_index" ON "public"."submissions" ("contest_id", "length");

CREATE INDEX IF NOT EXISTS "submissions_contest_id_point_index" ON "public"."submissions" ("contest_id", "point");

CREATE INDEX IF NOT EXISTS "submissions_execution_time_epoch_second_index" ON "public"."submissions" ("execution_time", "epoch_second");

CREATE INDEX IF NOT EXISTS "submissions_execution_time_length_index" ON "public"."submissions" ("execution_time", "length");

CREATE INDEX IF NOT EXISTS "submissions_execution_time_point_index" ON "public"."submissions" ("execution_time", "point");

CREATE UNIQUE INDEX IF NOT EXISTS "submissions_id_epoch_second_unique" ON "public"."submissions" ("id", "epoch_second");

CREATE INDEX IF NOT EXISTS "submissions_language_epoch_second_index" ON "public"."submissions" ("language", "epoch_second");

CREATE INDEX IF NOT EXISTS "submissions_language_execution_time_index" ON "public"."submissions" ("language", "execution_time");

CREATE INDEX IF NOT EXISTS "submissions_language_length_index" ON "public"."submissions" ("language", "length");

CREATE INDEX IF NOT EXISTS "submissions_language_point_index" ON "public"."submissions" ("language", "point");

CREATE INDEX IF NOT EXISTS "submissions_length_epoch_second_index" ON "public"."submissions" ("length", "epoch_second");

CREATE INDEX IF NOT EXISTS "submissions_length_execution_time_index" ON "public"."submissions" ("length", "execution_time");

CREATE INDEX IF NOT EXISTS "submissions_length_point_index" ON "public"."submissions" ("length", "point");

CREATE INDEX IF NOT EXISTS "submissions_point_epoch_second_index" ON "public"."submissions" ("point", "epoch_second");

CREATE INDEX IF NOT EXISTS "submissions_point_execution_time_index" ON "public"."submissions" ("point", "execution_time");

CREATE INDEX IF NOT EXISTS "submissions_point_length_index" ON "public"."submissions" ("point", "length");

CREATE INDEX IF NOT EXISTS "submissions_problem_id_epoch_second_index" ON "public"."submissions" ("problem_id", "epoch_second");

CREATE INDEX IF NOT EXISTS "submissions_problem_id_execution_time_index" ON "public"."submissions" ("problem_id", "execution_time");

CREATE INDEX IF NOT EXISTS "submissions_problem_id_length_index" ON "public"."submissions" ("problem_id", "length");

CREATE INDEX IF NOT EXISTS "submissions_problem_id_point_index" ON "public"."submissions" ("problem_id", "point");

CREATE INDEX IF NOT EXISTS "submissions_result_epoch_second_index" ON "public"."submissions" ("result", "epoch_second");

CREATE INDEX IF NOT EXISTS "submissions_result_execution_time_index" ON "public"."submissions" ("result", "execution_time");

CREATE INDEX IF NOT EXISTS "submissions_result_length_index" ON "public"."submissions" ("result", "length");

CREATE INDEX IF NOT EXISTS "submissions_result_point_index" ON "public"."submissions" ("result", "point");

CREATE INDEX IF NOT EXISTS "submissions_updated_at_index" ON "public"."submissions" ("epoch_second", "updated_at");

CREATE INDEX IF NOT EXISTS "submissions_user_id_epoch_second_index" ON "public"."submissions" ("user_id", "epoch_second");

CREATE INDEX IF NOT EXISTS "submissions_user_id_execution_time_index" ON "public"."submissions" ("user_id", "execution_time");

CREATE INDEX IF NOT EXISTS "submissions_user_id_length_index" ON "public"."submissions" ("user_id", "length");

CREATE INDEX IF NOT EXISTS "submissions_user_id_point_index" ON "public"."submissions" ("user_id", "point");

CREATE TABLE IF NOT EXISTS submissions_0_1304175600 partition of submissions FOR
VALUES
FROM
    (0) TO (1304175600);

CREATE TABLE IF NOT EXISTS submissions_1304175600_1321455600 partition of submissions FOR
VALUES
FROM
    (1304175600) TO (1321455600);

CREATE TABLE IF NOT EXISTS submissions_1321455600_1338735600 partition of submissions FOR
VALUES
FROM
    (1321455600) TO (1338735600);

CREATE TABLE IF NOT EXISTS submissions_1338735600_1356015600 partition of submissions FOR
VALUES
FROM
    (1338735600) TO (1356015600);

CREATE TABLE IF NOT EXISTS submissions_1356015600_1373295600 partition of submissions FOR
VALUES
FROM
    (1356015600) TO (1373295600);

CREATE TABLE IF NOT EXISTS submissions_1373295600_1390575600 partition of submissions FOR
VALUES
FROM
    (1373295600) TO (1390575600);

CREATE TABLE IF NOT EXISTS submissions_1390575600_1407855600 partition of submissions FOR
VALUES
FROM
    (1390575600) TO (1407855600);

CREATE TABLE IF NOT EXISTS submissions_1407855600_1425135600 partition of submissions FOR
VALUES
FROM
    (1407855600) TO (1425135600);

CREATE TABLE IF NOT EXISTS submissions_1425135600_1442415600 partition of submissions FOR
VALUES
FROM
    (1425135600) TO (1442415600);

CREATE TABLE IF NOT EXISTS submissions_1442415600_1459695600 partition of submissions FOR
VALUES
FROM
    (1442415600) TO (1459695600);

CREATE TABLE IF NOT EXISTS submissions_1459695600_1476975600 partition of submissions FOR
VALUES
FROM
    (1459695600) TO (1476975600);

CREATE TABLE IF NOT EXISTS submissions_1476975600_1494255600 partition of submissions FOR
VALUES
FROM
    (1476975600) TO (1494255600);

CREATE TABLE IF NOT EXISTS submissions_1494255600_1511535600 partition of submissions FOR
VALUES
FROM
    (1494255600) TO (1511535600);

CREATE TABLE IF NOT EXISTS submissions_1511535600_1528815600 partition of submissions FOR
VALUES
FROM
    (1511535600) TO (1528815600);

CREATE TABLE IF NOT EXISTS submissions_1528815600_1546095600 partition of submissions FOR
VALUES
FROM
    (1528815600) TO (1546095600);

CREATE TABLE IF NOT EXISTS submissions_1546095600_1563375600 partition of submissions FOR
VALUES
FROM
    (1546095600) TO (1563375600);

CREATE TABLE IF NOT EXISTS submissions_1563375600_1580655600 partition of submissions FOR
VALUES
FROM
    (1563375600) TO (1580655600);

CREATE TABLE IF NOT EXISTS submissions_1580655600_1597935600 partition of submissions FOR
VALUES
FROM
    (1580655600) TO (1597935600);

CREATE TABLE IF NOT EXISTS submissions_1597935600_1615215600 partition of submissions FOR
VALUES
FROM
    (1597935600) TO (1615215600);

CREATE TABLE IF NOT EXISTS submissions_1615215600_1632495600 partition of submissions FOR
VALUES
FROM
    (1615215600) TO (1632495600);

CREATE TABLE IF NOT EXISTS submissions_1632495600_1649775600 partition of submissions FOR
VALUES
FROM
    (1632495600) TO (1649775600);

CREATE TABLE IF NOT EXISTS submissions_1649775600_1667055600 partition of submissions FOR
VALUES
FROM
    (1649775600) TO (1667055600);

CREATE TABLE IF NOT EXISTS submissions_1667055600_1684335600 partition of submissions FOR
VALUES
FROM
    (1667055600) TO (1684335600);

CREATE TABLE IF NOT EXISTS submissions_1684335600_1701615600 partition of submissions FOR
VALUES
FROM
    (1684335600) TO (1701615600);

CREATE TABLE IF NOT EXISTS submissions_1701615600_1718895600 partition of submissions FOR
VALUES
FROM
    (1701615600) TO (1718895600);

CREATE TABLE IF NOT EXISTS submissions_1718895600_1736175600 partition of submissions FOR
VALUES
FROM
    (1718895600) TO (1736175600);

CREATE TABLE IF NOT EXISTS submissions_1736175600_1753455600 partition of submissions FOR
VALUES
FROM
    (1736175600) TO (1753455600);

CREATE TABLE IF NOT EXISTS submissions_1753455600_1770735600 partition of submissions FOR
VALUES
FROM
    (1753455600) TO (1770735600);

-- migrate:down
DROP TABLE IF EXISTS "public"."submissions";

DROP TABLE IF EXISTS submissions_0_1304175600;

DROP TABLE IF EXISTS submissions_1304175600_1321455600;

DROP TABLE IF EXISTS submissions_1321455600_1338735600;

DROP TABLE IF EXISTS submissions_1338735600_1356015600;

DROP TABLE IF EXISTS submissions_1356015600_1373295600;

DROP TABLE IF EXISTS submissions_1373295600_1390575600;

DROP TABLE IF EXISTS submissions_1390575600_1407855600;

DROP TABLE IF EXISTS submissions_1407855600_1425135600;

DROP TABLE IF EXISTS submissions_1425135600_1442415600;

DROP TABLE IF EXISTS submissions_1442415600_1459695600;

DROP TABLE IF EXISTS submissions_1459695600_1476975600;

DROP TABLE IF EXISTS submissions_1476975600_1494255600;

DROP TABLE IF EXISTS submissions_1494255600_1511535600;

DROP TABLE IF EXISTS submissions_1511535600_1528815600;

DROP TABLE IF EXISTS submissions_1528815600_1546095600;

DROP TABLE IF EXISTS submissions_1546095600_1563375600;

DROP TABLE IF EXISTS submissions_1563375600_1580655600;

DROP TABLE IF EXISTS submissions_1580655600_1597935600;

DROP TABLE IF EXISTS submissions_1597935600_1615215600;

DROP TABLE IF EXISTS submissions_1615215600_1632495600;

DROP TABLE IF EXISTS submissions_1632495600_1649775600;

DROP TABLE IF EXISTS submissions_1649775600_1667055600;

DROP TABLE IF EXISTS submissions_1667055600_1684335600;

DROP TABLE IF EXISTS submissions_1684335600_1701615600;

DROP TABLE IF EXISTS submissions_1701615600_1718895600;

DROP TABLE IF EXISTS submissions_1718895600_1736175600;

DROP TABLE IF EXISTS submissions_1736175600_1753455600;

DROP TABLE IF EXISTS submissions_1753455600_1770735600;
