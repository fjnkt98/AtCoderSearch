// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
)

const createBatchHistory = `-- name: CreateBatchHistory :one
INSERT INTO
    "batch_history" ("name", "started_at", "options")
VALUES
    ($1, NOW(), $2)
RETURNING
    "id"
`

type CreateBatchHistoryParams struct {
	Name    string `db:"name"`
	Options []byte `db:"options"`
}

func (q *Queries) CreateBatchHistory(ctx context.Context, arg CreateBatchHistoryParams) (int64, error) {
	row := q.db.QueryRow(ctx, createBatchHistory, arg.Name, arg.Options)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const createCrawlHistory = `-- name: CreateCrawlHistory :one
INSERT INTO
    "submission_crawl_history" ("started_at", "contest_id")
VALUES
    ($1, $2)
RETURNING
    "started_at",
    "contest_id"
`

type CreateCrawlHistoryParams struct {
	StartedAt int64  `db:"started_at"`
	ContestID string `db:"contest_id"`
}

type CreateCrawlHistoryRow struct {
	StartedAt int64  `db:"started_at"`
	ContestID string `db:"contest_id"`
}

func (q *Queries) CreateCrawlHistory(ctx context.Context, arg CreateCrawlHistoryParams) (CreateCrawlHistoryRow, error) {
	row := q.db.QueryRow(ctx, createCrawlHistory, arg.StartedAt, arg.ContestID)
	var i CreateCrawlHistoryRow
	err := row.Scan(&i.StartedAt, &i.ContestID)
	return i, err
}

const fetchCategories = `-- name: FetchCategories :many
SELECT DISTINCT
    "category"
FROM
    "contests"
ORDER BY
    "category" ASC
`

func (q *Queries) FetchCategories(ctx context.Context) ([]string, error) {
	rows, err := q.db.Query(ctx, fetchCategories)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			return nil, err
		}
		items = append(items, category)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const fetchContestIDs = `-- name: FetchContestIDs :many
SELECT
    "contest_id"
FROM
    "contests"
ORDER BY
    "start_epoch_second" DESC
`

func (q *Queries) FetchContestIDs(ctx context.Context) ([]string, error) {
	rows, err := q.db.Query(ctx, fetchContestIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var contest_id string
		if err := rows.Scan(&contest_id); err != nil {
			return nil, err
		}
		items = append(items, contest_id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const fetchContestIDsByCategory = `-- name: FetchContestIDsByCategory :many
SELECT
    "contest_id"
FROM
    "contests"
WHERE
    "category" = ANY ($1::TEXT[])
ORDER BY
    "start_epoch_second" DESC
`

func (q *Queries) FetchContestIDsByCategory(ctx context.Context, category []string) ([]string, error) {
	rows, err := q.db.Query(ctx, fetchContestIDsByCategory, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var contest_id string
		if err := rows.Scan(&contest_id); err != nil {
			return nil, err
		}
		items = append(items, contest_id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const fetchLanguages = `-- name: FetchLanguages :one
WITH
    "l" AS (
        SELECT
            "group",
            ARRAY_AGG(
                "language"
                ORDER BY
                    "language" ASC
            ) AS "language"
        FROM
            "languages"
        GROUP BY
            "group"
        ORDER BY
            "group" ASC
    )
SELECT
    JSON_AGG("l") AS "data"
FROM
    "l"
`

func (q *Queries) FetchLanguages(ctx context.Context) ([]byte, error) {
	row := q.db.QueryRow(ctx, fetchLanguages)
	var data []byte
	err := row.Scan(&data)
	return data, err
}

const fetchLanguagesByGroup = `-- name: FetchLanguagesByGroup :one
WITH
    "l" AS (
        SELECT
            "group",
            ARRAY_AGG(
                "language"
                ORDER BY
                    "language" ASC
            ) AS "language"
        FROM
            "languages"
        WHERE
            "group" = ANY ($1::TEXT[])
        GROUP BY
            "group"
        ORDER BY
            "group" ASC
    )
SELECT
    JSON_AGG("l") AS "data"
FROM
    "l"
`

func (q *Queries) FetchLanguagesByGroup(ctx context.Context, groups []string) ([]byte, error) {
	row := q.db.QueryRow(ctx, fetchLanguagesByGroup, groups)
	var data []byte
	err := row.Scan(&data)
	return data, err
}

const fetchLatestBatchHistory = `-- name: FetchLatestBatchHistory :one
SELECT
    "id",
    "started_at",
    "finished_at"
FROM
    "batch_history"
WHERE
    "name" = $1
    AND "status" = 'finished'
ORDER BY
    "started_at" DESC
LIMIT
    1
`

type FetchLatestBatchHistoryRow struct {
	ID         int64      `db:"id"`
	StartedAt  time.Time  `db:"started_at"`
	FinishedAt *time.Time `db:"finished_at"`
}

func (q *Queries) FetchLatestBatchHistory(ctx context.Context, name string) (FetchLatestBatchHistoryRow, error) {
	row := q.db.QueryRow(ctx, fetchLatestBatchHistory, name)
	var i FetchLatestBatchHistoryRow
	err := row.Scan(&i.ID, &i.StartedAt, &i.FinishedAt)
	return i, err
}

const fetchLatestCrawlHistory = `-- name: FetchLatestCrawlHistory :one
SELECT
    "started_at"
FROM
    "submission_crawl_history"
WHERE
    "contest_id" = $1
ORDER BY
    "started_at" DESC
LIMIT
    1
`

func (q *Queries) FetchLatestCrawlHistory(ctx context.Context, contestID string) (int64, error) {
	row := q.db.QueryRow(ctx, fetchLatestCrawlHistory, contestID)
	var started_at int64
	err := row.Scan(&started_at)
	return started_at, err
}

const fetchProblemIDs = `-- name: FetchProblemIDs :many
SELECT
    "problem_id"
FROM
    "problems"
ORDER BY
    "problem_id" ASC
`

func (q *Queries) FetchProblemIDs(ctx context.Context) ([]string, error) {
	rows, err := q.db.Query(ctx, fetchProblemIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var problem_id string
		if err := rows.Scan(&problem_id); err != nil {
			return nil, err
		}
		items = append(items, problem_id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const fetchProblemIDsByCategory = `-- name: FetchProblemIDsByCategory :many
SELECT
    "problem_id"
FROM
    "problems"
    LEFT JOIN "contests" USING ("contest_id")
WHERE
    "category" = ANY ($1::TEXT[])
ORDER BY
    "problem_id" ASC
`

func (q *Queries) FetchProblemIDsByCategory(ctx context.Context, category []string) ([]string, error) {
	rows, err := q.db.Query(ctx, fetchProblemIDsByCategory, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var problem_id string
		if err := rows.Scan(&problem_id); err != nil {
			return nil, err
		}
		items = append(items, problem_id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const fetchProblemIDsByContestID = `-- name: FetchProblemIDsByContestID :many
SELECT
    "problem_id"
FROM
    "problems"
WHERE
    "contest_id" = ANY ($1::TEXT[])
ORDER BY
    "problem_id" ASC
`

func (q *Queries) FetchProblemIDsByContestID(ctx context.Context, contestID []string) ([]string, error) {
	rows, err := q.db.Query(ctx, fetchProblemIDsByContestID, contestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var problem_id string
		if err := rows.Scan(&problem_id); err != nil {
			return nil, err
		}
		items = append(items, problem_id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const fetchRatingByUserID = `-- name: FetchRatingByUserID :one
SELECT
    "rating"
FROM
    "users"
WHERE
    "user_id" = $1
LIMIT
    1
`

func (q *Queries) FetchRatingByUserID(ctx context.Context, userID string) (int32, error) {
	row := q.db.QueryRow(ctx, fetchRatingByUserID, userID)
	var rating int32
	err := row.Scan(&rating)
	return rating, err
}

const insertProblem = `-- name: InsertProblem :execresult
INSERT INTO
    "problems" (
        "problem_id",
        "contest_id",
        "problem_index",
        "name",
        "title",
        "url",
        "html",
        "updated_at"
    )
VALUES
    ($1, $2, $3, $4, $5, $6, $7, NOW())
ON CONFLICT ("problem_id") DO
UPDATE
SET
    "contest_id" = EXCLUDED."contest_id",
    "problem_index" = EXCLUDED."problem_index",
    "name" = EXCLUDED."name",
    "title" = EXCLUDED."title",
    "url" = EXCLUDED."url",
    "html" = EXCLUDED."html",
    "updated_at" = NOW()
`

type InsertProblemParams struct {
	ProblemID    string `bulk:"unique" db:"problem_id"`
	ContestID    string `db:"contest_id"`
	ProblemIndex string `db:"problem_index"`
	Name         string `db:"name"`
	Title        string `db:"title"`
	Url          string `db:"url"`
	Html         string `db:"html"`
}

func (q *Queries) InsertProblem(ctx context.Context, arg InsertProblemParams) (pgconn.CommandTag, error) {
	return q.db.Exec(ctx, insertProblem,
		arg.ProblemID,
		arg.ContestID,
		arg.ProblemIndex,
		arg.Name,
		arg.Title,
		arg.Url,
		arg.Html,
	)
}

const updateBatchHistory = `-- name: UpdateBatchHistory :exec
UPDATE "batch_history"
SET
    "finished_at" = NOW(),
    "status" = $1
WHERE
    "id" = $2
`

type UpdateBatchHistoryParams struct {
	Status string `db:"status"`
	ID     int64  `db:"id"`
}

func (q *Queries) UpdateBatchHistory(ctx context.Context, arg UpdateBatchHistoryParams) error {
	_, err := q.db.Exec(ctx, updateBatchHistory, arg.Status, arg.ID)
	return err
}

const updateLanguages = `-- name: UpdateLanguages :execresult
WITH
    "l" AS (
        SELECT DISTINCT
            "language"
        FROM
            "submissions"
    )
INSERT INTO
    "languages" ("language", "group")
SELECT
    "language",
    CASE
        WHEN "language" = 'Ada2012 (GNAT 9.2.1)' THEN 'Ada'
        WHEN "language" = 'Awk (GNU Awk 4.1.4)' THEN 'AWK'
        WHEN "language" = 'AWK (GNU Awk 5.0.1)' THEN 'AWK'
        WHEN "language" = 'Awk (mawk 1.3.3)' THEN 'AWK'
        WHEN "language" = 'Bash (4.2.25)' THEN 'Bash'
        WHEN "language" = 'Bash (5.0.11)' THEN 'Bash'
        WHEN "language" = 'Bash (bash 5.2.2)' THEN 'Bash'
        WHEN "language" = 'Bash (GNU bash v4.3.11)' THEN 'Bash'
        WHEN "language" = 'bc (1.07.1)' THEN 'bc'
        WHEN "language" = 'Brainfuck (bf 20041219)' THEN 'Brainfuck'
        WHEN "language" = 'C# 11.0 AOT (.NET 7.0.7)' THEN 'C#'
        WHEN "language" = 'C# 11.0 (.NET 7.0.7)' THEN 'C#'
        WHEN "language" = 'C++11 (Clang++ 3.4)' THEN 'C++'
        WHEN "language" = 'C++11 (GCC 4.8.1)' THEN 'C++'
        WHEN "language" = 'C++11 (GCC 4.9.2)' THEN 'C++'
        WHEN "language" = 'C++14 (Clang++ 3.4)' THEN 'C++'
        WHEN "language" = 'C++14 (Clang 3.8.0)' THEN 'C++'
        WHEN "language" = 'C++14 (GCC 5.4.1)' THEN 'C++'
        WHEN "language" = 'C++ 17 (Clang 16.0.5)' THEN 'C++'
        WHEN "language" = 'C++ 17 (gcc 12.2)' THEN 'C++'
        WHEN "language" = 'C++ 20 (Clang 16.0.5)' THEN 'C++'
        WHEN "language" = 'C++ 20 (gcc 12.2)' THEN 'C++'
        WHEN "language" = 'C++ 23 (Clang 16.0.5)' THEN 'C++'
        WHEN "language" = 'C++ 23 (gcc 12.2)' THEN 'C++'
        WHEN "language" = 'Carp (Carp 0.5.5)' THEN 'Carp'
        WHEN "language" = 'C (Clang 10.0.0)' THEN 'C'
        WHEN "language" = 'C++ (Clang 10.0.0)' THEN 'C++'
        WHEN "language" = 'C++ (Clang 10.0.0 with AC Library)' THEN 'C++'
        WHEN "language" = 'C++ (Clang 10.0.0 with AC Library v1.0)' THEN 'C++'
        WHEN "language" = 'C++ (Clang 10.0.0 with AC Library v1.1)' THEN 'C++'
        WHEN "language" = 'C (Clang 3.4)' THEN 'C'
        WHEN "language" = 'C++ (Clang++ 3.4)' THEN 'C++'
        WHEN "language" = 'C (Clang 3.8.0)' THEN 'C'
        WHEN "language" = 'C++ (Clang 3.8.0)' THEN 'C++'
        WHEN "language" = 'Ceylon (1.2.1)' THEN 'Other'
        WHEN "language" = 'C++ (G++ 4.6.4)' THEN 'C++'
        WHEN "language" = 'C (gcc 12.2.0)' THEN 'C'
        WHEN "language" = 'C (GCC 4.4.7)' THEN 'C'
        WHEN "language" = 'C++ (GCC 4.4.7)' THEN 'C++'
        WHEN "language" = 'C (GCC 4.6.4)' THEN 'C'
        WHEN "language" = 'C (GCC 4.9.2)' THEN 'C'
        WHEN "language" = 'C++ (GCC 4.9.2)' THEN 'C++'
        WHEN "language" = 'C (GCC 5.4.1)' THEN 'C'
        WHEN "language" = 'C++ (GCC 5.4.1)' THEN 'C++'
        WHEN "language" = 'C (GCC 9.2.1)' THEN 'C'
        WHEN "language" = 'C++ (GCC 9.2.1)' THEN 'C++'
        WHEN "language" = 'C++ (GCC 9.2.1 with AC Library)' THEN 'C++'
        WHEN "language" = 'C++ (GCC 9.2.1 with AC Library v1.0)' THEN 'C++'
        WHEN "language" = 'C++ (GCC 9.2.1 with AC Library v1.1)' THEN 'C++'
        WHEN "language" = 'Clojure (1.10.1.536)' THEN 'Clojure'
        WHEN "language" = 'Clojure (1.1.0 + OpenJDK 1.7)' THEN 'Clojure'
        WHEN "language" = 'Clojure (1.8.0)' THEN 'Clojure'
        WHEN "language" = 'Clojure (babashka 1.3.181)' THEN 'Clojure'
        WHEN "language" = 'Clojure (clojure 1.11.1)' THEN 'Clojure'
        WHEN "language" = 'C# (Mono 2.10.8.1)' THEN 'C#'
        WHEN "language" = 'C# (Mono 3.2.1.0)' THEN 'C#'
        WHEN "language" = 'C# (Mono 4.6.2.0)' THEN 'C#'
        WHEN "language" = 'C# (Mono-csc 3.5.0)' THEN 'C#'
        WHEN "language" = 'C# (Mono-mcs 6.8.0.105)' THEN 'C#'
        WHEN "language" = 'C# (.NET Core 3.1.201)' THEN 'C#'
        WHEN "language" = 'COBOL - Fixed (OpenCOBOL 1.1.0)' THEN 'COBOL'
        WHEN "language" = 'COBOL (Free) (GnuCOBOL 3.1.2)' THEN 'COBOL'
        WHEN "language" = 'COBOL - Free (OpenCOBOL 1.1.0)' THEN 'COBOL'
        WHEN "language" = 'COBOL (GnuCOBOL(Fixed) 3.1.2)' THEN 'COBOL'
        WHEN "language" = 'Common Lisp (SBCL 1.0.55.0)' THEN 'Common Lisp'
        WHEN "language" = 'Common Lisp (SBCL 1.1.14)' THEN 'Common Lisp'
        WHEN "language" = 'Common Lisp (SBCL 2.0.3)' THEN 'Common Lisp'
        WHEN "language" = 'Common Lisp (SBCL 2.3.6)' THEN 'Common Lisp'
        WHEN "language" = 'Crystal (0.20.5)' THEN 'Crystal'
        WHEN "language" = 'Crystal (0.33.0)' THEN 'Crystal'
        WHEN "language" = 'Crystal (Crystal 1.9.1)' THEN 'Crystal'
        WHEN "language" = 'Cyber (Cyber v0.2-Latest)' THEN 'Cyber'
        WHEN "language" = 'Cython (0.29.16)' THEN 'Python'
        WHEN "language" = 'Dart (2.7.2)' THEN 'Dart'
        WHEN "language" = 'Dart (Dart 3.0.5)' THEN 'Dart'
        WHEN "language" = 'Dash (0.5.8)' THEN 'Other'
        WHEN "language" = 'dc (1.4.1)' THEN 'dc'
        WHEN "language" = 'D (DMD 2.060)' THEN 'D'
        WHEN "language" = 'D (DMD 2.066.1)' THEN 'D'
        WHEN "language" = 'D (DMD 2.091.0)' THEN 'D'
        WHEN "language" = 'D (DMD 2.104.0)' THEN 'D'
        WHEN "language" = 'D (DMD64 v2.070.1)' THEN 'D'
        WHEN "language" = 'D (GDC 4.9.4)' THEN 'D'
        WHEN "language" = 'D (GDC 9.2.1)' THEN 'D'
        WHEN "language" = 'D (LDC 0.17.0)' THEN 'D'
        WHEN "language" = 'D (LDC 1.20.1)' THEN 'D'
        WHEN "language" = 'D (LDC 1.32.2)' THEN 'D'
        WHEN "language" = 'ECLiPSe (ECLiPSe 7.1_13)' THEN 'ECLiPSe'
        WHEN "language" = 'Elixir (1.10.2)' THEN 'Elixir'
        WHEN "language" = 'Elixir (Elixir 1.15.2)' THEN 'Elixir'
        WHEN "language" = 'Emacs Lisp (Native Compile) (GNU Emacs 28.2)' THEN 'Emacs Lisp'
        WHEN "language" = 'Erlang (22.3)' THEN 'Erlang'
        WHEN "language" = 'F# 7.0 (.NET 7.0.7)' THEN 'F#'
        WHEN "language" = '><> (fishr 0.1.0)' THEN '><>'
        WHEN "language" = 'F# (Mono 10.2.3)' THEN 'F#'
        WHEN "language" = 'F# (Mono 4.0)' THEN 'F#'
        WHEN "language" = 'F# (.NET Core 3.1.201)' THEN 'F#'
        WHEN "language" = 'Forth (gforth 0.7.3)' THEN 'Forth'
        WHEN "language" = 'Fortran (gfortran 12.2)' THEN 'Fortran'
        WHEN "language" = 'Fortran (gfortran v4.8.4)' THEN 'Fortran'
        WHEN "language" = 'Fortran (GNU Fortran 9.2.1)' THEN 'Fortran'
        WHEN "language" = 'Fortran(GNU Fortran 9.2.1)' THEN 'Fortran'
        WHEN "language" = 'Go (1.14.1)' THEN 'Go'
        WHEN "language" = 'Go (1.4.1)' THEN 'Go'
        WHEN "language" = 'Go (1.6)' THEN 'Go'
        WHEN "language" = 'Go (go 1.20.6)' THEN 'Go'
        WHEN "language" = 'Haskell (GHC 7.10.3)' THEN 'Haskell'
        WHEN "language" = 'Haskell (GHC 7.4.1)' THEN 'Haskell'
        WHEN "language" = 'Haskell (GHC 8.8.3)' THEN 'Haskell'
        WHEN "language" = 'Haskell (GHC 9.4.5)' THEN 'Haskell'
        WHEN "language" = 'Haskell (Haskell Platform 2014.2.0.0)' THEN 'Haskell'
        WHEN "language" = 'Haxe (4.0.3); Java' THEN 'Haxe'
        WHEN "language" = 'Haxe (4.0.3); js' THEN 'Haxe'
        WHEN "language" = 'IOI-Style C++ (GCC 5.4.1)' THEN 'C++'
        WHEN "language" = 'Java7 (OpenJDK 1.7.0)' THEN 'Java'
        WHEN "language" = 'Java8 (OpenJDK 1.8.0)' THEN 'Java'
        WHEN "language" = 'Java (OpenJDK 11.0.6)' THEN 'Java'
        WHEN "language" = 'Java (OpenJDK 17)' THEN 'Java'
        WHEN "language" = 'Java (OpenJDK 1.7.0)' THEN 'Java'
        WHEN "language" = 'Java (OpenJDK 1.8.0)' THEN 'Java'
        WHEN "language" = 'JavaScript (Deno 1.35.1)' THEN 'JavaScript'
        WHEN "language" = 'JavaScript (Node.js 0.6.12)' THEN 'JavaScript'
        WHEN "language" = 'JavaScript (Node.js 12.16.1)' THEN 'JavaScript'
        WHEN "language" = 'JavaScript (Node.js 18.16.1)' THEN 'JavaScript'
        WHEN "language" = 'JavaScript (Node.js v0.10.36)' THEN 'JavaScript'
        WHEN "language" = 'JavaScript (node.js v5.12)' THEN 'JavaScript'
        WHEN "language" = 'jq (jq 1.6)' THEN 'jq'
        WHEN "language" = 'Julia (0.5.0)' THEN 'Julia'
        WHEN "language" = 'Julia (1.4.0)' THEN 'Julia'
        WHEN "language" = 'Julia (Julia 1.9.2)' THEN 'Julia'
        WHEN "language" = 'Koka (koka 2.4.0)' THEN 'Koka'
        WHEN "language" = 'Kotlin (1.0.0)' THEN 'Kotlin'
        WHEN "language" = 'Kotlin (1.3.71)' THEN 'Kotlin'
        WHEN "language" = 'Kotlin (Kotlin/JVM 1.8.20)' THEN 'Kotlin'
        WHEN "language" = 'LLVM IR (Clang 16.0.5)' THEN 'LLVM IR'
        WHEN "language" = 'Lua (5.3.2)' THEN 'Lua'
        WHEN "language" = 'LuaJIT (2.0.4)' THEN 'Lua'
        WHEN "language" = 'Lua (Lua 5.3.5)' THEN 'Lua'
        WHEN "language" = 'Lua (Lua 5.4.6)' THEN 'Lua'
        WHEN "language" = 'Lua (LuaJIT 2.1.0)' THEN 'Lua'
        WHEN "language" = 'Lua (LuaJIT 2.1.0-beta3)' THEN 'Lua'
        WHEN "language" = 'MoonScript (0.5.0)' THEN 'Other'
        WHEN "language" = 'Nibbles (literate form) (nibbles 1.01)' THEN 'Nibbles'
        WHEN "language" = 'Nim (0.13.0)' THEN 'Nim'
        WHEN "language" = 'Nim (1.0.6)' THEN 'Nim'
        WHEN "language" = 'Nim (Nim 1.6.14)' THEN 'Nim'
        WHEN "language" = 'Objective-C (Clang 10.0.0)' THEN 'Objective-C'
        WHEN "language" = 'Objective-C (Clang3.8.0)' THEN 'Objective-C'
        WHEN "language" = 'Objective-C (GCC 5.3.0)' THEN 'Objective-C'
        WHEN "language" = 'OCaml (3.12.1)' THEN 'Ocaml'
        WHEN "language" = 'OCaml (4.02.1)' THEN 'Ocaml'
        WHEN "language" = 'OCaml (4.02.3)' THEN 'Ocaml'
        WHEN "language" = 'OCaml (4.10.0)' THEN 'Ocaml'
        WHEN "language" = 'OCaml (ocamlopt 5.0.0)' THEN 'Ocaml'
        WHEN "language" = 'Octave (4.0.2)' THEN 'Octave'
        WHEN "language" = 'Octave (5.2.0)' THEN 'Octave'
        WHEN "language" = 'Octave (GNU Octave 8.2.0)' THEN 'Octave'
        WHEN "language" = 'Pascal (fpc 2.4.4)' THEN 'Pascal'
        WHEN "language" = 'Pascal (FPC 2.6.2)' THEN 'Pascal'
        WHEN "language" = 'Pascal (FPC 3.0.4)' THEN 'Pascal'
        WHEN "language" = 'Pascal (fpc 3.2.2)' THEN 'Pascal'
        WHEN "language" = 'Perl (5.14.2)' THEN 'Perl'
        WHEN "language" = 'Perl (5.26.1)' THEN 'Perl'
        WHEN "language" = 'Perl6 (rakudo-star 2016.01)' THEN 'Perl'
        WHEN "language" = 'Perl (perl  5.34)' THEN 'Perl'
        WHEN "language" = 'Perl (v5.18.2)' THEN 'Perl'
        WHEN "language" = 'PHP (5.6.30)' THEN 'PHP'
        WHEN "language" = 'PHP (7.4.4)' THEN 'PHP'
        WHEN "language" = 'PHP7 (7.0.15)' THEN 'PHP'
        WHEN "language" = 'PHP (PHP 5.3.10)' THEN 'PHP'
        WHEN "language" = 'PHP (PHP 5.5.21)' THEN 'PHP'
        WHEN "language" = 'PHP (php 8.2.8)' THEN 'PHP'
        WHEN "language" = 'PowerShell (PowerShell 7.3.1)' THEN 'PowerShell'
        WHEN "language" = 'Prolog (SWI-Prolog 8.0.3)' THEN 'Prolog'
        WHEN "language" = 'Prolog (SWI-Prolog 9.0.4)' THEN 'Prolog'
        WHEN "language" = 'PyPy2 (5.6.0)' THEN 'Python'
        WHEN "language" = 'PyPy2 (7.3.0)' THEN 'Python'
        WHEN "language" = 'PyPy3 (2.4.0)' THEN 'Python'
        WHEN "language" = 'PyPy3 (7.3.0)' THEN 'Python'
        WHEN "language" = 'Python2 (2.7.6)' THEN 'Python'
        WHEN "language" = 'Python (2.7.3)' THEN 'Python'
        WHEN "language" = 'Python (2.7.6)' THEN 'Python'
        WHEN "language" = 'Python3 (3.2.3)' THEN 'Python'
        WHEN "language" = 'Python3 (3.4.2)' THEN 'Python'
        WHEN "language" = 'Python3 (3.4.3)' THEN 'Python'
        WHEN "language" = 'Python (3.4.3)' THEN 'Python'
        WHEN "language" = 'Python (3.8.2)' THEN 'Python'
        WHEN "language" = 'Python (CPython 3.11.4)' THEN 'Python'
        WHEN "language" = 'Python (Cython 0.29.34)' THEN 'Python'
        WHEN "language" = 'Python (Mambaforge / CPython 3.10.10)' THEN 'Python'
        WHEN "language" = 'Python (PyPy 3.10-v7.3.12)' THEN 'Python'
        WHEN "language" = 'Racket (7.6)' THEN 'Other'
        WHEN "language" = 'Raku (Rakudo 2020.02.1)' THEN 'Raku'
        WHEN "language" = 'Raku (Rakudo 2023.06)' THEN 'Raku'
        WHEN "language" = 'ReasonML (reason 3.9.0)' THEN 'ReasonML'
        WHEN "language" = 'R (GNU R 4.2.1)' THEN 'R'
        WHEN "language" = 'Ruby (1.9.3)' THEN 'Ruby'
        WHEN "language" = 'Ruby (1.9.3p550)' THEN 'Ruby'
        WHEN "language" = 'Ruby (2.1.5p273)' THEN 'Ruby'
        WHEN "language" = 'Ruby (2.3.3)' THEN 'Ruby'
        WHEN "language" = 'Ruby (2.7.1)' THEN 'Ruby'
        WHEN "language" = 'Ruby (ruby 3.2.2)' THEN 'Ruby'
        WHEN "language" = 'Rust (1.15.1)' THEN 'Rust'
        WHEN "language" = 'Rust (1.42.0)' THEN 'Rust'
        WHEN "language" = 'Rust (rustc 1.70.0)' THEN 'Rust'
        WHEN "language" = 'SageMath (SageMath 9.5)' THEN 'SageMath'
        WHEN "language" = 'Scala (2.11.5)' THEN 'Scala'
        WHEN "language" = 'Scala (2.11.7)' THEN 'Scala'
        WHEN "language" = 'Scala (2.13.1)' THEN 'Scala'
        WHEN "language" = 'Scala (2.9.1)' THEN 'Scala'
        WHEN "language" = 'Scala 3.3.0 (Scala Native 0.4.14)' THEN 'Scala'
        WHEN "language" = 'Scala (Dotty 3.3.0)' THEN 'Scala'
        WHEN "language" = 'Scheme (Gauche 0.9.1)' THEN 'Scheme'
        WHEN "language" = 'Scheme (Gauche 0.9.12)' THEN 'Scheme'
        WHEN "language" = 'Scheme (Gauche 0.9.3.3)' THEN 'Scheme'
        WHEN "language" = 'Scheme (Gauche 0.9.9)' THEN 'Scheme'
        WHEN "language" = 'Scheme (Scheme 9.1)' THEN 'Scheme'
        WHEN "language" = 'Sed (4.4)' THEN 'Sed'
        WHEN "language" = 'Sed (GNU sed 4.2.2)' THEN 'Sed'
        WHEN "language" = 'Sed (GNU sed 4.8)' THEN 'Sed'
        WHEN "language" = 'Seed7 (Seed7 3.2.1)' THEN 'Seed7'
        WHEN "language" = 'Standard ML (MLton 20100608)' THEN 'Other'
        WHEN "language" = 'Standard ML (MLton 20130715)' THEN 'Other'
        WHEN "language" = 'Swift (5.2.1)' THEN 'Swift'
        WHEN "language" = 'Swift (swift-2.2-RELEASE)' THEN 'Swift'
        WHEN "language" = 'Swift (swift 5.8.1)' THEN 'Swift'
        WHEN "language" = 'Text (cat)' THEN 'Text'
        WHEN "language" = 'Text (cat 8.28)' THEN 'Text'
        WHEN "language" = 'Text (cat 8.32)' THEN 'Text'
        WHEN "language" = 'TypeScript (2.1.6)' THEN 'TypeScript'
        WHEN "language" = 'TypeScript (3.8)' THEN 'TypeScript'
        WHEN "language" = 'TypeScript 5.1 (Deno 1.35.1)' THEN 'TypeScript'
        WHEN "language" = 'TypeScript 5.1 (Node.js 18.16.1)' THEN 'TypeScript'
        WHEN "language" = 'Unlambda (0.1.3)' THEN 'Other'
        WHEN "language" = 'Unlambda (2.0.0)' THEN 'Other'
        WHEN "language" = 'Vim (8.2.0460)' THEN 'Vim'
        WHEN "language" = 'Visual Basic (Mono 2.10.8)' THEN 'Visual Basic'
        WHEN "language" = 'Visual Basic (Mono 4.0.1)' THEN 'Visual Basic'
        WHEN "language" = 'Visual Basic (.NET Core 3.1.101)' THEN 'Visual Basic'
        WHEN "language" = 'V (V 0.4)' THEN 'V'
        WHEN "language" = 'Whitespace (whitespacers 1.0.0)' THEN 'Whitespace'
        WHEN "language" = 'Zig (Zig 0.10.1)' THEN 'Zig'
        WHEN "language" = 'Zsh (5.4.2)' THEN 'Zsh'
        WHEN "language" = 'なでしこ (cnako3 3.4.20)' THEN 'なでしこ'
        WHEN "language" = 'プロデル (mono版プロデル 1.9.1182)' THEN 'プロデル'
        WHEN "language" ~* '^><>' THEN '><>'
        WHEN "language" ~* '^awk' THEN 'AWK'
        WHEN "language" ~* '^Ada' THEN 'Ada'
        WHEN "language" ~* '^bash' THEN 'Bash'
        WHEN "language" ~* '^Brainfuck' THEN 'Brainfuck'
        WHEN "language" ~* '^C ' THEN 'C'
        WHEN "language" ~* '^C#' THEN 'C#'
        WHEN "language" ~* '^C\+\+' THEN 'C++'
        WHEN "language" ~* '^COBOL' THEN 'COBOL'
        WHEN "language" ~* '^Carp' THEN 'Carp'
        WHEN "language" ~* '^Clojure' THEN 'Clojure'
        WHEN "language" ~* '^Common Lisp' THEN 'Common Lisp'
        WHEN "language" ~* '^Crystal' THEN 'Crystal'
        WHEN "language" ~* '^Cyber' THEN 'Cyber'
        WHEN "language" ~* '^D ' THEN 'D'
        WHEN "language" ~* '^Dart' THEN 'Dart'
        WHEN "language" ~* '^ECLiPSe' THEN 'ECLiPSe'
        WHEN "language" ~* '^Elixir' THEN 'Elixir'
        WHEN "language" ~* '^Emacs Lisp' THEN 'Emacs Lisp'
        WHEN "language" ~* '^Erlang' THEN 'Erlang'
        WHEN "language" ~* '^F#' THEN 'F#'
        WHEN "language" ~* '^Forth' THEN 'Forth'
        WHEN "language" ~* '^Fortran' THEN 'Fortran'
        WHEN "language" ~* '^Go' THEN 'Go'
        WHEN "language" ~* '^Haskell' THEN 'Haskell'
        WHEN "language" ~* '^Haxe' THEN 'Haxe'
        WHEN "language" ~* '^Java\d{0,2} ' THEN 'Java'
        WHEN "language" ~* '^JavaScript' THEN 'JavaScript'
        WHEN "language" ~* '^Julia' THEN 'Julia'
        WHEN "language" ~* '^Koka' THEN 'Koka'
        WHEN "language" ~* '^Kotlin' THEN 'Kotlin'
        WHEN "language" ~* '^LLVM IR' THEN 'LLVM IR'
        WHEN "language" ~* '^Lua' THEN 'Lua'
        WHEN "language" ~* '^Nibbles' THEN 'Nibbles'
        WHEN "language" ~* '^Nim' THEN 'Nim'
        WHEN "language" ~* '^Objective-C' THEN 'Objective-C'
        WHEN "language" ~* '^Ocaml' THEN 'Ocaml'
        WHEN "language" ~* '^Octave' THEN 'Octave'
        WHEN "language" ~* '^PHP' THEN 'PHP'
        WHEN "language" ~* '^Pascal' THEN 'Pascal'
        WHEN "language" ~* '^Perl' THEN 'Perl'
        WHEN "language" ~* '^PowerShell' THEN 'PowerShell'
        WHEN "language" ~* '^Prolog' THEN 'Prolog'
        WHEN "language" ~* '^([PC]ython|PyPy)' THEN 'Python'
        WHEN "language" ~* '^R ' THEN 'R'
        WHEN "language" ~* '^Raku' THEN 'Raku'
        WHEN "language" ~* '^Ruby' THEN 'Ruby'
        WHEN "language" ~* '^Rust' THEN 'Rust'
        WHEN "language" ~* '^SageMath' THEN 'SageMath'
        WHEN "language" ~* '^Scala' THEN 'Scala'
        WHEN "language" ~* '^Scheme' THEN 'Scheme'
        WHEN "language" ~* '^Sed' THEN 'Sed'
        WHEN "language" ~* '^Swift' THEN 'Swift'
        WHEN "language" ~* '^Text' THEN 'Text'
        WHEN "language" ~* '^TypeScript' THEN 'TypeScript'
        WHEN "language" ~* '^V ' THEN 'V'
        WHEN "language" ~* '^Vim' THEN 'Vim'
        WHEN "language" ~* '^Visual Basic' THEN 'Visual Basic'
        WHEN "language" ~* '^Whitespace' THEN 'Whitespace'
        WHEN "language" ~* '^Zig' THEN 'Zig'
        WHEN "language" ~* '^bc ' THEN 'bc'
        WHEN "language" ~* '^dc ' THEN 'dc'
        WHEN "language" ~* '^jq' THEN 'jq'
        WHEN "language" ~* '^なでしこ' THEN 'なでしこ'
        WHEN "language" ~* '^プロデル' THEN 'プロデル'
        ELSE 'Other'
    END AS "group"
FROM
    "l"
ON CONFLICT ("language") DO
UPDATE
SET
    "group" = EXCLUDED."group"
`

func (q *Queries) UpdateLanguages(ctx context.Context) (pgconn.CommandTag, error) {
	return q.db.Exec(ctx, updateLanguages)
}
