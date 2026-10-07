package postgres

import (
	"context"
	"fmt"
	"kino-notifier/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

func (r *UserRepository) Create(
	ctx context.Context,
	user domain.User,
) error {
	const query = `
		INSERT INTO users (
			telegram_user_id,
			telegram_chat_id,
			username,
			first_name,
			last_name,
			language_code,
			is_active,
			last_interaction_at
		)
		VALUES (
			$1, $2, $3, $4, $5, $6,
			TRUE,
			NOW(),
			NOW()
		)
		ON CONFLICT (telegram_user_id)
		DO UPDATE SET
			telegram_chat_id = EXCLUDED.telegram_chat_id,
			username = EXCLUDED.username,
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name,
			language_code = EXCLUDED.language_code,
			is_active = TRUE,
			last_interaction_at = NOW(),
			updated_at = NOW()
	`

	_, err := r.pool.Exec(
		ctx,
		query,
		user.TelegramUserID,
		user.TelegramChatID,
		user.Username,
		user.FirstName,
		user.LastName,
		user.LanguageCode,
	)
	if err != nil {
		return fmt.Errorf("error upsert user: %w", err)
	}

	return nil
}
