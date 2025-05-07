package http_middleware

import (
	"strings"
	"time"

	http_util "github.com/D4rk1ink/gin-hexagonal-example/internal/application/handler/http/util"
	custom_error "github.com/D4rk1ink/gin-hexagonal-example/internal/core/error"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/port"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/jwt"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/logger"
	random_util "github.com/D4rk1ink/gin-hexagonal-example/internal/util/random"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Middleware interface {
	Authentication() gin.HandlerFunc
	Logger() gin.HandlerFunc
	CorrelationId() gin.HandlerFunc
}

type middleware struct {
	userService port.UserService
	jwt         jwt.Jwt
}

func NewMiddleware(userService port.UserService, jwt jwt.Jwt) Middleware {
	return &middleware{
		userService: userService,
		jwt:         jwt,
	}
}

func (m *middleware) Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c := ctx.Request.Context()

		token := ctx.Request.Header.Get("Authorization")
		if token == "" {
			logger.Error("Authorization header is empty")
			http_util.AbortResponseError(ctx, custom_error.NewError(custom_error.ErrUnauthorized, nil), nil)
			return
		}

		bearer, ok := strings.CutPrefix(token, "Bearer ")
		if !ok {
			logger.Error("Authorization header is not Bearer")
			http_util.AbortResponseError(ctx, custom_error.NewError(custom_error.ErrUnauthorized, nil), nil)
			return
		}

		claims, err := m.jwt.ValidateAccessToken(bearer)
		if err != nil {
			logger.Error("Invalid access token")
			http_util.AbortResponseError(ctx, custom_error.NewError(custom_error.ErrUnauthorized, nil), nil)
			return
		}

		user, err := m.userService.GetById(c, claims.ID)
		if err != nil {
			logger.Error("Failed to get user by id")
			http_util.AbortResponseError(ctx, custom_error.NewError(custom_error.ErrInternalServerError, nil), nil)
			return
		}
		if user == nil {
			logger.Error("User not found")
			http_util.AbortResponseError(ctx, custom_error.NewError(custom_error.ErrUnauthorized, nil), nil)
			return
		}

		ctx.Set("user", user)
		ctx.Next()
	}
}

func (m *middleware) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		duration := time.Since(start)
		logger.Info(
			"",
			zap.String("correlation_id", c.Writer.Header().Get("X-Correlation-ID")),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
		)
	}
}

func (m *middleware) CorrelationId() gin.HandlerFunc {
	return func(c *gin.Context) {
		correlationID := c.Request.Header.Get("X-Correlation-ID")
		if correlationID == "" {
			correlationID = random_util.RandomCorrelationId()
			c.Request.Header.Set("X-Correlation-ID", correlationID)
		}
		c.Writer.Header().Set("X-Correlation-ID", correlationID)
		c.Next()
	}
}
