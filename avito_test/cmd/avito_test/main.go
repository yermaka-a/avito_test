package main

import (
	"avito_test/internal/app"
	"avito_test/internal/config"
	"avito_test/internal/logger"
	"avito_test/internal/storage/postgres"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config := config.MustLoad()
	log := logger.Setup(os.Stdout, config.LogLevel())

	ctx := context.Background()

	dbc := config.DBConfig
	storagePath := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", dbc.USER, dbc.PASS, dbc.HOST, dbc.PORT, dbc.DBName)
	storage, err := postgres.New(ctx, storagePath, log)
	if err != nil {
		fmt.Print(err)
	}
	defer storage.Close()

	app := app.New(*config, log, storage)
	go func() {
		app.MustRun()
	}()

	stopChan := make(chan os.Signal, 1)

	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)

	<-stopChan
	log.Info("server is stopped")
}
