package handlers

import (
	"avito_test/internal/services"
	"errors"
	"log/slog"
)

const (
	UnexpectedMsg = "unexpected error"
)

var (
	BadRequestMsg        = "bad request"
	IncorrectData        = "incorrect data"
	FailedToMarhsalJSON        = "failed to marshal json"
	PRExistsCode               = "PR_EXISTS"
	NotFoundCode               = "NOT_FOUND"
	PRMergedCode               = "PR_MERGED"
	NoAssignedCode             = "NOT_ASSIGNED"
	NoCandidateCode            = "NO_CANDIDATE"
	TeamExistsCode            = "TEAM_EXISTS"
	ErrResourceNotFound        = errors.New("resource not found")
	ErrCantReassignMergedPR    = errors.New("cannot reassign on merged PR")
	ErrNoActiveCandidates      = errors.New("no active replacement candidate in team")
	ErrReviewerNotAssignedToPR = errors.New("reviewer is not assigned to this PR")
)

type Handler struct {
	service *services.Service
	log     *slog.Logger
}

func New(log *slog.Logger, service *services.Service) *Handler {
	return &Handler{
		service: service,
		log:     log,
	}
}
