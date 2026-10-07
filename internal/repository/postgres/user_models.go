package postgres

import (
	"kino-notifier/internal/domain"
	"time"

	"github.com/google/uuid"
)

type userRow struct {
	ID uuid.UUID `db:"id"`

	TelegramUserID int64 `db:"telegram_user_id"`
	TelegramChatID int64 `db:"telegram_chat_id"`

	Username     *string `db:"username"`
	FirstName    *string `db:"first_name"`
	LastName     *string `db:"last_name"`
	LanguageCode *string `db:"language_code"`

	IsBot    bool `db:"is_bot"`
	IsActive bool `db:"is_active"`

	// Время последнего сообщения или команды
	LastInteractionAt time.Time `db:"last_interaction_at"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func userToDomain(user userRow) domain.User {
	return domain.User{
		ID: user.ID,

		TelegramUserID: user.TelegramUserID,
		TelegramChatID: user.TelegramChatID,

		Username:     user.Username,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		LanguageCode: user.LanguageCode,

		IsActive: user.IsActive,

		LastInteractionAt: user.LastInteractionAt,
	}
}

func domainToUser(user domain.User) userRow {
	return userRow{
		ID: user.ID,

		TelegramUserID: user.TelegramUserID,
		TelegramChatID: user.TelegramChatID,

		Username:     user.Username,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		LanguageCode: user.LanguageCode,

		IsActive: user.IsActive,

		LastInteractionAt: user.LastInteractionAt,
	}
}
