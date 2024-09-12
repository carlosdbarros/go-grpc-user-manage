package database

import (
	"database/sql"
	domainUser "github.com/carlosdbarros/go-grpc-user-manage/internal/domain/user"
	"github.com/go-faker/faker/v4"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UserDBTestSuite struct {
	suite.Suite
	db  *sql.DB
	sut domainUser.UserRepository
}

func (suite *UserDBTestSuite) SetupTest() {
	var err error
	suite.db, err = initSqliteInMemory()
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.sut = NewUserDB(suite.db)
}

func (suite *UserDBTestSuite) TearDownTest() {
	suite.db.Close()
}

func TestSuiteUserDB(t *testing.T) {
	suite.Run(t, new(UserDBTestSuite))
}

func (suite *UserDBTestSuite) TestUserDB_AddUser_ShouldAddUserToDatabase() {
	var (
		stmt          *sql.Stmt
		foundInstance domainUser.User
	)

	user, err := suite.sut.AddUser(makeUser(suite.T(), "", "", ""))
	assert.Nil(suite.T(), err)

	stmt, err = suite.db.Prepare("select id, name, email, password from users where id = $1")
	assert.Nil(suite.T(), err)
	row := stmt.QueryRow(user.ID)
	err = row.Scan(&foundInstance.ID, &foundInstance.Name, &foundInstance.Email, &foundInstance.Password)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), user.ID, foundInstance.ID)
	assert.Equal(suite.T(), user.Name, foundInstance.Name)
	assert.Equal(suite.T(), user.Email, foundInstance.Email)
	assert.NotEmpty(suite.T(), foundInstance.Password)
	assert.Equal(suite.T(), user.Password, foundInstance.Password)
}

func (suite *UserDBTestSuite) TestUserDB_FindUserByEmail_ShouldFindUserByEmail() {
	var (
		err           error
		foundInstance *domainUser.User
	)
	email := faker.Email()
	user := makeUser(suite.T(), "", email, "")

	user, err = suite.sut.AddUser(user)
	assert.Nil(suite.T(), err)
	foundInstance, err = suite.sut.FindUserByEmail(email)
	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), user.ID, foundInstance.ID)
	assert.Equal(suite.T(), user.Name, foundInstance.Name)
	assert.Equal(suite.T(), user.Email, foundInstance.Email)
	assert.NotEmpty(suite.T(), foundInstance.Password)
	assert.Equal(suite.T(), user.Password, foundInstance.Password)
}

func (suite *UserDBTestSuite) TestUserDB_FindAllUsers_ShouldFindAllUsers() {
	var (
		err              error
		userOne, userTwo *domainUser.User
		users            []*domainUser.User
	)

	userOne = makeUser(suite.T(), "", "", "")
	userOne, err = suite.sut.AddUser(userOne)
	assert.Nil(suite.T(), err)

	userTwo = makeUser(suite.T(), "", "", "")
	userTwo, err = suite.sut.AddUser(userTwo)
	assert.Nil(suite.T(), err)

	users, err = suite.sut.FindAllUsers()
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 2, len(users))

	for k, v := range users {
		assert.Equal(suite.T(), v.ID, users[k].ID)
		assert.Equal(suite.T(), v.Name, users[k].Name)
		assert.Equal(suite.T(), v.Email, users[k].Email)
		assert.NotEmpty(suite.T(), users[k].Password)
	}
}

func makeUser(t *testing.T, name, email, password string) *domainUser.User {
	if name == "" {
		name = faker.Name()
	}
	if email == "" {
		email = faker.Email()
	}
	if password == "" {
		password = faker.Password()
	}
	user, err := domainUser.NewUser(name, email, password)
	if err != nil {
		t.Fatal(err)
	}
	return user
}
