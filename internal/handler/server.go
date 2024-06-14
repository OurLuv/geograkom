package handler

import (
	"net/http"
	"time"

	"github.com/OurLuv/geograkom/internal/config"
	"github.com/gorilla/mux"
)

type Server struct {
	s *http.Server
}

func (s *Server) Start() error {
	err := s.s.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) ShutDown() error {
	err := s.s.Close()
	if err != nil {
		return err
	}
	return nil
}

func NewServer(cfg config.Config, r *mux.Router) *Server {
	s := &http.Server{
		Addr:           cfg.Server,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return &Server{
		s: s,
	}
}
