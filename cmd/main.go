package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "painaway_test/docs"
	"painaway_test/internal/app"

	"go.uber.org/zap"
)

// TODO: Покрыть swagger весь проект
// @title PainAway API
// @version 1.0
// @description API for PainAway
// @host localhost:8080
// @BasePath /
func main() {
	a, err := app.New()
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	// Запуск сервера
	go func() {
		if err := a.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.Logger.Fatal("server failed", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	a.Logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := a.Server.Shutdown(ctx); err != nil {
		a.Logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	a.Logger.Info("Server exited gracefully")
}
