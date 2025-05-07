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
	GetRouter() *gin.Engine
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

	wrapper := http_apigen.ServerInterfaceWrapper{
		Handler: h,
	}

	if config.Config.App.Env != "production" {
		logger.Info("Running in development mode")
		h.router.StaticFile("/swagger/doc.yaml", "./docs/server/doc.yaml")
		h.router.StaticFile("/swagger", "./docs/server/swagger.html")
	}

	auth := h.router.Group("/api/auth")
	auth.POST("/register", wrapper.Register)
	auth.POST("/login", wrapper.Login)

	users := h.router.Group("/api/users")
	users.Use(h.middleware.Authentication())
	users.POST("", wrapper.CreateUser)
	users.GET("", wrapper.GetUsers)
	users.GET("/:id", wrapper.GetUserById)
	users.PATCH("/:id", wrapper.UpdateUserById)
	users.DELETE("/:id", wrapper.DeleteUserById)

	return nil
}

func (h *httpHandler) GetRouter() *gin.Engine {
	return h.router
}

func (h *httpHandler) Listen() error {
	h.SetRouter()

	if err := h.router.Run(":8080"); err != nil {
		panic(err)
	}

	logger.Info(fmt.Sprintf("Server started on port %s", config.Config.App.Port))

	return nil
}
