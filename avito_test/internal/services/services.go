package services

import (
	"avito_test/internal/domain/models"
	"context"
	"log/slog"
)

type Service struct {
	teamRepo TeamRepository
	userRepo UserRepository
	prRepo   PRRepository
	log      *slog.Logger
}

func New(log *slog.Logger, teamRepo TeamRepository, userRepo UserRepository, prRepo PRRepository) *Service {

	return &Service{
		teamRepo: teamRepo,
		userRepo: userRepo,
		prRepo:   prRepo,
		log:      log,
	}
}

type PRRepository interface {
	CreatePR(ctx context.Context, pr *models.PullRequest) ([]*models.Reviewer, error)
	Reassign(ctx context.Context, prID, revID string) (*models.PRExtended, error)
	PRMarkAsMerged(ctx context.Context, prID string) (*models.PullRequest, error)
}

type TeamRepository interface {
	CreateTeam(ctx context.Context, team *models.Team) error
	GetTeamWithMembers(ctx context.Context, teamName string) (*models.Team, error)
}

type UserRepository interface {
	GetReview(ctx context.Context, userID string) ([]*models.PullRequest, error)
	SetIsActive(ctx context.Context, user *models.User) error
}
