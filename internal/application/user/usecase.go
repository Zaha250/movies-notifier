package user

import (
	"context"
	"kino-notifier/internal/domain"
)

type StartUsecase interface {
	Execute(ctx context.Context, user domain.User) error
}
