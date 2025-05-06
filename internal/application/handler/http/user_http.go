package http

import (
	"net/http"

	http_apigen "github.com/D4rk1ink/gin-hexagonal-example/internal/application/handler/http/apigen"
	http_mapper "github.com/D4rk1ink/gin-hexagonal-example/internal/application/handler/http/mapper"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/domain"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

func (h *httpHandler) GetUsers(ctx *gin.Context) {
	c := ctx.Request.Context()

	result, err := h.service.UserService.GetAll(c)
	if err != nil {
		h.ResponseError(ctx, err, nil)
		return
	}

	users := lo.Map(result, func(user *domain.User, _ int) http_apigen.User {
		return http_mapper.ToUserResponse(user)
	})
	ctx.JSON(http.StatusOK, http_apigen.UsersRes{
		Data: users,
	})
}
