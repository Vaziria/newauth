package models

type Quota struct {
	ID     uint `gorm:"primarykey"`
	TeamID uint `gorm:"uniqueIndex:quota_unique"`
	BotID  uint `gorm:"uniqueIndex:quota_unique"`
	Count  int  `json:"count"`
	Limit  int  `json:"limit"`

	Team *Team `json:"team"`
	Bot  *Bot  `json:"bot"`
}

func (m *Quota) TableName() string {
	return "quotas"
}
