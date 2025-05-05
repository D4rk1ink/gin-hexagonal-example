package port

import (
	"context"

	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/dto"
)

type AuthService interface {
	Login(ctx context.Context, payload dto.CredentialDto) (*dto.AccessTokenDto, error)
	Register(ctx context.Context, payload dto.UserRegisterDto) (*string, error)
}
