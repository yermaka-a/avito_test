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

func (h *Handler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var team models.Team
	if err := json.NewDecoder(r.Body).Decode(&team); err != nil {
		h.log.Debug(BadRequestMsg, slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(IncorrectData))
		return
	}

	err := h.service.CreateTeam(context.Background(), &team)
	if err != nil {

		if errors.Is(err, storage.ErrTeamAlreadyExists) {
			h.ErrResponse(w, storage.ErrTeamAlreadyExists, TeamExistsCode, http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(team)
	h.writeResponse(w, b, http.StatusCreated, err)
}

func (h *Handler) GetTeamWithMembers(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	teamName := params.Get("team_name")

	if teamName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(IncorrectData))
		return
	}
	team, err := h.service.GetTeamWithMembers(context.Background(), teamName)
	if err != nil {
		if errors.Is(err, storage.ErrTeamNotFound) {
			h.ErrResponse(w, ErrResourceNotFound, NotFoundCode, http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(team)
	if err != nil {
		h.log.Warn(FailedToMarhsalJSON, slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
