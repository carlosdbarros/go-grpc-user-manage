package entity

import "github.com/google/uuid"

type Permission struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Codename string `json:"codename"`
}

func NewPermission(name, codename string) *Permission {
	return &Permission{
		ID:       uuid.New().String(),
		Name:     name,
		Codename: codename,
	}
}
