package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kino-notifier/internal/config"
	"kino-notifier/internal/infrastructure/httpserver"
	"kino-notifier/internal/infrastructure/postgres"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	if err := run(ctx, logger); err != nil {
		logger.Error("приложение завершилось с ошибкой", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, logger *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("загрузка конфигурации: %w", err)
	}

	startupCtx, cancelStartup := context.WithTimeout(ctx, 10*time.Second)
	db, err := postgres.New(startupCtx, cfg.PostgresUrl)
	cancelStartup()
	if err != nil {
		return fmt.Errorf("ошибка подключения к PostgreSQL: %w", err)
	}
	defer db.Close()

	healthHandler := httpserver.NewHealthHandler(db.Pool)
	router := httpserver.NewRouter(healthHandler)

	server, err := httpserver.New(cfg.Port, router)
	if err != nil {
		return fmt.Errorf("создание HTTP-сервера: %w", err)
	}

	serverErrors := make(chan error, 1)

	logger.Info("запуск HTTP-сервера", "port", cfg.Port)

	go func() {
		serverErrors <- server.Run()
	}()

	var serverErr error
	serverStopped := false

	select {
	case serverErr = <-serverErrors:
		serverStopped = true
	case <-ctx.Done():
		logger.Info("получен сигнал остановки")
	}

	shutdownCtx, cancelShutdown := context.WithTimeout(
		context.Background(),
		15*time.Second,
	)
	defer cancelShutdown()

	shutdownErr := server.Shutdown(shutdownCtx)

	if !serverStopped {
		serverErr = <-serverErrors
	}

	return errors.Join(serverErr, shutdownErr)
}
