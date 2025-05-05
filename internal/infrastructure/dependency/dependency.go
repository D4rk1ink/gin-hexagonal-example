package dependency

import (
	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/port"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/service"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/database"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/jwt"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/repository"
)

type Service struct {
	AuthService port.AuthService
	UserService port.UserService
}

type Dependency struct {
	Service *Service
}

func NewDependency() *Dependency {
	db, err := database.NewMongodb()
	if err != nil {
		panic(err)
	}
	err = db.Connect()
	if err != nil {
		panic(err)
	}

	jwtImpl := jwt.NewJwt()

	userRepo := repository.NewUserRepository(*db)

	authService := service.NewAuthService(userRepo, jwtImpl)
	userService := service.NewUserService(userRepo)

	return &Dependency{
		Service: &Service{
			AuthService: authService,
			UserService: userService,
		},
	}
}
