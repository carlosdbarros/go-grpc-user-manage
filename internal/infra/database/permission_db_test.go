package database

import (
	"database/sql"
	"fmt"
	permissionDomain "github.com/carlosdbarros/go-grpc-user-manage/internal/domain/permission"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type PermissionDBTestSuite struct {
	suite.Suite
	db  *sql.DB
	sut permissionDomain.PermissionRepository
}

func (suite *PermissionDBTestSuite) SetupTest() {
	var err error
	suite.db, err = initSqliteInMemory()
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.sut = NewPermissionDB(suite.db)
}

func (suite *PermissionDBTestSuite) TearDownTest() {
	suite.db.Close()
}

func TestSuitePermissionDB(t *testing.T) {
	suite.Run(t, new(PermissionDBTestSuite))
}

func (suite *PermissionDBTestSuite) TestPermissionDB_AddPermission_ShouldAddPermissionToDatabase() {
	var (
		err             error
		stmt            *sql.Stmt
		foundPermission permissionDomain.Permission
	)
	permission := makePermission("", "")

	permission, err = suite.sut.AddPermission(permission)
	assert.Nil(suite.T(), err)

	stmt, err = suite.db.Prepare("select id, codename, name from permissions where id = $1")
	assert.Nil(suite.T(), err)
	row := stmt.QueryRow(permission.ID)
	err = row.Scan(&foundPermission.ID, &foundPermission.Codename, &foundPermission.Name)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), permission.ID, foundPermission.ID)
	assert.Equal(suite.T(), permission.Codename, foundPermission.Codename)
	assert.Equal(suite.T(), permission.Name, foundPermission.Name)
}

func (suite *PermissionDBTestSuite) TestPermissionDB_FindPermissionById_ShouldFindPermissionById() {
	var (
		err             error
		stmt            *sql.Stmt
		foundPermission permissionDomain.Permission
	)

	permission := makePermission("", "")
	permission, err = suite.sut.AddPermission(permission)
	assert.Nil(suite.T(), err)

	stmt, err = suite.db.Prepare("select id, codename, name from permissions where id = $1")
	assert.Nil(suite.T(), err)
	row := stmt.QueryRow(permission.ID)
	err = row.Scan(&foundPermission.ID, &foundPermission.Codename, &foundPermission.Name)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), permission.ID, foundPermission.ID)
	assert.Equal(suite.T(), permission.Codename, foundPermission.Codename)
	assert.Equal(suite.T(), permission.Name, foundPermission.Name)
}

func (suite *PermissionDBTestSuite) TestPermissionDB_DeletePermission_ShouldDeletePermissionFromDatabase() {
	permission := makePermission("", "")
	_, err := suite.sut.AddPermission(permission)
	assert.Nil(suite.T(), err)

	err = suite.sut.DeletePermission(permission.ID)
	assert.Nil(suite.T(), err)

	result, err := suite.sut.FindPermissionById(permission.ID)
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), result)
}

func (suite *PermissionDBTestSuite) TestPermissionDB_FindAllPermissions_ShouldFindAllPermissions() {
	var (
		err         error
		permissions []*permissionDomain.Permission
	)

	permissionOne := makePermission("", "")
	permissionOne, err = suite.sut.AddPermission(permissionOne)
	assert.Nil(suite.T(), err)

	permissionTwo := makePermission("", "")
	permissionTwo, err = suite.sut.AddPermission(permissionTwo)
	assert.Nil(suite.T(), err)

	permissions, err = suite.sut.FindAllPermissions()
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 2, len(permissions))

	for k, v := range permissions {
		assert.Equal(suite.T(), v.ID, permissions[k].ID)
		assert.Equal(suite.T(), v.Codename, permissions[k].Codename)
		assert.Equal(suite.T(), v.Name, permissions[k].Name)
	}
}

func makePermission(name, codename string) *permissionDomain.Permission {
	if name == "" {
		name = faker.Word()
	}
	if codename == "" {
		codename = fmt.Sprintf("%s.%s", faker.Word(), faker.Word())
	}
	return permissionDomain.NewPermission(name, codename)
}
