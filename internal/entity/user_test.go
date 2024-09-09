package entity

import (
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UserTestSuite struct {
	suite.Suite
	sut      *User
	err      error
	name     string
	email    string
	password string
}

func (suite *UserTestSuite) SetupTest() {
	suite.name = faker.Name()
	suite.email = faker.Email()
	suite.password = faker.Password()
	suite.sut, suite.err = NewUser(suite.name, suite.email, suite.password)
}

func (suite *UserTestSuite) TearDownTest() {
}

func (suite *UserTestSuite) TestNewUser_ShouldCreateANewUserWithCorrectParams() {
	assert.Nil(suite.T(), suite.err)
	assert.NotNil(suite.T(), suite.sut)
	assert.NotEmpty(suite.T(), suite.sut.ID)
	assert.NotEmpty(suite.T(), suite.sut.Password)
	assert.Equal(suite.T(), suite.name, suite.sut.Name)
	assert.Equal(suite.T(), suite.email, suite.sut.Email)
}

func (suite *UserTestSuite) TestNewUser_ShouldCreateANewUserWithValidPassword() {
	assert.Nil(suite.T(), suite.err)
	assert.NotNil(suite.T(), suite.sut)
	assert.True(suite.T(), suite.sut.ValidatePassword(suite.password))
	assert.False(suite.T(), suite.sut.ValidatePassword("123"))
	assert.NotEqual(suite.T(), suite.password, suite.sut.Password)
}

func TestSuiteUser(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
