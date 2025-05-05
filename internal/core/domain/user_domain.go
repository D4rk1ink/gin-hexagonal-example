package domain

import (
	"regexp"
	"time"

	err "github.com/D4rk1ink/gin-hexagonal-example/internal/core/error"
)

type User struct {
	ID        string
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(name, email, password string) (*User, error) {

	err := validateEmail(email)
	if err != nil {
		return nil, err
	}

	return &User{
		Name:      name,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func validateEmail(email string) error {
	regexpEmail := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !regexpEmail.MatchString(email) {
		return err.NewError(err.ErrAuthInvalidEmailFormat, nil)
	}
	return nil
}

func (u *User) SetId(id string) {
	u.ID = id
}

func (u *User) Update(name, email string) {
	u.Name = name
	u.Email = email
	u.UpdatedAt = time.Now()
}
