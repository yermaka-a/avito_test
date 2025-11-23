package postgres

import (
	"avito_test/internal/domain/models"
	"avito_test/internal/storage"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
)

func (s *Storage) CreatePR(ctx context.Context, pr *models.PullRequest) ([]*models.Reviewer, error) {

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		s.log.Debug(CantStartTransacMsg, slog.Any("error", err))
		return nil, storage.ErrCreatePR
	}
	defer tx.Rollback(ctx)
	query := fmt.Sprintf("SELECT TRUE FROM %s WHERE user_id = $1", UsersTable)
	exists := false
	err = tx.QueryRow(ctx, query, pr.AuthorID).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) && !exists {
			return nil, storage.ErrUserNotFound
		}
		s.log.Debug(FailedToScanRowMsg, slog.Any("error", err))
		return nil, err
	}
	query = fmt.Sprintf("INSERT INTO %s (pull_request_id, pull_request_name, author_id, status) VALUES ($1, $2, $3, $4) ON CONFLICT (pull_request_id) DO NOTHING", PRTable)

	tag, err := tx.Exec(ctx, query, pr.ID, pr.PullRequestName, pr.AuthorID, pr.Status)

	if err != nil {
		s.log.Debug("failed to insert pr value", slog.Any("error", err))
		return nil, err
	}

	if tag.RowsAffected() == 0 {
		return nil, storage.ErrPRIDAlreadyExists
	}

	teamName, err := s.teamNameByAuthorId(ctx, tx, pr.AuthorID)
	if err != nil {
		return nil, err
	}

	reviewers, err := s.reviewers(ctx, tx, pr.AuthorID, teamName)
	if err != nil {
		return nil, err
	}
	batch := &pgx.Batch{}
	query = fmt.Sprintf("INSERT INTO %s (pull_request_id, user_id) VALUES ($1, $2)", ReviewersTable)
	for _, r := range reviewers {
		batch.Queue(query, pr.ID, r.ReviewerID)
	}
	br := tx.SendBatch(ctx, batch)
	for range reviewers {
		if _, err := br.Exec(); err != nil {
			s.log.Debug("failed to read batch results", slog.Any("error", err))
			return nil, err
		}
	}
	br.Close()

	err = tx.Commit(ctx)
	if err != nil {
		s.log.Warn(FailedToCommitTransac, slog.Any("error", err))
		return nil, err
	}
	return reviewers, nil
}

func (s *Storage) reviewers(ctx context.Context, tx pgx.Tx, authorID string, teamName string) ([]*models.Reviewer, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id != $1 AND team_name = $2 AND is_active = true", UsersTable)
	rows, err := tx.Query(ctx, query, authorID, teamName)
	if err != nil {
		s.log.Debug(FailedToStartQuery, slog.Any("error", err))
		return nil, err
	}
	defer rows.Close()
	count := 0
	reviewers := make([]*models.Reviewer, 0, 2)

	for rows.Next() {
		if count == 2 {
			break
		}
		var user models.User
		if err := rows.Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, storage.ErrUsersNotFound
			}
			return nil, err
		}
		if count < 2 {
			reviewers = append(reviewers, &models.Reviewer{ReviewerID: user.UserID})
			count += 1
		}
	}
	return reviewers, nil
}

func (s *Storage) teamNameByAuthorId(ctx context.Context, tx pgx.Tx, authorID string) (string, error) {
	query := fmt.Sprintf("SELECT team_name FROM %s WHERE user_id = $1", UsersTable)
	row := tx.QueryRow(ctx, query, authorID)
	var teamName string

	if err := row.Scan(&teamName); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", storage.ErrUserNotFound
		}
		s.log.Debug("author is not found by id", slog.String("id", authorID))
		return "", err
	}
	return teamName, nil

}

func (s *Storage) PRMarkAsMerged(ctx context.Context, prID string) (*models.PullRequest, error) {
	tx, err := s.pool.Begin(ctx)

	if err != nil {
		s.log.Debug(CantStartTransacMsg, slog.Any("error", err))
		return nil, err
	}
	defer tx.Rollback(ctx)

	query := fmt.Sprintf("SELECT * FROM %s WHERE pull_request_id = $1", PRTable)
	row := tx.QueryRow(ctx, query, prID)
	var pr models.PullRequest
	err = row.Scan(&pr.ID, &pr.PullRequestName, &pr.AuthorID, &pr.Status, &pr.MergerdAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.log.Debug("pull request not found", slog.String("pr_id", prID))
			return nil, storage.ErrPRNotFound
		}

		s.log.Debug(FailedToScanRowMsg, slog.Any("error", err))
		return nil, err
	}

	if pr.Status == models.PRStatusOpen {
		query = fmt.Sprintf(`UPDATE %s SET status = 'MERGED', merged_at = $1 WHERE pull_request_id = $2`, PRTable)
		mergedAt := time.Now()
		_, err = tx.Exec(ctx, query, mergedAt, prID)
		if err != nil {
			s.log.Debug(FailedToUpdateValues, slog.Any("error", err))
			return nil, err
		}
		pr.Status = models.PRStatusMerged
		pr.MergerdAt = &mergedAt
	}

	err = tx.Commit(ctx)
	if err != nil {
		s.log.Debug(FailedToCommitTransac, slog.Any("error", err))
		return nil, err
	}
	return &pr, nil
}

func (s *Storage) Reassign(ctx context.Context, prID, revID string) (*models.PRExtended, error) {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})

	if err != nil {
		s.log.Debug(CantStartTransacMsg, slog.Any("error", err))
		return nil, err
	}
	defer tx.Rollback(ctx)

	query := fmt.Sprintf("SELECT TRUE FROM %s WHERE pull_request_id = $1 AND user_id = $2", ReviewersTable)
	exists := false
	err = tx.QueryRow(ctx, query, prID, revID).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrReviewerNotAssigned
		}
		s.log.Debug(FailedToUpdateValues, slog.Any("error", err))
		return nil, err
	}

	pr, err := s.getPR(ctx, tx, prID)
	if err != nil {
		return nil, err
	}

	teamName, err := s.teamName(ctx, tx, revID)
	if err != nil {
		return nil, err
	}

	user, err := s.randomUser(ctx, tx, teamName, revID)
	if err != nil {
		return nil, err
	}

	err = s.deleteOld(ctx, tx, prID, revID)
	if err != nil {
		return nil, err
	}

	err = s.insertNewReviewer(ctx, tx, user, prID)
	if err != nil {
		return nil, err
	}

	query = fmt.Sprintf("SELECT user_id FROM %s WHERE pull_request_id = $1", ReviewersTable)
	rows, err := tx.Query(ctx, query, prID)
	if err != nil {
		s.log.Debug(FailedToStartQuery, slog.Any("error", err))
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			s.log.Debug(FailedToScanRowsMsg, slog.Any("error", err))
			return nil, err
		}
		pr.AssignedReviewers = append(pr.AssignedReviewers, &models.Reviewer{ReviewerID: userID})
	}
	err = tx.Commit(ctx)
	if err != nil {
		s.log.Debug(FailedToCommitTransac, slog.Any("error", err))
		return nil, err
	}
	return &models.PRExtended{
		PullRequest: pr,
		ReplacedBy:  user.UserID,
	}, nil
}

func (s *Storage) deleteOld(ctx context.Context, tx pgx.Tx, prID string, revID string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE pull_request_id = $1 AND user_id = $2", ReviewersTable)
	_, err := tx.Exec(ctx, query, prID, revID)
	if err != nil {
		s.log.Debug(FailedToDeleteValues, slog.Any("error", err))
		return err
	}
	return nil
}

func (s *Storage) insertNewReviewer(ctx context.Context, tx pgx.Tx, user *models.User, prID string) error {
	query := fmt.Sprintf("INSERT INTO %s (pull_request_id, user_id)	VALUES ($1, $2) ON CONFLICT (pull_request_id, user_id) DO NOTHING", ReviewersTable)
	_, err := tx.Exec(ctx, query, user.UserID, prID)
	if err != nil {
		s.log.Debug(FailedToUpdateValues, slog.Any("error", err))
		return err
	}
	return nil
}

func (s *Storage) randomUser(ctx context.Context, tx pgx.Tx, teamName string, revID string) (*models.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE team_name = $1 AND user_id != $2 AND is_active = true ORDER BY RANDOM() LIMIT 1", UsersTable)
	row := tx.QueryRow(ctx, query, teamName, revID)

	var user models.User
	if err := row.Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrNoAvailableUsers
		}
		s.log.Debug(FailedToScanRowMsg, slog.Any("error", err))
		return nil, err
	}
	return &user, nil
}

func (s *Storage) teamName(ctx context.Context, tx pgx.Tx, revID string) (string, error) {
	query := fmt.Sprintf("SELECT team_name FROM %s WHERE user_id = $1", UsersTable)
	row := tx.QueryRow(ctx, query, revID)
	var teamName string
	if err := row.Scan(&teamName); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", storage.ErrUserNotFound
		}
		s.log.Debug(FailedToScanRowMsg, slog.Any("error", err))
		return "", err
	}
	return teamName, nil
}

func (s *Storage) getPR(ctx context.Context, tx pgx.Tx, prID string) (*models.PullRequest, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE pull_request_id = $1", PRTable)
	row := tx.QueryRow(ctx, query, prID)
	var pr models.PullRequest
	if err := row.Scan(&pr.ID, &pr.PullRequestName, &pr.AuthorID, &pr.Status, &pr.MergerdAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.log.Debug("pull request not found", slog.String("pr_id", prID))
			return nil, storage.ErrPRNotFound
		}
		s.log.Debug(FailedToScanRowMsg, slog.Any("error", err))
		return nil, err
	}
	if pr.Status == models.PRStatusMerged {
		return nil, storage.ErrPRAlreadyMerged
	}
	return &pr, nil

}
