package storage

import "errors"

var (
	ErrCreateTeam          = errors.New("team is not created")
	ErrTeamNotFound        = errors.New("team not found")
	ErrGetTeam             = errors.New("team is not got")
	ErrCreatePR            = errors.New("pull request is not created")
	ErrUsersNotFound       = errors.New("users not found in table users")
	ErrPRAlreadyMerged     = errors.New("pull request already merged")
	ErrPRIDAlreadyExists   = errors.New("PR id already exists")
	ErrFieldNotFound       = errors.New("field not found")
	ErrPRNotFound          = errors.New("pull request  not found")
	ErrTeamAlreadyExists   = errors.New("team already exists")
	ErrUserNotFound        = errors.New("user not found")
	ErrNoAvailableUsers    = errors.New("no available users")
	ErrReviewerNotAssigned = errors.New("reviewer is not assigned to this PR")
)
