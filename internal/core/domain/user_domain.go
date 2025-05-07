package domain

import (
	"regexp"
	"time"

	err "github.com/D4rk1ink/gin-hexagonal-example/internal/core/error"
	time_util "github.com/D4rk1ink/gin-hexagonal-example/internal/util/time"
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
		CreatedAt: time_util.Now(),
		UpdatedAt: time_util.Now(),
	}, nil
}

func validateEmail(email string) error {
	regexpEmail := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !regexpEmail.MatchString(email) {
		return err.NewError(err.ErrUserInvalidEmailFormat, nil)
	}
	return nil
}

func (u *User) SetName(name string) {
	u.Name = name
	u.UpdatedAt = time_util.Now()
}

func (u *User) SetEmail(email string) error {
	err := validateEmail(email)
	if err != nil {
		return err
	}

	u.Email = email
	u.UpdatedAt = time_util.Now()

	return nil
}
