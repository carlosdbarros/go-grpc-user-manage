package database

import (
	"database/sql"
	"github.com/carlosdbarros/go-grpc-user-manage/internal/entity"
)

type UserRepository interface {
	AddUser(user *entity.User) (*entity.User, error)
	FindUserByEmail(email string) (*entity.User, error)
}

type UserDBRepository struct {
	DB *sql.DB
}

func NewUserDBRepository(db *sql.DB) *UserDBRepository {
	return &UserDBRepository{DB: db}
}

func (repo *UserDBRepository) AddUser(input *entity.User) (*entity.User, error) {
	stmt, err := repo.DB.Prepare("insert into users(id, name, email, password) values ($1, $2, $3, $4)")
	if err != nil {
		return nil, err
	}
	_, err = stmt.Exec(input.ID, input.Name, input.Email, input.Password)
	if err != nil {
		return nil, err
	}
	return input, nil
}

func (repo *UserDBRepository) FindUserByEmail(email string) (*entity.User, error) {
	stmt, err := repo.DB.Prepare("select id, name, email, password from users where email = $1")
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRow(email)
	user := &entity.User{}
	err = row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
