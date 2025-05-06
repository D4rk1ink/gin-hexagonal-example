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
	GenerateTokenWithOptions(input *GenerateTokenInput, options *GenerateTokenOptions) (*string, *int64, error)
	ValidateAccessToken(token string) (*TokenPayload, error)
	ValidateTokenWithOptions(token string, options *ValidateTokenOptions) (*TokenPayload, error)
}

type jwt struct{}

type GenerateTokenInput struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type GenerateTokenOptions struct {
	Secret   string
	Duration string
}

type ValidateTokenOptions struct {
	Secret string
}

type TokenPayload struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Exp   int64  `json:"exp"`
	_jwt.RegisteredClaims
}

func NewJwt() Jwt {
	return &jwt{}
}

func (j *jwt) GenerateAccessToken(input *GenerateTokenInput) (*string, *int64, error) {
	return j.GenerateTokenWithOptions(input, &GenerateTokenOptions{
		Secret:   config.Config.Jwt.Secret,
		Duration: config.Config.Jwt.ExpiresIn,
	})
}

func (j *jwt) GenerateTokenWithOptions(input *GenerateTokenInput, options *GenerateTokenOptions) (*string, *int64, error) {
	duration, err := util.ParseDurationToSeconds(options.Duration)
	if err != nil {
		return nil, nil, err
	}

	payload := TokenPayload{
		ID:    input.ID,
		Email: input.Email,
		Exp:   time.Now().Add(time.Duration(duration) * time.Second).Unix(),
	}
	token := _jwt.NewWithClaims(_jwt.SigningMethodHS256, payload)
	tokenString, err := token.SignedString([]byte(options.Secret))
	if err != nil {
		return nil, nil, err
	}

	return &tokenString, &duration, nil
}

func (j *jwt) ValidateAccessToken(tokenStr string) (*TokenPayload, error) {
	return j.ValidateTokenWithOptions(tokenStr, &ValidateTokenOptions{
		Secret: config.Config.Jwt.Secret,
	})
}

func (j *jwt) ValidateTokenWithOptions(tokenStr string, options *ValidateTokenOptions) (*TokenPayload, error) {
	token, err := _jwt.ParseWithClaims(tokenStr, &TokenPayload{}, func(_token *_jwt.Token) (interface{}, error) {
		return []byte(options.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenPayload); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
