package repository

import "github.com/carlosdbarros/go-grpc-user-manage/internal/entity"

type PermissionRepository interface {
	AddPermission(permission *entity.Permission) (*entity.Permission, error)
	//FindPermissionByCodename(codename string) (*entity.Permission, error)
	//FindAllPermissions() ([]*entity.Permission, error)
	//DeletePermission(id string) error
}
