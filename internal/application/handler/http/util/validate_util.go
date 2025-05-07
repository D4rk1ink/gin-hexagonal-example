package http_util

import (
	"fmt"
	"reflect"
	"strings"

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
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		return strings.ToLower(fld.Tag.Get("json"))
	})
	if err := validate.Struct(payload); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			return fmt.Errorf("validation error: %s %s %s", e.Field(), e.Tag(), e.Param())
		}
	}
	return nil
}
