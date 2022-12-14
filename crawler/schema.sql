CREATE TABLE IF NOT EXISTS contest (
    id VARCHAR(255) PRIMARY KEY,
    start_epoch_second BIGINT NOT NULL,
    duration_second BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    rate_change VARCHAR(255) NOT NULL,
    category VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS problem (
    id VARCHAR(255) PRIMARY KEY,
    contest_id VARCHAR(255) NOT NULL,
    problem_index VARCHAR(20) NOT NULL,
    name VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    url VARCHAR(255) NOT NULL,
    html TEXT NOT NULL
);