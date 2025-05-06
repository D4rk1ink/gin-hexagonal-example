package auth_service

import (
	"context"

	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/domain"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/dto"
	custom_error "github.com/D4rk1ink/gin-hexagonal-example/internal/core/error"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/port"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/hash"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/jwt"
)

type authService struct {
	userRepo port.UserRepository
	jwt      jwt.Jwt
	hash     hash.Hash
}

func NewAuthService(userRepo port.UserRepository, jwt jwt.Jwt, hash hash.Hash) port.AuthService {
	return &authService{
		userRepo: userRepo,
		jwt:      jwt,
		hash:     hash,
	}
}

func (s *authService) Register(ctx context.Context, payload dto.UserRegisterDto) (*string, error) {
	if payload.Password != payload.ConfirmPassword {
		return nil, custom_error.NewError(custom_error.ErrAuthInvalidConfirmPassword, nil)
	}

	existsUser, err := s.userRepo.GetByEmail(ctx, payload.Email)
	if err != nil {
		return nil, err
	}
	if existsUser != nil {
		return nil, custom_error.NewError(custom_error.ErrAuthEmailAlreadyExists, nil)
	}

	hashed, err := s.hash.HashPassword(ctx, payload.Password)
	if err != nil {
		return nil, err
	}
	user, err := domain.NewUser(payload.Name, payload.Email, *hashed)
	if err != nil {
		return nil, err
	}

	id, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return id, nil
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
