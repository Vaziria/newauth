package models

import "gorm.io/gorm"

type Device struct {
	gorm.Model

	Hostname string

	Platform  string
	Suspended bool
}

type Interface struct {
	gorm.Model

	Mac  string
	Addr string
}
