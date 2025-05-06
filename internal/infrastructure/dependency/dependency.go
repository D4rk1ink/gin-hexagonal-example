package dependency

import (
	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/port"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/service"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/database"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/hash"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/jwt"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/repository"
)

type Infrastructure struct {
	Database database.MongoDb
	Jwt      jwt.Jwt
}

type Service struct {
	AuthService port.AuthService
	UserService port.UserService
}

type Dependency struct {
	Service        *Service
	Infrastructure *Infrastructure
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

	jwt := jwt.NewJwt()
	hash := hash.NewHash()

	userRepo := repository.NewUserRepository(db)

	authService := service.NewAuthService(userRepo, jwt, hash)
	userService := service.NewUserService(userRepo)

	return &Dependency{
		Service: &Service{
			AuthService: authService,
			UserService: userService,
		},
		Infrastructure: &Infrastructure{
			Database: db,
			Jwt:      jwt,
		},
	}
}
