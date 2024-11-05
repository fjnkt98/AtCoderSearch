-- migrate:up
CREATE TABLE IF NOT EXISTS "public"."languages" ("language" TEXT NOT NULL, "group" TEXT NULL, PRIMARY KEY ("language"));

CREATE INDEX IF NOT EXISTS "languages_group_index" ON "public"."languages" ("group");

-- migrate:down
DROP TABLE IF EXISTS "public"."languages";
