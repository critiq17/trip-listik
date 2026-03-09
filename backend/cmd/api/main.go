package main

import (
	"context"
	"log"
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
	cfg := config.LoadConfig()
	if err := cfg.ValidateAPI(); err != nil {
		log.Fatalf("config error: %v", err)
	}

	gormDB, sqlDB, err := db.Connect(cfg.DSN)
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}
	defer sqlDB.Close()

	app := server.New(cfg, store.New(gormDB))

	go func() {
		if err := app.Listen(cfg.HTTPAddr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
}
