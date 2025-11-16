package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/Wucop228/avito-PullRequest/internal/app"
	"github.com/Wucop228/avito-PullRequest/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatalf("failed to init app: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		if err := application.RunHTTP(); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Printf("shutdown signal received, shutting down gracefully")
	case err := <-errCh:
		if err != nil {
			log.Printf("server error: %v", err)
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := application.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	} else {
		log.Println("server stopped gracefully")
	}
}
