package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	tgbot "github.com/go-telegram/bot"
)

const initializationTimeout = 10 * time.Second

// Client инкапсулирует Telegram SDK для входящих команд и исходящих уведомлений.
type Client struct {
	bot      *tgbot.Bot
	username string
}

func NewClient(ctx context.Context, token string, logger *slog.Logger) (*Client, error) {
	if strings.TrimSpace(token) == "" {
		return nil, fmt.Errorf("токен Telegram не задан")
	}
	if logger == nil {
		return nil, fmt.Errorf("логгер Telegram не задан")
	}

	client, err := tgbot.New(token,
		tgbot.WithSkipGetMe(),
		tgbot.WithAllowedUpdates(tgbot.AllowedUpdates{"message"}),
		tgbot.WithWorkers(1),
		tgbot.WithNotAsyncHandlers(),
		tgbot.WithDefaultHandler(ignoreUpdate),
		tgbot.WithErrorsHandler(func(err error) {
			logger.Error("ошибка Telegram", "error", redactTelegramError(err, token))
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("создание Telegram-клиента: %s", redactTelegramError(err, token))
	}

	initializationCtx, cancel := context.WithTimeout(ctx, initializationTimeout)
	defer cancel()

	botUser, err := client.GetMe(initializationCtx)
	if err != nil {
		return nil, fmt.Errorf("проверка токена Telegram: %s", redactTelegramError(err, token))
	}
	return &Client{bot: client, username: botUser.Username}, nil
}

// Токен входит в URL Bot API и может оказаться в сетевой ошибке.
func redactTelegramError(err error, token string) string {
	return strings.ReplaceAll(err.Error(), token, "[REDACTED]")
}
