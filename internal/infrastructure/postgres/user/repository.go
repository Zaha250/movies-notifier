package user

import (
	"context"
	"fmt"

	userapp "kino-notifier/internal/application/user"
	"kino-notifier/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

var _ userapp.UserRepository = (*Repository)(nil)

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Upsert сохраняет состояние пользователя, сохраняя исходные ID и CreatedAt.
func (r *Repository) Upsert(ctx context.Context, user domain.User) error {
	const query = `
		INSERT INTO users (
			id,
			telegram_user_id,
			telegram_chat_id,
			username,
			first_name,
			last_name,
			language_code,
			is_bot,
			is_active,
			last_interaction_at,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (telegram_user_id)
		DO UPDATE SET
			telegram_chat_id = EXCLUDED.telegram_chat_id,
			username = EXCLUDED.username,
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name,
			language_code = EXCLUDED.language_code,
			is_bot = EXCLUDED.is_bot,
			is_active = EXCLUDED.is_active,
			last_interaction_at = EXCLUDED.last_interaction_at,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.pool.Exec(ctx, query,
		user.ID,
		user.TelegramUserID,
		user.TelegramChatID,
		user.Username,
		user.FirstName,
		user.LastName,
		user.LanguageCode,
		user.IsBot,
		user.IsActive,
		user.LastInteractionAt,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("upsert пользователя в PostgreSQL: %w", err)
	}
	return nil
}
