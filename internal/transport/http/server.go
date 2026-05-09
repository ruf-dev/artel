package http

import (
	"context"
	"errors"
	"net/http"

	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/service"
)

type Server struct {
	router *http.ServeMux
	addr   string
	server *http.Server
}

func New(addr string, vs service.VaultService) *Server {
	router := http.NewServeMux()
	h := NewVaultHandler(vs)
	router.HandleFunc("POST /api/vault", h.CreateVault)

	return &Server{
		router: router,
		addr:   addr,
	}
}

func (s *Server) Start() error {
	s.server = &http.Server{
		Addr:    s.addr,
		Handler: s.router,
	}

	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return rerrors.Wrap(err, "failed to start server")
	}

	return nil
}

func (s *Server) Stop() error {
	if s.server != nil {
		err := s.server.Shutdown(context.Background())
		if err != nil {
			return rerrors.Wrap(err, "failed to shutdown server")
		}
	}
	return nil
}
