package hash

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -package mock_hash -source=hash.go -destination=mock/mock_hash.go *
type Hash interface {
	HashPassword(ctx context.Context, password string) (*string, error)
	ComparePassword(ctx context.Context, password, hashedPassword string) error
}

type hash struct {
}

func NewHash() Hash {
	return &hash{}
}

func (s *hash) HashPassword(ctx context.Context, password string) (*string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	hashedPassword := string(hashed)

	return &hashedPassword, nil
}

func (s *hash) ComparePassword(ctx context.Context, password string, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}

	return nil
}
