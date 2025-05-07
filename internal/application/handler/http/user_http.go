package http

import (
	"fmt"
	"net/http"

	http_apigen "github.com/D4rk1ink/gin-hexagonal-example/internal/application/handler/http/apigen"
	http_mapper "github.com/D4rk1ink/gin-hexagonal-example/internal/application/handler/http/mapper"
	http_util "github.com/D4rk1ink/gin-hexagonal-example/internal/application/handler/http/util"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/domain"
	custom_error "github.com/D4rk1ink/gin-hexagonal-example/internal/core/error"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
	"github.com/guregu/null"
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

func (h *httpHandler) CreateUser(ctx *gin.Context, params http_apigen.CreateUserParams) {
	c := ctx.Request.Context()

	var body http_apigen.CreateUserJSONRequestBody
	if err := http_util.ValidateJSON(ctx, &body); err != nil {
		err = custom_error.NewError(custom_error.ErrBadRequest, null.StringFrom(err.Error()).Ptr())
		http_util.ResponseError(ctx, err, nil)
		return
	}

	result, err := h.service.UserService.Create(c, http_mapper.ToUserCreateDto(body))
	if err != nil {
		logger.Error(fmt.Sprintf("Register error: %v", err))
		http_util.ResponseError(ctx, err, nil)
		return
	}

	ctx.JSON(http.StatusCreated, http_apigen.UserCreateRes{
		Data: http_mapper.ToUserResponse(result),
	})
}

func (h *httpHandler) UpdateUserById(ctx *gin.Context, id string, params http_apigen.UpdateUserByIdParams) {
	c := ctx.Request.Context()

	var body http_apigen.UpdateUserByIdJSONRequestBody
	if err := http_util.ValidateJSON(ctx, &body); err != nil {
		err = custom_error.NewError(custom_error.ErrBadRequest, null.StringFrom(err.Error()).Ptr())
		http_util.ResponseError(ctx, err, nil)
		return
	}

	payload := http_mapper.ToUserUpdateDto(id, body)
	result, err := h.service.UserService.Update(c, payload)
	if err != nil {
		logger.Error(fmt.Sprintf("UpdateUserById error: %v", err))
		http_util.ResponseError(ctx, err, nil)
		return
	}

	ctx.JSON(http.StatusOK, http_apigen.UserRes{
		Data: http_mapper.ToUserResponse(result),
	})
}

func (h *httpHandler) DeleteUserById(ctx *gin.Context, id string, params http_apigen.DeleteUserByIdParams) {
	c := ctx.Request.Context()

	err := h.service.UserService.Delete(c, id)
	if err != nil {
		logger.Error(fmt.Sprintf("DeleteUserById error: %v", err))
		http_util.ResponseError(ctx, err, nil)
		return
	}

	ctx.JSON(http.StatusOK, http_apigen.RegisterRes{
		Success: true,
	})
}
