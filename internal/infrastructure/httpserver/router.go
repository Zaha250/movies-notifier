package httpserver

import "net/http"

func NewRouter(healthHandler, usersHandler http.Handler) http.Handler {
	const apiPrefix = "/api"

	apiMux := http.NewServeMux()
	apiMux.Handle("GET /health", healthHandler)
	apiMux.Handle("GET /users", usersHandler)

	mux := http.NewServeMux()
	mux.Handle(apiPrefix+"/", http.StripPrefix(apiPrefix, apiMux))

	return mux
}
