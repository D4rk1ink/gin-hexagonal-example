package http

import (
	"fmt"
	"net/http"

	"github.com/D4rk1ink/gin-hexagonal-example/internal/application/handler/http/apigen"
	custom_error "github.com/D4rk1ink/gin-hexagonal-example/internal/core/error"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/config"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/dependency"
	"github.com/gin-gonic/gin"
	"github.com/guregu/null"
)

type HttpHandler interface {
	apigen.ServerInterface
	Listen() error
}

type httpHandler struct {
	router  *gin.Engine
	service *dependency.Service
}

func NewHttpHandler(
	service *dependency.Service,
) HttpHandler {
	router := gin.Default()

	return &httpHandler{
		router:  router,
		service: service,
	}
}

func (h *httpHandler) Listen() error {
	if config.Config.Env != "production" {
		fmt.Println("Swagger is enabled")
		h.router.StaticFile("/swagger/doc.yaml", "./docs/server/doc.yaml")
		h.router.StaticFile("/swagger", "./docs/server/swagger.html")
	}

	apigen.RegisterHandlers(h.router, h)

	if err := h.router.Run(":8080"); err != nil {
		panic(err)
	}

	return nil
}

func (h *httpHandler) ResponseError(ctx *gin.Context, err error, httpCode *int) {
	switch e := err.(type) {
	case custom_error.CustomError:
		ctx.JSON(e.GetHttpCode(), apigen.ErrorRes{
			Error: apigen.ErrorBody{
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
		ctx.JSON(*httpCode, apigen.ErrorRes{
			Error: apigen.ErrorBody{
				Code:    custom_error.ErrInternalServerError,
				Message: null.StringFrom("Internal server error").Ptr(),
			},
		})
		return
	}
}
