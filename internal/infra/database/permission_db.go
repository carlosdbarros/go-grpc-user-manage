package database

import (
	"database/sql"
	"github.com/carlosdbarros/go-grpc-user-manage/internal/entity"
)

type PermissionDBRepository struct {
	DB *sql.DB
}

func NewPermissionDBRepository(db *sql.DB) *PermissionDBRepository {
	return &PermissionDBRepository{DB: db}
}

func (repo *PermissionDBRepository) AddPermission(permission *entity.Permission) (*entity.Permission, error) {
	stmt, err := repo.DB.Prepare("insert into permissions(id, codename, name) values ($1, $2, $3)")
	if err != nil {
		return nil, err
	}
	_, err = stmt.Exec(permission.ID, permission.Codename, permission.Name)
	if err != nil {
		return nil, err
	}
	return permission, nil
}
