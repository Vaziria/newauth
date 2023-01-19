package models

import "gorm.io/gorm"

type Team struct {
	gorm.Model
	Name        string
	Description string
	Users       []*User `gorm:"many2many:user_teams;"`
}
