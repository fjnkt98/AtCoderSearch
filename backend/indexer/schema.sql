CREATE TABLE IF NOT EXISTS contests (
    id TEXT PRIMARY KEY,
    start_epoch_second BIGINT NOT NULL,
    duration_second BIGINT NOT NULL,
    title TEXT NOT NULL,
    rate_change TEXT NOT NULL,
    category TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS problems (
    id TEXT PRIMARY KEY,
    contest_id TEXT NOT NULL REFERENCES contests (id) ON DELETE CASCADE,
    problem_index TEXT NOT NULL,
    name TEXT NOT NULL,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    html TEXT NOT NULL,
    difficulty INTEGER NOT NULL
);

CREATE INDEX contest_id_index ON problems (contest_id);