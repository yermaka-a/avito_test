-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS teams (
    team_name VARCHAR(100) PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS users (
    user_id VARCHAR(100) PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    team_name VARCHAR(100) NOT NULL REFERENCES teams(team_name) ON DELETE CASCADE,
    is_active BOOLEAN DEFAULT true
);

CREATE TABLE IF NOT EXISTS pull_requests (
    pull_request_id VARCHAR(100) PRIMARY KEY,
    pull_request_name VARCHAR(200) NOT NULL,
    author_id VARCHAR(100) NOT NULL REFERENCES users(user_id),
    status VARCHAR(20) DEFAULT 'OPEN' CHECK (status IN ('OPEN', 'MERGED')),
    FOREIGN KEY (author_id) REFERENCES users(user_id),
    merged_at TIMESTAMP WITH TIME ZONE NULL
);

CREATE TABLE IF NOT EXISTS pr_reviewers (
    pull_request_id VARCHAR(100) REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE,
    user_id VARCHAR(100) REFERENCES users(user_id) ON DELETE CASCADE,
    PRIMARY KEY (pull_request_id, user_id)
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pr_reviewers;
DROP TABLE IF EXISTS pull_requests;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS teams;
-- +goose StatementEnd
