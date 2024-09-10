package repository

import "github.com/carlosdbarros/go-grpc-user-manage/internal/entity"

type UserRepository interface {
	AddUser(user *entity.User) (*entity.User, error)
	FindUserByEmail(email string) (*entity.User, error)
	FindAllUsers() ([]*entity.User, error)
}
