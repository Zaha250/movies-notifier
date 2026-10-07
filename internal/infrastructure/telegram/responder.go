package telegram

import (
	"context"
	"log/slog"
	"strconv"

	"kino-notifier/internal/application/notification"
)

// Responder отправляет ответы любых Telegram-команд через общий порт уведомлений.
type Responder struct {
	notifier notification.Notifier
	logger   *slog.Logger
}

func NewResponder(notifier notification.Notifier, logger *slog.Logger) *Responder {
	return &Responder{notifier: notifier, logger: logger}
}

// Reply преобразует ID чата в адрес получателя и логирует ошибку доставки.
func (r *Responder) Reply(ctx context.Context, chatID int64, text string) {
	if ctx.Err() != nil {
		return
	}

	err := r.notifier.Send(ctx, notification.Message{
		Recipient: strconv.FormatInt(chatID, 10),
		Text:      text,
	})
	if err != nil && ctx.Err() == nil {
		r.logger.Error("не удалось отправить ответ Telegram",
			"telegram_chat_id", chatID,
			"error", err,
		)
	}
}
