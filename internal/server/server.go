package server

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"
)

type Config struct {
	Addr              string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ReadHeaderTimeout time.Duration
}

type Server struct {
	httpServer *http.Server
	shutdownOnce sync.Once
}

func NewServer(cfg Config, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    cfg.Addr,
			Handler: handler,
			ReadTimeout: cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout: cfg.IdleTimeout,
			ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		},
	}
}

func (s *Server) Start() error {
	log.Printf("HTTP server starting on %s\n", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	var err error

	s.shutdownOnce.Do(func() {
		log.Println("HTTP server shutting down...")
		err = s.httpServer.Shutdown(ctx)
	})
	return err
}
