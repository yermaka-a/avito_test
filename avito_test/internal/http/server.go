package http

import (
	"avito_test/internal/http/handlers"
	"avito_test/internal/http/routes"
	"net/http"
)

type Server struct {
	router *http.ServeMux
}

func NewServer(handler *handlers.Handler) *Server {
	r := http.NewServeMux()
	routes.SetupRoutes(r, handler)
	return &Server{router: r}
}

func (s *Server) Run(port string) error {
	return http.ListenAndServe(port, s.router)
}
