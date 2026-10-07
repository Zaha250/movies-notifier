package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidTelegramUserID = errors.New("некорректный Telegram ID пользователя")
	ErrInvalidTelegramChatID = errors.New("некорректный ID личного Telegram-чата")
	ErrBotUserNotSupported   = errors.New("регистрация ботов не поддерживается")
)

type User struct {
	ID             uuid.UUID
	TelegramUserID int64
	TelegramChatID int64
	Username       *string
	FirstName      *string
	LastName       *string
	LanguageCode   *string
	IsBot          bool
	IsActive       bool
	// Время последнего сообщения или команды
	LastInteractionAt time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// Validate проверяет данные пользователя для работы в личном чате.
func (u User) Validate() error {
	if u.TelegramUserID <= 0 {
		return ErrInvalidTelegramUserID
	}
	if u.TelegramChatID <= 0 {
		return ErrInvalidTelegramChatID
	}
	if u.IsBot {
		return ErrBotUserNotSupported
	}
	return nil
}
