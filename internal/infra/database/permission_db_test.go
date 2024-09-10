package database

import (
	"database/sql"
	"fmt"
	"github.com/carlosdbarros/go-grpc-user-manage/internal/entity"
	"github.com/carlosdbarros/go-grpc-user-manage/internal/repository"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type PermissionDBTestSuite struct {
	suite.Suite
	db   *sql.DB
	repo repository.PermissionRepository
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
		foundPermission entity.Permission
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
		foundPermission entity.Permission
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
		permission *entity.Permission
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

func makePermission(name, codename string) *entity.Permission {
	if name == "" {
		name = faker.Word()
	}
	if codename == "" {
		codename = fmt.Sprintf("%s.%s", faker.Word(), faker.Word())
	}
	return entity.NewPermission(name, codename)
}
