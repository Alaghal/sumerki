package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sumerki/backend/internal/config"
	"sumerki/backend/internal/db"
	httpserver "sumerki/backend/internal/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	database, err := db.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database open error: %v", err)
	}
	defer database.Close()

	server := httpserver.New(database)

	errCh := make(chan error, 1)
	go func() {
		log.Printf("starting backend on port %s", cfg.BackendPort)
		errCh <- server.Start(":" + cfg.BackendPort)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		log.Printf("received signal %s, shutting down", sig)
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}
}
