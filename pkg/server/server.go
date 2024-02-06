package server

import (
	"context"
	"movieMatchMaker/internal/config"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func Init() *Server {
	return new(Server)
}

func (s *Server) Run(cfg config.Config, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:         cfg.Address,
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
