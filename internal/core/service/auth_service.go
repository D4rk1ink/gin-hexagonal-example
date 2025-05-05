package service

import (
	"context"

	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/domain"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/dto"
	custom_error "github.com/D4rk1ink/gin-hexagonal-example/internal/core/error"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/port"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/jwt"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo port.UserRepository
	jwt      jwt.Jwt
}

func NewAuthService(userRepo port.UserRepository, customJwt jwt.Jwt) port.AuthService {
	return &authService{
		userRepo: userRepo,
		jwt:      customJwt,
	}
}

func (s *authService) hashPassword(ctx context.Context, password string) (*string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	hashedPassword := string(hashed)

	return &hashedPassword, nil
}

func (s *authService) comparePassword(ctx context.Context, password string, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return custom_error.NewError(custom_error.ErrAuthInvalidCredentials, nil)
	}

	return nil
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

	hashed, err := s.hashPassword(ctx, payload.Password)
	if err != nil {
		return nil, err
	}
	user, err := domain.NewUser(payload.Name, payload.Email, *hashed)
	if err != nil {
		return nil, err
	}

	id, err := s.userRepo.Create(ctx, user)
	if err != nil {
		println(err.Error())
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
	err = s.comparePassword(ctx, payload.Password, user.Password)
	if err != nil {
		return nil, err
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
