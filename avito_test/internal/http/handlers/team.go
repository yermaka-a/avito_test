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
		h.log.Debug("bad request", slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Некорректное тело запроса"))
		return
	}

	err := h.service.CreateTeam(context.Background(), &team)
	if err != nil {

		if errors.Is(err, storage.ErrTeamAlreadyExists) {
			w.WriteHeader(http.StatusBadRequest)
			h.ErrResponse(w, storage.ErrTeamAlreadyExists, TeamExistsCode)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Команда создана"))
}

func (h *Handler) GetTeamWithMembers(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	teamName := params.Get("team_name")

	if teamName == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Некорректное имя команды"))
		return
	}
	team, err := h.service.GetTeamWithMembers(context.Background(), teamName)
	if err != nil {
		if errors.Is(err, storage.ErrTeamNotFound) {
			w.WriteHeader(http.StatusNotFound)
			h.ErrResponse(w, ErrResourceNotFound, NotFoundCode)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(team)
	if err != nil {
		h.log.Warn("failed to marshal json", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
