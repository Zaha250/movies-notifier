package domain

import (
	"time"

	"github.com/google/uuid"
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
