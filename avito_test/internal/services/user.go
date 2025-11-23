package services

import (
	"avito_test/internal/domain/models"
	"avito_test/internal/storage"
	"context"
	"errors"
	"log/slog"
)

func (s *Service) GetReview(ctx context.Context, userID string) ([]*models.PullRequest, error) {

	prs, err := s.userRepo.GetReview(ctx, userID)
	if err != nil {
		s.log.Warn(UnexpectedMsg, slog.Any("error", err))
		return nil, ErrUnexpectedError
	}

	return prs, nil

}

func (s *Service) SetIsActive(ctx context.Context, userID string, isActive bool) (*models.User, error) {
	user := models.NewUser(userID, "", "", isActive)
	err := s.userRepo.SetIsActive(ctx, user)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, storage.ErrUserNotFound
		}
		s.log.Warn(UnexpectedMsg, slog.Any("error", err))
		return nil, ErrUnexpectedError
	}

	return user, nil

}
