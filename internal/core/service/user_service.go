package service

import (
	"context"

	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/domain"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/port"
)

type userService struct {
	userRepo port.UserRepository
}

func NewUserService(userRepo port.UserRepository) port.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetAll(ctx context.Context) ([]*domain.User, error) {
	users, err := s.userRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *userService) GetById(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.userRepo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
