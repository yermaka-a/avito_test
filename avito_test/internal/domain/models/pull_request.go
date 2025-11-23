package models

import "time"

type Reviewer struct {
	ReviewerID string
}

type PullRequest struct {
	ID                string      `json:"pull_request_id"`
	PullRequestName   string      `json:"pull_request_name"`
	AuthorID          string      `json:"author_id"`
	Status            string      `json:"status"`
	AssignedReviewers []*Reviewer `json:"reviewers"`
	MergerdAt         *time.Time  `json:"merged_at"`
}

type PRExtended struct {
	*PullRequest
	ReplacedBy string `json:"replaced_by"`
}

type PRStatus = string

const (
	PRStatusOpen   PRStatus = "OPEN"
	PRStatusMerged PRStatus = "MERGED"
)

func NewPullRequest(id, title, authorID string) *PullRequest {
	return &PullRequest{
		ID:              id,
		PullRequestName: title,
		AuthorID:        authorID,
		Status:          PRStatusOpen,
	}
}

func NewPRExtended(id, title, authorID, replacedBy string, assignedReviewers []*Reviewer) *PRExtended {
	prEx := &PRExtended{
		PullRequest: NewPullRequest(id, title, authorID),
		ReplacedBy:  replacedBy,
	}
	prEx.AssignedReviewers = assignedReviewers
	return prEx

}
