package notification

import (
	"context"
	"errors"
	"strings"
)

var (
	ErrInvalidRecipient = errors.New("адрес получателя не задан или некорректен")
	ErrEmptyText        = errors.New("текст уведомления не задан")
	ErrDeliveryFailed   = errors.New("не удалось доставить уведомление")
)

// Message содержит данные уведомления без типов конкретного канала доставки.
// Формат Recipient определяется адаптером: например, chat ID или email.
type Message struct {
	Recipient string
	Text      string
}

func (m Message) Validate() error {
	if strings.TrimSpace(m.Recipient) == "" {
		return ErrInvalidRecipient
	}
	if strings.TrimSpace(m.Text) == "" {
		return ErrEmptyText
	}
	return nil
}

// Notifier — выходной порт отправки уведомлений.
type Notifier interface {
	Send(ctx context.Context, message Message) error
}
