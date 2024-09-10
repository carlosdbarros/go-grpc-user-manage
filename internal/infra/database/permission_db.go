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

func (repo *PermissionDBRepository) FindPermissionById(id string) (*entity.Permission, error) {
	stmt, err := repo.DB.Prepare("select id, codename, name from permissions where id = $1")
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRow(id)
	permission := &entity.Permission{}
	err = row.Scan(&permission.ID, &permission.Codename, &permission.Name)
	if err != nil {
		return nil, err
	}
	return permission, nil
}

func (repo *PermissionDBRepository) DeletePermission(id string) error {
	stmt, err := repo.DB.Prepare("delete from permissions where id = $1")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func (repo *PermissionDBRepository) FindAllPermissions() ([]*entity.Permission, error) {
	rows, err := repo.DB.Query("select id, codename, name from permissions")
	if err != nil {
		return nil, err
	}
	var permissions []*entity.Permission
	for rows.Next() {
		permission := &entity.Permission{}
		err = rows.Scan(&permission.ID, &permission.Codename, &permission.Name)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	return permissions, nil
}
