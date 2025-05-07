package http_mapper

import (
	http_apigen "github.com/D4rk1ink/gin-hexagonal-example/internal/application/handler/http/apigen"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/domain"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/dto"
)

func ToCredentialDto(payload http_apigen.LoginJSONRequestBody) dto.CredentialDto {
	return dto.CredentialDto{
		Email:    string(payload.Email),
		Password: payload.Password,
	}
}

func ToAccessTokenResponse(token *dto.AccessTokenDto) http_apigen.LoginRes {
	return http_apigen.LoginRes{
		AccessToken: token.AccessToken,
		TokenType:   token.TokenType,
		ExpiresIn:   int(token.ExpiresIn),
	}
}

func ToUserRegisterDto(payload http_apigen.RegisterJSONRequestBody) dto.UserRegisterDto {
	return dto.UserRegisterDto{
		Name:            payload.Name,
		Email:           string(payload.Email),
		Password:        payload.Password,
		ConfirmPassword: payload.ConfirmPassword,
	}
}

func ToUserResponse(user *domain.User) http_apigen.User {
	return http_apigen.User{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
