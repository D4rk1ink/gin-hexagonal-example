package http_util

import (
	"net/http"

	http_apigen "github.com/D4rk1ink/gin-hexagonal-example/internal/application/handler/http/apigen"
	custom_error "github.com/D4rk1ink/gin-hexagonal-example/internal/core/error"
	"github.com/gin-gonic/gin"
	"github.com/guregu/null"
)

func ResponseError(ctx *gin.Context, err error, httpCode *int) {
	switch e := err.(type) {
	case custom_error.CustomError:
		ctx.JSON(e.GetHttpCode(), http_apigen.ErrorRes{
			Error: http_apigen.ErrorBody{
				Code:    e.GetCode(),
				Message: e.GetMessage(),
			},
		})
		return
	default:
		if httpCode == nil {
			code := http.StatusInternalServerError
			httpCode = &code
		}
		ctx.JSON(*httpCode, http_apigen.ErrorRes{
			Error: http_apigen.ErrorBody{
				Code:    custom_error.ErrInternalServerError,
				Message: null.StringFrom("Internal server error").Ptr(),
			},
		})
		return
	}
}

func AbortResponseError(ctx *gin.Context, err error, httpCode *int) {
	switch e := err.(type) {
	case custom_error.CustomError:
		ctx.AbortWithStatusJSON(e.GetHttpCode(), http_apigen.ErrorRes{
			Error: http_apigen.ErrorBody{
				Code:    e.GetCode(),
				Message: e.GetMessage(),
			},
		})
		return
	default:
		if httpCode == nil {
			code := http.StatusInternalServerError
			httpCode = &code
		}
		ctx.AbortWithStatusJSON(*httpCode, http_apigen.ErrorRes{
			Error: http_apigen.ErrorBody{
				Code:    custom_error.ErrInternalServerError,
				Message: null.StringFrom("Internal server error").Ptr(),
			},
		})
		return
	}
}
