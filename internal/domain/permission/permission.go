package permission

import (
	"github.com/google/uuid"
)

type Permission struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Codename string `json:"codename"`
}

func NewPermission(name, codename string) *Permission {
	return &Permission{
		ID:       uuid.NewString(),
		Name:     name,
		Codename: codename,
	}
}
