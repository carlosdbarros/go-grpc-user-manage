package entity

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type PermissionTestSuite struct {
	suite.Suite
	sut      *Permission
	name     string
	codename string
}

func (suite *PermissionTestSuite) SetupTest() {
	suite.name = "Criar Todo"
	suite.codename = "todo.add"
	suite.sut = NewPermission(suite.name, suite.codename)
}

func (suite *PermissionTestSuite) TearDownTest() {
}

func (suite *PermissionTestSuite) TestPermission_NewPermission_ShouldCreateNewPermission() {
	assert.NotNil(suite.T(), suite.sut)
	assert.NotEmpty(suite.T(), suite.sut.ID)
	assert.Equal(suite.T(), suite.name, suite.sut.Name)
	assert.Equal(suite.T(), suite.codename, suite.sut.Codename)
}

func TestSuitePermission(t *testing.T) {
	suite.Run(t, new(PermissionTestSuite))
}
