package http_util

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ValidateQuery(ctx *gin.Context, payload any) error {
	if err := ctx.ShouldBindQuery(payload); err != nil {
		return err
	}
	if err := validator.New().Struct(payload); err != nil {
		return err
	}
	return nil
}

func ValidateJSON(ctx *gin.Context, payload any) error {
	if err := ctx.ShouldBindJSON(payload); err != nil {
		return err
	}
	if err := validator.New().Struct(payload); err != nil {
		return err
	}
	return nil
}
