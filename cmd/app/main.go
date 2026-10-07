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

	userapp "kino-notifier/internal/application/user"
	"kino-notifier/internal/config"
	"kino-notifier/internal/infrastructure/httpserver"
	"kino-notifier/internal/infrastructure/postgres/connection"
	postgresuser "kino-notifier/internal/infrastructure/postgres/user"
	"kino-notifier/internal/infrastructure/telegram"
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
	database, err := connection.Open(startupCtx, cfg.PostgresURL)
	cancelStartup()
	if err != nil {
		return fmt.Errorf("ошибка подключения к PostgreSQL: %w", err)
	}
	defer database.Close()

	userRepository := postgresuser.NewRepository(database.Pool)
	registerUser := userapp.NewRegisterUser(userRepository)
	telegramClient, err := telegram.NewClient(ctx, cfg.Telegram.Token, logger)
	if err != nil {
		return fmt.Errorf("инициализация Telegram-клиента: %w", err)
	}
	notifier := telegram.NewNotifier(telegramClient)
	telegramBot, err := telegram.NewBot(ctx, telegramClient, registerUser, notifier, logger)
	if err != nil {
		return fmt.Errorf("инициализация Telegram-бота: %w", err)
	}

	healthHandler := httpserver.NewHealthHandler(database.Pool)
	router := httpserver.NewRouter(healthHandler)

	server, err := httpserver.New(cfg.Port, router)
	if err != nil {
		return fmt.Errorf("создание HTTP-сервера: %w", err)
	}

	runtimeCtx, cancelRuntime := context.WithCancel(ctx)
	defer cancelRuntime()

	serverErrors := make(chan error, 1)
	botDone := make(chan struct{})

	logger.Info("запуск HTTP-сервера", "port", cfg.Port)

	go func() {
		serverErrors <- server.Run()
	}()
	go func() {
		defer close(botDone)
		telegramBot.Run(runtimeCtx)
	}()
	logger.Info("запуск Telegram-бота", "mode", "long_polling")

	var serverErr error
	serverStopped := false

	select {
	case serverErr = <-serverErrors:
		serverStopped = true
	case <-botDone:
		if ctx.Err() == nil {
			serverErr = fmt.Errorf("Telegram-бот неожиданно завершил работу")
		}
	case <-ctx.Done():
		logger.Info("получен сигнал остановки")
	}

	cancelRuntime()

	shutdownCtx, cancelShutdown := context.WithTimeout(
		context.Background(),
		15*time.Second,
	)
	defer cancelShutdown()

	shutdownErr := server.Shutdown(shutdownCtx)

	if !serverStopped {
		serverErr = errors.Join(serverErr, <-serverErrors)
	}
	select {
	case <-botDone:
	case <-shutdownCtx.Done():
		shutdownErr = errors.Join(shutdownErr,
			fmt.Errorf("ожидание остановки Telegram-бота: %w", shutdownCtx.Err()),
		)
	}

	return errors.Join(serverErr, shutdownErr)
}
