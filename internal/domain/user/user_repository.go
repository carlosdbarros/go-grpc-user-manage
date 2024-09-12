package user

import "errors"

type UserRepository interface {
	AddUser(user *User) (*User, error)
	FindUserByEmail(email string) (*User, error)
	FindAllUsers() ([]*User, error)
}

var ErrUserNotFound = errors.New("user not found")
