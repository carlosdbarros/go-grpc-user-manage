package database

import (
	"database/sql"
	"github.com/carlosdbarros/go-grpc-user-manage/internal/entity"
	"github.com/carlosdbarros/go-grpc-user-manage/internal/repository"
	"github.com/go-faker/faker/v4"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UserDBTestSuite struct {
	suite.Suite
	db   *sql.DB
	repo repository.UserRepository
}

func (suite *UserDBTestSuite) SetupTest() {
	var err error
	suite.db, err = initSqliteInMemory()
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.repo = NewUserDBRepository(suite.db)
}

func (suite *UserDBTestSuite) TearDownTest() {
	suite.db.Close()
}

func TestSuiteUserDB(t *testing.T) {
	suite.Run(t, new(UserDBTestSuite))
}

func (suite *UserDBTestSuite) TestUserDBRepository_AddUser_ShouldAddUserToDatabase() {
	var (
		err       error
		stmt      *sql.Stmt
		foundUser entity.User
	)
	user := makeUser(suite.T(), "", "", "")

	user, err = suite.repo.AddUser(user)
	assert.Nil(suite.T(), err)

	stmt, err = suite.db.Prepare("select id, name, email, password from users where id = $1")
	assert.Nil(suite.T(), err)
	row := stmt.QueryRow(user.ID)
	err = row.Scan(&foundUser.ID, &foundUser.Name, &foundUser.Email, &foundUser.Password)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), user.ID, foundUser.ID)
	assert.Equal(suite.T(), user.Name, foundUser.Name)
	assert.Equal(suite.T(), user.Email, foundUser.Email)
	assert.NotEmpty(suite.T(), foundUser.Password)
	assert.Equal(suite.T(), user.Password, foundUser.Password)
}

func (suite *UserDBTestSuite) TestUserDBRepository_FindUserByEmail_ShouldFindUserByEmail() {
	var (
		err       error
		foundUser *entity.User
	)
	email := faker.Email()
	user := makeUser(suite.T(), "", email, "")

	user, err = suite.repo.AddUser(user)
	assert.Nil(suite.T(), err)
	foundUser, err = suite.repo.FindUserByEmail(email)
	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), user.ID, foundUser.ID)
	assert.Equal(suite.T(), user.Name, foundUser.Name)
	assert.Equal(suite.T(), user.Email, foundUser.Email)
	assert.NotEmpty(suite.T(), foundUser.Password)
	assert.Equal(suite.T(), user.Password, foundUser.Password)
}

func makeUser(t *testing.T, name, email, password string) *entity.User {
	if name == "" {
		name = faker.Name()
	}
	if email == "" {
		email = faker.Email()
	}
	if password == "" {
		password = faker.Password()
	}
	user, err := entity.NewUser(name, email, password)
	if err != nil {
		t.Fatal(err)
	}
	return user
}
