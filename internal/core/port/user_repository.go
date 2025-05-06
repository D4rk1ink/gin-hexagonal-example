package port

import (
	"context"

	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/domain"
)

//go:generate mockgen -package mock_port -source=user_repository.go -destination=mock/mock_user_repository.go *
type UserRepository interface {
	GetAll(ctx context.Context) ([]*domain.User, error)
	GetById(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	Create(ctx context.Context, payload *domain.User) (*string, error)
	Update(ctx context.Context, payload *domain.User) error
	Delete(ctx context.Context, id string) error
}
