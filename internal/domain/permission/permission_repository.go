package permission

type PermissionRepository interface {
	AddPermission(permission *Permission) (*Permission, error)
	FindPermissionById(id string) (*Permission, error)
	DeletePermission(id string) error
	FindAllPermissions() ([]*Permission, error)
}
