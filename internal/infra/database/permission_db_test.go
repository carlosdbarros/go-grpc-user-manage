package database

import (
	"database/sql"
	"fmt"
	"github.com/carlosdbarros/go-grpc-user-manage/internal/domain/permission"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type PermissionDBTestSuite struct {
	suite.Suite
	db   *sql.DB
	repo permission.PermissionRepository
}

func (suite *PermissionDBTestSuite) SetupTest() {
	var err error
	suite.db, err = initSqliteInMemory()
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.repo = NewPermissionDBRepository(suite.db)
}

func (suite *PermissionDBTestSuite) TearDownTest() {
	suite.db.Close()
}

func TestSuitePermissionDB(t *testing.T) {
	suite.Run(t, new(PermissionDBTestSuite))
}

func (suite *PermissionDBTestSuite) TestPermissionDBRepository_AddPermission_ShouldAddPermissionToDatabase() {
	var (
		err             error
		stmt            *sql.Stmt
		foundPermission permission.Permission
	)
	permission := makePermission("", "")

	permission, err = suite.repo.AddPermission(permission)
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

func (suite *PermissionDBTestSuite) TestPermissionDBRepository_FindPermissionById_ShouldFindPermissionById() {
	var (
		err             error
		stmt            *sql.Stmt
		foundPermission permission.Permission
	)

	permission := makePermission("", "")
	permission, err = suite.repo.AddPermission(permission)
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

func (suite *PermissionDBTestSuite) TestPermissionDBRepository_DeletePermission_ShouldDeletePermissionFromDatabase() {
	var (
		err        error
		stmt       *sql.Stmt
		permission *permission.Permission
	)

	permission = makePermission("", "")
	_, err = suite.repo.AddPermission(permission)
	assert.Nil(suite.T(), err)

	err = suite.repo.DeletePermission(permission.ID)
	assert.Nil(suite.T(), err)

	stmt, err = suite.db.Prepare("select id, codename, name from permissions where id = $1")
	assert.Nil(suite.T(), err)
	row := stmt.QueryRow(permission.ID)
	err = row.Scan(&permission.ID, &permission.Codename, &permission.Name)
	assert.NotNil(suite.T(), err)
}

func (suite *PermissionDBTestSuite) TestPermissionDBRepository_FindAllPermissions_ShouldFindAllPermissions() {
	var (
		err         error
		permissions []*permission.Permission
	)

	permission1 := makePermission("", "")
	permission2 := makePermission("", "")
	permission1, err = suite.repo.AddPermission(permission1)
	assert.Nil(suite.T(), err)
	permission2, err = suite.repo.AddPermission(permission2)
	assert.Nil(suite.T(), err)

	permissions, err = suite.repo.FindAllPermissions()
	assert.Nil(suite.T(), err)

	assert.Equal(suite.T(), 2, len(permissions))

	assert.Equal(suite.T(), permission1.ID, permissions[0].ID)
	assert.Equal(suite.T(), permission1.Codename, permissions[0].Codename)
	assert.Equal(suite.T(), permission1.Name, permissions[0].Name)

	assert.Equal(suite.T(), permission2.ID, permissions[1].ID)
	assert.Equal(suite.T(), permission2.Codename, permissions[1].Codename)
	assert.Equal(suite.T(), permission2.Name, permissions[1].Name)
}

func makePermission(name, codename string) *permission.Permission {
	if name == "" {
		name = faker.Word()
	}
	if codename == "" {
		codename = fmt.Sprintf("%s.%s", faker.Word(), faker.Word())
	}
	return permission.NewPermission(name, codename)
}
