package app

import (
	"avito_test/internal/config"
	"avito_test/internal/http"
	"avito_test/internal/http/handlers"
	"avito_test/internal/services"
	"avito_test/internal/storage/postgres"
	"fmt"
	"log/slog"
)

type App struct {
	cfg    config.Config
	logger *slog.Logger
	server *http.Server
}

func New(cfg config.Config, logger *slog.Logger, storage *postgres.Storage) *App {

	service := services.New(logger, storage, storage, storage)
	handler := handlers.New(logger, service)
	server := http.NewServer(handler)

	return &App{
		cfg:    cfg,
		logger: logger,
		server: server,
	}
}

func (a *App) MustRun() error {
	a.logger.Info(fmt.Sprintf("starting server on port %s", a.cfg.Port()))
	return a.server.Run(fmt.Sprintf(":%s", a.cfg.Port()))
}
