package user

import (
	"context"

	"kino-notifier/internal/domain"
)

// UserRepository описывает необходимый сценарию контракт хранения пользователей.
type UserRepository interface {
	Upsert(ctx context.Context, user domain.User) error
}

type UserReader interface {
	List(ctx context.Context) ([]domain.User, error)
}
