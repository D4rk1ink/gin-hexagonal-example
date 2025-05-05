package port

import (
	"context"

	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/domain"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/dto"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]*domain.User, error)
	GetById(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Create(ctx context.Context, payload *domain.User) (*string, error)
	Update(ctx context.Context, payload dto.UserUpdateDto) error
	Delete(ctx context.Context, id string) error
}
