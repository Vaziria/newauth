package models

import (
	"github.com/PDC-Repository/newauth/newauth/authorize"
)

type Team struct {
	ID          uint    `gorm:"primarykey" json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Users       []*User `gorm:"many2many:user_teams;" json:"users"`
}

func (user *Team) TableName() string {
	return "teams"
}

type UserTeam struct {
	UserID uint `gorm:"index:unique_user,unique"`
	TeamID uint `gorm:"index:unique_user,unique"`
	Role   authorize.RoleEnum
}

func (user *UserTeam) TableName() string {
	return "user_teams"
}
