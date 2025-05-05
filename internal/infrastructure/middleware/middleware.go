package middleware

import "github.com/gin-gonic/gin"

type Middleware struct{}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
