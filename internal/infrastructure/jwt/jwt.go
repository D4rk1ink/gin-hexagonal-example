package jwt

import (
	"errors"
	"time"

	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/config"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/util"
	_jwt "github.com/golang-jwt/jwt/v5"
)

type Jwt interface {
	GenerateAccessToken(input *GenerateTokenInput) (*string, *int64, error)
	ValidateAccessToken(token string) (*TokenPayload, error)
	ParseAccessToken(token string) (string, error)
}

type jwt struct{}

type GenerateTokenInput struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type TokenPayload struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Exp   int64  `json:"exp"`
	_jwt.Claims
}

func NewJwt() Jwt {
	return &jwt{}
}

func (j *jwt) GenerateAccessToken(input *GenerateTokenInput) (*string, *int64, error) {
	duration, err := util.ParseDurationToSeconds(config.Config.Jwt.ExpiresIn)
	if err != nil {
		return nil, nil, err
	}

	payload := TokenPayload{
		ID:    input.ID,
		Email: input.Email,
		Exp:   time.Now().Add(time.Duration(duration) * time.Second).Unix(),
	}
	token := _jwt.NewWithClaims(_jwt.SigningMethodHS256, payload)
	tokenString, err := token.SignedString([]byte(config.Config.Jwt.SecretKey))
	if err != nil {
		return nil, nil, err
	}

	return &tokenString, &duration, nil
}

func (j *jwt) ValidateAccessToken(tokenStr string) (*TokenPayload, error) {
	token, err := _jwt.ParseWithClaims(tokenStr, &TokenPayload{}, func(_token *_jwt.Token) (interface{}, error) {
		return []byte(config.Config.Jwt.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenPayload); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func (j *jwt) ParseAccessToken(token string) (string, error) {
	return "", nil
}
