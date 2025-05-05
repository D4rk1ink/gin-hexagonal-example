package mapper

import (
	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/domain"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/repository/model"
)

func ToUserModel(user *domain.User) *model.UserModel {
	return &model.UserModel{
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToUserDomain(userModel *model.UserModel) *domain.User {
	return &domain.User{
		ID:        userModel.ID.Hex(),
		Name:      userModel.Name,
		Email:     userModel.Email,
		Password:  userModel.Password,
		CreatedAt: userModel.CreatedAt,
		UpdatedAt: userModel.UpdatedAt,
	}
}
