package handlers

import (
	"avito_test/internal/domain/models"
	"avito_test/internal/storage"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

func (h *Handler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	var uSt UserStateRequest
	if err := json.NewDecoder(r.Body).Decode(&uSt); err != nil {
		h.log.Debug(BadRequestMsg, slog.Any("error", err))
		h.ErrResponse(w, ErrResourceNotFound, NotFoundCode, http.StatusBadRequest)
		return
	}

	user, err := h.service.SetIsActive(context.Background(), uSt.UserId, uSt.IsActive)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			h.ErrResponse(w, ErrResourceNotFound, NotFoundCode, http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(user)
	if err != nil {
		h.log.Warn(FailedToMarhsalJSON, slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (h *Handler) GetReview(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	userId := params.Get("user_id")
	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(IncorrectData))
		return
	}
	pr, err := h.service.GetReview(context.Background(), userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	prsDTO := h.prsToDTO(pr)
	prsDTO.UserId = userId
	b, err := json.Marshal(prsDTO)

	if err != nil {
		h.log.Warn(FailedToMarhsalJSON, slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (h *Handler) prsToDTO(prs []*models.PullRequest) *ReviewerPRsResponse {
	var rPRs ReviewerPRsResponse
	rPRs.PRs = make([]*PRShortResponse, 0)
	for _, v := range prs {
		rPRs.PRs = append(rPRs.PRs, &PRShortResponse{ID: v.ID, Name: v.PullRequestName, Author: v.AuthorID, Status: v.Status})
	}
	return &rPRs
}
