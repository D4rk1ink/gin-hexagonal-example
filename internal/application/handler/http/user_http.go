package http

import (
	"fmt"
	"net/http"

	http_apigen "github.com/D4rk1ink/gin-hexagonal-example/internal/application/handler/http/apigen"
	http_mapper "github.com/D4rk1ink/gin-hexagonal-example/internal/application/handler/http/mapper"
	http_util "github.com/D4rk1ink/gin-hexagonal-example/internal/application/handler/http/util"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/domain"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

func (h *httpHandler) GetUsers(ctx *gin.Context, params http_apigen.GetUsersParams) {
	c := ctx.Request.Context()

	result, err := h.service.UserService.GetAll(c)
	if err != nil {
		logger.Error(fmt.Sprintf("GetUsers error: %v", err))
		http_util.ResponseError(ctx, err, nil)
		return
	}

	users := lo.Map(result, func(user *domain.User, _ int) http_apigen.User {
		return http_mapper.ToUserResponse(user)
	})
	ctx.JSON(http.StatusOK, http_apigen.UsersRes{
		Data: users,
	})
}

func (h *httpHandler) GetUserById(ctx *gin.Context, id string, params http_apigen.GetUserByIdParams) {
	c := ctx.Request.Context()

	result, err := h.service.UserService.GetById(c, id)
	if err != nil {
		logger.Error(fmt.Sprintf("GetUsers error: %v", err))
		http_util.ResponseError(ctx, err, nil)
		return
	}

	ctx.JSON(http.StatusOK, http_apigen.UserRes{
		Data: http_mapper.ToUserResponse(result),
	})
}
