package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func New(port int, handler http.Handler) (*Server, error) {
	if handler == nil {
		return nil, fmt.Errorf("HTTP handler is required")
	}

	return &Server{
		httpServer: &http.Server{
			Addr:              net.JoinHostPort("", strconv.Itoa(port)),
			Handler:           handler,
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       10 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       60 * time.Second,
		},
	}, nil
}

func (s *Server) Run() error {
	err := s.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http server: %w", err)
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		closeErr := s.httpServer.Close()

		return errors.Join(
			fmt.Errorf("ошибка завершения HTTP-сервера: %w", err),
			closeErr,
		)
	}
	return nil
}
