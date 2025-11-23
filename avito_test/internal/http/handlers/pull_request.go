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

func (h *Handler) writeResponse(w http.ResponseWriter, b []byte, status int, err error) {
	if err != nil {
		h.log.Warn(FailedToMarhsalJSON, slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
func (h *Handler) CreatePR(w http.ResponseWriter, r *http.Request) {
	var prDTO PRRequest
	if err := json.NewDecoder(r.Body).Decode(&prDTO); err != nil {
		h.log.Debug(BadRequestMsg, slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(IncorrectData))
		return
	}
	pr := models.NewPullRequest(prDTO.ID, prDTO.Name, prDTO.AuthorId)
	err := h.service.CreatePR(context.Background(), pr)
	if err != nil {
		if errors.Is(err, storage.ErrPRIDAlreadyExists) {
			h.ErrResponse(w, storage.ErrPRIDAlreadyExists, PRExistsCode, http.StatusConflict)
			return
		}

		if errors.Is(err, storage.ErrUserNotFound) {
			h.ErrResponse(w, storage.ErrUserNotFound, NotFoundCode, http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(pr)
	h.writeResponse(w, b, http.StatusCreated, err)

}

func (h *Handler) PRMarkAsMerged(w http.ResponseWriter, r *http.Request) {
	var merge merge
	if err := json.NewDecoder(r.Body).Decode(&merge); err != nil {
		h.log.Debug(BadRequestMsg, slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(IncorrectData))
		return
	}
	pr, err := h.service.PRMarkAsMerged(context.Background(), merge.PrID)
	if err != nil {
		if errors.Is(err, storage.ErrPRNotFound) {

			h.ErrResponse(w, ErrResourceNotFound, NotFoundCode, http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(pr)
	h.writeResponse(w, b, http.StatusOK, err)

}

func (h *Handler) Reassign(w http.ResponseWriter, r *http.Request) {
	var rRevws ReassignReviewers
	if err := json.NewDecoder(r.Body).Decode(&rRevws); err != nil {
		h.log.Debug(BadRequestMsg, slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(IncorrectData))
		return
	}

	prEx, err := h.service.Reassign(context.Background(), rRevws.PrID, rRevws.OldID)
	if err != nil {
		if errors.Is(err, storage.ErrPRNotFound) || errors.Is(err, storage.ErrUserNotFound) {
			h.ErrResponse(w, ErrResourceNotFound, NotFoundCode, http.StatusNotFound)
			return
		}

		if errors.Is(err, storage.ErrPRAlreadyMerged) {
			h.ErrResponse(w, ErrCantReassignMergedPR, PRMergedCode, http.StatusConflict)
			return
		}

		if errors.Is(err, storage.ErrNoAvailableUsers) {
			h.ErrResponse(w, ErrNoActiveCandidates, NoCandidateCode, http.StatusConflict)
			return
		}

		if errors.Is(err, storage.ErrReviewerNotAssigned) {
			h.ErrResponse(w, ErrReviewerNotAssignedToPR, NoAssignedCode, http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(prEx)
	if err != nil {
		h.log.Warn(FailedToMarhsalJSON, slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (h *Handler) ErrResponse(w http.ResponseWriter, err error, code string, status int) {
	resp := NewErrResponse(err, code)
	b, err := json.Marshal(resp)
	h.writeResponse(w, b, status, err)
}
