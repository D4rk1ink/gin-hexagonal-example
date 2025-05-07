package http_middleware

import (
	"strings"

	http_util "github.com/D4rk1ink/gin-hexagonal-example/internal/application/handler/http/util"
	custom_error "github.com/D4rk1ink/gin-hexagonal-example/internal/core/error"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/port"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/jwt"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
)

type Middleware interface {
	Authentication() gin.HandlerFunc
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
