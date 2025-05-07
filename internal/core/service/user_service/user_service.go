package user_service

import (
	"context"

	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/domain"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/dto"
	custom_error "github.com/D4rk1ink/gin-hexagonal-example/internal/core/error"
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

func (s *userService) Update(ctx context.Context, userUpdate dto.UserUpdateDto) (*domain.User, error) {
	if userUpdate.Name == nil && userUpdate.Email == nil {
		return nil, custom_error.NewError(custom_error.ErrUserInvalidateUpdateInput, nil)
	}

	if userUpdate.Email != nil {
		existingUser, err := s.userRepo.GetByEmail(ctx, *userUpdate.Email)
		if err != nil {
			return nil, err
		}
		if existingUser != nil && existingUser.Email == *userUpdate.Email {
			return nil, custom_error.NewError(custom_error.ErrUserEmailAlreadyExists, nil)
		}
	}

	user, err := s.userRepo.GetById(ctx, userUpdate.ID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, custom_error.NewError(custom_error.ErrUserNotFound, nil)
	}

	if userUpdate.Name != nil {
		user.SetName(*userUpdate.Name)
	}
	if userUpdate.Email != nil {
		err := user.SetEmail(*userUpdate.Email)
		if err != nil {
			return nil, err
		}
	}

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
