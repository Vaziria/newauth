package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name string

	Permissions []*Permission `gorm:"many2many:role_permission;"`
}

type Permission struct {
	ID  uint   `gorm:"primarykey"`
	Key string `gorm:"uniqueIndex"`
	// Role []*Role `gorm:"many2many:role_accesses;"`
}

func (*Permission) TableName() string {
	return "permissions"
}

type UserPermission struct {
	UserID       uint
	PermissionID uint
}
