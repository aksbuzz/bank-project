package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aksbuzz/library-project/config"
	"github.com/aksbuzz/library-project/helper/logger"
	"github.com/aksbuzz/library-project/server"
	"github.com/aksbuzz/library-project/store"
	"github.com/jackc/pgx/v5"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	logger := logger.New(cfg)
	slog.SetDefault(logger)

	ctx := context.Background()

	logger.Info("Connecting to database")
	conn, err := pgx.Connect(ctx, fmt.Sprintf("%s sslmode=disable", cfg.DSN))
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)
	logger.Info("Connected to database")

	store := store.New(conn)

	s := server.New(logger)
	s.Add(server.NewHttpServer(store, cfg, logger))
	s.Add(server.NewTaskScheduler(store, logger))

	err = s.Run(ctx)
	if err != nil {
		panic(err)
	}
}
