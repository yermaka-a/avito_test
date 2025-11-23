package services

import (
	"avito_test/internal/domain/models"
	"avito_test/internal/storage"
	"context"
	"errors"
	"log/slog"
)

const (
	UnexpectedMsg = "unexpected error"
)

func (s *Service) CreateTeam(ctx context.Context, team *models.Team) error {

	err := s.teamRepo.CreateTeam(ctx, team)
	if err != nil {
		if errors.Is(err, storage.ErrTeamAlreadyExists) {
			return storage.ErrTeamAlreadyExists
		}
		s.log.Warn(UnexpectedMsg, slog.Any("error", err))
		return ErrUnexpectedError
	}

	return nil

}

func (s *Service) GetTeamWithMembers(ctx context.Context, teamName string) (*models.Team, error) {

	team, err := s.teamRepo.GetTeamWithMembers(ctx, teamName)
	if err != nil {
		if errors.Is(err, storage.ErrTeamNotFound) {
			return nil, storage.ErrTeamNotFound
		}

		s.log.Warn(UnexpectedMsg, slog.Any("error", err))
		return nil, ErrUnexpectedError
	}

	return team, nil

}
