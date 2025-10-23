package server

import (
	"context"
	"net/http"
	"user-subscriptions/internal/config"
)

type Server struct {
	standardServer http.Server
}

func New(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		standardServer: http.Server{
			Addr:    ":" + cfg.HTTP.Port,
			Handler: handler,
		},
	}
}

func (s *Server) Run() error {
	return s.standardServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.standardServer.Shutdown(ctx)
}
