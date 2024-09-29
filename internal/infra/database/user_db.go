package database

import (
	"database/sql"
	userDomain "github.com/carlosdbarros/go-grpc-user-manage/internal/domain/user"
)

type UserDB struct {
	DB *sql.DB
}

func NewUserDB(db *sql.DB) *UserDB {
	return &UserDB{DB: db}
}

func (repo *UserDB) AddUser(input *userDomain.User) (*userDomain.User, error) {
	//stmt, err := repo.DB.Prepare("insert into users(id, name, email, password) values ($1, $2, $3, $4)")
	//if err != nil {
	//	return nil, err
	//}
	//_, err = stmt.Exec(input.ID, input.Name, input.Email, input.Password)
	//if err != nil {
	//	return nil, err
	//}
	return input, nil
}

func (repo *UserDB) FindUserByEmail(email string) (*userDomain.User, error) {
	stmt, err := repo.DB.Prepare("select id, name, email, password from users where email = $1")
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRow(email)
	user := &userDomain.User{}
	err = row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *UserDB) FindAllUsers() ([]*userDomain.User, error) {
	rows, err := repo.DB.Query("select id, name, email, password from users")
	if err != nil {
		return nil, err
	}
	var users []*userDomain.User
	for rows.Next() {
		user := &userDomain.User{}
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
