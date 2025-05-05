package http

import (
	"context"
	"net/http"

	"github.com/D4rk1ink/gin-hexagonal-example/internal/application/handler/http/apigen"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/application/handler/http/mapper"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/domain"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

func (h *httpHandler) GetUsers(ctx *gin.Context) {
	c := context.Background()

	result, err := h.service.UserService.GetAll(c)
	if err != nil {
		h.ResponseError(ctx, err, nil)
		return
	}

	users := lo.Map(result, func(user *domain.User, _ int) apigen.User {
		return mapper.ToUserResponse(user)
	})
	ctx.JSON(http.StatusOK, apigen.UsersRes{
		Data: users,
	})
}
