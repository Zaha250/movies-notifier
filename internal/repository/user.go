package repository

import (
	"context"
	"kino-notifier/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
}
