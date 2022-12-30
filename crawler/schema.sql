CREATE TABLE IF NOT EXISTS contests (
    id VARCHAR(255) PRIMARY KEY,
    start_epoch_second BIGINT NOT NULL,
    duration_second BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    rate_change VARCHAR(255) NOT NULL,
    category VARCHAR(255) NOT NULL
);

CREATE INDEX contest_id_index ON contests (id);

CREATE TABLE IF NOT EXISTS problems (
    id VARCHAR(255) PRIMARY KEY,
    contest_id VARCHAR(255) NOT NULL REFERENCES contests(id) ON DELETE CASCADE,
    problem_index VARCHAR(20) NOT NULL,
    name VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL,
    html TEXT NOT NULL,
    difficulty INTEGER NOT NULL
);

CREATE INDEX problem_id_index ON problems (id);