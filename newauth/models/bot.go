package models

type Bot struct {
	ID   uint   `gorm:"primarykey" json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`

	Quotas []Quota `gorm:"constraint:OnUpdate:CASCADE, OnDelete:CASCADE, foreignKey:BotID;"`
}
