package http_mapper

import (
	"github.com/D4rk1ink/gin-hexagonal-example/internal/application/handler/http/apigen"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/domain"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/dto"
)

func ToCredentialDto(payload apigen.LoginJSONRequestBody) dto.CredentialDto {
	return dto.CredentialDto{
		Email:    payload.Email,
		Password: payload.Password,
	}
}

func ToAccessTokenResponse(token *dto.AccessTokenDto) apigen.LoginRes {
	return apigen.LoginRes{
		AccessToken: token.AccessToken,
		TokenType:   token.TokenType,
		ExpiresIn:   int(token.ExpiresIn),
	}
}

func ToUserRegisterDto(payload apigen.RegisterJSONRequestBody) dto.UserRegisterDto {
	return dto.UserRegisterDto{
		Name:            payload.Name,
		Email:           payload.Email,
		Password:        payload.Password,
		ConfirmPassword: payload.ConfirmPassword,
	}
}

func ToUserResponse(user *domain.User) apigen.User {
	return apigen.User{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
