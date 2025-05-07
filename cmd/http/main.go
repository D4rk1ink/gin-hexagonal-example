package main

import (
	"context"
	_http "net/http"
	"os/signal"
	"syscall"

	http_handler "github.com/D4rk1ink/gin-hexagonal-example/internal/application/handler/http"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/dependency"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/logger"
)

func main() {
	dep := dependency.NewDependency()

	handler := http_handler.NewHttpHandler(dep)

	go func() {
		if err := handler.ListenAndServe(); err != nil && err != _http.ErrServerClosed {
			panic(err)
		}
	}()

	shutdown, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-shutdown.Done()

	logger.Info("Received shutdown signal, shutting down server...")
	if err := handler.Shutdown(context.Background()); err != nil {
		panic(err)
	}
	if err := dep.Infrastructure.Database.Disconnect(context.Background()); err != nil {
		panic(err)
	}
	logger.Info("Server shutdown gracefully")
}
