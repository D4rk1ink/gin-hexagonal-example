package http_middleware

import (
	"strings"

	http_util "github.com/D4rk1ink/gin-hexagonal-example/internal/application/handler/http/util"
	custom_error "github.com/D4rk1ink/gin-hexagonal-example/internal/core/error"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/port"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/jwt"
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
			http_util.AbortResponseError(ctx, custom_error.NewError(custom_error.ErrAuthUnauthenticated, nil), nil)
		}

		bearer, ok := strings.CutPrefix(token, "Bearer ")
		if !ok {
			http_util.AbortResponseError(ctx, custom_error.NewError(custom_error.ErrAuthUnauthenticated, nil), nil)
		}

		claims, err := m.jwt.ValidateAccessToken(bearer)
		if err != nil {
			http_util.AbortResponseError(ctx, custom_error.NewError(custom_error.ErrAuthUnauthenticated, nil), nil)
		}

		user, err := m.userService.GetById(c, claims.ID)
		if err != nil {
			http_util.AbortResponseError(ctx, custom_error.NewError(custom_error.ErrInternalServerError, nil), nil)
		}
		if user == nil {
			http_util.AbortResponseError(ctx, custom_error.NewError(custom_error.ErrAuthUnauthenticated, nil), nil)
		}

		ctx.Set("user", user)
		ctx.Next()
	}
}
