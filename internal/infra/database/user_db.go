package database

import (
	"database/sql"
	"github.com/carlosdbarros/go-grpc-user-manage/internal/domain/user"
)

type UserDBRepository struct {
	DB *sql.DB
}

func NewUserDBRepository(db *sql.DB) *UserDBRepository {
	return &UserDBRepository{DB: db}
}

func (repo *UserDBRepository) AddUser(input *user.User) (*user.User, error) {
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

func (repo *UserDBRepository) FindUserByEmail(email string) (*user.User, error) {
	stmt, err := repo.DB.Prepare("select id, name, email, password from users where email = $1")
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRow(email)
	user := &user.User{}
	err = row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *UserDBRepository) FindAllUsers() ([]*user.User, error) {
	rows, err := repo.DB.Query("select id, name, email, password from users")
	if err != nil {
		return nil, err
	}
	var users []*user.User
	for rows.Next() {
		user := &user.User{}
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
