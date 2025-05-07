package http_util

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Validate(ctx *gin.Context, payload any) error {
	if err := ctx.ShouldBindJSON(payload); err != nil {
		return err
	}
	if err := validator.New().Struct(payload); err != nil {
		return err
	}
	return nil
}
