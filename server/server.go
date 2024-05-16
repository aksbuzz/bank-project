package server

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type server struct {
	logger  *slog.Logger
	servers []Server
}

func New(logger *slog.Logger, servers []Server) *server {
	return &server{logger: logger, servers: servers}
}

type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
}

func (s *server) Run(ctx context.Context) error {
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	for _, sv := range s.servers {
		go func(sv Server) {
			if err := sv.Start(ctx); err != nil {
				s.logger.ErrorContext(ctx, "Failed to start server", slog.Any("error", err))
			}
		}(sv)
	}

	select {
	case <-signals:
		s.logger.Info("Received signal, stopping...")
	case <-ctx.Done():
		s.logger.Info("Context is done, stopping...")
	}

	for _, sv := range s.servers {
		if err := sv.Stop(ctx); err != nil {
			s.logger.ErrorContext(ctx, "Failed to stop server", slog.Any("error", err))
		}
	}

	return nil
}
