package auth_service

import (
	"context"

	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/dto"
	custom_error "github.com/D4rk1ink/gin-hexagonal-example/internal/core/error"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/port"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/hash"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/jwt"
)

type authService struct {
	userService port.UserService
	userRepo    port.UserRepository
	jwt         jwt.Jwt
	hash        hash.Hash
}

func NewAuthService(userService port.UserService, userRepo port.UserRepository, jwt jwt.Jwt, hash hash.Hash) port.AuthService {
	return &authService{
		userService: userService,
		userRepo:    userRepo,
		jwt:         jwt,
		hash:        hash,
	}
}

func (s *authService) Register(ctx context.Context, payload dto.UserRegisterDto) (*string, error) {
	user, err := s.userService.Create(ctx, dto.UserCreateDto(payload))
	if err != nil {
		return nil, err
	}

	return &user.ID, nil
}

func (s *authService) Login(ctx context.Context, payload dto.CredentialDto) (*dto.AccessTokenDto, error) {
	user, err := s.userRepo.GetByEmail(ctx, payload.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, custom_error.NewError(custom_error.ErrAuthInvalidCredentials, nil)
	}
	err = s.hash.ComparePassword(ctx, payload.Password, user.Password)
	if err != nil {
		return nil, custom_error.NewError(custom_error.ErrAuthInvalidCredentials, nil)
	}

	accessToken, expiresIn, err := s.jwt.GenerateAccessToken(&jwt.GenerateTokenInput{
		ID:    user.ID,
		Email: user.Email,
	})
	if err != nil {
		return nil, err
	}
	return &dto.AccessTokenDto{
		AccessToken: *accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   *expiresIn,
	}, nil
}
