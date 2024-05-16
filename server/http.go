package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/aksbuzz/library-project/config"
	"github.com/aksbuzz/library-project/service"
	"github.com/aksbuzz/library-project/store"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type HttpServer struct {
	logger *slog.Logger
	store  *store.Store
	server *http.Server
}

func NewHttpServer(store *store.Store, cfg *config.Config, logger *slog.Logger) *HttpServer {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("👋"))
	})

	service.New(store, logger).Register(r)

	s := &HttpServer{
		logger: logger.With("service", "http-server"),
		store:  store,
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort),
			Handler: r,
		},
	}

	return s
}

func (h *HttpServer) Start(ctx context.Context) error {
	h.logger.Info("Starting http server", slog.String("address", h.server.Addr))
	if err := h.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http server start error: %w", err)
	}

	return nil
}

func (h *HttpServer) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("http server stop error: %w", err)
	}

	h.logger.Info("http server stopped")

	return nil
}
