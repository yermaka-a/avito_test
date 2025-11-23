package handlers

type UserStateRequest struct {
	UserId   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type PRShortResponse struct {
	ID     string `json:"pull_request_id"`
	Name   string `json:"pull_request_naem"`
	Author string `json:"author_id"`
	Status string `json:"status"`
}

type ReviewerPRsResponse struct {
	UserId string             `json:"user_id"`
	PRs    []*PRShortResponse `json:"pull_requests"`
}

type PRRequest struct {
	ID       string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorId string `json:"author_id"`
}

type ErrorDetails struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorDetails `json:"error"`
}

func NewErrResponse(err error, code string) *ErrorResponse {

	return &ErrorResponse{
		Error: ErrorDetails{
			Code:    code,
			Message: err.Error(),
		},
	}
}

type ReassignReviewers struct {
	PrID  string `json:"pull_request_id"`
	OldID string `json:"old_user_id"`
}

type merge struct {
	PrID string `json:"pull_request_id"`
}
