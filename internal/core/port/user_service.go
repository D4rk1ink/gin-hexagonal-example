package port

import (
	"context"

	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/domain"
)

type UserService interface {
	GetAll(ctx context.Context) ([]*domain.User, error)
	GetById(ctx context.Context, id string) (*domain.User, error)
}
