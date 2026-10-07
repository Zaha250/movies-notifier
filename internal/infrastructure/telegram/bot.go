package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"kino-notifier/internal/application/notification"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// Bot управляет получением событий Telegram через long polling.
type Bot struct {
	client *tgbot.Bot
}

func NewBot(
	ctx context.Context,
	client *Client,
	userRegistration UserRegistration,
	notifier notification.Notifier,
	logger *slog.Logger,
) (*Bot, error) {
	if client == nil || client.bot == nil {
		return nil, fmt.Errorf("Telegram-клиент не задан")
	}
	if userRegistration == nil {
		return nil, fmt.Errorf("сценарий регистрации пользователя не задан")
	}
	if notifier == nil {
		return nil, fmt.Errorf("отправитель уведомлений не задан")
	}
	if logger == nil {
		return nil, fmt.Errorf("логгер Telegram не задан")
	}

	initializationCtx, cancel := context.WithTimeout(ctx, initializationTimeout)
	defer cancel()

	webhook, err := client.bot.GetWebhookInfo(initializationCtx)
	if err != nil {
		return nil, fmt.Errorf("проверка режима Telegram: %s", redactTelegramError(err, client.bot.Token()))
	}
	if webhook.URL != "" {
		return nil, fmt.Errorf("у бота настроен webhook: отключите его перед запуском long polling")
	}

	responder := NewResponder(notifier, logger)
	handler := newStartHandler(userRegistration, responder, logger)
	client.bot.RegisterHandlerMatchFunc(func(update *models.Update) bool {
		return matchesStartCommand(update, client.username)
	}, handler.handle)

	return &Bot{client: client.bot}, nil
}

// Run блокируется до отмены контекста и выхода текущего обработчика.
func (b *Bot) Run(ctx context.Context) {
	b.client.Start(ctx)
}

func matchesStartCommand(update *models.Update, botUsername string) bool {
	if update == nil || update.Message == nil {
		return false
	}
	commandParts := strings.Fields(update.Message.Text)
	if len(commandParts) == 0 {
		return false
	}
	command := commandParts[0]
	return command == "/start" || strings.EqualFold(command, "/start@"+botUsername)
}

func ignoreUpdate(_ context.Context, _ *tgbot.Bot, _ *models.Update) {}
