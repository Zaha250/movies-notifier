package user

import (
	"context"
	"fmt"
	"kino-notifier/internal/domain"
)

type ListUsers struct {
	users UserReader
}

func NewUserList(users UserReader) *ListUsers {
	return &ListUsers{users: users}
}

func (uc *ListUsers) Execute(ctx context.Context) ([]domain.User, error) {
	list, err := uc.users.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("получение списка пользователей: %w", err)
	}
	return list, nil
}
