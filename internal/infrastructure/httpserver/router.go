package httpserver

import "net/http"

func NewRouter(healthHandler http.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /health", healthHandler)

	return mux
}
