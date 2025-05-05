package jwt

import (
	"time"

	"encoding/json"

	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/config"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/util"
	_jwt "github.com/golang-jwt/jwt/v5"
)

type Jwt interface {
	GenerateAccessToken(input *GenerateTokenInput) (*string, *int64, error)
	ValidateAccessToken(token string) (string, error)
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
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, err
	}
	var claims _jwt.MapClaims
	err = json.Unmarshal(b, &claims)
	if err != nil {
		return nil, nil, err
	}

	println("config.Config.Jwt.SecretKey", config.Config.Jwt.SecretKey)

	token := _jwt.NewWithClaims(_jwt.SigningMethodHS256, _jwt.MapClaims{
		"exp": "ss",
		"iat": time.Now().Unix(),
	})
	tokenString, err := token.SignedString([]byte(config.Config.Jwt.SecretKey))
	if err != nil {
		return nil, nil, err
	}

	return &tokenString, &duration, nil
}

func (j *jwt) ValidateAccessToken(token string) (string, error) {
	return "", nil
}

func (j *jwt) ParseAccessToken(token string) (string, error) {
	return "", nil
}
