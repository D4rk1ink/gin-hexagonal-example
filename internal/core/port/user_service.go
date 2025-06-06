package port

import (
	"context"

	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/domain"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/dto"
)

//go:generate mockgen -package mock_port -source=user_service.go -destination=mock/mock_user_service.go *
type UserService interface {
	Count(ctx context.Context) (int64, error)
	GetAll(ctx context.Context) ([]*domain.User, error)
	GetById(ctx context.Context, id string) (*domain.User, error)
	Create(ctx context.Context, userCreate dto.UserCreateDto) (*domain.User, error)
	Update(ctx context.Context, userUpdate dto.UserUpdateDto) (*domain.User, error)
	Delete(ctx context.Context, id string) error
}
