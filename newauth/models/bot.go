package models

type Bot struct {
	ID   uint   `gorm:"primarykey"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}
