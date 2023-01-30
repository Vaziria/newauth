package models

import "time"

type BotToken struct {
	ID uint `gorm:"primarykey" json:"id"`

	BotID    uint  `json:"bot_id"`
	DeviceID *uint `json:"device_id"`
	UserID   uint  `json:"user_id"`
	TeamID   uint  `json:"team_id"`

	secretPwd string

	LastLog   time.Time `json:"last_login"`
	Device    *Device   `json:"device"`
	CreatedAt time.Time `json:"created_at"`

	User User `json:"user"`
	Bot  Bot  `json:"bot"`
	Team Team `json:"team"`
}

func (t *BotToken) CheckPwd(pwd string) bool {
	return t.secretPwd == HashPassword(pwd)
}

func (t *BotToken) SetPwd(pwd string) {
	t.secretPwd = HashPassword(pwd)
}
