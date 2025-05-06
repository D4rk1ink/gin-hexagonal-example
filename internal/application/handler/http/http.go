package http

import (
	"fmt"

	http_apigen "github.com/D4rk1ink/gin-hexagonal-example/internal/application/handler/http/apigen"
	http_middleware "github.com/D4rk1ink/gin-hexagonal-example/internal/application/handler/http/middleware"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/config"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/dependency"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
)

type HttpHandler interface {
	http_apigen.ServerInterface
	SetRouter() error
	Listen() error
}

type httpHandler struct {
	router         *gin.Engine
	service        *dependency.Service
	infrastructure *dependency.Infrastructure
	middleware     http_middleware.Middleware
}

func NewHttpHandler(
	dep *dependency.Dependency,
) HttpHandler {
	router := gin.Default()
	middleware := http_middleware.NewMiddleware(dep.Service.UserService, dep.Infrastructure.Jwt)

	return &httpHandler{
		router:         router,
		service:        dep.Service,
		infrastructure: dep.Infrastructure,
		middleware:     middleware,
	}
}

func (h *httpHandler) SetRouter() error {
	h.router.Use(gin.Recovery())
	h.router.Use(gin.Logger())

	if config.Config.App.Env != "production" {
		logger.Info("Running in development mode")
		h.router.StaticFile("/swagger/doc.yaml", "./docs/server/doc.yaml")
		h.router.StaticFile("/swagger", "./docs/server/swagger.html")
	}

	http_apigen.RegisterHandlers(h.router, h)

	return nil
}

func (h *httpHandler) Listen() error {
	h.SetRouter()

	if err := h.router.Run(":8080"); err != nil {
		panic(err)
	}

	logger.Info(fmt.Sprintf("Server started on port %s", config.Config.App.Port))

	return nil
}
