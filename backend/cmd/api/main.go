package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/critiq17/tripListik/internal/config"
	"github.com/critiq17/tripListik/internal/db"
	"github.com/critiq17/tripListik/internal/server"
	"github.com/critiq17/tripListik/internal/store"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg := config.LoadConfig()
	if err := cfg.ValidateAPI(); err != nil {
		slog.Error("config error", "err", err)
		os.Exit(1)
	}

	gormDB, sqlDB, err := db.Connect(cfg.DSN)
	if err != nil {
		slog.Error("db connect error", "err", err)
		os.Exit(1)
	}
	defer sqlDB.Close()

	app := server.New(cfg, store.New(gormDB))

	go func() {
		if err := app.Listen(cfg.HTTPAddr); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		slog.Error("shutdown error", "err", err)
		os.Exit(1)
	}
}
