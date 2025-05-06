package service

import (
	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/service/auth_service"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/service/user_service"
)

var NewAuthService = auth_service.NewAuthService
var NewUserService = user_service.NewUserService
