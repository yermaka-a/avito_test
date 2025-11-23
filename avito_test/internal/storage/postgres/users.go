package postgres

import (
	"avito_test/internal/domain/models"
	"avito_test/internal/storage"
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

func (s *Storage) SetIsActive(ctx context.Context, user *models.User) error {

	query := fmt.Sprintf("UPDATE %s SET is_active = $1 WHERE user_id = $2 RETURNING  username, team_name", UsersTable)
	row := s.pool.QueryRow(ctx, query, user.IsActive, user.UserID)
	if err := row.Scan(&user.Username, &user.TeamName); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.log.Debug("user not found for update", slog.String("user_id", user.UserID))
			return storage.ErrUserNotFound
		}

		s.log.Debug(FailedToScanRowMsg, slog.Any("error", err))
		return err
	}

	return nil
}

func (s *Storage) GetReview(ctx context.Context, userID string) ([]*models.PullRequest, error) {

	query := fmt.Sprintf(`
        SELECT p.pull_request_id, p.pull_request_name, p.author_id, p.status 
        FROM %s p
        JOIN %s r ON p.pull_request_id = r.pull_request_id
        WHERE r.user_id = $1`,
		PRTable, ReviewersTable)

	rows, err := s.pool.Query(ctx, query, userID)
	if err != nil {
		s.log.Debug(FailedToStartQuery, slog.Any("error", err))
		return nil, err
	}
	defer rows.Close()

	pullRequests := make([]*models.PullRequest, 0)
	for rows.Next() {
		var pr models.PullRequest
		if err := rows.Scan(&pr.ID, &pr.PullRequestName, &pr.AuthorID, &pr.Status); err != nil {
			s.log.Debug(FailedToScanRowsMsg, slog.Any("error", err))
			return nil, err
		}
		pullRequests = append(pullRequests, &pr)
	}

	if err := rows.Err(); err != nil {
		s.log.Warn("error during rows iteration", slog.Any("error", err))
		return nil, err
	}

	return pullRequests, nil
}
