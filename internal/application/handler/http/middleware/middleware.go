package middleware

import (
	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/port"
	"github.com/gin-gonic/gin"
)

type Middleware struct{}

func NewMiddleware(userRepo port.UserRepository) *Middleware {
	return &Middleware{}
}

func (m *Middleware) Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
