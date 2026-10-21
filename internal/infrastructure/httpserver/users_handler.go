package httpserver

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"kino-notifier/internal/domain"

	"github.com/google/uuid"
)

// UserLister описывает сценарий получения пользователей, необходимый HTTP-обработчику.
type UserLister interface {
	Execute(ctx context.Context) ([]domain.User, error)
}

type userResponse struct {
	ID                uuid.UUID `json:"id"`
	TelegramUserID    int64     `json:"telegram_user_id"`
	TelegramChatID    int64     `json:"telegram_chat_id"`
	Username          *string   `json:"username"`
	FirstName         *string   `json:"first_name"`
	LastName          *string   `json:"last_name"`
	LanguageCode      *string   `json:"language_code"`
	IsBot             bool      `json:"is_bot"`
	IsActive          bool      `json:"is_active"`
	LastInteractionAt time.Time `json:"last_interaction_at"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type usersResponse struct {
	List []userResponse `json:"list"`
}

type usersErrorResponse struct {
	Error string `json:"error"`
}

// NewUsersHandler создаёт обработчик получения списка всех пользователей.
func NewUsersHandler(users UserLister, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		list, err := users.Execute(ctx)
		if err != nil {
			logger.ErrorContext(r.Context(), "не удалось получить список пользователей", "error", err)
			writeUsersJSON(w, r, logger, http.StatusInternalServerError, usersErrorResponse{
				Error: "Не удалось получить список пользователей",
			})
			return
		}

		response := usersResponse{List: make([]userResponse, 0, len(list))}
		for _, user := range list {
			response.List = append(response.List, userResponse{
				ID:                user.ID,
				TelegramUserID:    user.TelegramUserID,
				TelegramChatID:    user.TelegramChatID,
				Username:          user.Username,
				FirstName:         user.FirstName,
				LastName:          user.LastName,
				LanguageCode:      user.LanguageCode,
				IsBot:             user.IsBot,
				IsActive:          user.IsActive,
				LastInteractionAt: user.LastInteractionAt,
				CreatedAt:         user.CreatedAt,
				UpdatedAt:         user.UpdatedAt,
			})
		}

		writeUsersJSON(w, r, logger, http.StatusOK, response)
	})
}

func writeUsersJSON(w http.ResponseWriter, r *http.Request, logger *slog.Logger, status int, response any) {
	body, err := json.Marshal(response)
	if err != nil {
		logger.ErrorContext(r.Context(), "не удалось сериализовать HTTP-ответ", "error", err)
		status = http.StatusInternalServerError
		body = []byte(`{"error":"Не удалось сформировать ответ"}`)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if _, err := w.Write(body); err != nil {
		logger.ErrorContext(r.Context(), "не удалось записать HTTP-ответ", "error", err)
	}
}
