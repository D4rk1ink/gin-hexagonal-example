package user_service

import (
	"context"

	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/domain"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/dto"
	custom_error "github.com/D4rk1ink/gin-hexagonal-example/internal/core/error"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/port"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/hash"
)

type userService struct {
	userRepo port.UserRepository
	hash     hash.Hash
}

func NewUserService(userRepo port.UserRepository, hash hash.Hash) port.UserService {
	return &userService{
		userRepo: userRepo,
		hash:     hash,
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

func (s *userService) Create(ctx context.Context, payload dto.UserCreateDto) (*domain.User, error) {
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

	user, err = s.userRepo.GetById(ctx, *id)
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
