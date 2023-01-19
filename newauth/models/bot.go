package models

import "gorm.io/gorm"

type Bot struct {
	gorm.Model
	Name string
	Desc string
}
