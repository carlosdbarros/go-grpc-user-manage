package database

import (
	"database/sql"
	permissionDomain "github.com/carlosdbarros/go-grpc-user-manage/internal/domain/permission"
)

type PermissionDBR struct {
	DB *sql.DB
}

func NewPermissionDB(db *sql.DB) *PermissionDBR {
	return &PermissionDBR{DB: db}
}

func (repo *PermissionDBR) AddPermission(input *permissionDomain.Permission) (*permissionDomain.Permission, error) {
	stmt, err := repo.DB.Prepare("insert into permissions(id, codename, name) values ($1, $2, $3)")
	if err != nil {
		return nil, err
	}
	_, err = stmt.Exec(input.ID, input.Codename, input.Name)
	if err != nil {
		return nil, err
	}
	return input, nil
}

func (repo *PermissionDBR) FindPermissionById(id string) (*permissionDomain.Permission, error) {
	stmt, err := repo.DB.Prepare("select id, codename, name from permissions where id = $1")
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRow(id)
	permission := &permissionDomain.Permission{}
	err = row.Scan(&permission.ID, &permission.Codename, &permission.Name)
	if err != nil {
		return nil, err
	}
	return permission, nil
}

func (repo *PermissionDBR) DeletePermission(id string) error {
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

func (repo *PermissionDBR) FindAllPermissions() ([]*permissionDomain.Permission, error) {
	rows, err := repo.DB.Query("select id, codename, name from permissions")
	if err != nil {
		return nil, err
	}
	var permissions []*permissionDomain.Permission
	for rows.Next() {
		permission := &permissionDomain.Permission{}
		err = rows.Scan(&permission.ID, &permission.Codename, &permission.Name)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	return permissions, nil
}
