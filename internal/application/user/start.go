package user

import (
	"context"
	"fmt"
	"kino-notifier/internal/domain"
	"kino-notifier/internal/repository"
)

type StartInteractor struct {
	userRepo repository.UserRepository
}

func NewStartInteractor(userRepo repository.UserRepository) *StartInteractor {
	return &StartInteractor{userRepo: userRepo}
}

func (i *StartInteractor) Execute(ctx context.Context, user domain.User) error {
	return fmt.Errorf("not implemented")
}
