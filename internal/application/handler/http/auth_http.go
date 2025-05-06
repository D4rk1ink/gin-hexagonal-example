package http

import (
	"fmt"
	"net/http"

	http_apigen "github.com/D4rk1ink/gin-hexagonal-example/internal/application/handler/http/apigen"
	http_mapper "github.com/D4rk1ink/gin-hexagonal-example/internal/application/handler/http/mapper"
	http_util "github.com/D4rk1ink/gin-hexagonal-example/internal/application/handler/http/util"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
)

func (h *httpHandler) Register(ctx *gin.Context) {
	c := ctx.Request.Context()

	var body http_apigen.RegisterJSONRequestBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		httpCode := http.StatusBadRequest
		http_util.ResponseError(ctx, err, &httpCode)
		return
	}

	_, err := h.service.AuthService.Register(c, http_mapper.ToUserRegisterDto(body))
	if err != nil {
		logger.Error(fmt.Sprintf("Register error: %v", err))
		http_util.ResponseError(ctx, err, nil)
		return
	}

	ctx.JSON(http.StatusOK, http_apigen.RegisterRes{
		Success: true,
	})
}

func (h *httpHandler) Login(ctx *gin.Context) {
	c := ctx.Request.Context()

	var body http_apigen.LoginJSONRequestBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		httpCode := http.StatusBadRequest
		http_util.ResponseError(ctx, err, &httpCode)
		return
	}

	result, err := h.service.AuthService.Login(c, http_mapper.ToCredentialDto(body))
	if err != nil {
		logger.Error(fmt.Sprintf("Login error: %v", err))
		http_util.ResponseError(ctx, err, nil)
		return
	}

	ctx.JSON(http.StatusOK, http_mapper.ToAccessTokenResponse(result))
}
