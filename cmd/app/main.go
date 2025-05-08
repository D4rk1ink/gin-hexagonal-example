package main

import (
	"context"
	"fmt"
	_http "net/http"
	"os/signal"
	"syscall"
	"time"

	http_handler "github.com/D4rk1ink/gin-hexagonal-example/internal/application/handler/http"
	scheduler_handler "github.com/D4rk1ink/gin-hexagonal-example/internal/application/handler/scheduler"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/dependency"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/logger"
)

func main() {
	dep := dependency.NewDependency()

	httpHandler := http_handler.NewHttpHandler(dep)
	scheduleHandler := scheduler_handler.NewScheduler(dep)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	err := scheduleHandler.Start(ctx)
	if err != nil {
		panic(err)
	}
	go func() {
		if err := httpHandler.ListenAndServe(); err != nil && err != _http.ErrServerClosed {
			panic(err)
		}
	}()

	<-ctx.Done()
	logger.Info("Received shutdown signal, shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := scheduleHandler.Shutdown(shutdownCtx); err != nil {
		logger.Error(fmt.Sprintf("Scheduler shutdown error: %v", err))
	}
	if err := httpHandler.Shutdown(shutdownCtx); err != nil {
		logger.Error(fmt.Sprintf("HTTP server shutdown error: %v", err))
	}
	if err := dep.Infrastructure.Database.Disconnect(shutdownCtx); err != nil {
		logger.Error(fmt.Sprintf("Database shutdown error: %v", err))
	}
	logger.Info("Server shutdown gracefully")
}
