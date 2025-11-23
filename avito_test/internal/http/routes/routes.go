package routes

import (
	"avito_test/internal/http/handlers"
	"net/http"
)

func SetupRoutes(r *http.ServeMux, handler *handlers.Handler) {
	r.HandleFunc("POST /team/add", handler.CreateTeam)
	r.HandleFunc("GET /team/get", handler.GetTeamWithMembers)
	r.HandleFunc("POST /users/setIsActive", handler.SetIsActive)
	r.HandleFunc("GET /users/getReview", handler.GetReview)
	r.HandleFunc("POST /pullRequest/create", handler.CreatePR)
	r.HandleFunc("POST /pullRequest/merge", handler.PRMarkAsMerged)
	r.HandleFunc("POST /pullRequest/reassign", handler.Reassign)
}
