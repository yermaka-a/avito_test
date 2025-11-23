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

func (h *Handler) writeResponse(w http.ResponseWriter, b []byte, err error) {
	if err != nil {
		h.log.Warn(FailedToMarhsalJSON, slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
func (h *Handler) CreatePR(w http.ResponseWriter, r *http.Request) {
	var prDTO PRRequest
	if err := json.NewDecoder(r.Body).Decode(&prDTO); err != nil {
		h.log.Debug("bad request", slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Некорректное тело запроса"))
		return
	}
	pr := models.NewPullRequest(prDTO.ID, prDTO.Name, prDTO.AuthorId)
	err := h.service.CreatePR(context.Background(), pr)
	if err != nil {
		if errors.Is(err, storage.ErrPRIDAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			h.ErrResponse(w, storage.ErrPRIDAlreadyExists, PRExistsCode)
			return
		}

		if errors.Is(err, storage.ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			h.ErrResponse(w, storage.ErrUserNotFound, NotFoundCode)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(pr)
	w.WriteHeader(http.StatusCreated)
	h.writeResponse(w, b, err)

}

func (h *Handler) PRMarkAsMerged(w http.ResponseWriter, r *http.Request) {
	var merge merge
	if err := json.NewDecoder(r.Body).Decode(&merge); err != nil {
		h.log.Debug("bad request", slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Некорректное тело запроса"))
		return
	}
	pr, err := h.service.PRMarkAsMerged(context.Background(), merge.PrID)
	if err != nil {
		if errors.Is(err, storage.ErrPRNotFound) {

			w.WriteHeader(http.StatusNotFound)
			h.ErrResponse(w, ErrResourceNotFound, NotFoundCode)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(pr)
	w.WriteHeader(http.StatusOK)
	h.writeResponse(w, b, err)

}

func (h *Handler) Reassign(w http.ResponseWriter, r *http.Request) {
	var rRevws ReassignReviewers
	if err := json.NewDecoder(r.Body).Decode(&rRevws); err != nil {
		h.log.Debug("bad request", slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Некорректное тело запроса"))
		return
	}

	prEx, err := h.service.Reassign(context.Background(), rRevws.PrID, rRevws.OldID)
	if err != nil {
		if errors.Is(err, storage.ErrPRNotFound) || errors.Is(err, storage.ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			h.ErrResponse(w, ErrResourceNotFound, NotFoundCode)
			return
		}

		if errors.Is(err, storage.ErrPRAlreadyMerged) {
			w.WriteHeader(http.StatusConflict)
			h.ErrResponse(w, ErrCantReassignMergedPR, PRMergedCode)
			return
		}

		if errors.Is(err, storage.ErrNoAvailableUsers) {
			w.WriteHeader(http.StatusConflict)
			h.ErrResponse(w, ErrNoActiveCandidates, NoCandidateCode)
			return
		}

		if errors.Is(err, storage.ErrReviewerNotAssigned) {
			w.WriteHeader(http.StatusConflict)
			h.ErrResponse(w, ErrReviewerNotAssignedToPR, NoAssignedCode)
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

func (h *Handler) ErrResponse(w http.ResponseWriter, err error, code string) {
	resp := NewErrResponse(err, code)
	b, err := json.Marshal(resp)
	h.writeResponse(w, b, err)
}
