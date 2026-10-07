package user

import (
	"context"
	"fmt"
	"time"

	"kino-notifier/internal/domain"

	"github.com/google/uuid"
)

// RegisterUserInput не зависит от формата событий Telegram SDK.
type RegisterUserInput struct {
	TelegramUserID int64
	TelegramChatID int64
	Username       string
	FirstName      string
	LastName       string
	LanguageCode   string
	IsBot          bool
}

type RegisterUser struct {
	userRepository UserRepository
}

func NewRegisterUser(userRepository UserRepository) *RegisterUser {
	return &RegisterUser{userRepository: userRepository}
}

// Execute регистрирует пользователя или обновляет и активирует существующего.
func (uc *RegisterUser) Execute(ctx context.Context, input RegisterUserInput) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	user := domain.User{
		TelegramUserID: input.TelegramUserID,
		TelegramChatID: input.TelegramChatID,
		Username:       optionalString(input.Username),
		FirstName:      optionalString(input.FirstName),
		LastName:       optionalString(input.LastName),
		LanguageCode:   optionalString(input.LanguageCode),
		IsBot:          input.IsBot,
	}
	if err := user.Validate(); err != nil {
		return fmt.Errorf("регистрация пользователя: %w", err)
	}

	now := time.Now().UTC()
	user.ID = uuid.New()
	user.IsActive = true
	user.LastInteractionAt = now
	user.CreatedAt = now
	user.UpdatedAt = now

	if err := uc.userRepository.Upsert(ctx, user); err != nil {
		return fmt.Errorf("сохранение пользователя: %w", err)
	}
	return nil
}

func optionalString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
