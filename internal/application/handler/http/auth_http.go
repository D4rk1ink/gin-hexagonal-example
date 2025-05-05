package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/D4rk1ink/gin-hexagonal-example/internal/application/handler/http/apigen"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/application/handler/http/mapper"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
)

func (h *httpHandler) Register(ctx *gin.Context) {
	c := context.Background()

	var body apigen.RegisterJSONRequestBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		httpCode := http.StatusBadRequest
		h.ResponseError(ctx, err, &httpCode)
		return
	}

	_, err := h.service.AuthService.Register(c, mapper.ToUserRegisterDto(body))
	if err != nil {
		logger.Error(fmt.Sprintf("Register error: %v", err))
		h.ResponseError(ctx, err, nil)
		return
	}

	ctx.JSON(http.StatusOK, apigen.RegisterRes{
		Success: true,
	})
}

func (h *httpHandler) Login(ctx *gin.Context) {
	c := context.Background()

	var body apigen.LoginJSONRequestBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		httpCode := http.StatusBadRequest
		h.ResponseError(ctx, err, &httpCode)
		return
	}

	result, err := h.service.AuthService.Login(c, mapper.ToCredentialDto(body))
	if err != nil {
		logger.Error(fmt.Sprintf("Login error: %v", err))
		h.ResponseError(ctx, err, nil)
		return
	}

	ctx.JSON(http.StatusOK, mapper.ToAccessTokenResponse(result))
}
