-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS teams (
    name TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    team_name TEXT NOT NULL REFERENCES teams(name) ON UPDATE CASCADE ON DELETE RESTRICT,
    is_active BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS idx_users_team_name ON users(team_name);
CREATE INDEX IF NOT EXISTS idx_users_team_name_is_active ON users(team_name, is_active);

CREATE TABLE IF NOT EXISTS pull_requests (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    author_id TEXT NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    status TEXT NOT NULL,
    reviewer1_id TEXT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    reviewer2_id TEXT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    merged_at TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS idx_pull_requests_reviewer1_id ON pull_requests(reviewer1_id);
CREATE INDEX IF NOT EXISTS idx_pull_requests_reviewer2_id ON pull_requests(reviewer2_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS pull_requests;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS teams;

-- +goose StatementEnd
