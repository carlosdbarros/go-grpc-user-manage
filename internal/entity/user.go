package entity

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

func NewUser(name, email, password string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &User{
		ID:       uuid.NewString(),
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
	}
	err = user.Validate()
	if err != nil {
		return nil, err
	}
	return user, nil
}

var (
	ErrNameIsRequired     = errors.New("name is required")
	ErrEmailIsRequired    = errors.New("email is required")
	ErrPasswordIsRequired = errors.New("password is required")
)

func (u *User) Validate() error {
	if u.Name == "" {
		return ErrNameIsRequired
	}
	if u.Email == "" {
		return ErrEmailIsRequired
	}
	if u.Password == "" {
		return ErrPasswordIsRequired
	}
	return nil
}

func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
