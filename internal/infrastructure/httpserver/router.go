package httpserver

import "net/http"

func NewRouter(healthHandler http.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /api/health", healthHandler)

	return mux
}
