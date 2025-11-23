-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_users_team_active ON users(team_name, is_active);
CREATE INDEX IF NOT EXISTS idx_prs_status ON pull_requests(status);
CREATE INDEX IF NOT EXISTS idx_pr_reviewers_user ON pr_reviewers(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_users_team_active;
DROP INDEX IF EXISTS idx_prs_status;
DROP INDEX IF EXISTS idx_pr_reviewers_user;
-- +goose StatementEnd
