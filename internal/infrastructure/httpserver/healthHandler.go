package httpserver

import (
	"context"
	"net/http"
	"time"
)

// DatabasePinger проверяет доступность базы данных.
type DatabasePinger interface {
	Ping(ctx context.Context) error
}

// NewHealthHandler создаёт обработчик состояния приложения.
func NewHealthHandler(db DatabasePinger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), time.Second)
		defer cancel()

		if err := db.Ping(ctx); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
