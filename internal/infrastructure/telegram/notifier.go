package telegram

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"kino-notifier/internal/application/notification"

	tgbot "github.com/go-telegram/bot"
)

const deliveryTimeout = 5 * time.Second

// Notifier реализует отправку уведомлений через Telegram Bot API.
type Notifier struct {
	client *Client
}

var _ notification.Notifier = (*Notifier)(nil)

func NewNotifier(client *Client) *Notifier {
	return &Notifier{client: client}
}

func (n *Notifier) Send(ctx context.Context, message notification.Message) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := message.Validate(); err != nil {
		return err
	}

	chatID, err := strconv.ParseInt(message.Recipient, 10, 64)
	if err != nil || chatID == 0 {
		return fmt.Errorf("%w: Telegram ожидает ненулевой числовой chat ID", notification.ErrInvalidRecipient)
	}

	deliveryCtx, cancel := context.WithTimeout(ctx, deliveryTimeout)
	defer cancel()

	_, err = n.client.bot.SendMessage(deliveryCtx, &tgbot.SendMessageParams{
		ChatID: chatID,
		Text:   message.Text,
	})
	if err == nil {
		return nil
	}
	if contextErr := deliveryCtx.Err(); contextErr != nil {
		return fmt.Errorf("%w: %w", notification.ErrDeliveryFailed, contextErr)
	}
	return fmt.Errorf("%w: %s", notification.ErrDeliveryFailed,
		redactTelegramError(err, n.client.bot.Token()),
	)
}
