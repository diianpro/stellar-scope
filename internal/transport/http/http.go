package http

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	public       *http.Server
	publicRouter *chi.Mux

	handler *Handler
}

func New(handler *Handler) *Server {
	return &Server{
		publicRouter: chi.NewRouter(),

		handler: handler,
	}
}

func (s *Server) ServePublic(addr string) error {
	s.registerPublicRoutes()

	s.public = &http.Server{
		Addr:         addr,
		Handler:      s.publicRouter,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	return s.public.ListenAndServe()
}

func (s *Server) ShutdownPublic(ctx context.Context) error {
	if err := s.public.Shutdown(ctx); err != nil {
		return s.public.Close()
	}
	return nil
}

func (s *Server) registerPublicRoutes() {
	s.publicRouter.Get("/_/ready", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("OK"))
	})

	s.publicRouter.Route("/v1", func(r chi.Router) {
		r.Get("/picture/by/{date}", s.handler.GetByDate)
		r.Route("/list", func(r chi.Router) {
			r.Get("/picture", s.handler.GetAll)
		})
	})
}
