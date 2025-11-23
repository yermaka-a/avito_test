package services

import (
	"avito_test/internal/domain/models"
	"avito_test/internal/storage"
	"context"
	"errors"
	"log/slog"
)

func (s *Service) CreatePR(ctx context.Context, pr *models.PullRequest) error {
	reviewers, err := s.prRepo.CreatePR(ctx, pr)
	if err != nil {
		if errors.Is(err, storage.ErrPRIDAlreadyExists) || errors.Is(err, storage.ErrUserNotFound) {
			return err
		}
		s.log.Warn(UnexpectedMsg, slog.Any("error", err))
		return ErrUnexpectedError
	}
	pr.AssignedReviewers = reviewers
	return nil

}

func (s *Service) PRMarkAsMerged(ctx context.Context, prID string) (*models.PullRequest, error) {
	pr, err := s.prRepo.PRMarkAsMerged(ctx, prID)
	if err != nil {
		if errors.Is(err, storage.ErrPRNotFound) {
			return nil, storage.ErrPRNotFound
		}
		s.log.Warn(UnexpectedMsg, slog.Any("error", err))
		return nil, ErrUnexpectedError
	}
	return pr, nil

}

func (s *Service) Reassign(ctx context.Context, prID, revID string) (*models.PRExtended, error) {
	pr, err := s.prRepo.Reassign(ctx, prID, revID)
	if err != nil {
		if errors.Is(err, storage.ErrPRNotFound) || errors.Is(err, storage.ErrPRAlreadyMerged) || errors.Is(err, storage.ErrNoAvailableUsers) || errors.Is(err, storage.ErrUserNotFound) || errors.Is(err, storage.ErrReviewerNotAssigned) {
			return nil, err
		}

		s.log.Warn(UnexpectedMsg, slog.Any("error", err))
		return nil, ErrUnexpectedError
	}
	return pr, nil

}
